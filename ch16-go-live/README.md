# Chapter 16: Go Live

Code for this chapter: the app as it stands at the end of Chapter 16.

Bogie version: v0.6.8 (created with `bogie new depot`, module name `depot`).

## Run it

You need Go 1.26+, Docker with Compose v2, and Bogie v0.6.8 (`go install github.com/bogie-go/bogie@v0.6.8`).

```sh
cd ch16-go-live
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
- Deployment config only; nothing was deployed (no server, domain or registry account exists for this book). Every value is a placeholder: the server is `203.0.113.10` (the reserved documentation address the chapter uses), the hosts are `staging-depot.bogie-go.com` and `depot.bogie-go.com`, the image is `roonglit/bogie-depot` and the registry user is `your-docker-hub-username`. Replace them with your own.
- `config/deploy.yml` has the Postgres accessory uncommented and pinned to the server with `host:`; `.kamal/secrets` and `.kamal/secrets.production` each gain a `POSTGRES_PASSWORD` line that reads it from the encrypted credentials. No secret value is in any file: the credentials (`config/credentials/*.yml.enc`) and their keys (`*.key`) are not committed, so create yours with `make credentials ENV=staging` and `make credentials ENV=production`, as the chapter says.
- Checked on the box: `kamal config` and `kamal config -d production` (Kamal 2.12.0, with dummy secrets in a scratch copy) accept both files, and `docker build .` produces the image (about 38 MB).
- Migration timestamps continue after Chapter 8's `20261104094600_create_price_rules`: users `20261104094700`, orders `20261104094800`, order items `20261104094801`, coupons `20261104094900`, `role` `20261104095000`, `archived_at` `20261104095100`, categories `20261104095200`, product categories `20261104095300`. Earlier chapters in this list are omitted from a folder that does not have them yet. The book prints the `20261031…` names (`20261031072640_create_users`, `20261031073237_create_orders`, `20261031073941_add_role_to_users`, `20261031074654_add_archived_at_to_products`, `20261031075305_create_product_categories`, and the rest). The SQL is the book's. Bogie refuses a migration older than one already applied.
- `bin/ci` is unchanged from Chapter 15 (no Go code changed).
