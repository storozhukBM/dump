# Dump
One function library to simplify local print debug.

## Motivation
How many times during local code debug you've written huge `fmt.Printf`s similar to that one?

```go
fmt.Printf("Tx commit. body: %v; hash: %v; isValid: %v; offset: %v;", body, hashCode, codeIsValid, offset)
```

When you in the middle of complex code debug this printfs can be annoying and sometimes painful to maintain if you introduce new variables or rename things as you go.

But now our struggles are over and you can just:

```go
dump.Dump("Tx commit. ", body, hashCode, codeIsValid, offset)
```

And you will see something like that in stdout:

```
[DEBUG] /dump/example_test.go:23: Tx commit. body: `txBody`; hashCode: `94876`; codeIsValid: `false`; offset: `{TxName:Final idx:34 deadline:160}`
```

**Notice** that we automatically find the file and line number of `Dump` call, and we also resolve all the variable names you have in your code.

## How to use it
You can import it as a library via Go modules or copy-paste `dump.go` file it to your project and add it to `.gitignore` if you don't want to pollute your dependencies.

# Important
This library will work only for debug scenarios in your local development environment. Please don't try to use it in production.

It also won't handle 100% of different cases for the sake of simplicity.
It won't handle multiline `Dump` statements, so please create more consecutive statements if necessary.