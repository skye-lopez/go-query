# go-query
Better query tools for sql/psql

not stable yet

`go get github.com/skye-lopez/go-query`

### Roadmap
- More testing
- Errors
- More featuresss

## Usage
```golang
package main

import (
    "database/sql"
    _ "github.com/lib/pq"
    "github.com/skye-lopez/go-query"
)

func main() {
    conn := YourAwesomeDBConnLogic() // get your db conn
    gq := NewGoQuery(conn)

    // Add all your query files to the map for ease of use.
    gq.AddQueriesToMap("queriesDir")

    // You can add specific args to a query if you want it to be used every call.
    gq.AddDefaultArgsToQuery("queriesDir/queryName", args) // args is just a []any

    // Preform a query
    rows, err := gq.Query("queries/queryName", ...args) // args are fully dynamic.

    // do whatever you want with rows...
    for _, row := range rows {
        for _, col := range row.([]interface{}) {

        }
    }

    // You can also do it with a queryString if you dont want to write a query file for it.
    rows, err := gq.QueryString("SELECT awesome FROM table;")
}
```
