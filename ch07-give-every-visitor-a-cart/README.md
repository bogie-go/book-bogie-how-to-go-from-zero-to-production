# Chapter 7: Give Every Visitor a Cart

Code for this chapter: the app as it stands at the end of Chapter 7.

Bogie version: v1.0.0 (created with `bogie new depot`, module name `depot`).

## Run it

You need Go 1.26+, Docker with Compose v2, and Bogie v1.0.0 (`go install github.com/bogie-go/bogie@v1.0.0`).

```sh
cd ch07-give-every-visitor-a-cart
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
- Every visitor gets their own cart: `POST /guest_sessions` mints a signed guest token, and `/cart` and `/cart_items` sit behind `WithUser`. The cart is found from the token's subject (`middlewares.Subject`).
- Try it: `TOKEN=$(curl -s -X POST localhost:8080/guest_sessions | jq -r .token)` then `curl -H "Authorization: Bearer $TOKEN" localhost:8080/cart`.
- Set a token secret of at least 32 bytes in the environment when running without credentials: `export TOKEN_SECRET=$(openssl rand -hex 32)`. The tests mint their own tokens and don't need it.
- The "Someone else's cart item" specs (`guest:alice` / `guest:mallory`) and the scoped queries `GetCartItemInCart` / `DeleteCartItemInCart` are the bug and fix from the chapter.
- `public/*.html` are the updated pages the chapter says to copy from the companion repo (this chapter's `public/`): they get a guest token once, keep it in local storage and send it with each cart request. The book does not print them, so they are written for this repo.
- The cart migrations use the same filenames as Chapter 6, matching the book's listings: `20261031071536_create_carts.sql`, `20261031071537_create_cart_items.sql`, and `20261031071548_add_unique_product_per_cart.sql`.
- `bin/ci` passes.
