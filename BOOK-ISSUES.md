# Book issues found while building the companion

While following *Bogie: How to Go from Zero to Production* against **Bogie v0.6.8**, these places needed a working fix that the printed listing does not show (or shows incorrectly). Each entry is: chapter, listing or claim, what differs, and what the companion does.

**Re-checked 2026-10-31** against the updated chapters in Drive (Introduction, Chapters 1–23, Appendices A–D). Each entry below now has a **Status** line. A summary:

| # | Item | Status |
|---|---|---|
| 1 | Ch7–16: UI pages "copy from the companion" | **Fixed.** Each chapter names the exact `chNN-…/public/` folder and files; no `chNN_NN` checkpoint tags remain after Ch5; the Introduction explains the `public/` folders and `UI-PAGES.md`. |
| 2 | Ch10 / 12: 401-retry in `api()` | **Fixed.** Ch10 now describes the retry (get a guest token, save it, try again) and the `api()` helper the pages carry. |
| 3 | Ch14: embedding `public/admin/` | **Fixed.** Ch14 now shows `//go:embed *.html admin/*.html images` in `public/public.go`. |
| 4 | Ch16: catalog UI and admin checkboxes | **Fixed** (copy the whole `public/`, keep your `public.go`). Small gap: the book does not mention that an archived product's checkboxes are disabled. |
| 5 | Ch17: SMTP in encrypted credentials | **Not a book bug** (a limit of the companion, which commits no credentials). The Introduction now says the repo never holds credentials and the README shows how to create them. |
| 6 | Ch17: worker test not printed | **Fixed.** The text now points to `ch17-email/app/jobs/send_order_confirmation_test.go`. |
| 7 | Ch18: OpenAPI and contract tests | **Fixed.** The text points to `ch18-contract/docs/openapi.yaml` and `ch18-contract/test/requests/contract_test.go`. |
| 8 | Ch18: 404 vs 422 for a missing product | **Not a book bug.** Ch9's test expects `models.ErrNotFound`, which is a 404, and the book never says 422 for that case. The Ch18 prose ("a 422 for an empty cart") matches the companion's contract test. |
| 9 | Ch19: credentials and deploy secrets | **Not a book bug** (companion placeholders only). See the new item below: the book itself hard-codes `depot.bogie-go.com`. |
| 10 | Ch8: deliberate `PriceCents < 0` | **Not a book bug** (unchanged). |

New findings from this pass (Introduction, Ch19, Ch23, Appendices B and C) are in the section **"Round 2: findings from the updated Drive chapters"** at the end of this file. The detailed entries follow.

## Chapter 7–16: UI pages are “copy from the companion”

**Claim.** From Chapter 7 on, the book says to copy `public/*.html` (and later `public/admin/*.html`) from a companion checkpoint (`ch07_05`, `ch09_05`, `ch10_04`, `ch11_06`, `ch12_06`, `ch13_07`, `ch14_05`, `ch15_03`, `ch16_07`). The HTML is never printed.

**What differs.** A reader without the companion has no pages to copy.

**Fix.** Every page is in this repo under each chapter's `public/` folder; [`UI-PAGES.md`](UI-PAGES.md) lists, per chapter, each file, what it is, and whether it is new or changed since the previous chapter (so you know what to print). This repo writes those pages (inline CSS/JS; Bogie embeds only `*.html` and `images` from `public/`). They match the chapter’s described behaviour; they are not a verbatim dump of an unpublished checkpoint.

## Chapter 10 / 12: 401-retry in the UI `api()`

**Claim.** Chapter 12 describes forgetting an expired token, minting a fresh guest token, and retrying once.

**What differs.** The Chapter 10 pages were first committed without that retry. Chapter 12’s text assumes it.

**Fix.** Every chapter from 10 onward now uses the same `api(method, path, body, retried)` helper. Chapters 10 and 11 were regenerated so they match.

## Chapter 14: embedding `public/admin/`

**Claim.** Copy `public/admin/*.html` from the companion; open `/admin/products.html`.

