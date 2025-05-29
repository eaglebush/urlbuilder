# urlbuilder

`urlbuilder` is a Go package for building URLs in a flexible, composable, and readable way. It supports schemes, hosts, ports, paths, IDs, query parameters (with deduplication strategies), user credentials, and fragments.

## Features

- Compose URLs using functional options (`UrlPart`)
- Support for basic auth credentials
- Support for query deduplication modes:
  - `QModeArray`: keep all query parameters (duplicates allowed)
  - `QModeLast`: keep only the last value for duplicate names
  - `QModeError`: fail on duplicate query names
- Cloneable URL builders
- Minimal allocations (optimized with `strings.Builder`)
- Support for independent query string generation

## Installation

```bash
go get github.com/eaglebush/urlbuilder
```

## Usage

### Basic URL

```go
import "github.com/eaglebush/urlbuilder"

url := urlbuilder.New(
    urlbuilder.Sch("https"),
    urlbuilder.Host("example.com"),
    urlbuilder.Path("api"),
    urlbuilder.Path("v1"),
    urlbuilder.ID(123),
    urlbuilder.Query("q", "go"),
).Build()

fmt.Println(url)
// Output: https://example.com/api/v1/123?q=go
```

### Simple Helper

```go
url := urlbuilder.NewUrlWithPath("example.com", "status").Build()
fmt.Println(url)
// Output: https://example.com/status/
```

### With Authentication

```go
url := urlbuilder.New(
    urlbuilder.Sch("https"),
    urlbuilder.Host("secure.example.com"),
    urlbuilder.UsrPwd("admin", "secret"),
    urlbuilder.Path("dashboard"),
).Build()

fmt.Println(url)
// Output: https://admin:secret@secure.example.com/dashboard/
```

### Query Deduplication Modes

```go
url := urlbuilder.New(
    urlbuilder.Host("example.com"),
    urlbuilder.Path("search"),
    urlbuilder.Mode(urlbuilder.QModeLast),
    urlbuilder.Query("q", "first"),
    urlbuilder.Query("q", "last"),
).Build()

fmt.Println(url)
// Output: https://example.com/search/?q=last
```

### Cloning and Modifying

```go
base := urlbuilder.NewUrl("example.com", "items", urlbuilder.Query("sort", "asc"))
filtered := base.Clone(urlbuilder.Query("category", "books"))

fmt.Println(filtered.Build())
// Output: https://example.com/items/?sort=asc&category=books
```
### UrlBuilder with QueryString Sample
This allows the UrlBuilder to accept a query from a ready built encoded query string.
```go
url := urlbuilder.NewUrlWithPath(
	"example.com",
	"protected",
	urlbuilder.Query(
		"q",
		MyEncyptedEncoder(
			urlbuilder.NewQueryString(urlbuilder.QModeArray,
			urlbuilder.Nv("secret", "mysecretvalue"),
			urlbuilder.Nv("for", "sherlock"),
			urlbuilder.Nv("execute", "order-66"),
			urlbuilder.Nv("clean", true),
		).String(),
	)),
).Build()
fmt.Println(url)
// Output: https://example.com/protected/?q=<garbled>
```

## API Reference

### Constructor

- `New(...UrlPart) *UrlBuilder`
- `NewUrl(host, path string, ...UrlPart) *UrlBuilder` (Convenience function)
- `NewUrlWithPath(host, path string, id any, ...UrlPart) *UrlBuilder` (Convenience function)
- `NewUrlWithID(host, path string, id any, ...UrlPart) *UrlBuilder` (Convenience function)
- `Clone(ub *UrlBuilder, ...UrlPart) *UrlBuilder`
- `NewQueryString(mode QueryMode, qp ...QueryPart) *QueryString` (Independent query string creation)

### UrlPart Functions

- `Sch(string)` - set scheme
- `Host(string)` - set host
- `Port(uint)` - set port
- `Path(string)` - add path segment
- `ID(any)` - set ID segment
- `Usr(string)` - set username
- `Pwd(string)` - set password
- `UsrPwd(string, string)` - set both username and password
- `Query(string, any)` - add query parameter
- `Frag(string)` - set fragment
- `Mode(QueryMode)` - set query deduplication mode


### UrlBuilder Methods

- `(*UrlBuilder) Build() string` - build the URL
- `(*UrlBuilder) String() string` - alias for `Build()`
- `(*UrlBuilder) Clone(...UrlPart) *UrlBuilder`
- `(*UrlBuilder) Err() error`

### QueryPart Functions

- `Nv(string, any)` - add query parameter

### QueryString Methods

- `(*UrlBuilder) Build() string` - build the query string
- `(*UrlBuilder) String() string` - alias for `Build()`


## License

MIT
