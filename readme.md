# urlbuilder

**urlbuilder** is a composable and expressive URL builder library for Go. It provides a fluent API to construct URLs with minimal boilerplate and full flexibility.

---

## Features

- Composable URL parts using function-based configuration
- Support for all common URL components: scheme, host, user/password, path, query, fragment
- Query handling modes: array, last value wins, error on duplicates
- Built-in cloning and shortcut constructors
- Clean, idiomatic Go interface

---

## Installation

```bash
go get github.com/eaglebush/urlbuilder
```

---

## Basic Usage

```go
package main

import (
    "fmt"
    "github.com/eaglebush/urlbuilder"
)

func main() {
    url := urlbuilder.New(
        urlbuilder.Sch("https"),
        urlbuilder.Host("example.com"),
        urlbuilder.Path("api"),
        urlbuilder.Path("users"),
        urlbuilder.ID(42),
        urlbuilder.Query("active", true),
        urlbuilder.Fragment("details"),
    ).Build()

    fmt.Println(url)
    // Output: https://example.com/api/users/42?active=true#details
}
```

---

## Query Modes

### ▶ QModeArray

Allows repeated parameters:

```go
url := urlbuilder.New(
    urlbuilder.Host("example.com"),
    urlbuilder.Query("x", 1),
    urlbuilder.Query("x", 2),
    urlbuilder.Mode(urlbuilder.QModeArray),
).Build()

fmt.Println(url)
// Output: https://example.com/?x=1&x=2
```

### ▶ QModeLast

Keeps only the last value for each query key:

```go
url := urlbuilder.New(
    urlbuilder.Host("example.com"),
    urlbuilder.Query("x", 1),
    urlbuilder.Query("x", 2),
    urlbuilder.Mode(urlbuilder.QModeLast),
).Build()

fmt.Println(url)
// Output: https://example.com/?x=2
```

### ▶ QModeError

Throws an error on duplicate query keys:

```go
ub := urlbuilder.New(
    urlbuilder.Host("example.com"),
    urlbuilder.Query("x", 1),
    urlbuilder.Query("x", 2),
    urlbuilder.Mode(urlbuilder.QModeError),
)

url := ub.Build()

if err := ub.Err(); err != nil {
    fmt.Println("Error:", err)
}
// Output: Error: duplicate query name found
```

---

## Cloning URLs

Clone an existing builder and add more parts:

```go
base := urlbuilder.New(
    urlbuilder.Host("example.com"),
    urlbuilder.Path("api"),
)

extended := base.Clone(
    urlbuilder.Path("products"),
    urlbuilder.ID(88),
)

fmt.Println(extended.Build())
// Output: https://example.com/api/products/88
```

---

## Convenience Functions

### ▶ NewSimpleUrl

```go
url := urlbuilder.NewSimpleUrl("example.com", "dashboard").Build()
fmt.Println(url)
// Output: https://example.com/dashboard
```

### ▶ NewSimpleUrlWithID

```go
url := urlbuilder.NewSimpleUrlWithID("example.com", "user", 789).Build()
fmt.Println(url)
// Output: https://example.com/user/789
```

---

## String Conversion

```go
ub := urlbuilder.New(
    urlbuilder.Host("example.com"),
    urlbuilder.Path("login"),
)

fmt.Println(ub.String())
// Output: https://example.com/login
```

---

## License

MIT License

---

## Author

[eaglebush]
