package handlers

import (
	"cache-demo/internal/templates"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type GridData struct {
	Properties struct {
		GridID string `json:"gridId"`
		GridX  int64  `json:"gridX"`
		GridY  int64  `json:"gridY"`
	} `json:"properties"`
}

type Forecast struct {
	Properties struct {
		Period []Period `json:"periods"`
	} `json:"properties"`
}

type Period struct {
	Renderable               bool
	Name                     string `json:"name"`
	Temperature              int64  `json:"temperature"`
	TempU                    string `json:"temperatureUnit"`
	ProbabilityPrecipitation struct {
		Value int64 `json:"value"`
	} `json:"probabilityOfPrecipitation"`
	WindSpeed     string `json:"windSpeed"`
	WindDirection string `json:"windDirection"`
	Short         string `json:"shortForecast"`
	Detailed      string `json:"detailedForecast"`
}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "web/templates/index.html")
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing response", http.StatusInternalServerError)
		log.Println("Error parsing response")
		return
	}
	client := &http.Client{Timeout: 8 * time.Second}
	latitude := r.FormValue("latitude")
	latitude = strings.TrimSpace(latitude)
	longitude := r.FormValue("longitude")
	longitude = strings.TrimSpace(longitude)
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.weather.gov/points/%s,%s", latitude, longitude), nil)
	if err != nil {
		http.Error(w, "Error talking to weather api", http.StatusInternalServerError)
		return
	}
	req.Header.Set("User-Agent", "cache-demo hi@gmail.com")
	req.Header.Set("Accept", "application/geo+json")
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error talking to weather api", http.StatusInternalServerError)
		log.Println("Error accessing api")
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "upstream error", http.StatusBadGateway)
		log.Println("Error retreving data")
		return
	}
	var gridData GridData
	if err = json.NewDecoder(resp.Body).Decode(&gridData); err != nil {
		http.Error(w, "Error decoding data", http.StatusInternalServerError)
		log.Println("Error decoding data")
		return
	}
	log.Println(gridData.Properties.GridID, gridData.Properties.GridX, gridData.Properties.GridY)
	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.weather.gov/gridpoints/%s/%d,%d/forecast",
		gridData.Properties.GridID, gridData.Properties.GridX, gridData.Properties.GridY), nil)
	if err != nil {
		http.Error(w, "Error talking to weather api", http.StatusInternalServerError)
		return
	}
	req.Header.Set("User-Agent", "cache-demo hi@gmail.com")
	req.Header.Set("Accept", "application/geo+json")
	resp, err = client.Do(req)
	if err != nil {
		http.Error(w, "Error talking to weather api", http.StatusInternalServerError)
		log.Println("Error accessing api")
		return
	}
	defer resp.Body.Close()
	var forecast Forecast
	json.NewDecoder(resp.Body).Decode(&forecast)
	log.Println(forecast.Properties.Period[0].Name, forecast.Properties.Period[0].Temperature, forecast.Properties.Period[0].TempU,
		forecast.Properties.Period[0].ProbabilityPrecipitation.Value, forecast.Properties.Period[0].WindSpeed,
		forecast.Properties.Period[0].WindDirection, forecast.Properties.Period[0].Short, forecast.Properties.Period[0].Detailed)
	forecast.Properties.Period[0].Renderable = true
	type View struct{ Period Period }
	view := View{
		Period: forecast.Properties.Period[0],
	}
	templates.Render(w, view, "web/templates/index.html")
}
