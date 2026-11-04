# Chapter 1: Hello, Bogie Depot

Code for this chapter: the app as it stands at the end of Chapter 1.

Bogie version: v0.6.8 (created with `bogie new depot`, module name `depot`).

## Run it

You need Go 1.26+, Docker with Compose v2, and Bogie v0.6.8 (`go install github.com/bogie-go/bogie@v0.6.8`).

```sh
cd ch01-hello-bogie-depot
make up                              # postgres on :5440
bogie db:prepare                     # create and migrate
bogie s                              # http://localhost:8080
```

The two products in the chapter are inserted by hand in `make psql`. They are not part of this folder. Seeds arrive in Chapter 3.

Run the tests:

```sh
DEPOT_ENV=test bogie db:prepare      # once, creates the test database
bogie test
```

## Notes

- This folder is a complete Bogie app, generated with `bogie new depot`, then `bogie g model product name description:text price_cents:integer stock:integer image_url` and `bogie g scaffold_controller products --only index,show`. The generated `README.md` of the app is kept as `README-bogie.md`.
- Credentials are not committed: no `config/master.key` and no `config/credentials.yml.enc`. Without them the app reads its settings from environment variables and its built-in development defaults, so it runs as is. Run `make credentials` to create your own encrypted settings.
- The migration filename's timestamp is whatever `bogie g model` wrote. The book's listing uses a different timestamp, and the chapter says yours will differ.
