// internal/api/router.go
package api

import (
	"github.com/gin-gonic/gin"
)

// NewRouter configura as rotas da nossa API.
func NewRouter() *gin.Engine {
	router := gin.Default()

	// Define um limite para o tamanho do upload para evitar abusos.
	router.MaxMultipartMemory = 8 << 20 // 8 MB

	// Nossa rota principal para compressão.
	router.POST("/compress", CompressPDFHandler)

	return router
}