# Chapter 13: Task E: Let Customers Place Orders

Code for this chapter: the app as it stands at the end of Chapter 13.

Bogie version: v0.6.8 (created with `bogie new depot`, module name `depot`).

## Run it

You need Go 1.26+, Docker with Compose v2, and Bogie v0.6.8 (`go install github.com/bogie-go/bogie@v0.6.8`).

```sh
cd ch13-*/
make up                              # postgres on :5440
bogie db:prepare                     # create, migrate, seed
bogie s                              # http://localhost:8080
```

Run the tests:

```sh
DEPOT_ENV=test bogie db:prepare      # once, creates the test database
make test
```

## Notes

- This folder is a complete Bogie app, a snapshot of Bogie Depot at the end of this chapter. The generated `README.md` of the app is kept as `README-bogie.md`.
- Credentials are not committed: no `config/master.key` and no `config/credentials.yml.enc`. Without them the app reads its settings from environment variables and its built-in development defaults, so it runs as is. Run `make credentials` to create your own encrypted settings.
- Customers place orders with `POST /orders`: one transaction takes stock, writes the order and its lines (prices saved as they were), and empties the cart. Guests get a 403; `GET /orders/:id` shows a customer their own receipt only.
- Coupons are a `coupons` table pointing at a price rule; send `{"coupon_code": "LAUNCH5"}` in the body. A rule with a coupon only applies to orders that use its code. The book suggests adding a coupon to the seeds to try it in the browser; none is seeded here, so create one in `make psql` (insert a price rule, then a coupon with `price_rule_id`).
- `db/adapter/orders_extra_test.go` holds the tests that race several checkouts at once. Remove `LockCart` or the sort by product ID from `PlaceOrder` and they fail.
- `public/checkout.html`, `public/order.html` and the "Check out" button in `public/cart.html` are the pages the chapter says to copy from the companion repo (checkpoint `ch13_07`). The book does not print them, so they are written for this repo. The chapter's "Playtime" exercises (`GET /orders`, moving the cart-merge deletes into `withTx`) are not part of this code.
- `bin/ci` passes.
