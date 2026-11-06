# Chapter 2: Just Enough Go

Code for this chapter: the Chapter 1 app, unchanged, plus the scratch programs the chapter runs.

Bogie version: v0.6.8. The shop is the same tree as [`ch01-hello-bogie-depot`](../ch01-hello-bogie-depot/). The chapter deletes `scratch/` at the end, so this folder does not contain one.

## Run the shop

Same as Chapter 1:

```sh
cd ch02-just-enough-go
make up
bogie db:prepare
bogie s
```

## Checkpoints

Each step the chapter labels is a `package main` under `checkpoints/`. From this folder:

```sh
go run ./checkpoints/ch02_01/scratch
go run ./checkpoints/ch02_02/scratch
go run ./checkpoints/ch02_03/scratch
go run ./checkpoints/ch02_04/scratch
go run ./checkpoints/ch02_05/scratch
go run ./checkpoints/ch02_06/scratch
go run ./checkpoints/ch02_07/scratch
```

| Checkpoint | What it prints |
| --- | --- |
| `ch02_01` | the zero `models.Product` |
| `ch02_02` | a mug's name, price, stock `0`, and `Validate` (`<nil>`) |
| `ch02_03` | `parsePrice` on `"1400"` and `"cheap"` |
| `ch02_04` | `0.30000000000000004`, then `$14.00 $3.05` and `$21474836.47` |
| `ch02_05` | `150` — `Take` on a value receiver does not change the product |
| `ch02_06` | a slice, and a nil slice versus an empty one |
| `ch02_07` | sold out, took 3, found, no keycaps, stock `2 497` |

`ch02_01` and `ch02_02` import `depot/app/models`, so they are run from this module. `ch02_05` is the value-receiver program the chapter labels. The pointer receiver and the `range` trap are edits the chapter makes in place; they are not separate checkpoints. `ch02_06` is the slice program the chapter labels. The map example that follows is an edit to that same `main`.
