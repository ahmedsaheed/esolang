## Eso Samples

### Hello, World

```js
print("Hello, World!");
println("Hello" + "," + " " + "World");
```

### Variables

```js
let x = 10;
let y = 20;
let sum = x + y;
let correct_sum = sum == 20;

let version = "0.1.0";
let message = "Welcome to EsoLang";
```

### Conditionals

```js
let x = 10;
let y = 20;

if (x < y) {
  println("x is less than y");
} else {
  println("x is greater than y");
}
```

### Loops

```js
let x = 0;

while (x < 10) {
  println(x);
  let x = x + 1;
}
```

### Functions

```js
let min = fn(a, b) {
  if (a < b) {
    return a;
  } else {
    return b;
  }
};

let result = min(10, 20);
```

### Arrays

```js
let arr = [1, 2, 3, 4, 5];
// or
let arr = array_new(1, 2, 3, 4, 5);
arr[0] = 1;
arr[1] = 2;
```

### Fibonacci

```js
let fibonacci = fn(x) {
   if (x < 2) {x} else {
    fibonacci(x - 1) + fibonacci(x - 2);
  }
}

let result = fibonacci(10);
```

### Array Sum

```js

let arr = [1, 2, 3, 4, 5];
// using built-in array_reduce
let sum = array_reduce(arr, 0, fn(accum, el) { accum + el });

```

Many more examples can be found in the [esolang playground](https://esolang.onrender.com/). Go have fun!
