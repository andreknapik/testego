package main

import (
    "context"
    "fmt"
    "github.com/jackc/pgx/v4"
)

func main() {
    // String de conexão
    connStr := "postgresql://postgres:unRnsdTHJVefMLSrKFWiQjqblaUnRqOH@viaduct.proxy.rlwy.net:14523/railway"

    // Abrindo conexão com o banco de dados (use pgx.Connect)
    conn, err := pgx.Connect(context.Background(), connStr)
    if err != nil {
        panic(err)
    }
    defer conn.Close(context.Background())

    // Executando uma consulta
    rows, err := conn.Query(context.Background(), "SELECT id, nome FROM usuarios")
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    // Varrendo os resultados da consulta
    for rows.Next() {
        var id int
        var nome string
        err := rows.Scan(&id, &nome)
        if err != nil {
            panic(err)
        }
        fmt.Println("ID:", id, "Nome:", nome)
    }
}
