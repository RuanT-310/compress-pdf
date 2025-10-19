// internal/compressor/ghostscript.go
package compressor

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// Define diferentes níveis de qualidade para compressão.
type Quality string

const (
	ScreenQuality Quality = "/screen" // Menor tamanho, menor qualidade.
	EbookQuality  Quality = "/ebook"  // Bom equilíbrio.
	PrinterQuality Quality = "/printer" // Alta qualidade.
)

// Compress processa um PDF de um leitor (reader) e escreve o resultado comprimido
// em um escritor (writer), usando a qualidade especificada.
func Compress(input io.Reader, output io.Writer, quality Quality) error {
	// Ghostscript precisa de arquivos no disco. Criamos arquivos temporários.
	inputFile, err := os.CreateTemp("", "input-*.pdf")
	if err != nil {
		return fmt.Errorf("falha ao criar arquivo de entrada temporário: %w", err)
	}
	defer os.Remove(inputFile.Name()) // Garante a limpeza do arquivo

	outputFile, err := os.CreateTemp("", "output-*.pdf")
	if err != nil {
		return fmt.Errorf("falha ao criar arquivo de saída temporário: %w", err)
	}
	defer os.Remove(outputFile.Name())

	// Copia o conteúdo do input (que vem da requisição HTTP) para o arquivo temporário.
	if _, err := io.Copy(inputFile, input); err != nil {
		return fmt.Errorf("falha ao escrever no arquivo de entrada temporário: %w", err)
	}
	inputFile.Close() // Fecha para garantir que o Ghostscript possa ler

	// Monta e executa o comando Ghostscript.
	cmd := exec.Command("gs",
		"-sDEVICE=pdfwrite",
		"-dCompatibilityLevel=1.4",
		fmt.Sprintf("-dPDFSETTINGS=%s", quality),
		"-dNOPAUSE",
		"-dQUIET",
		"-dBATCH",
		fmt.Sprintf("-sOutputFile=%s", outputFile.Name()),
		inputFile.Name(),
	)

	if cmdOutput, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("erro ao executar o Ghostscript: %s: %w", string(cmdOutput), err)
	}

	// Lê o arquivo de saída comprimido.
	compressedFile, err := os.Open(outputFile.Name())
	if err != nil {
		return fmt.Errorf("falha ao abrir o arquivo de saída comprimido: %w", err)
	}
	defer compressedFile.Close()

	// Copia o resultado para o writer de saída (que será a resposta HTTP).
	if _, err := io.Copy(output, compressedFile); err != nil {
		return fmt.Errorf("falha ao copiar o resultado para a saída: %w", err)
	}

	return nil
}