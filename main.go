package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v4"
)

type Usuario struct {
	ID   int    `json:"id"`
	Nome string `json:"nome"`
}

func main() {
	// Configuração da conexão com o banco de dados
	connStr := "postgresql://postgres:unRnsdTHJVefMLSrKFWiQjqblaUnRqOH@viaduct.proxy.rlwy.net:14523/railway"

	// Criação do pool de conexões
	config, err := pgx.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("Erro ao fazer parsing da string de conexão: %v", err)
	}
	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     *config,
		MaxConnections: 5,
	})
	if err != nil {
		log.Fatalf("Erro ao criar pool de conexões: %v", err)
	}
	defer pool.Close()

	// Manipuladores de rota
	http.HandleFunc("/usuarios", func(w http.ResponseWriter, r *http.Request) {
		// Executa a consulta no banco de dados
		rows, err := pool.Query(context.Background(), "SELECT id, nome FROM usuarios")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Itera sobre os resultados e constrói uma lista de usuários
		var usuarios []Usuario
		for rows.Next() {
			var usuario Usuario
			err := rows.Scan(&usuario.ID, &usuario.Nome)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			usuarios = append(usuarios, usuario)
		}

		// Serializa a lista de usuários como JSON e envia como resposta
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(usuarios)
	})

	// Inicia o servidor HTTP na porta 8080
	fmt.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
