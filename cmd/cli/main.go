// cmd/cli/main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	// Importamos o nosso pacote de lógica central!
	"compress-pdf/internal/compressor" 
)

func main() {
	// 1. Definir e parsear os argumentos (flags) da linha de comando
	inputPath := flag.String("i", "", "Caminho do arquivo de entrada (obrigatório)")
	outputPath := flag.String("o", "", "Caminho do arquivo de saída (obrigatório)")
	quality := flag.String("q", "ebook", "Nível de qualidade: screen, ebook, printer")

	flag.Parse()

	// 2. Validar os argumentos
	if *inputPath == "" || *outputPath == "" {
		fmt.Println("Erro: Os caminhos de entrada (-i) e saída (-o) são obrigatórios.")
		flag.Usage() // Mostra como usar os comandos
		os.Exit(1)
	}

	// 3. Converter a string de qualidade para o nosso tipo
	var qualityLevel compressor.Quality
	switch *quality {
	case "screen":
		qualityLevel = compressor.ScreenQuality
	case "ebook":
		qualityLevel = compressor.EbookQuality
	case "printer":
		qualityLevel = compressor.PrinterQuality
	default:
		log.Fatalf("Qualidade inválida: %s. Use 'screen', 'ebook' ou 'printer'.", *quality)
	}

	// 4. Abrir o arquivo de entrada
	inputFile, err := os.Open(*inputPath)
	if err != nil {
		log.Fatalf("Erro ao abrir arquivo de entrada %s: %v", *inputPath, err)
	}
	defer inputFile.Close()

	// 5. Criar o arquivo de saída
	outputFile, err := os.Create(*outputPath)
	if err != nil {
		log.Fatalf("Erro ao criar arquivo de saída %s: %v", *outputPath, err)
	}
	defer outputFile.Close()

	// 6. Chamar a lógica de negócio!
	// Note que é a *mesma função* que a API usa.
	fmt.Printf("Comprimindo %s para %s com qualidade %s...\n", *inputPath, *outputPath, *quality)
	err = compressor.Compress(inputFile, outputFile, qualityLevel)
	if err != nil {
		log.Fatalf("Erro durante a compressão: %v", err)
	}

	fmt.Println("PDF comprimido com sucesso!")
}