# Chapter 7: Task A: Show Our Merch

Code for this chapter: the app as it stands at the end of Chapter 7.

Bogie version: v0.6.8 (created with `bogie new depot`, module name `depot`).

## Run it

You need Go 1.26+, Docker with Compose v2, and Bogie v0.6.8 (`go install github.com/bogie-go/bogie@v0.6.8`).

```sh
cd ch07-*/
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
- Seed data: `bogie db:seed` (also run by `bogie db:prepare`). It is safe to run again: seven products, not fourteen.
- `public/index.html`, `public/product.html` and `public/images/products/*.webp` are the UI files the chapter says to copy from the companion repo (checkpoint `ch07_05`). The book does not print them, so they are written for this repo: plain HTML that fetches `GET /products` and `GET /products/:id`, with the CSS inline, because Bogie embeds only `*.html` and `images` from `public/`. The product images are generated placeholders, not photos.
- The tests in this chapter are the ones `bogie g` generated. Chapter 8 writes the real ones.
