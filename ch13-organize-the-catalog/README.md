# Chapter 13: Organize the Catalog

Code for this chapter: the app as it stands at the end of Chapter 13.

Bogie version: v1.0.0 (created with `bogie new depot`, module name `depot`).

## Run it

You need Go 1.26+, Docker with Compose v2, and Bogie v1.0.0 (`go install github.com/bogie-go/bogie@v1.0.0`).

```sh
cd ch13-organize-the-catalog
make up                              # postgres on :5440
bogie db:prepare                     # create, migrate, seed
bogie s                              # http://localhost:8080
```

Run the tests:

```sh
DEPOT_ENV=test bogie db:prepare      # once, creates the test database
make test
```


## If you're following along

- In "Products in Many Categories" (p.265–266), stop `bogie s` before you run `bogie g model product_category` and edit the migration. The server applies new migrations as it reloads, so the unedited file can run first. If it already has, run `bogie db:rollback`, add `ON DELETE CASCADE` and the unique pair index, then `bogie db:migrate`.
- In "Count with a LEFT JOIN" (p.277), the printed `Describe("GET /categories")` puts its `BeforeEach` at Describe level. If you keep the generated empty-list example in the same Describe, that setup will make it fail. Put the new examples in a `Context`, or replace the whole generated Describe.

## Notes

- This folder is a complete Bogie app, a snapshot of Bogie Depot at the end of this chapter. The generated `README.md` of the app is kept as `README-bogie.md`.
- Credentials are not committed: no `config/master.key` and no `config/credentials.yml.enc`. Without them the app reads its settings from environment variables and its built-in development defaults, so it runs as is. Run `make credentials` to create your own encrypted settings.
- Categories: public `GET /categories` (index only, with a `product_count` per category, archived products not counted) and a full staff resource under `/admin/categories`.
- Products and categories are many-to-many through the `product_categories` join table (`ON DELETE CASCADE`, unique pair). `CategoriesForProduct` is the JOIN; `SetProductCategories` replaces a product's whole set in one transaction (`PUT /admin/products/:id/categories`).
- `GET /products?category_id=N&page=N&per_page=N` filters with an optional `sqlc.narg` and `EXISTS`, and pages with `LIMIT`/`OFFSET` plus a `CountCatalogProducts` total (default 12 per page, at most 48). `GET /products/:id` lists the product's categories.
- The UI pages in `public/` are written for this repo (the pages are in this chapter's `public/`; the book does not print them): the catalog has a category menu with counts and numbered pages; the admin sidebar gains a Categories link; the product form (`admin/product.html`) has category checkboxes (saved with the PUT; an archived product's boxes are disabled because the public product page does not show archived products); and `public/admin/categories.html` and `category.html` are new.
- Migration filenames match the book and stay the same from chapter to chapter (`20261031071211_create_products.sql` through `20261031075305_create_product_categories.sql`). Earlier files in the chain are omitted from a folder that does not have them yet. If a database still has this repo's older `20261104…` versions, re-prepare it: Bogie refuses a migration older than one already applied.
- `bin/ci` passes.
