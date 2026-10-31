# Chapter 5: A Web Server with Go's Standard Library

Code for this chapter: the app as it stands at the end of Chapter 5.

This is the `depot-playground` product API built with `net/http` only. The folder holds the end state of the chapter (checkpoint `ch05_06`). It needs Go 1.26 or newer.

## Run it

```sh
cd ch05-web-server-standard-library
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

Products live in memory, so a restart throws them away. That is the point of the chapter's closing section, and the reason Chapter 6 brings in Bogie and a database.

## Checkpoints

Each one is its own module named `depot-playground`, in [`checkpoints/`](checkpoints/):

| Checkpoint | Listing |
| --- | --- |
| [`ch05_01`](checkpoints/ch05_01/) | the smallest server (`GET /hello`) |
| [`ch05_02`](checkpoints/ch05_02/) | `GET /products` |
| [`ch05_03`](checkpoints/ch05_03/) | `GET /products/{id}` |
| [`ch05_04`](checkpoints/ch05_04/) | `POST /products` |
| [`ch05_05`](checkpoints/ch05_05/) | `writeJSON` and `writeError` |
| [`ch05_06`](checkpoints/ch05_06/) | the `logRequests` middleware |

Where a listing shows only part of a file, the checkpoint completes it with the code the chapter built earlier.
