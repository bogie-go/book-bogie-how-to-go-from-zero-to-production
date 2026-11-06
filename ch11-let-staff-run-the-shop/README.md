# Chapter 11: Let Staff Run the Shop

Code for this chapter: the app as it stands at the end of Chapter 11.

Bogie version: v0.6.8 (created with `bogie new depot`, module name `depot`).

## Run it

You need Go 1.26+, Docker with Compose v2, and Bogie v0.6.8 (`go install github.com/bogie-go/bogie@v0.6.8`).

```sh
cd ch11-let-staff-run-the-shop
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
- Staff can manage the shop from `/admin`: products, sales (`price_rules`) and coupons, behind `RequireStaff` (it asks the database on every request, so removing a role takes effect at once).
- Make the first staff member by hand: sign up, then in `make psql` run `UPDATE users SET role = 'staff' WHERE email = 'ada@example.com';`.
- `public/admin/products.html`, `price_rules.html` and `coupons.html` are the pages the chapter says to copy from the companion repo (this chapter's `public/admin/`), plus an "Admin" link in the header for staff. The book does not print them, so they are written for this repo. Open them at `/admin/products.html` and so on: `public/` serves no directory index.
- `public/public.go` changed from `//go:embed *.html images` to `//go:embed *.html admin/*.html images`, because the default pattern would not embed `public/admin/`. The book does not mention this.
- Migration timestamps continue after Chapter 8's `20261104094600_create_price_rules`: users `20261104094700`, orders `20261104094800`, order items `20261104094801`, coupons `20261104094900`, `role` `20261104095000`, `archived_at` `20261104095100`, categories `20261104095200`, product categories `20261104095300`. Earlier chapters in this list are omitted from a folder that does not have them yet. The book prints the `20261031…` names (`20261031072640_create_users`, `20261031073237_create_orders`, `20261031073941_add_role_to_users`, `20261031074654_add_archived_at_to_products`, `20261031075305_create_product_categories`, and the rest). The SQL is the book's. Bogie refuses a migration older than one already applied.
- `bin/ci` passes.
