// cmd/api/main.go
package main

import (
	"log"
	"compress-pdf/internal/api" // Importa nosso pacote de API
)

func main() {
	router := api.NewRouter() // Cria o roteador com as nossas rotas

	log.Println("Servidor iniciado na porta :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Não foi possível iniciar o servidor: %v", err)
	}
}