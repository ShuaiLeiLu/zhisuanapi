# Task Specification

## Task Shape

- **Shape**: `single-full`

## Goals

- Implement `docs/prd/invitation-rebate-checkin-amount-prd.md`.
- Configure invitation rewards as registration fixed USD rewards, first top-up rebates, or every top-up rebates.
- Express check-in rewards and affiliate transfer inputs as USD amounts while preserving internal quota accounting.
- Preserve legacy API fields and old configuration compatibility.

## Non-Goals

- Do not migrate the global balance unit away from `quota`.
- Do not change model billing or quota deduction semantics.
- Do not trigger invitation rebates for subscription purchases.
- Do not repurpose `aff_count`; it remains invitation registration count.

## Constraints

- Backend must remain compatible with SQLite, MySQL, and PostgreSQL.
- JSON marshal/unmarshal in business code must use `common` wrappers.
- Frontend package manager and scripts use Bun under `web/default/`.
- Frontend user-facing text must be translated in `en`, `zh`, `fr`, `ja`, `ru`, and `vi`.
- Existing protected project identity/branding must not be removed or renamed.

## Environment

- **Project root**: `D:\workspace\zhisuanapi`
- **Language/runtime**: Go 1.22+, React 19 + TypeScript
- **Package manager**: Go modules, Bun for `web/default`
- **Test framework**: Go test, frontend build/type/i18n scripts
- **Build command**: to be confirmed from repository scripts
- **Existing test count**: to be discovered during implementation

## Risk Assessment

- [ ] Payment completion paths must be identified and hooked without changing recharge semantics.
- [ ] Rebate idempotency must hold across repeated webhooks and manual completion.
- [ ] Epay path has a special transaction/lock boundary.
- [ ] Legacy quota-based settings must remain readable after upgrade.

## Deliverables

- Backend settings, models, migrations, services, controllers, and tests for invitation rebate and check-in amount behavior.
- Frontend admin settings, wallet affiliate transfer/check-in UI, and locale updates.
- PRD retained as implementation reference.

## Done-When

- [ ] Backend supports `registration`, `first_topup`, and `every_topup` invitation reward modes.
- [ ] Top-up rebates use `original_pay_amount_usd`, not provider-specific `Money` / `Amount` semantics.
- [ ] Rebate logs prevent duplicate rebates for repeated success handling.
- [ ] `aff_count` remains invitation count and `aff_rebate_count` tracks rebate count.
- [ ] Check-in configuration and responses support amount fields with legacy fallback.
- [ ] Affiliate transfer accepts amount semantics while preserving old quota fields.
- [ ] Frontend uses amount wording and translations for all new user-facing strings.
- [ ] Relevant Go tests and frontend checks pass or any unavailable check is documented.

## Final Validation Command

```bash
go test ./... && cd web/default && bun run i18n:sync && bun run build
```

