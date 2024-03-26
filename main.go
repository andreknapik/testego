package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	_ "github.com/lib/pq"
)

type Usuario struct {
	ID   int    `json:"id"`
	Nome string `json:"nome"`
}

func main() {
	// String de conexão com o banco de dados
	connStr := "postgresql://postgres:unRnsdTHJVefMLSrKFWiQjqblaUnRqOH@viaduct.proxy.rlwy.net:14523/railway"

	// Abre a conexão com o banco de dados
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Verifica se a conexão com o banco de dados está disponível
	err = db.Ping()
	if err != nil {
		log.Fatalf("Erro ao fazer ping no banco de dados: %v", err)
	}

	// Manipulador de rota para /usuarios
	http.HandleFunc("/usuarios", func(w http.ResponseWriter, r *http.Request) {
		// Executa a consulta no banco de dados
		rows, err := db.QueryContext(context.Background(), "SELECT id, nome FROM usuarios")
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

		// Verifica se houve algum erro durante a iteração
		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Serializa a lista de usuários como JSON e envia como resposta
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(usuarios)
	})

	// Inicia o servidor HTTP na porta 8080
	fmt.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
