package query

import (
    "database/sql"
    "os"
    "strings"
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