**What differs.** Generated `public/public.go` is `//go:embed *.html images`. That pattern does **not** include files under `admin/`, so the admin pages would 404 at runtime even after you copy them in.

**Fix.** From Chapter 14 on, `public/public.go` is:

```go
//go:embed *.html admin/*.html images
```

The book never mentions this change.

## Chapter 16: catalog UI and admin category checkboxes

**Claim.** Copy the `ch16_07` checkpoint for the category menu, paging, admin checkboxes, and admin categories page.

**What differs.** Same as Chapters 7–15: the checkpoint is not printed.

**Fix.** Written for this repo. An archived product’s category checkboxes are disabled on the admin products page, because `GET /products/:id` (used to load the current set) does not return archived products.

## Chapter 17: SMTP in encrypted credentials

**Claim.** `bogie g config smtp …` puts Mailpit’s address in the encrypted development credentials.

**What differs.** This companion commits **no** `config/master.key`, `config/credentials.yml.enc`, or `config/credentials/*.key` (and must not). Without credentials, those Mailpit values would be missing.

**Fix.** `SmtpConfig.defaults` sets `localhost:1025` and the book’s `From` for every non-production environment. Production still requires `SMTP_ADDRESS`, `SMTP_PASSWORD`, and `SMTP_FROM` (and optionally username/port) from the environment or production credentials. No real SMTP secrets are committed.

## Chapter 17: worker test not printed

**Claim.** After wiring `Confirm` into the worker, “give it a fake `Confirm` that records the order ID” and check `Work` passes `job.Args.OrderID`.

**What differs.** The test body is described, not listed.

**Fix.** `app/jobs/send_order_confirmation_test.go` in this repo holds that fake.

## Chapter 18: full OpenAPI document and contract tests

**Claim.** Copy the full `docs/openapi.yaml` and the rest of `test/requests/contract_test.go` from companion checkpoints `ch18_01` / `ch18_02`.

**What differs.** The book only prints the product paths and the first two contract examples; the rest is “copy from the companion”.

**Fix.** Written here from the real controllers and views (every public route and documented status). Playtime (`public/docs.html` with Redoc) is not included.

## Chapter 18: nonexistent product when adding to the cart

**Claim.** (Inferred from the pattern of documenting every error.) A `POST /cart_items` for a product that does not exist should be a documented client error.

**What differs.** `AddCartItem` maps a missing product to `models.ErrNotFound`, and `failCart` answers **404**, not 422.

**Fix.** The contract documents and tests **404** for that case. A reader who expects 422 from a generic “invalid” path will be surprised.

## Chapter 19: credentials and deploy secrets

**Claim.** `make credentials ENV=staging` / `production`, fill real values, deploy with Kamal.

**What differs.** This companion must not commit secrets. Deploy files still need real hosts, a registry, and credentials keys to deploy.

**Fix.** Placeholders only: server `203.0.113.10` (the documentation address the chapter uses), hosts `staging-depot.bogie-go.com` / `depot.bogie-go.com`, image `roonglit/bogie-depot`, registry user `your-docker-hub-username`. `.kamal/secrets*` only shell out to `go tool credentials get` / `cat …key`. Create your own credentials locally; nothing was deployed from this repo.

## Chapter 8: deliberate `PriceCents < 0` mistake

**Claim.** The chapter introduces a wrong rule (`PriceCents < 0`), then fixes it to `<= 0` after the test fails.

**What differs.** Not a book bug. The companion folder for Chapter 8 holds the **fixed** rule (`<= 0`) and the passing tests, which is the state at the *end* of the chapter.

## Not issues (verified against the book)

