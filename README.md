# Bogie: How to Go from Zero to Production: Companion Code

Companion code for the book **Bogie: How to Go from Zero to Production**.

**Pinned Bogie version: v0.6.8.** Every chapter folder targets this version.

## How to use

Each chapter folder holds the app as it stands at the end of that chapter. To follow along or check your work, `cd` into the chapter you are reading:

```sh
git clone git@github.com:bogie-go/book-bogie-how-to-go-from-zero-to-production.git
cd book-bogie-how-to-go-from-zero-to-production/ch03-show-the-catalog
```

Each chapter folder has its own README with how to run that chapter's code.

Chapters 1 and 3 are Bogie Depot, generated with Bogie v0.6.8 and left as those chapters leave the app. Chapter 2 is that same app plus `checkpoints/` for the scratch programs (the chapter deletes `scratch/` at the end). Chapter 18 is a standalone `net/http` playground, also with `checkpoints/`. Chapters 4–17 are full Bogie Depot snapshots. Chapters 19 and 20 are README only (no new code in those chapters).

Product images (one master copy; the Depot chapters carry their own copies): [`images/`](images/). Book feedback from building this repo lives in [`BOOK-ISSUES.md`](BOOK-ISSUES.md). Which UI pages to print is [`UI-PAGES.md`](UI-PAGES.md).

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

## Layout notes

Chapters 1–3 were generated fresh against Bogie v0.6.8. From Chapter 4 on, each folder is the old chapter's code moved to its new number. Two of those moves are the closest snapshot rather than a clean cut:

- **Chapter 4** is the end of the old Chapter 7. The shop's `Validate` rules are here. The hand-written model tests Chapter 4 adds were not saved on their own; they first appear in Chapter 5.
- **Chapter 5** is the end of the old Chapter 8, so it has those model tests and the rest of the testing chapter.
- **Chapter 18** is the old Chapter 5 `net/http` playground. Its step folders are `checkpoints/ch18_01` through `checkpoints/ch18_06`.

The old setting-up chapter, the two Go-essentials chapters, the empty "meet Bogie" app, and the tools chapter are not folders anymore. Chapter 2 replaces the essentials playgrounds, Chapter 1 replaces the empty app, and the tools chapter moves to Appendix C.
