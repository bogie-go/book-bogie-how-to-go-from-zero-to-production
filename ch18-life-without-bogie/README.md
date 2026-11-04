# Chapter 18: Life Without Bogie

Code for this chapter: the closest snapshot, the old Chapter 5 playground.

This is the `depot-playground` product API built with `net/http` only. The folder holds the end state of that chapter (checkpoint `ch18_06`). It needs Go 1.26 or newer. The new Chapter 18 text is not written yet; this is the old standard-library chapter, kept because that is what Chapter 18 retells.

## Run it

```sh
cd ch18-life-without-bogie
go run .      # or: make run
```

In a second terminal:

```sh
curl localhost:8080/products
curl localhost:8080/products/1
curl -i localhost:8080/products/99
curl -X POST localhost:8080/products -d '{"name": "Bogie Keycap", "price_cents": 900, "stock": 10}'
curl -X POST localhost:8080/products -d '{"price_cents": 900}'
```

Products live in memory, so a restart throws them away. That is the point of the chapter's closing section. Chapter 1 is where Bogie and a database take over.

## Checkpoints

Each one is its own module named `depot-playground`, in [`checkpoints/`](checkpoints/). The directories were `ch05_01` … `ch05_06`; they are renamed to this chapter's number. The programs are unchanged.

| Checkpoint | Listing |
| --- | --- |
| [`ch18_01`](checkpoints/ch18_01/) | the smallest server (`GET /hello`) |
| [`ch18_02`](checkpoints/ch18_02/) | `GET /products` |
| [`ch18_03`](checkpoints/ch18_03/) | `GET /products/{id}` |
| [`ch18_04`](checkpoints/ch18_04/) | `POST /products` |
| [`ch18_05`](checkpoints/ch18_05/) | `writeJSON` and `writeError` |
| [`ch18_06`](checkpoints/ch18_06/) | the `logRequests` middleware |

Where a listing shows only part of a file, the checkpoint completes it with the code the chapter built earlier.
