// internal/api/handler.go
package api

import (
	"fmt"
	"net/http"
	"compress-pdf/internal/compressor" // Importa o nosso pacote de compressão.
	"github.com/gin-gonic/gin"
)

// CompressPDFHandler lida com a requisição de compressão de PDF.
func CompressPDFHandler(c *gin.Context) {
	// 1. Obter o arquivo da requisição (multipart/form-data)
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nenhum arquivo enviado. Use o campo 'file'."})
		return
	}

	// 2. Obter o parâmetro de qualidade (opcional)
	// Ex: /compress?quality=ebook
	qualityQuery := c.DefaultQuery("quality", "ebook")
	var qualityLevel compressor.Quality
	switch qualityQuery {
	case "screen":
		qualityLevel = compressor.ScreenQuality
	case "ebook":
		qualityLevel = compressor.EbookQuality
	case "printer":
		qualityLevel = compressor.PrinterQuality
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Qualidade inválida. Use 'screen', 'ebook' ou 'printer'."})
		return
	}

	// 3. Abrir o arquivo enviado
	srcFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível abrir o arquivo enviado."})
		return
	}
	defer srcFile.Close()

	// 4. Chamar a nossa lógica de negócio (o compressor)
	// O gin.Context.Writer é um io.Writer, então podemos passar diretamente!
	err = compressor.Compress(srcFile, c.Writer, qualityLevel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Falha na compressão: %v", err)})
		return
	}

	// 5. Configurar os headers da resposta para indicar que é um PDF para download
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=comprimido-%s", file.Filename))
	c.Header("Content-Type", "application/pdf")
}