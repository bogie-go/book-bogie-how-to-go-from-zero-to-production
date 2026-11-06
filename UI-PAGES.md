# UI pages to print in the book

Chapters 3 and 6–13 use plain HTML pages in `public/`. Each page is one self-contained file (inline CSS and JavaScript, no build step, no framework). The book tells the reader to copy them from the chapter folder, but never prints them. This file lists, per chapter, every page the chapter needs and what changed since the previous chapter, so you know what to print.

Each chapter folder holds **all** pages the app has at the end of that chapter, so a reader can copy the whole `public/` folder from any chapter. `public/404.html`, `public/500.html` and `public/images/` come from `bogie new` (the product photos are added in Chapter 3) and are not part of the UI to print.

Chapters 14–17 do not change `public/`: their pages are identical to Chapter 13's. Chapters 1, 2, and 18–20 have no shop UI pages to print (Chapters 1 and 2 have only Bogie's generated `404.html`/`500.html`).

Chapters 4 and 5 do not change the pages either. Chapter 4's folder still has `index.html` and `product.html` because that snapshot is the old catalog chapter, which already included them. They match Chapter 3.

**Status key:** NEW = first appears this chapter. CHANGED = differs from the previous chapter's copy. same = byte-identical to the previous chapter (no need to print again).

**Embed pattern.** Bogie serves only what `public/public.go` embeds. Chapters 3–10 use the generated `//go:embed *.html images`. From Chapter 11 on, change it to `//go:embed *.html admin/*.html images` (otherwise `/admin/*.html` is a 404). See `BOOK-ISSUES.md`.

## Chapter 3: Show the Catalog (`ch03-show-the-catalog/public/`)

| File | What it is | Status |
| --- | --- | --- |
| `index.html` | Catalog page. Fetches `GET /products` and draws a card per product. | NEW |
| `product.html` | One product (`?id=`). Fetches `GET /products/:id`; 404 shows "No such product." | NEW |

A missing image file stays a broken image: the page sets `img.src` to `image_url` and does not hide the element or swap in a placeholder. The browser shows its broken-image icon, with the product name as the alt text. The rest of the card or product page still renders.

## Chapter 4: Make the Catalog Safe (`ch04-make-the-catalog-safe/public/`)

No UI change in the new chapter. `index.html`, `product.html`: same as Chapter 3. (This folder is the old catalog snapshot, so the files are already here.)

## Chapter 5: Test the Whole App (`ch05-test-the-whole-app/public/`)

No UI change. `index.html`, `product.html`: same as Chapter 4.

## Chapter 6: Add a Cart (`ch06-add-a-cart/public/`)

| File | What it is | Status |
| --- | --- | --- |
| `cart.html` | Cart page: quantity stepper per line, remove link, total. Uses `GET /cart`, `PATCH`/`DELETE /cart_items/:id`. | NEW |
| `index.html` | Adds a "Cart (n)" badge in the nav; adds a small `api()` helper. | CHANGED |
| `product.html` | Adds the cart badge and an "Add to cart" button (`POST /cart_items`). | CHANGED |

## Chapter 7: Give Every Visitor a Cart (`ch07-give-every-visitor-a-cart/public/`)

The shared script block now gets a guest token (`POST /guest_sessions`), keeps it in `localStorage`, sends it as a bearer token on cart calls, and retries once with a fresh token on a 401. The same block is in every page.

| File | What it is | Status |
| --- | --- | --- |
| `cart.html` | Cart page; token handling and 401 retry added. | CHANGED |
| `index.html` | Catalog; same token/`api(method, path, body, retried)` block. | CHANGED |
| `product.html` | Product page; same token/`api()` block. | CHANGED |

## Chapter 8: Run a Sale (`ch08-run-a-sale/public/`)

| File | What it is | Status |
| --- | --- | --- |
| `cart.html` | Shows Subtotal, one row per applied rule (`$68.00 → $54.40`, from `price_before_cents` and `price_after_cents`), and Total. | CHANGED |
| `index.html`, `product.html` | No change from Chapter 7. | same |

