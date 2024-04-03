# Esolang

Esolang is a beginner friendly dynamically typed scripting language quick prototyping and learning functional programming.

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

| Page               | Description                                |
| ------------------ | ------------------------------------------ |
| Syntax & Semantics | Overview of the esolang syntax & semantics |
| Examples           | Code samples in esolang                    |
| Internals          | Overview of the inner workings of esolang  |
| Built-in Functions | List of built-in functions in esolang      |

## Inspiration

This scripting language is highly inspired by the book [Writing an Interpreter in Go](https://interpreterbook.com/) by Thorsten Ball and [Crafting Interpreters](https://craftinginterpreters.com/) by Bob Nystrom. Also some takeaway from the [Structural and Interpretation of Computer Programs - SICP](https://web.mit.edu/6.001/6.037/sicp.pdf) by Harold Abelson and Gerald Jay Sussman. And [loads of blog posts](https://journal.stuffwithstuff.com/category/parsing/) from Bob Nystrom.

Also some concepts are gather from primary languages like golang, python, javascript, and basic arithmetic.

List of features & Sources

| Features               | Inspration       | Description                                        |
| ---------------------- | ---------------- | -------------------------------------------------- |
| Variable Bindings      | Ocaml & JS       | Variable bindings are done using `let` keyword     |
| Conditionals           | C-Family, Golang | Conditionals are done using `if` keyword           |
| Loops                  | C-Family, Golang | Loops are done using `for` keyword                 |
| Function Literals      | Golang           | Function literals are done using `fn` keyword      |
| Higher Order Functions | JavaScript       | Higher order functions are done using `fn` keyword |
| Closures               | Golang           | Closures are done using `fn` keyword               |
| Lexer Error Messages   | Rust             | Rigid error messages just like rust                |

## Performance

The performance could be improved. At current its mediocre. Evaluation stage seems to be hindering speed as its a recursive decent which traverses each AST node and interprets what its sees. I'm looking into a bytecode compilable solution which runs on a VM. This would improve the performance significantly.

| Command                             |     Mean [ms] | Min [ms] | Max [ms] |      Relative |
| :---------------------------------- | ------------: | -------: | -------: | ------------: |
| `esolang ./benchmarkings/bench.eso` | 607.4 ± 326.3 |    487.3 |   1532.3 | 24.58 ± 17.50 |
| `python3 ./benchmarkings/bench.py`  |   24.7 ± 11.5 |     21.0 |    107.2 |          1.00 |
