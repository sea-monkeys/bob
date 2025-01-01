# Rust Programming Tutorial for Beginners

Welcome to the Rust programming tutorial! Rust is a powerful, safe, and concurrent systems programming language. It's known for its performance and memory safety without a garbage collector. This tutorial will guide you through the basics of Rust, including syntax, control flow, functions, ownership, error handling, and more. We'll also build some simple projects to put your knowledge into practice.

## Setting Up Your Development Environment

Before we dive into coding, you'll need to set up your Rust development environment. Follow these steps:

1. **Install Rust**: Go to the [official Rust website](https://www.rust-lang.org/tools/install) and follow the instructions to install Rust using `rustup`. This will install the Rust compiler, `cargo` (Rust's package manager and build system), and other necessary tools.

2. **Verify Installation**: Open a terminal and run `rustc --version` and `cargo --version` to ensure Rust and Cargo are installed correctly.

3. **Set Up an IDE/Editor**: You can use any text editor, but an IDE with Rust support can be helpful. Popular choices include Visual Studio Code with the Rust Analyzer extension, IntelliJ IDEA with the Rust plugin, or CLion.

## 1. Basic Syntax and Concepts

### Variables and Data Types

In Rust, variables are immutable by default. You can make them mutable by using the `mut` keyword.

```rust
fn main() {
    // Immutable variable
    let x = 5;
    println!("The value of x is: {}", x);

    // Mutable variable
    let mut y = 10;
    y = 20;
    println!("The value of y is: {}", y);
}
```

Rust has several data types, including integers (`i32`, `u32`, etc.), floating-point numbers (`f32`, `f64`), booleans (`bool`), and characters (`char`).

### Control Flow

#### If/Else Statements

```rust
fn main() {
    let number = 3;

    if number < 5 {
        println!("Condition was true");
    } else {
        println!("Condition was false");
    }
}
```

#### Loops

Rust has three types of loops: `loop`, `while`, and `for`.

```rust
fn main() {
    // Infinite loop
    loop {
        println!("This loop runs forever!");
        break; // Use break to exit the loop
    }

    // While loop
    let mut counter = 0;
    while counter < 5 {
        println!("Counter: {}", counter);
        counter += 1;
    }

    // For loop
    for number in 1..5 {
        println!("Number: {}", number);
    }
}
```

### Functions

Functions in Rust are defined using the `fn` keyword.

```rust
fn main() {
    let result = add(5, 3);
    println!("The result is: {}", result);
}

fn add(a: i32, b: i32) -> i32 {
    a + b
}
```

### Ownership and Borrowing Basics

Rust's ownership system ensures memory safety without a garbage collector. Here's a simple example:

```rust
fn main() {
    let s1 = String::from("hello");
    let s2 = s1; // Ownership of s1 is moved to s2

    // println!("{}", s1); // This will cause a compile-time error
    println!("{}", s2);
}
```

To avoid moving ownership, you can borrow using references:

```rust
fn main() {
    let s1 = String::from("hello");
    let len = calculate_length(&s1); // Borrowing s1

    println!("The length of '{}' is {}.", s1, len);
}

fn calculate_length(s: &String) -> usize {
    s.len()
}
```

### Simple Error Handling

Rust uses the `Result` and `Option` types for error handling.

```rust
use std::fs::File;
use std::io::{self, Read};

fn main() {
    let filename = "hello.txt";

    let mut f = File::open(filename).expect("Failed to open file");

    let mut contents = String::new();
    f.read_to_string(&mut contents)
        .expect("Failed to read file");

    println!("File contents:\n{}", contents);
}
```

## 2. Sample Projects

### Hello World

```rust
fn main() {
    println!("Hello, world!");
}
```

### Basic Calculator

```rust
fn main() {
    let a = 10;
    let b = 5;

    println!("Addition: {}", add(a, b));
    println!("Subtraction: {}", subtract(a, b));
    println!("Multiplication: {}", multiply(a, b));
    println!("Division: {}", divide(a, b));
}

fn add(x: i32, y: i32) -> i32 {
    x + y
}

fn subtract(x: i32, y: i32) -> i32 {
    x - y
}

fn multiply(x: i32, y: i32) -> i32 {
    x * y
}

fn divide(x: i32, y: i32) -> i32 {
    if y == 0 {
        panic!("Division by zero!");
    }
    x / y
}
```

### Number Guessing Game

```rust
use std::io;
use std::cmp::Ordering;
use rand::Rng;

fn main() {
    println!("Guess the number!");

    let secret_number = rand::thread_rng().gen_range(1..=100);

    loop {
        println!("Please input your guess.");

        let mut guess = String::new();
        io::stdin()
            .read_line(&mut guess)
            .expect("Failed to read line");

        let guess: u32 = match guess.trim().parse() {
            Ok(num) => num,
            Err(_) => continue,
        };

        println!("You guessed: {}", guess);

        match guess.cmp(&secret_number) {
            Ordering::Less => println!("Too small!"),
            Ordering::Greater => println!("Too big!"),
            Ordering::Equal => {
                println!("You win!");
                break;
            }
        }
    }
}
```

### Simple File Operations

```rust
use std::fs::File;
use std::io::{self, Read};

fn main() {
    let filename = "hello.txt";

    match File::open(filename) {
        Ok(mut file) => {
            let mut contents = String::new();
            file.read_to_string(&mut contents)
                .expect("Failed to read file");
            println!("File contents:\n{}", contents);
        }
        Err(e) => println!("Failed to open file: {}", e),
    }
}
```

## 3. Exercises

1. **Modify the calculator to handle floating-point numbers.**
2. **Add a loop to the number guessing game that allows the user to guess multiple times until they win.**
3. **Create a program that reads a file and counts the number of words.**

## 4. Common Pitfalls and Solutions

- **Ownership Errors**: Always ensure that you're not using a variable after it has been moved. Use references (`&`) to borrow variables.
- **Mutable Borrowing**: You can only have one mutable borrow or multiple immutable borrows at a time.
- **Error Handling**: Always check for errors and handle them appropriately using `Result` and `Option`.

## 5. Tips for Further Learning

- **Rust Book**: The [official Rust book](https://doc.rust-lang.org/book/) is an excellent resource for learning Rust.
- **Rust by Example**: [Rust by Example](https://doc.rust-lang.org/rust-by-example/) provides practical examples of Rust concepts.
- **Rustlings**: [Rustlings](https://github.com/rust-lang/rustlings) is a collection of small exercises to help you practice Rust.

That's it for this beginner-friendly Rust tutorial! Rust can be challenging at first, but its safety features make it a powerful language for systems programming. Keep practicing and experimenting, and you'll become proficient in no time. Happy coding!