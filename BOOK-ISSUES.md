# Book issues found while building the companion

While following *Bogie: How to Go from Zero to Production* against **Bogie v0.6.8**, these places needed a working fix that the printed listing does not show (or shows incorrectly). Each entry is: chapter, listing or claim, what differs, and what the companion does.

Chapter numbers below are the **20-chapter layout**. The review on 2026-10-31 was written against the old 23-chapter manuscript; each old number is translated here (old 7 → split across 1, 3, and 4; old 8 → 4 and 5; old 9–20 → 6–17; old 5 → 18; old 18 → 15; old 19 → 16; old 20 → 17; old 21 → 19; old 23 → 20).

**Re-checked 2026-10-31** against the updated chapters in Drive (Introduction, the old Chapters 1–23, Appendices A–D). Each entry below now has a **Status** line. A summary:

| # | Item | Status |
|---|---|---|
| 1 | Ch3, Ch6–13: UI pages "copy from the companion" | **Fixed.** Each chapter's pages live in that chapter's `public/`; [`UI-PAGES.md`](UI-PAGES.md) lists them. The old `chNN_NN` checkpoint tags are gone from this repo. |
| 2 | Ch7 / 9: 401-retry in `api()` | **Fixed.** Chapter 7 describes the retry (get a guest token, save it, try again) and the `api()` helper the pages carry. |
| 3 | Ch11: embedding `public/admin/` | **Fixed.** Chapter 11 shows `//go:embed *.html admin/*.html images` in `public/public.go`. |
| 4 | Ch13: catalog UI and admin checkboxes | **Fixed** (copy the whole `public/`, keep your `public.go`). Small gap: the book does not mention that an archived product's checkboxes are disabled. |
| 5 | Ch14: SMTP in encrypted credentials | **Not a book bug** (a limit of the companion, which commits no credentials). The Introduction says the repo never holds credentials and the README shows how to create them. |
| 6 | Ch14: worker test not printed | **Fixed.** The text points at `ch14-send-email-in-the-background/app/jobs/send_order_confirmation_test.go`. |
| 7 | Ch15: OpenAPI and contract tests | **Fixed.** The text points at `ch15-publish-an-api-contract/docs/openapi.yaml` and `ch15-publish-an-api-contract/test/requests/contract_test.go`. |
| 8 | Ch15: 404 vs 422 for a missing product | **Not a book bug.** Chapter 6's test expects `models.ErrNotFound`, which is a 404, and the book never says 422 for that case. The Chapter 15 prose ("a 422 for an empty cart") matches the companion's contract test. |
| 9 | Ch16: credentials and deploy secrets | **Not a book bug** (companion placeholders only). See the new item below: the book itself hard-codes `depot.bogie-go.com`. |
| 10 | Ch5: deliberate `PriceCents < 0` | **Not a book bug** (unchanged). The Chapter 5 folder holds the **fixed** rule (`<= 0`) and the passing tests, which is the state at the *end* of the testing chapter. The same fixed rule is already in the Chapter 4 snapshot. |

New findings from the 2026-10-31 pass are in **"Round 2"** at the end of this file, restated with the new numbers. The detailed entries follow.

## Chapters 3 and 6–13: UI pages are “copy from the companion”

**Claim.** The book says to copy `public/*.html` (and later `public/admin/*.html`) from the companion. The HTML is never printed.

**What differs.** A reader without the companion has no pages to copy.

**Fix.** Every page is in this repo under each chapter's `public/` folder; [`UI-PAGES.md`](UI-PAGES.md) lists, per chapter, each file, what it is, and whether it is new or changed since the previous chapter (so you know what to print). This repo writes those pages (inline CSS/JS; Bogie embeds only `*.html` and `images` from `public/`). They match the chapter’s described behaviour; they are not a verbatim dump of an unpublished checkpoint.

## Chapters 7 / 9: 401-retry in the UI `api()`

**Claim.** Chapter 9 describes forgetting an expired token, minting a fresh guest token, and retrying once.

**What differs.** The Chapter 7 pages were first committed without that retry. Chapter 9’s text assumes it.

**Fix.** Every chapter from 7 onward now uses the same `api(method, path, body, retried)` helper. Chapters 7 and 8 were regenerated so they match.

## Chapter 11: embedding `public/admin/`

**Claim.** Copy `public/admin/*.html` from the companion; open `/admin/products.html`.

**What differs.** Generated `public/public.go` is `//go:embed *.html images`. That pattern does **not** include files under `admin/`, so the admin pages would 404 at runtime even after you copy them in.

**Fix.** From Chapter 11 on, `public/public.go` is:

```go
//go:embed *.html admin/*.html images
```

The book never mentions this change.

## Chapter 13: catalog UI and admin category checkboxes

**Claim.** Copy the companion's `public/` for the category menu, paging, admin checkboxes, and admin categories page.

**What differs.** Same as Chapters 3 and 6–12: the pages are not printed.

**Fix.** Written for this repo. An archived product’s category checkboxes are disabled on the admin products page, because `GET /products/:id` (used to load the current set) does not return archived products.

