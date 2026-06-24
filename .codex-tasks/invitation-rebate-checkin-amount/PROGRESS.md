# Progress Log

---

## Session Start

- **Date**: 2026-06-24
- **Task name**: `invitation-rebate-checkin-amount`
- **Task dir**: `.codex-tasks/invitation-rebate-checkin-amount/`
- **Spec**: See `SPEC.md`
- **Plan**: See `TODO.csv` (6 milestones)
- **Environment**: Go + React/TypeScript + Bun

---

## Context Recovery Block

- **Current milestone**: #6 — Complete i18n and final verification
- **Current status**: DONE
- **Last completed**: #6 — Complete i18n and final verification
- **Current artifact**: `.codex-tasks/invitation-rebate-checkin-amount/TODO.csv`
- **Key context**: Backend code now has amount/quota helpers, new quota/check-in settings, `original_pay_amount_usd`, `topup_rebate_logs`, `aff_rebate_count`, `has_first_topup_rebate`, registration reward handling, top-up rebate hook calls, check-in amount fields, amount-based affiliate transfer compatibility, frontend amount-setting forms, wallet amount transfer UI, check-in amount display, and complete locale keys.
- **Known issues**: `.workbuddy/` is untracked and unrelated; do not touch it. Local Go/Bun dependencies are absent, so verification was run through Docker.
- **Next action**: Review and commit the completed feature changes.

---

## Milestone 1-4: Backend Core

- **Status**: DONE (verification blocked by missing Go runtime)
- **Started**: 2026-06-24
- **Completed**: 2026-06-24
- **What was done**:
  - Added amount/quota conversion helpers and amount-aware settings.
  - Added user/top-up/rebate log schema fields.
  - Added top-up affiliate rebate service with idempotent logs.
  - Hooked Epay, Stripe, Creem, Waffo, Waffo Pancake, and manual completion paths.
  - Changed registration rewards, check-in rewards, and affiliate transfer to amount semantics with legacy compatibility.
- **Validation**: `go test ./model` could not run because `go` is not available in PATH.
- **Next step**: Milestone 5 — frontend UI and i18n.

---

## Milestone 5: Frontend UI and i18n

- **Status**: DONE (build verification blocked by missing frontend dependencies)
- **Started**: 2026-06-24
- **Completed**: 2026-06-24
- **What was done**:
  - Replaced invitation reward quota inputs with trigger mode plus USD amount / percentage fields.
  - Replaced check-in quota inputs with USD min/max amount fields.
  - Changed wallet affiliate card and transfer dialog to amount semantics and rebate count display.
  - Changed check-in calendar/toasts/statistics to prefer amount fields with quota fallback.
  - Added amount/rebate fields to wallet, profile, auth, users, and usage-log types.
  - Completed en/zh/fr/ja/ru/vi locale entries, including pre-existing missing `t()` keys found during scanning.
  - Tightened settings compatibility so new amount options override legacy quota runtime values when present.
- **Validation**:
  - `node scripts/sync-i18n.mjs` passed.
  - Custom `t()` key scanner found all source keys in `en.json`.
  - `npm run typecheck` failed because `tsgo` is unavailable; `web/default/node_modules` does not exist.
  - `npm run build` failed because `rsbuild` is unavailable; `web/default/node_modules` does not exist.
- **Next step**: Milestone 6 — final static review and environment-limited verification.

---

## Milestone 6: Static Review Pass

- **Status**: DONE
- **Started**: 2026-06-24
- **What was done**:
  - Added frontend validation so maximum check-in USD reward cannot be lower than the minimum.
  - Added locale coverage for the new validation message across en/zh/fr/ja/ru/vi.
  - Replaced remaining JSON request parsing in the touched top-up/user controllers with `common.DecodeJson` / `common.Unmarshal`.
  - Added registration reward tests and a first-top-up zero-percent opportunity-consumption test.
  - Added `original_pay_amount_usd` completion fallbacks for pending orders completed through callback or admin补单.
  - Changed check-in reward selection to randomize the USD amount first, then convert that amount to internal quota for入账.
  - Added backend check-in amount range validation so `max_amount` cannot be saved below `min_amount`.
  - Adjusted the frontend check-in settings save order so simultaneous min/max edits do not create an invalid intermediate state.
- **Validation**:
  - `node scripts/sync-i18n.mjs` passed.
  - Locale sync report shows 0 missing, 0 extra, and 0 untranslated keys for en/zh/fr/ja/ru/vi.
  - `git diff --check` passed.
  - `docker run ... golang:1.25-bookworm go test -timeout=120s ./...` passed.
  - `docker run ... oven/bun:1 ... bun run typecheck; bun run build` passed in a temporary copy of `web/`, without changing workspace lockfiles or `node_modules`.
- **Next step**: Review and commit the completed feature changes.

---
