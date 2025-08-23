package types

import "time"

type GridData struct {
	Properties struct {
		GridID string `json:"gridId"`
		GridX  int64  `json:"gridX"`
		GridY  int64  `json:"gridY"`
	} `json:"properties"`
}

func (g *GridData) ParseToKey() *GridKey {
	return &GridKey{
		GridID: g.Properties.GridID,
		GridX:  g.Properties.GridX,
		GridY:  g.Properties.GridY,
	}
}

func (f *Forecast) ParseToVal() *ForecastVals {
	return &ForecastVals{
		Periods:   f.Properties.Period,
		ExpiresAt: time.Now().Add(6 * time.Hour),
	}
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

type GridKey struct {
	GridID       string
	GridX, GridY int64
}

type ForecastVals struct {
	Periods   []Period
	ExpiresAt time.Time
}
