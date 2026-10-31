# Chapter 2: Setting Up Go

Code for this chapter: the app as it stands at the end of Chapter 2.

This is the `depot-playground` project: a small standalone Go program and the Makefile from the chapter. It needs Go 1.26 or newer and `make`.

## Run it

```sh
cd ch02-setting-up
make run      # go run .
make check    # go fmt ./... and go vet ./...
make build    # go build -o bin/depot-playground .
./bin/depot-playground
```

Expected output: `Welcome to Bogie Depot`

## Checkpoints

The book labels each listing with a checkpoint. They are in [`checkpoints/`](checkpoints/), and each one is its own Go module named `depot-playground`, so you can `cd` into it and run it:

| Checkpoint | Listing |
| --- | --- |
| [`ch02_01`](checkpoints/ch02_01/) | `main.go`: the first program |
| [`ch02_02`](checkpoints/ch02_02/) | `Makefile` (plus `main.go`) |

The end state of the chapter, in this folder, is the same as `ch02_02`.
