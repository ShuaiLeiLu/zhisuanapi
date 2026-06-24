# PRD 审查报告：邀请返利与签到奖励金额化

> 审查对象：`docs/prd/invitation-rebate-checkin-amount-prd.md`
> 审查日期：2026-06-24

---

## 一、总体结论

PRD 设计思路清晰，三层奖励模式（注册/首充/每次充值）的抽象是合理的。但它存在 **6 个需要修正的关键问题**，其中 **支付逻辑侵入性** 和 **返利基数不统一** 是两个必须解决的核心风险。以下是逐条分析。

---

## 二、你关心的三个问题

### 问题 1：能不能完成"充值返利"机制？

**能实现，但有几个坑必须填。**

#### 1.1 返利基数 `TopUp.Money` 在不同渠道含义不一致（P0）

PRD 第 78 行说"返利基数为订单实际支付金额 `TopUp.Money`"。但当前代码中：

| 支付渠道 | 额度计算方式 | `Money` 字段含义 |
|----------|-------------|------------------|
| Stripe | `Money * QuotaPerUnit` | 已按分组倍率/折扣换算后的美元数 |
| Creem | 直接用 `Amount`（不使用 `Money`） | 原始支付金额（仅记日志用） |
| Waffo | `Amount * QuotaPerUnit` | 原始支付金额（仅记日志用） |
| Waffo Pancake | `Amount * QuotaPerUnit` | 原始支付金额（仅记日志用） |
| Epay | `Amount * QuotaPerUnit` | 原始支付金额（仅记日志用） |
| 管理员补单 | Stripe 用 `Money`，其他用 `Amount` | 见上 |

**风险**：以 `Money` 作为统一返利基数会导致 Stripe 渠道的返利金额偏高或偏低（取决于分组倍率），不同渠道用户拿到的返利不一致。

**建议修改 PRD**：返利基数改为各渠道**用户实际支付的原始金额**（即前端传给支付网关的金额），而非 `TopUp.Money`。这需要在 TopUp 表新增字段或在订单创建时额外存储。

#### 1.2 返利逻辑会侵入所有 6 个支付完成路径（P0）

当前支付完成分布在 6 个不同位置：

| # | 函数 | 文件 | 位置 |
|---|------|------|------|
| 1 | `Recharge()` | `model/topup.go:109` | Model 层 |
| 2 | `RechargeCreem()` | `model/topup.go:392` | Model 层 |
| 3 | `RechargeWaffo()` | `model/topup.go:467` | Model 层 |
| 4 | `RechargeWaffoPancake()` | `model/topup.go:530` | Model 层 |
| 5 | `EpayNotify()` | `controller/topup.go:310` | **Controller 层** |
| 6 | `ManualCompleteTopUp()` | `model/topup.go:320` | Model 层 |

PRD 的充值返利需要在每一个充值成功后触发。这意味着 **必须修改这 6 处核心支付代码**。最理想的实现方式是抽取统一的"充值成功后处理"回调/钩子，而不是在每个函数里内联返利逻辑。

**强烈建议**：不要在每个 `Recharge*` 函数里直接写返利代码。改为：
1. 创建一个 `AfterTopUpSuccessHook(topUp *TopUp)` 统一钩子
2. 在 6 个完成点调用该钩子
3. 钩子内部判读 `affiliate_reward_trigger` 配置决定是否发返利

这样可以最小化对原有支付代码的改动量。

#### 1.3 Epay 回调不在数据库事务内（P1）

`EpayNotify` 是 6 个路径中唯一不使用数据库事务的。它使用的是 Go 内存锁 `LockOrder`（基于 `sync.Map` + `sync.Mutex`），然后分两步：
1. `topUp.Update()` — 更新订单状态
2. `model.IncreaseUserQuota()` — 加余额

这两步不在同一个 DB 事务里。如果在其中插入返利逻辑（加 `aff_quota`），一旦失败，订单已标记成功但返利没发，数据就脏了。

**建议**：Epay 的 Webhook 回调也应重构为事务内处理，或者返利逻辑做成可补偿的重试任务。

#### 1.4 `first_topup` 模式缺少判断依据（P1）

