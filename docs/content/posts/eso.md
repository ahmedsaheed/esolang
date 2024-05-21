+++
title = "Esolang"
date = "2024-05-22T00:01:24+01:00"
author = "Ahmed Saheed"
cover = "img/hello.jpg"
showFullContent = false
readingTime = false
hideComments = false
+++

Esolang is a beginner friendly dynamically typed scripting language for quick prototyping and learning functional programming.

## Installation

```bash
go get github.com/ahmedsaheed/esolang
```

or via homebrew

```bash
brew tap ahmedsaheed/esolang
brew install esolang
```

## Usage

Try it out now in a sandbox environment on your browser via the [esolang playground](https://esolang.onrender.com/).

Eso Expressions can evaluated on your terminal using the built-in repl
by running the command `esolang -repl` or via a file on your preferred text editor with the `.eso` extension and running the command `esolang -file <filename>`.

## Overview

| Page                          | Description                                |
| ----------------------------- | ------------------------------------------ |
| [Language Overview](./eso.md) | Overview of the esolang syntax & semantics |
| [Examples](./code_samples.md) | Code samples in esolang                    |
| Internals                     | Overview of the inner workings of esolang  |
| Built-in Functions            | List of built-in functions in esolang      |

## Inspiration

This scripting language draws its core concepts from the aclaimed works of Thorson Ball's [(Writing an Interpreter in Go)](https://interpreterbook.com/), Bob Nystrom [(Crafting Interpretes)](https://craftinginterpreters.com/), and the foundational [SICP](https://web.mit.edu/6.001/6.037/sicp.pdf).

> To make the language more approachable it borrows familiar concepts from popular languages like JavaScript, Golang, Ocaml as well as basic arithmetics principals

### List of features

| Features                   | Inspration       | Description                                                      |
| -------------------------- | ---------------- | ---------------------------------------------------------------- |
| Variable Bindings          | Ocaml & JS       | Variable bindings are done using `let` keyword                   |
| Conditionals               | C-Family, Golang | Conditionals are done using `if` keyword                         |
| Loops                      | C-Family, Golang | Loops are done using `while` keyword                             |
| Arrays                     | JavaScript       | Easy initialiasation of dictionary using `[]` or `array_new()`   |
| Hash                       | JavaScript       | A key value pair data structure `student = {"name": "John Doe"}` |
| Function Literals          | Golang           | Function literals are done using `fn` keyword                    |
| Higher Order Functions     | JavaScript       | Higher order functions are done using `fn` keyword               |
| Closures                   | Golang           | Closures are done using `fn` keyword                             |
| Error Messages             | Rust             | Rigid error messages                                             |
| Tons of built-in functions | Golang           | Battries includes functions to get you up an running ASAP        |

## Performance

The performance could be improved. At current its mediocre. Evaluation stage seems to be hindering speed as its a recursive decent which traverses each AST node and interprets what its sees. Bytecode compilation solution which runs on a VM would improve the performance significantly.

Current benchmark of `fib(20)` using [hyperfine](https://github.com/sharkdp/hyperfine) is as follows

```bash
hyperfine 'esolang.go ./benchmarkings/bench.eso' 'python3 ./benchmarkings/bench.py'
```

| Command                             |  Mean [ms] | Min [ms] | Max [ms] |    Relative |
| :---------------------------------- | ---------: | -------: | -------: | ----------: |
| `esolang ./benchmarkings/bench.eso` | 14.6 ± 1.5 |     14.0 |     28.4 |        1.00 |
| `python3 ./benchmarkings/bench.py`  | 22.3 ± 3.8 |     21.1 |     50.1 | 1.52 ± 0.30 |
