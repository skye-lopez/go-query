package query

import (
    "database/sql"
    "os"
    "strings"
    "reflect"
)

type GoQuery struct {
    Conn *sql.DB
    QueryMap map[string]Query
}

type Query struct {
    Query string
    Args []any
}

func NewGoQuery(conn *sql.DB) (GoQuery) {
    return GoQuery{
        Conn: conn,
        QueryMap: make(map[string]Query),
    }
}

func (gq *GoQuery) AddDefaultArgsToQuery(queryName string, args []any) (error) {
    if query, ok := gq.QueryMap[queryName]; ok {
        query.Args = args
        gq.QueryMap[queryName] = query
        return nil
    }
    return nil // TODO: Doesnt exist error
}

func (gq *GoQuery) QueryString(q string, args ...any) ([]any, error) {
    var result []any
    rows, err := gq.Conn.Query(q, args...)
    if err != nil {
        return result, err
    }
    defer rows.Close()

    for rows.Next() {
        types, err := rows.ColumnTypes()
        if err != nil {
            return result, err
        }

        values := make([]any, len(types))
        refs := make([]any, len(types))

        for i, t := range types {
            values[i] = reflect.New(t.ScanType())
            refs[i] = &values[i]
        }

        err = rows.Scan(refs...)
        if err != nil {
            return result, err
        }

        result = append(result, values)
    }

    return result, nil
}

func (gq *GoQuery) Query(queryName string, args ...any) ([]any, error) {
    var q Query
    if query, ok := gq.QueryMap[queryName]; ok {
        if len(query.Args) > 0 {
            args = query.Args
        }
        q = query
    } else {
        // TODO: ERROR
        return make([]any, 0), nil
    }

    result, err := gq.QueryString(q.Query, args...)
    return result, err
}

func (gq *GoQuery) AddQueriesToMap(dirPath string) (error) {
    dirs := []string{ dirPath }
    for len(dirs) > 0 {
        dir := pop(&dirs)

        queryFiles, err := os.ReadDir(dir)
        if err != nil {
            return err
        }

        for _, file := range queryFiles {
            if (file.IsDir()) {
                dirs = append(dirs, (dir + "/" + file.Name()))
                continue
            }

            data, err := os.ReadFile(dir + "/" + file.Name())
            if err != nil {
                return err
            }

            fileExtension := strings.Split(file.Name(), ".")
            if fileExtension[len(fileExtension)-1] != "sql" {
                continue
            }

            q := Query{
                Query: string(data),
                Args: make([]any, 0),
            }
            gq.QueryMap[dir + "/" + fileExtension[0]] = q
        }
    }
    return nil
}

// UTIL
func pop(l *[]string) string {
    f := len(*l)
    rv := (*l)[f-1]
    *l = (*l)[:f-1]
    return rv
}
