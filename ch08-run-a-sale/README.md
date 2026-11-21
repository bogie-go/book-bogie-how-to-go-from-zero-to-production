# Chapter 8: Run a Sale

Code for this chapter: the app as it stands at the end of Chapter 8.

Bogie version: v1.0.0 (created with `bogie new depot`, module name `depot`). Built on the Chapter 7 app: every visitor has a cart of their own. The chapter adds a pricing pipeline: price rules live in one table, each kind of rule is a struct behind a `Rule` interface, and the cart's price flows through the rules running now, one after another, with every step recorded.

The cart and guest sessions this folder still contains are the end of Chapters 6 and 7. Those chapters have their own folders now. Chapter 5's request tests are in `ch05-test-the-whole-app`.

## Run it

You need Go 1.26+, Docker with Compose v2, and Bogie v1.0.0 (`go install github.com/bogie-go/bogie@v1.0.0`).

```sh
cd ch08-run-a-sale
make up                              # postgres on :5440
bogie db:prepare                     # create, migrate, seed
bogie s                              # http://localhost:8080
```

Run the tests:

```sh
DEPOT_ENV=test bogie db:prepare      # once, creates the test database
bogie test
```

`bogie test` (or `make test`, or `bin/ci`) is the way to run them. Plain `go test ./...` runs packages in parallel against the one test database.

## Notes

- This folder is a complete Bogie app, a snapshot of Bogie Depot at the end of this chapter. The generated `README.md` of the app is kept as `README-bogie.md`.
- Credentials are not committed: no `config/master.key` and no `config/credentials.yml.enc`. Without them the app reads its settings from environment variables and its built-in development defaults, so it runs as is. Run `make credentials` to create your own encrypted settings. Set `TOKEN_SECRET` (at least 32 bytes) when you run the server without credentials: `export TOKEN_SECRET=$(openssl rand -hex 32)`. The tests mint their own tokens.
- Price rules live in one `price_rules` table: a `kind` and its settings as `jsonb` `properties`. `app/models/pricing.go` holds `PriceContext`, `RuleInfo` (embedded by each kind with `json:"-"`), the `Rule` interface (`Apply`, `Validate`, `Info`), `FreeItem`, the three kinds (`PercentageDiscount`, `BuyXGetY`, and `MinimumPurchase`, added at the end of the chapter), the registry `NewRule(pr)`, and `PriceCart`, the pipeline. `ListActivePriceRules` builds each row with `NewRule` and returns `(rules, skipped, err)`, so broken rows are skipped and logged. `Apply` takes the cart as it stands and returns it as the rule leaves it: each rule records its own step and names any free items it gives.
- `GET /cart` and every cart change answer with `subtotal_cents`, `discount_cents`, `total_cents`, `applied_rules`, and `free_items`. `applied_rules` lists every rule that took effect, with `price_before_cents`, `price_after_cents` and `discount_cents` (`0` for `BuyXGetY`). `free_items` has one entry per item a rule gave away, with `product_id`, `name`, `image_url`, `quantity` and `rule` (the rule's name); it is an empty list when no rule gave anything.
- `bogie db:seed` loads the launch rules (safe to run twice): "Launch weekend: 20% off everything" now for three days (priority 1), "Keycaps: buy 2, get 1 free" starting next month (priority 2), and "$5 off orders of $40 or more" now for a month (priority 3). The chapter first seeds the first two, then adds the third.
- To try the keycap deal early, as the chapter does: `make psql`, then `UPDATE price_rules SET starts_at = now() WHERE name LIKE 'Keycaps%';`. Change `priority` to see the order of the rules change the price.
- `public/cart.html` is the updated page the chapter says to copy from this folder: it shows the free items below the customer's own lines (marked free, with the rule that gave them), then the subtotal, one row per applied rule with the price before and after it, and the total. The rest of `public/` is as in Chapter 7. The book does not print these pages, so they are written for this repo.
- The migrations are numbered in the order they run: `create_products` (Chapter 3), then the cart tables (Chapters 6 and 7), then `create_price_rules`. Bogie refuses a migration older than one already applied, so a database that ran an earlier chapter's folder can take this one's new migration.
- `bin/ci` passes.
