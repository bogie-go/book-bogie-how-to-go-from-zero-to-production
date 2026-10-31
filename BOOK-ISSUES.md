# Book issues found while building the companion

While following *Bogie: How to Go from Zero to Production* against **Bogie v0.6.8**, these places needed a working fix that the printed listing does not show (or shows incorrectly). Each entry is: chapter, listing or claim, what differs, and what the companion does.

## Chapter 7–16: UI pages are “copy from the companion”

**Claim.** From Chapter 7 on, the book says to copy `public/*.html` (and later `public/admin/*.html`) from a companion checkpoint (`ch07_05`, `ch09_05`, `ch10_04`, `ch11_06`, `ch12_06`, `ch13_07`, `ch14_05`, `ch15_03`, `ch16_07`). The HTML is never printed.

**What differs.** A reader without the companion has no pages to copy.

**Fix.** This repo writes those pages (inline CSS/JS; Bogie embeds only `*.html` and `images` from `public/`). They match the chapter’s described behaviour; they are not a verbatim dump of an unpublished checkpoint.

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
