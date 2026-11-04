# Chapter 3: Show the Catalog

Code for this chapter: the app as it stands at the end of Chapter 3.

Bogie version: v0.6.8. Built on the Chapter 1 app: timestamps removed from `ProductResponse`, `GetProductByName` and `FindOrCreateProductByName`, the seven-product seed, and `public/index.html`, `public/product.html`, and the product images.

## Run it

You need Go 1.26+, Docker with Compose v2, and Bogie v0.6.8 (`go install github.com/bogie-go/bogie@v0.6.8`).

```sh
cd ch03-show-the-catalog
make up                              # postgres on :5440
bogie db:prepare                     # create, migrate, seed
bogie s                              # http://localhost:8080
```

`bogie db:prepare` runs the seeds. Running `bogie db:seed` again still leaves seven products.

Run the tests:

```sh
DEPOT_ENV=test bogie db:prepare      # once, creates the test database
bogie test
```

## Notes

- This folder is a complete Bogie app. The generated `README.md` of the app is kept as `README-bogie.md`.
- Credentials are not committed: no `config/master.key` and no `config/credentials.yml.enc`. Without them the app reads its settings from environment variables and its built-in development defaults, so it runs as is. Run `make credentials` to create your own encrypted settings.
- `public/index.html` and `public/product.html` are the catalog pages. They fetch `GET /products` and `GET /products/:id`. A product whose image file is missing still renders; the `<img>` keeps that `image_url`, and the browser shows its broken-image icon (and the product name as alt text). The pages do not hide the image or substitute a placeholder.
- The seven `public/images/products/*.webp` files are the same placeholders as in [`images/`](../images/).
