# structquery
Query data structures using string path in Go

**Note: This project is incomplete, not stable, and most features don't work. I am not working on it at the moment. I am shelving it for work later.**

Using an Xpath-like query language to get, add, remove, or replace values in a Go data structure. Great for dynamic data structures.

## Example
```go
package main

import (
  "fmt"
  "github.com/bryanaustin/structquery"
)

func main() {
  x := map[string]int{"one":1, "two":2, "ten":10}
  y, _ := structquery.Get(x, "two")
  fmt.Println(y[0])
  // Output: 2
}
```
