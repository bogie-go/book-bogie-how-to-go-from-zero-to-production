# Chapter 4: Go Essentials, Part 2

Code for this chapter: the app as it stands at the end of Chapter 4.

This is the `depot-playground` module with the `catalog` package. The folder holds the end state of the chapter (checkpoint `ch04_06`). It needs Go 1.26 or newer.

## Run it

```sh
cd ch04-go-essentials-part-2
go run .      # or: make run
```

Expected output:

```
context deadline exceeded
Bogie Mug 13
Bogie Sticker 10
Bogie Desk Mat 5
```

## Checkpoints

Each one is its own module named `depot-playground`, in [`checkpoints/`](checkpoints/):

| Checkpoint | Topic |
| --- | --- |
| [`ch04_01`](checkpoints/ch04_01/) | the `catalog` package |
| [`ch04_02`](checkpoints/ch04_02/) | a `String` method |
| [`ch04_03`](checkpoints/ch04_03/) | an interface declared where it is used (`Priced`) |
| [`ch04_04`](checkpoints/ch04_04/) | embedding (`Timestamps` in `Order`) |
| [`ch04_05`](checkpoints/ch04_05/) | functions as values, closures |
| [`ch04_06`](checkpoints/ch04_06/) | `context` with a timeout |

Where a listing shows only part of a file, the checkpoint completes it with the code the chapter built earlier.
