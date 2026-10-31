# Chapter 11: Task C: Run a Sale

Code for this chapter: the app as it stands at the end of Chapter 11.

Bogie version: v0.6.8 (created with `bogie new depot`, module name `depot`).

## Run it

You need Go 1.26+, Docker with Compose v2, and Bogie v0.6.8 (`go install github.com/bogie-go/bogie@v0.6.8`).

```sh
cd ch11-*/
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
- Price rules live in one `price_rules` table (`kind` plus `properties` as `jsonb`). `app/models/pricing.go` holds the `Rule` interface, the three kinds (`PercentOff`, `BuyXGetY`, and `AmountOff` added at the end of the chapter), the registry `PriceRule.Rule()` and `PriceCart`.
- `GET /cart` and every cart change now answer with `subtotal_cents`, `discount_cents`, `total_cents` and `applied_rules`.
- Load the launch rules with `bogie db:seed` (safe to run twice). The seed starts "20% off everything" now for three days, the keycap deal next month, and "$5 off orders of $40 or more" now for a month.
- To try the keycap deal early, as in the chapter: `make psql`, then `UPDATE price_rules SET starts_at = now() WHERE name LIKE 'Keycaps%';`.
- `public/cart.html` is the updated page the chapter says to copy from the companion repo (checkpoint `ch11_06`): it lists each applied rule above the total. The book does not print it, so it is written for this repo.
- `bin/ci` passes.
