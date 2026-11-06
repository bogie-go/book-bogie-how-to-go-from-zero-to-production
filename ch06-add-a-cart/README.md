# Chapter 6: Add a Cart

Code for this chapter: the app as it stands at the end of Chapter 6.

Bogie version: v0.6.8 (created with `bogie new depot`, module name `depot`).

## Run it

You need Go 1.26+, Docker with Compose v2, and Bogie v0.6.8 (`go install github.com/bogie-go/bogie@v0.6.8`).

```sh
cd ch06-add-a-cart
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
- The cart is one shared cart for every visitor at this point (the shared token is an empty string). Chapter 7 gives each visitor their own.
- Try it: `curl -X POST localhost:8080/cart_items -d '{"product_id": 3, "quantity": 1}'` and `curl localhost:8080/cart` (use a product id from `GET /products`).
- `public/index.html`, `public/product.html` and `public/cart.html` are the updated pages the chapter says to copy from the companion repo (this chapter's `public/`): a cart badge, an "Add to cart" button, and a cart page with a quantity stepper, a remove link and the total. The book does not print them, so they are written for this repo as plain HTML with inline CSS and script.
- The cart migrations are `20261104094512_create_carts`, `20261104094513_create_cart_items`, and `20261104094524_add_unique_product_per_cart`, after Chapter 4's `create_products`. The book prints `20261031071536`, `20261031071537`, and `20261031071548`. The SQL is the book's. Bogie refuses a migration older than one already applied.
- `bin/ci` passes.
