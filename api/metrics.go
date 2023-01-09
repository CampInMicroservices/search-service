package api

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// @BasePath /search-service/v1

// Metrics godoc
// @Summary Metrics
// @Schemes
// @Description Metrics
// @Tags Metrics
// @Accept json
// @Produce json
// @Success 200 {string} helloworld
// @Router /metrics [get]
func (server *Server) Metrics(ctx *gin.Context) {
	handler := promhttp.Handler()
	handler.ServeHTTP(ctx.Writer, ctx.Request)
}
