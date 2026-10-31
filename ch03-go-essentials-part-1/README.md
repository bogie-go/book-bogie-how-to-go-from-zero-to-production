# Chapter 3: Go Essentials, Part 1

Code for this chapter: the app as it stands at the end of Chapter 3.

This is the `depot-playground` module from Chapter 2, with everything in one `main.go`. The folder holds the end state of the chapter (checkpoint `ch03_11`). It needs Go 1.26 or newer.

## Run it

```sh
cd ch03-go-essentials-part-1
go run .      # or: make run
```

Expected output: `Sorry: selling 10 of "Bogie Mug": out of stock`

## Checkpoints

Each listing in the chapter replaces `main.go`. The versions as they stand after each listing are in [`checkpoints/`](checkpoints/), and each one is its own module named `depot-playground`:

| Checkpoint | Topic |
| --- | --- |
| [`ch03_01`](checkpoints/ch03_01/) | declaring values |
| [`ch03_02`](checkpoints/ch03_02/) | zero values |
| [`ch03_03`](checkpoints/ch03_03/) | functions with more than one result |
| [`ch03_04`](checkpoints/ch03_04/) | structs and methods |
| [`ch03_05`](checkpoints/ch03_05/) | struct tags and JSON |
| [`ch03_06`](checkpoints/ch03_06/) | a value receiver (the stock does not change) |
| [`ch03_07`](checkpoints/ch03_07/) | a pointer receiver |
| [`ch03_08`](checkpoints/ch03_08/) | slices and `for` |
| [`ch03_09`](checkpoints/ch03_09/) | maps |
| [`ch03_10`](checkpoints/ch03_10/) | errors as values (`Validate`) |
| [`ch03_11`](checkpoints/ch03_11/) | sentinel errors (`Sell`, `ErrOutOfStock`) |

The chapter's listings collapse unchanged code to `// ...`, so each checkpoint is the listing completed with the code the chapter built earlier, so that it compiles and runs.
