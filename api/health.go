package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RecommendationServicePingResponse struct {
	Status string `json:"status"`
}

func (server *Server) Live(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "UP"})
}

func (server *Server) Ready(ctx *gin.Context) {

	dbConnectionStatus := "UP"
	recommendationServiceConnectionStatus := "UP"

	// Check connection with database.
	dbErr := server.store.PingDB()
	if dbErr != nil {
		dbConnectionStatus = "DOWN"
	}

	url := fmt.Sprintf("http://%s/health/live", server.config.RecommendationServiceAddress)
	_, err := PingRecommendationService(url)

	if err != nil {
		recommendationServiceConnectionStatus = "DOWN"
	}

	status := gin.H{"status": gin.H{
		"db_connection":          dbConnectionStatus,
		"recommendation_service": recommendationServiceConnectionStatus,
	}}

	if dbConnectionStatus == "DOWN" || recommendationServiceConnectionStatus == "DOWN" {
		ctx.JSON(http.StatusServiceUnavailable, status)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, status)
}

func PingRecommendationService(url string) (*RecommendationServicePingResponse, error) {

	req, _ := http.NewRequest("GET", url, nil)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Cannot unmarshal Response")
		return nil, errors.New("Recommendation service unavailable!")
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var r RecommendationServicePingResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		fmt.Println("Cannot unmarshal Response")
		return nil, errors.New("Recommendation service unavailable!")
	}

	if res.StatusCode != 200 || r.Status != "UP" {
		return &r, errors.New("Recommendation service unavailable!")
	}

	return &r, nil
}
