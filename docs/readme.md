# Documentation - Esolang

Esolang is a minimal interpreted scripting language that syntactically looks like F# programming language. It is a simple language that is designed to be easy to learn and use. It offers familiar programming language features such as variable bindings, conditionals, and loops, as well function literals, higher order functions and closures.

It's implementation can be expressions passed into a repl or a file with the extension `.esolang` and run with the command `esolang <file>.esolang`.

Its a portable interpreted language with a Pratt parser and a tree-walking interpreter. It is written in Golang and can be easily extended to be compiled and ran in a virtual machine.

## Overview

| Page               | Description                                                            |
| ------------------ | ---------------------------------------------------------------------- |
| Syntax & Semantics | Overview of the esolang syntax & semantics                             |
| Examples           | Code samples in esolang                                                |
| Internals          | Overview of the inner workings and the code of interpreter / evaluator |

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

FizzBuzz in Esolang

```js
let fizzBuzz = fn(num, res){
  let i = 1;
  while(i < num) {
    let divided = false;
    if(i % 15 == 0){array_append(res, "FizzBuzz"); let divided = true; };
    if(i % 3 == 0){array_append(res, "Fizz"); let divided = true;};
    if(i % 5 == 0){array_append(res, "Bizz"); let divided = true;};
    if(dev == false){array_append(res, i);} let i = i + 1;}
  return res;
}

let res = array_new();
let result = fizzBuzz(100, res);
println(result); // println(fizzBuzz(100, res))
```