These claims were checked on the box against Bogie v0.6.8 and matched the book: factory-defaults tests failing until defaults are set; removing `LockCart` or the product-ID sort breaking the concurrency tests (including the deadlock); swapping `GetProductForUpdate` for `GetProduct` in `EditProduct` failing the Chapter 15 race; dropping the `guest:` check failing the “never merges another customer’s cart” test; sign-in timing for unknown / wrong / correct password; guest-token cart merge; cart totals with the launch rules ($68 / $18 / $13.60); Chapter 16’s archive check must sit in `ON` (moving it to `WHERE` drops a category whose only product is archived); Chapter 18’s two deliberate breaks (`price_cents` → `price`, and `bogie g controller wishlists`) fail as described; Chapter 17’s transactional outbox (job rolls back with a failed `after`); Mailpit delivery and retry after Mailpit is stopped.

## Round 2: findings from the updated Drive chapters (2026-10-31)

Checked by running Bogie v0.6.8 (`bogie new`, every generator command in Ch7–17 and 20, `make help`, `bogie routes`, `bogie doctor`, `bogie help`) and by comparing every labelled listing with this repo. Code listings, command outputs and file paths matched; the items below are what did not.

**Introduction**

- *"Each chapter starts with a request from me, as the shop's owner"* is not true for Ch8 (a reviewer's question), Ch10 (a confused customer), Ch12 (what "the regulars asked for"), Ch14 ("the request from my team") or Ch18 (a friend who builds mobile apps). Suggested: "Each chapter starts with a request or a problem report."
- *"every chapter after it builds its feature from the outside in, starting with the test"* is not true: Ch9, 10, 12, 15, 17 and 20 open with migrations, models or code; the first test listing is the 4th (Ch9, 12, 20), 7th (Ch17) or 9th (Ch10, 15) listing. Suggested: "from Chapter 8 on, every feature comes with its tests." Ch23's "outside-in testing for every feature after it" and "Outside in, as always" need the same softening.
- *"a real online shop that sells real Bogie merchandise… It's not a pretend shop… live at `depot.bogie-go.com`… I'm opening this shop for real"* is still an oversell: Ch23 says the shop "never charges anyone", and `depot.bogie-go.com` does not resolve today. Either put the site up before publishing, or reword to "a complete online shop… you deploy your own copy".
- *"the hundred decisions a Go web app needs"* contradicts Ch1's "adding up to dozens". Use one number, or "dozens".
- The Introduction never says the backend is a JSON API; one sentence in "What You'll Build" would set the expectation.
- "Getting Help" says "compare your code with the checkpoint"; from Ch6 on the unit is the chapter's folder.
- Fixed since the last review: listings point to chapter folders (not tags); "Four appendices" now matches Appendices A–D.

**Chapter 19** (also Ch17, 18, 20): the book tells the reader to use their own server address, but hard-codes `depot.bogie-go.com` and `staging-depot.bogie-go.com` ("Customers will type…") with no "use your own domain here and below".

**Chapter 23**: "Bogie Depot is a real shop, but a small one" — align with the Introduction fix.

**Appendix C** (and the same rows in Appendix B)

- "`bogie help`, and `bogie help <command>`": `bogie help g` just prints the general help. Per-command help is `bogie g -h` or `bogie app:update --help`.
- `bogie version` is listed under Chapter 14; it first appears in Chapter 6. `bogie routes` (7) first appears in Chapter 6.
- Rows whose chapter does not mention the command: `bogie lint` (8), `bogie ci` (8), `bogie db:migrate:status` (7), `credentials:edit` and `credentials:show` (6), `credentials:get KEY` (19), `bogie g controller admin/reports index` (14), the `:index` modifier (9). The commands all work in v0.6.8; the chapters just use `bin/ci`, `make credentials` or `go tool credentials`. Either drop the numbers, or add a line saying which commands the book never runs.

**Cross-cutting**

- Seven `[SCREENSHOT NEEDED]` placeholders remain: Ch6 `ch06-welcome.png`, Ch7 `ch07-catalog.png`, Ch9 `ch09-add-to-cart.png`, Ch11 `ch11-sale.png`, Ch14 `ch14-admin.png`, Ch16 `ch16-categories.png`, Ch17 `ch17-mailpit.png`.
- No chapter points to an appendix; the Introduction is the only mention.
