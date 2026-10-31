# UI pages to print in the book

Chapters 7–16 use plain HTML pages in `public/`. Each page is one self-contained file (inline CSS and JavaScript, no build step, no framework). The book tells the reader to copy them from a checkpoint, but never prints them. This file lists, per chapter, every page the chapter needs and what changed since the previous chapter, so you know what to print.

Each chapter folder holds **all** pages the app has at the end of that chapter, so a reader can copy the whole `public/` folder from any chapter. `public/404.html`, `public/500.html` and `public/images/` come from `bogie new` and are not part of the UI to print.

Chapters 17–20 do not change `public/`: their pages are identical to Chapter 16's. Chapters 1–6 and 21–23 have no UI pages to print (Chapter 6 has only Bogie's generated `404.html`/`500.html`).

**Status key:** NEW = first appears this chapter. CHANGED = differs from the previous chapter's copy. same = byte-identical to the previous chapter (no need to print again).

**Embed pattern.** Bogie serves only what `public/public.go` embeds. Chapters 7–13 use the generated `//go:embed *.html images`. From Chapter 14 on, change it to `//go:embed *.html admin/*.html images` (otherwise `/admin/*.html` is a 404). See `BOOK-ISSUES.md`.

**Book checkpoint names.** The book refers to these as `ch07_05`, `ch09_05`, `ch10_04`, `ch11_06`, `ch12_06`, `ch13_07`, `ch14_05`, `ch15_03`, `ch16_07`. In this repo, each chapter folder is that checkpoint's `public/`.

## Chapter 7: Show Our Merch (`ch07-show-our-merch/public/`)

| File | What it is | Status |
| --- | --- | --- |
| `index.html` | Catalog page. Fetches `GET /products` and draws a card per product. | NEW |
| `product.html` | One product (`?id=`). Fetches `GET /products/:id`; 404 shows "No such product." | NEW |

## Chapter 8: Testing (`ch08-testing/public/`)

No UI change. `index.html`, `product.html`: same as Chapter 7.

## Chapter 9: Add to Cart (`ch09-add-to-cart/public/`)

| File | What it is | Status |
| --- | --- | --- |
| `cart.html` | Cart page: quantity stepper per line, remove link, total. Uses `GET /cart`, `PATCH`/`DELETE /cart_items/:id`. | NEW |
| `index.html` | Adds a "Cart (n)" badge in the nav; adds a small `api()` helper. | CHANGED |
| `product.html` | Adds the cart badge and an "Add to cart" button (`POST /cart_items`). | CHANGED |

## Chapter 10: Cart for Every Visitor (`ch10-cart-for-every-visitor/public/`)

The shared script block now gets a guest token (`POST /guest_sessions`), keeps it in `localStorage`, sends it as a bearer token on cart calls, and retries once with a fresh token on a 401. The same block is in every page.

| File | What it is | Status |
| --- | --- | --- |
| `cart.html` | Cart page; token handling and 401 retry added. | CHANGED |
| `index.html` | Catalog; same token/`api(method, path, body, retried)` block. | CHANGED |
| `product.html` | Product page; same token/`api()` block. | CHANGED |

## Chapter 11: Run a Sale (`ch11-run-a-sale/public/`)

| File | What it is | Status |
| --- | --- | --- |
| `cart.html` | Shows the applied sale rules (name and discount) next to the total, from `applied_rules` in `GET /cart`. | CHANGED |
| `index.html`, `product.html` | No change from Chapter 10. | same |

## Chapter 12: Accounts (`ch12-accounts/public/`)

| File | What it is | Status |
| --- | --- | --- |
| `signin.html` | Sign-in form. `POST /sessions`, stores the token and goes to the cart (the guest cart merges on the server). | NEW |
| `signup.html` | Sign-up form. `POST /users`, which returns a session; stores the token and goes to the cart. | NEW |
| `cart.html` | Adds the account area to the nav (email, Sign out, or Sign in/Sign up links). | CHANGED |
| `index.html` | Same nav change and `renderAccount()`. | CHANGED |
| `product.html` | Same nav change and `renderAccount()`. | CHANGED |

## Chapter 13: Place Orders (`ch13-place-orders/public/`)

| File | What it is | Status |
| --- | --- | --- |
| `checkout.html` | Checkout: order summary, optional coupon code, "Place order" (`POST /orders`). | NEW |
| `order.html` | Order confirmation / detail (`?id=`), `GET /orders/:id`. | NEW |
| `cart.html` | Adds a "Check out" button, shown only when the cart has items. | CHANGED |
| `index.html`, `product.html`, `signin.html`, `signup.html` | No change from Chapter 12. | same |

## Chapter 14: Staff (`ch14-staff/public/`)

| File | What it is | Status |
| --- | --- | --- |
| `admin/products.html` | Staff list/create/edit/delete for products (`/admin/products`). | NEW |
| `admin/price_rules.html` | Staff screen for sale rules (`/admin/price_rules`). | NEW |
| `admin/coupons.html` | Staff screen for coupons (`/admin/coupons`). | NEW |
| `cart.html`, `checkout.html`, `index.html`, `order.html`, `product.html`, `signin.html`, `signup.html` | Each gets an "Admin" link in the nav for staff, found by probing `GET /admin/products` once after sign-in. | CHANGED |

The three admin pages share one script (a table driven by a `RES` config for slug and fields). Also change `public/public.go` (see above).

## Chapter 15: Archive (`ch15-archive/public/`)

| File | What it is | Status |
| --- | --- | --- |
| `admin/products.html` | Adds an Archive / Restore link per product (`POST /admin/products/:id/archive` or `/restore`); archived rows are dimmed. | CHANGED |
| `admin/price_rules.html`, `admin/coupons.html` | Same shared script, so the file changes, but the archive code does nothing for them (it only runs when `RES.slug === "products"`). | CHANGED |
| `cart.html`, `checkout.html`, `index.html`, `order.html`, `product.html`, `signin.html`, `signup.html` | No change from Chapter 14. | same |

## Chapter 16: Categories (`ch16-categories/public/`)

| File | What it is | Status |
| --- | --- | --- |
| `admin/categories.html` | Staff screen for categories (`/admin/categories`). | NEW |
| `admin/products.html` | Adds category checkboxes saved with `PUT /admin/products/:id/categories`. They are disabled for an archived product. | CHANGED |
| `admin/price_rules.html`, `admin/coupons.html` | Shared admin script gets the `onEdit`/`saveCats` hooks and a Categories nav link. | CHANGED |
| `index.html` | Category menu (`GET /categories`), `?category_id=` filter and paging (`?page=`). | CHANGED |
| `cart.html`, `checkout.html`, `order.html`, `product.html`, `signin.html`, `signup.html` | One-line fix in `api()`: ignore the query string when deciding whether a path needs a token. | CHANGED |

## What to print

Print the whole page once, where it is NEW. For CHANGED pages that only repeat the shared script block (Chapters 10, 12, 14, 16 on most pages), print the changed block or a short diff, not the whole file. If you prefer the book not to print full pages at all, point readers to `chNN-…/public/` in this repo instead of the `chNN_NN` checkpoint names.
