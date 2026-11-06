# Bogie: How to Go from Zero to Production: Companion Code

Companion code for the book **Bogie: How to Go from Zero to Production**.

**Pinned Bogie version: v0.6.8.** Every chapter folder targets this version.

## How to use

To follow along or check your work, `cd` into the chapter you are reading:

```sh
git clone git@github.com:bogie-go/book-bogie-how-to-go-from-zero-to-production.git
cd book-bogie-how-to-go-from-zero-to-production/ch03-show-the-catalog
```

Each folder holds the app (or the programs) as it stands at the end of that chapter, and its README says how to run it.

Chapters 1 and 3 through 17 are Bogie Depot, one snapshot per chapter, on Bogie v0.6.8. Chapter 2 is the Chapter 1 app plus `checkpoints/` for the scratch programs (the chapter deletes `scratch/` at the end). Chapters 5 through 7 and 9 through 17 continue that shop: tests, a cart, guest sessions, the sale from Chapter 8, accounts, orders, staff, archive, categories, background email, the API contract, deploy config, and `/status`. Chapters 12, 13, and 15 are complete apps on purpose: later chapters tell a reader who skipped them to start from those folders.

Chapter 18 is not another Depot snapshot. It is the standard-library playground, with `checkpoints/ch18_01` through `ch18_06`. Chapters 19 and 20 have no code of their own. Chapter 19 is read with the finished shop in Chapter 17 open beside it.

Product images: [`images/`](images/) is the master copy, and the catalog chapters carry them in `public/images/products/`. Book feedback from building this repo lives in [`BOOK-ISSUES.md`](BOOK-ISSUES.md). Which UI pages to print is [`UI-PAGES.md`](UI-PAGES.md).

Migration filenames in Chapters 5 through 17 use timestamps that sort after Chapter 4's `create_products` (and, from Chapter 9 on, after Chapter 8's `create_price_rules`). Bogie refuses a migration older than one already applied, so a database prepared in an earlier chapter can take the next chapter's new migrations. The book prints earlier timestamps. The SQL matches the book. Each chapter README notes the names it uses.

## Chapters

| Chapter | Title | Folder |
| --- | --- | --- |
| 1 | Hello, Bogie Depot | [`ch01-hello-bogie-depot/`](ch01-hello-bogie-depot/) |
| 2 | Just Enough Go | [`ch02-just-enough-go/`](ch02-just-enough-go/) |
| 3 | Show the Catalog | [`ch03-show-the-catalog/`](ch03-show-the-catalog/) |
| 4 | Make the Catalog Safe | [`ch04-make-the-catalog-safe/`](ch04-make-the-catalog-safe/) |
| 5 | Test the Whole App | [`ch05-test-the-whole-app/`](ch05-test-the-whole-app/) |
| 6 | Add a Cart | [`ch06-add-a-cart/`](ch06-add-a-cart/) |
| 7 | Give Every Visitor a Cart | [`ch07-give-every-visitor-a-cart/`](ch07-give-every-visitor-a-cart/) |
| 8 | Run a Sale | [`ch08-run-a-sale/`](ch08-run-a-sale/) |
| 9 | Add Customer Accounts | [`ch09-add-customer-accounts/`](ch09-add-customer-accounts/) |
| 10 | Take Orders | [`ch10-take-orders/`](ch10-take-orders/) |
| 11 | Let Staff Run the Shop | [`ch11-let-staff-run-the-shop/`](ch11-let-staff-run-the-shop/) |
| 12 | Archive, Don't Delete | [`ch12-archive-dont-delete/`](ch12-archive-dont-delete/) |
| 13 | Organize the Catalog | [`ch13-organize-the-catalog/`](ch13-organize-the-catalog/) |
| 14 | Send Email in the Background | [`ch14-send-email-in-the-background/`](ch14-send-email-in-the-background/) |
| 15 | Publish an API Contract | [`ch15-publish-an-api-contract/`](ch15-publish-an-api-contract/) |
| 16 | Go Live | [`ch16-go-live/`](ch16-go-live/) |
| 17 | Keep the Shop Healthy | [`ch17-keep-the-shop-healthy/`](ch17-keep-the-shop-healthy/) |
| 18 | Life Without Bogie | [`ch18-life-without-bogie/`](ch18-life-without-bogie/) |
| 19 | What Bogie Generates, and Why | [`ch19-what-bogie-generates/`](ch19-what-bogie-generates/) |
| 20 | When to Break the Conventions | [`ch20-break-the-conventions/`](ch20-break-the-conventions/) |
