+++
title = "Data Types"
date = "2024-05-22T00:01:24+01:00"
author = "Ahmed Saheed"
showFullContent = false
readingTime = false
hideComments = false
+++

## Datatypes

| **Datatype** | **Description**                                                   | **Examples**                     |
| ------------ | ----------------------------------------------------------------- | -------------------------------- |
| Integer      | 64Bit int number                                                  | 1, 100, 1000                     |
| String       | Text, multiple and single characters                              | `"John, Doe"`, `"A"`, `"!!!"`    |
| Bool         | Booleans                                                          | `true`, `false`                  |
| Array        | Loosely typed list which contains any of the above type           | [1,2,3], [a,b,1,2]               |
| Hash         | Loosely typed Key Value Pair which contains any of the above type | {"name": "John Doe", "age": 100} |

## Printing

The `println` or `print` keyword prints all arguements to the standard outstream. Our current knowledge of esolang can be applied to create an `Hello, World` example.

```js
println("Hello, World");
```

## Evaluation Expressions

Invoking the interpreter can easily be done via two ways

- The builtin repl
- A file with the `.eso` extension.

#### The Repl

The repl accepts all valid eso expressions. Running the aforemetioned `Hello, World` can be done in the repl like so:

```bash
$ esolang -repl

Hello username! Welcome to esolang's repl
Feel free to type in commands
Type '.help' for assistance

>>print("Hello, World")
INFO:  Hello, World
>>
```

#### Eso File

These are recognisable files to esolang. They usually have the `.eso` file extension. Running `Hello, World` in an esofile would look something like this:

```bash
touch ./hello.eso
echo "println("Hello, World")" > ./hello.eso

esolang ./hello.eso
"Hello, World"
```

## Comments

Single line comments are supported in esolang

```js
// this is a comment - parser skips this line
```

## Arithmetic Operations

Arithmentic expressions are evaluated with the operation precedence in mind.

| Operator | Description                 | Example                     |
| -------- | --------------------------- | --------------------------- |
| `+`      | Sums left & right           | let sum = 1 + 1 -> 2        |
| `-`      | Differences of left & right | let diff = 10 - 5 -> 5      |
| `*`      | Product of left & right     | let product = 2 \* 2 -> 4   |
| `/`      | Quotient of left & right    | let quotient = 10 / 2 -> 5  |
| `%`      | Remainder of left & right   | let remainder = 10 % 3 -> 1 |

## Comparison Operations

Comparison operations are used to compare two values. They return a boolean value.

| Operator | Description                | Example        |
| -------- | -------------------------- | -------------- |
| `==`     | Equality of left & right   | 1 == 1 -> true |
| `!=`     | Inequality of left & right | 1 != 2 -> true |
| `>`      | Greater than left & right  | 2 > 1 -> true  |
| `<`      | Less than left & right     | 1 < 2 -> true  |

## Logical Operations

Logical operations are used to combine multiple boolean values.

| Operator | Description                 | Example                    |
| -------- | --------------------------- | -------------------------- | ---- | --- | ------------- |
| `&&`     | Logical AND of left & right | true && true -> true       |
| `-       | `                           | Logical OR of left & right | true |     | false -> true |
| `!`      | Logical NOT of right        | !true -> false             |

## Variables

Variables are used to store values. They are declared using the `let` keyword.

```js
let name = "John Doe";
let age = 100;
let isAdult = true;
let arr = [1, 2, 3];
let hash = { name: "John Doe", age: 100 };
```

## Control Flow

Control flow statements are used to control the flow of the program. They include:

- `if` statements
- `else` statements
- `while` loops

### If Statements

If statements are used to execute a block of code if a condition is true.

```js
let age = 100;

if (age > 18) {
  println("You are an adult");
}
```

### Else Statements

Else statements are used to execute a block of code if the condition in the if statement is false.

```js
let age = 10;

if (age > 18) {
  println("You are an adult");
} else {
  println("You are not an adult");
}
```

### While Loops

While loops are used to execute a block of code as long as a condition is true.

```js
let i = 0;

while (i < 10) {
  println(i);
  let i = i + 1;
}
```

## Functions

Functions are used to group code into reusable blocks. They are declared using the `fn` keyword.

```js
let add = fn(a, b) {
  return a + b;
}

let result = add(1, 2);
println(result); // 3
```

### Hashes

Hashes are used to store key-value pairs. They are declared using curly braces `{}`.

```js
let person = { name: "John Doe", age: 100, siblings: ["Jane Doe", "Jack Doe"] };
```

Accessing hashes and array elements

Elements in a hash can be accessed using the key and elements in an array can be accessed using the index.

```js
let person = { name: "John Doe", age: 100, siblings: ["Jane Doe", "Jack Doe"] };

println(person["name"]); // John Doe
println(person["age"]); // 100
println(person["siblings"][0]); // Jane Doe

let arr = [1, 2, 3];
println(arr[0]); // 1
```
