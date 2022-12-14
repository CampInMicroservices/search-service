package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RecommendationServicePingResponse struct {
	Status string `json:"status"`
}

// @BasePath /search-service/v1

// Health godoc
// @Summary Liveness
// @Schemes
// @Description Liveness
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {array} api.RecommendationServicePingResponse
// @Router /health/live [get]
func (server *Server) Live(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "UP"})
}

// @BasePath /search-service/v1

// Health godoc
// @Summary Readiness
// @Schemes
// @Description Readiness
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {string} api.RecommendationServicePingResponse
// @Router /health/ready [get]
func (server *Server) Ready(ctx *gin.Context) {

	// Check connection with database.
	err := server.store.PingDB()
	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"status": "DOWN"})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "UP"})
}

func PingRecommendationService(url string) (*RecommendationServicePingResponse, error) {

	req, _ := http.NewRequest("GET", url, nil)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Cannot unmarshal Response")
		return nil, errors.New("Recommendation service unavailable!")
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var r RecommendationServicePingResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, errors.New("Recommendation service unavailable!")
	}

	if res.StatusCode != 200 || r.Status != "UP" {
		return &r, errors.New("Recommendation service unavailable!")
	}

	return &r, nil
}
