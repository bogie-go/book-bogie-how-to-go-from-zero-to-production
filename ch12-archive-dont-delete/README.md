# Chapter 12: Archive, Don't Delete

Code for this chapter: the app as it stands at the end of Chapter 12.

Bogie version: v0.6.8 (created with `bogie new depot`, module name `depot`).

## Run it

You need Go 1.26+, Docker with Compose v2, and Bogie v0.6.8 (`go install github.com/bogie-go/bogie@v0.6.8`).

```sh
cd ch12-archive-dont-delete
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
- Products are archived, not deleted: `POST /admin/products/:id/archive` and `/restore` change only `archived_at` (a nullable column added with `bogie model:sync products`, a `*time.Time` in the model).
- `PATCH /admin/products/:id` now edits inside a lock on the product's row (`EditProduct`), so a staff edit can't undo stock a checkout took. `db/adapter/products_extra_test.go` races edits against checkouts.
- Archived products are out of the catalog (`ListCatalogProducts`, `GetCatalogProduct`), the cart, and checkout (`TakeStock` skips them).
- `public/admin/products.html` is the updated page the chapter says to copy from the companion repo (this chapter's `public/admin/products.html`): Archive and Restore buttons, archived rows greyed out. The book does not print it, so it is written for this repo.
- Migration timestamps continue after Chapter 8's `20261104094600_create_price_rules`: users `20261104094700`, orders `20261104094800`, order items `20261104094801`, coupons `20261104094900`, `role` `20261104095000`, `archived_at` `20261104095100`, categories `20261104095200`, product categories `20261104095300`. Earlier chapters in this list are omitted from a folder that does not have them yet. The book prints the `20261031…` names (`20261031072640_create_users`, `20261031073237_create_orders`, `20261031073941_add_role_to_users`, `20261031074654_add_archived_at_to_products`, `20261031075305_create_product_categories`, and the rest). The SQL is the book's. Bogie refuses a migration older than one already applied.
- `bin/ci` passes.