## Chapter 14: SMTP in encrypted credentials

**Claim.** `bogie g config smtp …` puts Mailpit’s address in the encrypted development credentials.

**What differs.** This companion commits **no** `config/master.key`, `config/credentials.yml.enc`, or `config/credentials/*.key` (and must not). Without credentials, those Mailpit values would be missing.

**Fix.** `SmtpConfig.defaults` sets `localhost:1025` and the book’s `From` for every non-production environment. Production still requires `SMTP_ADDRESS`, `SMTP_PASSWORD`, and `SMTP_FROM` (and optionally username/port) from the environment or production credentials. No real SMTP secrets are committed.

## Chapter 14: worker test not printed

**Claim.** After wiring `Confirm` into the worker, “give it a fake `Confirm` that records the order ID” and check `Work` passes `job.Args.OrderID`.

**What differs.** The test body is described, not listed.

**Fix.** `app/jobs/send_order_confirmation_test.go` in `ch14-send-email-in-the-background` holds that fake.

## Chapter 15: full OpenAPI document and contract tests

**Claim.** Copy the full `docs/openapi.yaml` and the rest of `test/requests/contract_test.go` from the companion.

**What differs.** The book only prints the product paths and the first two contract examples; the rest is “copy from the companion”.

**Fix.** Written here from the real controllers and views (every public route and documented status). Playtime (`public/docs.html` with Redoc) is not included.

## Chapter 15: nonexistent product when adding to the cart

**Claim.** (Inferred from the pattern of documenting every error.) A `POST /cart_items` for a product that does not exist should be a documented client error.

**What differs.** `AddCartItem` maps a missing product to `models.ErrNotFound`, and `failCart` answers **404**, not 422.

**Fix.** The contract documents and tests **404** for that case. A reader who expects 422 from a generic “invalid” path will be surprised.

## Chapter 16: credentials and deploy secrets

**Claim.** `make credentials ENV=staging` / `production`, fill real values, deploy with Kamal.

**What differs.** This companion must not commit secrets. Deploy files still need real hosts, a registry, and credentials keys to deploy.

**Fix.** Placeholders only: server `203.0.113.10` (the documentation address the chapter uses), hosts `staging-depot.bogie-go.com` / `depot.bogie-go.com`, image `roonglit/bogie-depot`, registry user `your-docker-hub-username`. `.kamal/secrets*` only shell out to `go tool credentials get` / `cat …key`. Create your own credentials locally; nothing was deployed from this repo.

## Chapter 5: deliberate `PriceCents < 0` mistake

**Claim.** The testing chapter introduces a wrong rule (`PriceCents < 0`), then fixes it to `<= 0` after the test fails.

**What differs.** Not a book bug. The Chapter 5 folder holds the **fixed** rule (`<= 0`) and the passing tests, which is the state at the *end* of the chapter. Chapter 4's snapshot already has the same fixed rule, without those hand-written tests.

## Not issues (verified against the book)

These claims were checked on the box against Bogie v0.6.8 and matched the book: factory-defaults tests failing until defaults are set; removing `LockCart` or the product-ID sort breaking the concurrency tests (including the deadlock); swapping `GetProductForUpdate` for `GetProduct` in `EditProduct` failing the Chapter 12 race; dropping the `guest:` check failing the “never merges another customer’s cart” test; sign-in timing for unknown / wrong / correct password; guest-token cart merge; cart totals with the launch rules ($68 / $18 / $13.60); Chapter 13’s archive check must sit in `ON` (moving it to `WHERE` drops a category whose only product is archived); Chapter 15’s two deliberate breaks (`price_cents` → `price`, and `bogie g controller wishlists`) fail as described; Chapter 14’s transactional outbox (job rolls back with a failed `after`); Mailpit delivery and retry after Mailpit is stopped.

## Round 2: findings from the updated Drive chapters (2026-10-31)

Checked by running Bogie v0.6.8 (`bogie new`, every generator command in the shop chapters, `make help`, `bogie routes`, `bogie doctor`, `bogie help`) and by comparing every labelled listing with this repo. Code listings, command outputs and file paths matched; the items below are what did not. Numbers are the new layout; the Drive files reviewed still used the old ones.

**Introduction**