## Chapter 9: Add Customer Accounts (`ch09-add-customer-accounts/public/`)

| File | What it is | Status |
| --- | --- | --- |
| `signin.html` | Sign-in form. `POST /sessions`, stores the token and goes to the cart (the guest cart merges on the server). | NEW |
| `signup.html` | Sign-up form. `POST /users`, which returns a session; stores the token and goes to the cart. | NEW |
| `cart.html` | Adds the account area to the nav (email, Sign out, or Sign in/Sign up links). | CHANGED |
| `index.html` | Same nav change and `renderAccount()`. | CHANGED |
| `product.html` | Same nav change and `renderAccount()`. | CHANGED |

## Chapter 10: Take Orders (`ch10-take-orders/public/`)

| File | What it is | Status |
| --- | --- | --- |
| `checkout.html` | Checkout: order summary, optional coupon code, "Place order" (`POST /orders`). | NEW |
| `order.html` | Order confirmation / detail (`?id=`), `GET /orders/:id`. | NEW |
| `cart.html` | Adds a "Check out" button, shown only when the cart has items. | CHANGED |
| `index.html`, `product.html`, `signin.html`, `signup.html` | No change from Chapter 9. | same |

## Chapter 11: Let Staff Run the Shop (`ch11-let-staff-run-the-shop/public/`)

| File | What it is | Status |
| --- | --- | --- |
| `admin/products.html` | Staff list/create/edit/delete for products (`/admin/products`). | NEW |
| `admin/price_rules.html` | Staff screen for sale rules (`/admin/price_rules`). | NEW |
| `admin/coupons.html` | Staff screen for coupons (`/admin/coupons`). | NEW |
| `cart.html`, `checkout.html`, `index.html`, `order.html`, `product.html`, `signin.html`, `signup.html` | Each gets an "Admin" link in the nav for staff, found by probing `GET /admin/products` once after sign-in. | CHANGED |

The three admin pages share one script (a table driven by a `RES` config for slug and fields). Also change `public/public.go` (see above).

## Chapter 12: Archive, Don't Delete (`ch12-archive-dont-delete/public/`)

| File | What it is | Status |
| --- | --- | --- |
| `admin/products.html` | Adds an Archive / Restore link per product (`POST /admin/products/:id/archive` or `/restore`); archived rows are dimmed. | CHANGED |
| `admin/price_rules.html`, `admin/coupons.html` | Same shared script, so the file changes, but the archive code does nothing for them (it only runs when `RES.slug === "products"`). | CHANGED |
| `cart.html`, `checkout.html`, `index.html`, `order.html`, `product.html`, `signin.html`, `signup.html` | No change from Chapter 11. | same |

## Chapter 13: Organize the Catalog (`ch13-organize-the-catalog/public/`)

| File | What it is | Status |
| --- | --- | --- |
| `admin/categories.html` | Staff screen for categories (`/admin/categories`). | NEW |
| `admin/products.html` | Adds category checkboxes saved with `PUT /admin/products/:id/categories`. They are disabled for an archived product. | CHANGED |
| `admin/price_rules.html`, `admin/coupons.html` | Shared admin script gets the `onEdit`/`saveCats` hooks and a Categories nav link. | CHANGED |
| `index.html` | Category menu (`GET /categories`), `?category_id=` filter and paging (`?page=`). | CHANGED |
| `cart.html`, `checkout.html`, `order.html`, `product.html`, `signin.html`, `signup.html` | One-line fix in `api()`: ignore the query string when deciding whether a path needs a token. | CHANGED |

## What to print

Print the whole page once, where it is NEW. For CHANGED pages that only repeat the shared script block (Chapters 7, 9, 11, 13 on most pages), print the changed block or a short diff, not the whole file. If you prefer the book not to print full pages at all, point readers at `chNN-…/public/` in this repo.