PRD 第 79 行说"first_topup 只对受邀用户第一笔成功钱包充值发放"，但没有明确如何判断"第一笔"。

当前方案选项：
- **方案 A**：查询 topups 表，看该用户是否有过 status=success 的充值记录 → 性能差且存在时间窗口竞态
- **方案 B**：在 User 表新增 `has_received_first_topup_rebate` 布尔标记 → 推荐，清晰且可原子更新
- **方案 C**：利用 `aff_count` 字段 → 但语义冲突（见下文）

**建议**：采用方案 B，并在事务内原子标记。

#### 1.5 `aff_count` 语义变更存在兼容风险（P1）

当前 `aff_count` 在 `model/user.go:344` 的定义是"**邀请了的人数**"（注册时 +1）。PRD 第 81 行说 `every_topup` 模式下 `aff_count` 表示返利次数，每笔充值 +1。

这意味着：
- `aff_count` 的语义取决于 `affiliate_reward_trigger` 配置值
- 切换模式后，`aff_count` 的历史含义和新含义混在一起
- 前端如果展示了"邀请了 X 人"，在不同模式下含义完全不同

**建议**：新增独立字段如 `aff_rebate_count` 表示充值返利次数，保持 `aff_count` 语义不变。

---

### 问题 2：支付逻辑不要被修改

**PRD 的设计确实会修改支付逻辑**——这是必然的，因为充值返利必须挂钩到充值完成事件上。

但可以通过架构设计把影响控制在可接受范围内：

#### 2.1 理论上最小化修改的方案

```
充值成功 -> AfterTopUpSuccessHook (新增)
                |
                +-> 判断 affiliate_reward_trigger
                +-> 如果是 first_topup/every_topup
                        +-> 查找 inviter_id
                        +-> 计算返利
                        +-> 增加 inviter.AffQuota
                        +-> 增加 inviter.AffHistoryQuota
                        +-> 如果是 every_topup: inviter.AffRebateCount++
                        +-> 记录返利日志
```

支付完成函数只需加一行 `AfterTopUpSuccessHook(topUp)` 调用，逻辑集中在一个地方。

#### 2.2 事务边界建议

**不要把返利放进充值事务里面**。理由：
- 返利是附加逻辑，不应阻塞充值
- 返利失败不应回滚充值（用户充了钱，必须到账）
- 邀请人可能已被删除或禁用

**推荐方案**：
- 充值事务完成后，同步调用返利钩子
- 返利钩子使用独立事务
- 钩子内部幂等（按 trade_no + inviter_id 防重复发放）
- 极低概率的钩子失败由错误日志和手动补偿覆盖

#### 2.3 订阅购买已正确排除

PRD 第 23 行和第 83 行明确排除了订阅购买触发返利，这是正确的。当前 `CompleteSubscriptionOrder()`（`model/subscription.go:553`）使用独立的 `SubscriptionOrder` 模型，不会经过充值补单逻辑，不会被影响。

---

### 问题 3：不能影响现有业务

#### 3.1 现有业务影响分析

| 影响范围 | 风险等级 | 说明 |
|----------|---------|------|
| 注册奖励流程 | 低 | `affiliate_reward_trigger=registration` 保持现有行为不变（PRD 第 48 行） |
| 签到流程 | 低 | 只改配置值和前端展示，签到本身逻辑不变 |
| 余额转出 | 低 | 只改前端输入方式，后端逻辑变化小 |
| 计费/扣费 | 无 | PRD 第 21 行的非目标明确排除 |
| 支付回调 | **中** | 需在 6 处加钩子调用 |
| 老数据兼容 | 低 | 保留旧 quota 字段，PRD 第 17 行 |

#### 3.2 老配置兼容需明确（P2）

PRD 第 231 行提到兼容读取老配置，但不够具体：

当前 `CheckinSetting` 结构：
```go
type CheckinSetting struct {
    Enabled  bool `json:"enabled"`   
    MinQuota int  `json:"min_quota"` 
    MaxQuota int  `json:"max_quota"` 
}
```

PRD 新增 `min_amount` / `max_amount` 后，需要明确优先级：
- 新字段有值时用新字段（按 `amount * QuotaPerUnit` 换算为 quota）
- 新字段为空时回退到老字段 `min_quota` / `max_quota`
- 升级时自动从老配额反算为金额写入新字段