- *"Each chapter starts with a request from me, as the shop's owner"* is not true for Chapter 5 (a reviewer's question; Chapter 4 is too), Chapter 7 (a confused customer), Chapter 9 (what "the regulars asked for"), Chapter 11 ("the request from my team") or Chapter 15 (a friend who builds mobile apps). Suggested: "Each chapter starts with a request or a problem report."
- *"every chapter after it builds its feature from the outside in, starting with the test"* is not true: Chapters 6, 7, 9, 12, 14 and 17 open with migrations, models or code; the first test listing is well into the chapter. Suggested: "from Chapter 5 on, every feature comes with its tests." Chapter 20's "outside-in testing for every feature after it" and "Outside in, as always" need the same softening.
- *"a real online shop that sells real Bogie merchandise… It's not a pretend shop… live at `depot.bogie-go.com`… I'm opening this shop for real"* is still an oversell: Chapter 20 says the shop "never charges anyone", and `depot.bogie-go.com` does not resolve today. Either put the site up before publishing, or reword to "a complete online shop… you deploy your own copy".
- *"the hundred decisions a Go web app needs"* contradicts Chapter 1's "adding up to dozens". Use one number, or "dozens".
- The Introduction never says the backend is a JSON API; one sentence in "What You'll Build" would set the expectation.
- "Getting Help" says "compare your code with the checkpoint"; from Chapter 1 on the unit is the chapter's folder.
- The Introduction still needs the outline's update: Part I is two chapters, Part III is three, and there are five appendices.

**Chapter 16** (also Chapters 14, 15, 17): the book tells the reader to use their own server address, but hard-codes `depot.bogie-go.com` and `staging-depot.bogie-go.com` ("Customers will type…") with no "use your own domain here and below".

**Chapter 20**: "Bogie Depot is a real shop, but a small one" — align with the Introduction fix.

**Appendix C** (and the same rows in Appendix B)

- "`bogie help`, and `bogie help <command>`": `bogie help g` just prints the general help. Per-command help is `bogie g -h` or `bogie app:update --help`.
- `bogie version` was listed under the old Chapter 14 (now 11); it first appears in Chapter 1. `bogie routes` was listed under the old Chapter 7; it first appears in Chapter 1.
- Rows whose chapter does not mention the command. After the old-to-new renumber, those rows land on: `bogie lint` and `bogie ci` (Chapter 5, old 8), `bogie db:migrate:status` (the old Chapter 7 row; that chapter is now split, and Chapters 1, 3, and 4 do not run it), `credentials:edit` and `credentials:show` (Chapter 1, old 6), `credentials:get KEY` (Chapter 16, old 19), `bogie g controller admin/reports index` (Chapter 11, old 14), the `:index` modifier (Chapter 6, old 9). The commands all work in v0.6.8; the chapters just use `bin/ci`, `make credentials` or `go tool credentials`. Either drop the numbers, or add a line saying which commands the book never runs.

**Cross-cutting**

- Seven `[SCREENSHOT NEEDED]` placeholders, now named for the new chapters: Chapter 1 `ch01-welcome.png`, Chapter 3 `ch03-catalog.png`, Chapter 6 `ch06-add-to-cart.png`, Chapter 8 `ch08-sale.png`, Chapter 11 `ch11-admin.png`, Chapter 13 `ch13-categories.png`, Chapter 14 `ch14-mailpit.png`. The files are in [`screenshots/`](screenshots/).

## Chapters 1–3, checked against Bogie v0.6.8 (2026-11-04)

Regenerated with `bogie new depot`, `bogie g model product name description:text price_cents:integer stock:integer image_url`, `bogie g scaffold_controller products --only index,show`, and `bogie g seed products`. `go build`, `go vet`, and `bogie test` pass for Chapters 1 and 3. Chapter 2's six scratch programs print what the chapter shows.

What matched: generator file lists for `g model` and `scaffold_controller` (the migration timestamp differs, and the chapter says it will); `GET /healthz` is `ok`; an empty catalog is `{"products":[]}`; the two inserted products come back newest first, and `GET /products/1` matches the printed JSON, including `created_at` and `updated_at`; after Chapter 3's view edit those two fields are gone; `GET /products/nope` and `GET /products/999` are both `{"status":"error","error":"no such product"}` with status 404; seeding twice leaves `count(*) = 7`, with Chapter 1's two rows at the end of the newest-first list; `\d products` matches the printed table; the welcome page shows the app name, `development`, `:8080`, and the Go version.

What does not match the printed sample, or only matches in a particular setup:

- **`bogie routes`** prints two configuration lines, then Gin's debug-mode warning (`Running in "debug" mode` …), then the three routes. The book's sample shows only the three routes, and with slightly tighter spacing. The method, path, and handler on each route match.
- **`bogie db:prepare`** logs `database created` when `depot_development` does not exist yet. The book says you will most likely see `database already exists`, because Docker's `POSTGRES_DB` creates it before Bogie runs. That is the Docker path. A Postgres that has not pre-created the database gets `database created`.
- **The migration file's column alignment** is whatever `bogie g model` wrote. The SQL matches the listing; the spaces used to line the types up do not. `\d products` still matches the book.
- **The migration filename** in the Chapter 3 listing is `20261031071211_create_products.sql`. A fresh `bogie g model` writes a new timestamp. This repo keeps the generated name. The chapter already says the timestamp will differ.
- **A missing product photo** does what the chapter says. `index.html` and `product.html` set `img.src` to `image_url` and leave it there. Chrome shows the broken-image icon and the alt text (the product name). `naturalWidth` is 0. The card and the product page still render the name and the price. The pages do not hide the image or substitute a placeholder. With the `.webp` files in `public/images/products/`, the same pages show the pictures (HTTP 200, `image/webp`). A path that is not in `public/` is HTTP 404, `text/html` (the app's `404.html`), which is why the browser treats it as a broken image.
