# Book images

The shop's product images, in one place. Every Depot chapter from 3 to 17 carries the same seven files in its own `public/images/products/` (Bogie embeds only `public/`, so each chapter needs its own copy); this folder is the master set for the book. Chapters 1 and 2 have only the 404 and 500 images from `bogie new`. Chapter 18 is the standard-library playground and does not use these files.

| File | Product (seed row) |
| --- | --- |
| `products/bogie-logo-sticker.webp` | Bogie Logo Sticker |
| `products/bogie-sticker-pack.webp` | Bogie Sticker Pack |
| `products/bogie-ceramic-mug.webp` | Bogie Ceramic Mug |
| `products/bogie-desk-mat.webp` | Bogie Desk Mat |
| `products/bogie-keycap.webp` | Bogie Keycap |
| `products/bogie-poster.webp` | Bogie Poster |
| `products/bogie-notebook.webp` | Bogie Notebook |
| `bogie-404.webp`, `bogie-500.webp` | Bogie mascot for the 404 and 500 pages (also in every chapter's `public/images/`) |

The product images are 800×533 WebP **placeholders** (a flat colour, the product name and "BOGIE DEPOT / placeholder image"), copied from the Bogie Depot app in [`bogie-go/zero-to-production`](https://github.com/bogie-go/zero-to-production) (`public/images/products/`, the images `db/seeds/product.go` points at). They are original artwork made for the project, not stock photos or third-party images. That repo has no LICENSE file, so add one (or state the licence here) before the book is published.
