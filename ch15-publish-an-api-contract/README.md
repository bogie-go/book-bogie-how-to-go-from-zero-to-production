# Chapter 15: Publish an API Contract

Code for this chapter: the app as it stands at the end of Chapter 15.

Bogie version: v0.6.8 (created with `bogie new depot`, module name `depot`).

## Run it

You need Go 1.26+, Docker with Compose v2, and Bogie v0.6.8 (`go install github.com/bogie-go/bogie@v0.6.8`).

```sh
cd ch15-publish-an-api-contract
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
- `docs/openapi.yaml` describes the whole public API (catalog, categories, guest sessions, sign up and sign in, cart, cart items, orders), every documented status included, with `additionalProperties: false` on every body. `/admin` and `/healthz` are not part of it.
- `test/requests/contract_test.go` sends real requests through the app and checks each response against the document with `kin-openapi`, and checks that the routes served are exactly the routes documented.
- The book shows only the start of `docs/openapi.yaml` and says to copy the rest from the companion repo; this file is written for this repo from the real controllers and views. The book also shows only the first two contract tests; the remaining ones (sessions, accounts, cart items, orders, categories, error responses) are written here.
- Both breaks from the chapter were run on the box: renaming the `price_cents` tag fails the catalog test (`property "price" is unsupported`), and `bogie g controller wishlists index` fails the routes test (`GET /wishlists` served, not documented). Both were reverted.
- The Playtime (`public/docs.html` with Redoc) is not included.
- Migration timestamps continue after Chapter 8's `20261104094600_create_price_rules`: users `20261104094700`, orders `20261104094800`, order items `20261104094801`, coupons `20261104094900`, `role` `20261104095000`, `archived_at` `20261104095100`, categories `20261104095200`, product categories `20261104095300`. Earlier chapters in this list are omitted from a folder that does not have them yet. The book prints the `20261031…` names (`20261031072640_create_users`, `20261031073237_create_orders`, `20261031073941_add_role_to_users`, `20261031074654_add_archived_at_to_products`, `20261031075305_create_product_categories`, and the rest). The SQL is the book's. Bogie refuses a migration older than one already applied.
- `bin/ci` passes.
