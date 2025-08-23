package handlers

import (
	"cache-demo/internal/cache"
	"cache-demo/internal/templates"
	"cache-demo/internal/types"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func Home(w http.ResponseWriter, r *http.Request, cache *cache.Cache) {
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
	req, err := http.NewRequestWithContext(
		r.Context(),
		http.MethodGet,
		fmt.Sprintf("https://api.weather.gov/points/%s,%s", latitude, longitude),
		nil)
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
	var gridData types.GridData
	if err = json.NewDecoder(resp.Body).Decode(&gridData); err != nil {
		http.Error(w, "Error decoding data", http.StatusInternalServerError)
		log.Println("Error decoding data")
		return
	}
	gridKey := gridData.ParseToKey()
	weatherReport, exists := cache.Get(*gridKey)
	if !exists {
		log.Println(gridData.Properties.GridID, gridData.Properties.GridX, gridData.Properties.GridY)
		req, err = http.NewRequestWithContext(r.Context(), http.MethodGet, fmt.Sprintf("https://api.weather.gov/gridpoints/%s/%d,%d/forecast",
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
		var forecast types.Forecast
		json.NewDecoder(resp.Body).Decode(&forecast)
		forecastValue := forecast.ParseToVal()
		log.Println(forecast.Properties.Period[0].Name, forecast.Properties.Period[0].Temperature, forecast.Properties.Period[0].TempU,
			forecast.Properties.Period[0].ProbabilityPrecipitation.Value, forecast.Properties.Period[0].WindSpeed,
			forecast.Properties.Period[0].WindDirection, forecast.Properties.Period[0].Short, forecast.Properties.Period[0].Detailed)
		forecast.Properties.Period[0].Renderable = true
		type View struct{ Period types.Period }
		view := View{
			Period: forecastValue.Periods[0],
		}
		cache.Set(*gridKey, *forecastValue)
		templates.Render(w, view, "web/templates/index.html")
	} else {
		type View struct{ Period types.Period }
		view := View{
			Period: weatherReport.Periods[0],
		}
		templates.Render(w, view, "web/templates/index.html")
	}

}
