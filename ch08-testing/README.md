# Chapter 8: Testing: Proving Task A Works

Code for this chapter: the app as it stands at the end of Chapter 8.

Bogie version: v0.6.8 (created with `bogie new depot`, module name `depot`).

## Run it

You need Go 1.26+, Docker with Compose v2, and Bogie v0.6.8 (`go install github.com/bogie-go/bogie@v0.6.8`).

```sh
cd ch08-*/
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
- Tests: `bin/ci` runs every check (it needs golangci-lint: `go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest`). `bin/ci` passes at the end of this chapter.
- The suite runs with `make test` or `bogie test`, not plain `go test ./...` (the packages share one test database).
- The deliberate mistake the chapter makes (`PriceCents < 0`) is fixed here: the rule reads `p.PriceCents <= 0`, and the model test catches the mistake.