#### 3.3 幂等性（P0）

PRD 第 82 行说"重复 webhook 或重复补单不得重复发放同一订单返利"。当前各渠道的幂等实现方式不同：

- Stripe/Creem/Waffo/Waffo Pancake：事务内 `FOR UPDATE` + status 检查
- Epay：`LockOrder()` 内存锁
- 管理员补单：事务内 `FOR UPDATE` + status 检查

**建议**：新增一张 `topup_rebate_log` 表，以 `trade_no + inviter_id` 作为唯一索引，天然保证幂等。返利发放前 `INSERT` 该记录，冲突即跳过。

---

## 三、PRD 具体内容修正建议

### 3.1 第 78 行 — 返利基数

原文：
> 返利基数为订单实际支付金额 `TopUp.Money`。

建议改为：
```
返利基数 = 用户在支付渠道实际支付的原始金额（单位 USD）。
- Stripe/Creem/Waffo 等：由订单创建时从支付网关获取并存入 TopUp.original_pay_amount_usd 字段
- Epay：从 webhook 回调参数提取
- 管理员补单：由管理员手动输入
```

### 3.2 第 81 行 — aff_count 语义

原文：
> `every_topup` 模式下，`aff_count` 表示返利次数，每笔成功充值返利都加 1。

建议改为：
```
新增 AffRebateCount 字段，专用于记录充值返利次数。
aff_count 保持"邀请注册人数"语义不变。
```

### 3.3 第 82 行 — 幂等性

建议补充：
```
新增 topup_rebate_log 表：
- trade_no + inviter_id 作为联合唯一索引
- rebate_amount、rebate_quota、created_at
- 返利发放前 INSERT 此记录，冲突即幂等跳过
```

### 3.4 第 83 行 — 触发范围

原文：
> 触发范围包括 Epay、Stripe、Creem、Waffo、Waffo Pancake、管理员补单。

建议补充：
```
所有充值完成路径统一通过 AfterTopUpSuccessHook 钩子触发返利逻辑，
支付完成函数本身只增加一行钩子调用，不内联返利业务。
```

### 3.5 事务边界（新增建议）

建议在 PRD 第 82 行后增加：
```
**事务边界**：
- 充值事务与返利事务分离。充值成功是返利的前提，但返利失败不应回滚充值。
- 返利使用独立事务，内部幂等。
- Epay 回调需先重构为事务内处理，或将返利逻辑放在 LockOrder 保护区内。
```

### 3.6 first_topup 判断（新增建议）

建议在 PRD 5.3 节补充：
```
"首充"判断规则：
- 在 users 表新增 has_first_topup_rebate 布尔字段（默认 false）
- 第一次成功充值发放返利后，在同一个事务内将该字段设为 true
- 后续充值查到此字段为 true 时跳过返利
```

---

## 四、实现优先级建议

| 阶段 | 任务 | 风险 |
|------|------|------|
| P0 前置 | 解决返利基数不统一（统一使用原始支付金额） | 高 |
| P0 前置 | Epay 回调重构为事务内处理 | 高 |
| P0 | 创建统一钩子 `AfterTopUpSuccessHook` | 中 |
| P1 | 实现返利逻辑 + 幂等表 | 中 |
| P1 | 新增 AffRebateCount 代替复用 aff_count | 低 |
| P1 | first_topup 判断标记字段 | 低 |
| P2 | 签到配置兼容迁移 | 低 |
| P2 | 前端改造 + i18n | 低 |

---

## 五、总结

| 评估维度 | 结论 |
|----------|------|
| 能否实现充值返利？ | **能**，但需要解决 6 个前置问题 |
| 支付逻辑会被修改吗？ | **会被修改**（6 处加钩子调用），但可用架构设计最小化侵入 |
| 会影响现有业务吗？ | **可控**，但 Epay 回调和 `aff_count` 语义需要特别处理 |

**核心建议**：批准进行，但要求 PRD 补充上述 6 处修正（尤其是返利基数、事务边界、幂等方案），技术方案评审时重点审查"支付完成钩子"和"Epay 事务化"两个部分。
