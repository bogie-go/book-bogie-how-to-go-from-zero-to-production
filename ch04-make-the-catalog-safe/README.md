# Chapter 4: Make the Catalog Safe

Code for this chapter: the app as it stands at the end of Chapter 4.

Bogie version: v0.6.8. Built on the Chapter 3 app: `Product.Validate` refuses a blank name, a price that is not more than zero, and negative stock. A stock of zero is still valid. The model's tests cover each rule. The store already calls `Validate` on save, so a seed cannot skip the rules.

## Run it

You need Go 1.26+, Docker with Compose v2, and Bogie v0.6.8 (`go install github.com/bogie-go/bogie@v0.6.8`).

```sh
cd ch04-make-the-catalog-safe
make up                              # postgres on :5440
bogie db:prepare                     # create, migrate, seed
bogie s                              # http://localhost:8080
```

The seed is the same seven products as Chapter 3. The chapter adds a zero-price "Bogie Freebie", watches `bogie db:seed` refuse it with `depot: invalid: price must be more than zero`, and then removes that row from the list. This folder is the app after that removal.

Run the tests:

```sh
DEPOT_ENV=test bogie db:prepare      # once, creates the test database
bogie test
```

`bogie test` (or `make test`, or `bin/ci`) is the way to run them. Plain `go test ./...` runs packages in parallel against the one test database.

## Notes

- This folder is a complete Bogie app. The generated `README.md` of the app is kept as `README-bogie.md`.
- Credentials are not committed: no `config/master.key` and no `config/credentials.yml.enc`. Without them the app reads its settings from environment variables and its built-in development defaults, so it runs as is. Run `make credentials` to create your own encrypted settings.
- `public/index.html`, `public/product.html`, and the product images are unchanged from Chapter 3.
