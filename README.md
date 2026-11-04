# Bogie: How to Go from Zero to Production: Companion Code

Companion code for the book **Bogie: How to Go from Zero to Production**.

**Pinned Bogie version: v0.6.8.** Every chapter folder targets this version.

## How to use

Only Chapters 1 to 3 have code so far. To follow along or check your work, `cd` into the chapter you are reading:

```sh
git clone git@github.com:bogie-go/book-bogie-how-to-go-from-zero-to-production.git
cd book-bogie-how-to-go-from-zero-to-production/ch03-show-the-catalog
```

Each of those folders holds the app as it stands at the end of that chapter, and its README says how to run it. Chapters 1 and 3 are Bogie Depot, generated with Bogie v0.6.8. Chapter 2 is that same app plus `checkpoints/` for the scratch programs (the chapter deletes `scratch/` at the end).

Chapters 4 to 20 are placeholders. Each folder is a README until that chapter is written. Earlier snapshots of those chapters stay in git history.

Product images for the written chapters: [`images/`](images/) is the master copy, and Chapter 3 carries them in `public/images/products/`. Book feedback from building this repo lives in [`BOOK-ISSUES.md`](BOOK-ISSUES.md). Which UI pages to print is [`UI-PAGES.md`](UI-PAGES.md).

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
