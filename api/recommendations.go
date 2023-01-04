package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	gclient "github.com/machinebox/graphql"
)

type LocationsResponse struct {
	Cities struct {
		Data []struct {
			City         string  `json:"city"`
			Latitude     float64 `json:"latitude"`
			Longitude    float64 `json:"longitude"`
			Population   int     `json:"population"`
			TempC        int     `json:"tempC"`
			Weather      string  `json:"weather"`
			WeatherShort string  `json:"weatherShort"`
		} `json:"data"`
		Error string `json:"error"`
	} `json:"cities"`
}

type weatherApiResponse struct {
	Success  bool        `json:"success"`
	Error    interface{} `json:"error"`
	Response struct {
		ID         string `json:"id"`
		DataSource string `json:"dataSource"`
		Loc        struct {
			Long float64 `json:"long"`
			Lat  float64 `json:"lat"`
		} `json:"loc"`
		Place struct {
			Name    string `json:"name"`
			City    string `json:"city"`
			State   string `json:"state"`
			Country string `json:"country"`
		} `json:"place"`
		Profile struct {
			Tz       string `json:"tz"`
			Tzname   string `json:"tzname"`
			Tzoffset int    `json:"tzoffset"`
			IsDST    bool   `json:"isDST"`
			ElevM    int    `json:"elevM"`
			ElevFT   int    `json:"elevFT"`
		} `json:"profile"`
		ObTimestamp int       `json:"obTimestamp"`
		ObDateTime  time.Time `json:"obDateTime"`
		Ob          struct {
			Type                string      `json:"type"`
			Timestamp           int         `json:"timestamp"`
			DateTimeISO         time.Time   `json:"dateTimeISO"`
			RecTimestamp        int         `json:"recTimestamp"`
			RecDateTimeISO      time.Time   `json:"recDateTimeISO"`
			TempC               int         `json:"tempC"`
			TempF               int         `json:"tempF"`
			DewpointC           int         `json:"dewpointC"`
			DewpointF           int         `json:"dewpointF"`
			Humidity            int         `json:"humidity"`
			PressureMB          int         `json:"pressureMB"`
			PressureIN          float64     `json:"pressureIN"`
			SpressureMB         int         `json:"spressureMB"`
			SpressureIN         float64     `json:"spressureIN"`
			AltimeterMB         int         `json:"altimeterMB"`
			AltimeterIN         float64     `json:"altimeterIN"`
			WindKTS             float64     `json:"windKTS"`
			WindKPH             float64     `json:"windKPH"`
			WindMPH             float64     `json:"windMPH"`
			WindSpeedKTS        float64     `json:"windSpeedKTS"`
			WindSpeedKPH        float64     `json:"windSpeedKPH"`
			WindSpeedMPH        float64     `json:"windSpeedMPH"`
			WindDirDEG          float64     `json:"windDirDEG"`
			WindDir             string      `json:"windDir"`
			WindGustKTS         interface{} `json:"windGustKTS"`
			WindGustKPH         interface{} `json:"windGustKPH"`
			WindGustMPH         interface{} `json:"windGustMPH"`
			FlightRule          string      `json:"flightRule"`
			VisibilityKM        float64     `json:"visibilityKM"`
			VisibilityMI        float64     `json:"visibilityMI"`
			Weather             string      `json:"weather"`
			WeatherShort        string      `json:"weatherShort"`
			WeatherCoded        string      `json:"weatherCoded"`
			WeatherPrimary      string      `json:"weatherPrimary"`
			WeatherPrimaryCoded string      `json:"weatherPrimaryCoded"`
			CloudsCoded         string      `json:"cloudsCoded"`
			Icon                string      `json:"icon"`
			HeatindexC          float64     `json:"heatindexC"`
			HeatindexF          int         `json:"heatindexF"`
			WindchillC          float64     `json:"windchillC"`
			WindchillF          int         `json:"windchillF"`
			FeelslikeC          float64     `json:"feelslikeC"`
			FeelslikeF          int         `json:"feelslikeF"`
			IsDay               bool        `json:"isDay"`
			Sunrise             int         `json:"sunrise"`
			SunriseISO          time.Time   `json:"sunriseISO"`
			Sunset              int         `json:"sunset"`
			SunsetISO           time.Time   `json:"sunsetISO"`
			SnowDepthCM         interface{} `json:"snowDepthCM"`
			SnowDepthIN         interface{} `json:"snowDepthIN"`
			PrecipMM            int         `json:"precipMM"`
			PrecipIN            int         `json:"precipIN"`
			SolradWM2           int         `json:"solradWM2"`
			SolradMethod        string      `json:"solradMethod"`
			CeilingFT           int         `json:"ceilingFT"`
			CeilingM            float64     `json:"ceilingM"`
			Light               int         `json:"light"`
			Uvi                 interface{} `json:"uvi"`
			Qc                  string      `json:"QC"`
			QCcode              int         `json:"QCcode"`
			TrustFactor         int         `json:"trustFactor"`
			Sky                 int         `json:"sky"`
		} `json:"ob"`
		Raw        string `json:"raw"`
		RelativeTo struct {
			Lat        float64 `json:"lat"`
			Long       float64 `json:"long"`
			Bearing    int     `json:"bearing"`
			BearingENG string  `json:"bearingENG"`
			DistanceKM float64 `json:"distanceKM"`
			DistanceMI float64 `json:"distanceMI"`
		} `json:"relativeTo"`
	} `json:"response"`
}

func (server *Server) GetRecommendedLocations(ctx *gin.Context) {

	graphqlClient := gclient.NewClient(server.config.RecommendationServiceAddress + "/v1/locations")

	query := `{
				cities {
					data {
						city
						latitude
						longitude
						population
					}
					error
				}
			}`

	graphqlRequest := gclient.NewRequest(query)

	var graphqlResponse interface{}
	if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		log.Fatalf("Error querying GeoDB, error: %v", err)
	}

	// Convert map to json string
	rJSON, err := json.Marshal(graphqlResponse)
	if err != nil {
		log.Panic("Cannot marshal graphqlResponse to JSON")
	}

	// Convert struct
	var locationsResponse LocationsResponse
	err = json.Unmarshal(rJSON, &locationsResponse)
	if err != nil {
		log.Panic("Cannot unmarshal LocationResponse")
	}

	for key, value := range locationsResponse.Cities.Data {
		fmt.Println(key)
		fmt.Println(value)

		url := server.config.RecommendationServiceAddress + "/v1/weather"

		req, _ := http.NewRequest("GET", url, nil)

		params := req.URL.Query()
		params.Add("lat", fmt.Sprintf("%v", value.Latitude))
		params.Add("long", fmt.Sprintf("%v", value.Longitude))
		req.URL.RawQuery = params.Encode()

		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)

		var r weatherApiResponse
		err := json.Unmarshal(body, &r)
		if err != nil {
			log.Panic("Cannot unmarshal Response")
		}

		locationsResponse.Cities.Data[key].TempC = r.Response.Ob.TempC
		locationsResponse.Cities.Data[key].Weather = r.Response.Ob.Weather
		locationsResponse.Cities.Data[key].WeatherShort = r.Response.Ob.WeatherShort
	}

	ctx.JSON(http.StatusOK, locationsResponse)
}
