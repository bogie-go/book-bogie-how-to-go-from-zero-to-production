# Chapter 18: Life Without Bogie

Code for this chapter: a small catalog API built with the standard library, not another snapshot of Bogie Depot.

Bogie version: v0.6.8. The shop is finished by Chapter 17. This folder is the playground the chapter builds beside it. The reader creates `without-bogie` next to `depot` and runs `go mod init without-bogie`. The listings in the book are labeled under this folder's `checkpoints/`.

## Run it

You need Go 1.26+.

```sh
cd ch18-life-without-bogie
go run .      # or: make run
```

In a second terminal:

```sh
curl localhost:8080/products
curl localhost:8080/products/2
curl -i localhost:8080/products/99
curl -X POST localhost:8080/products -d '{"name": "Bogie Keycap", "price_cents": 1800, "stock": 80}'
curl -X POST localhost:8080/products -d '{"price_cents": 1800}'
```

`main.go` is the program at the end of the chapter (the same as `checkpoints/ch18_06`). Products live in memory, so a restart throws them away.

## Checkpoints

Each listing is `checkpoints/ch18_NN/main.go`. Where the book shows `// ...`, the checkpoint is the whole file as it stands after that step. From this folder:

```sh
go run ./checkpoints/ch18_01
go run ./checkpoints/ch18_02
go run ./checkpoints/ch18_03
go run ./checkpoints/ch18_04
go run ./checkpoints/ch18_05
go run ./checkpoints/ch18_06
```

| Checkpoint | What it is |
| --- | --- |
| `ch18_01` | `GET /hello` answers `Hello from Bogie Depot` |
| `ch18_02` | `GET /products` lists the three products as JSON |
| `ch18_03` | `GET /products/{id}`; a missing id is `404` and plain text `product not found` |
| `ch18_04` | `POST /products` with `Validate`; a blank name is plain text `name is required` |
| `ch18_05` | `writeJSON` and `writeError`; a refused product is `{"error":"name is required"}` |
| `ch18_06` | `logRequests` wraps the mux |
