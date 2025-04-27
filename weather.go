package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Example return data from https://openweathermap.org/api/one-call-3
// {
//     "lat": 40.0389,
//     "lon": -75.4483,
//     "timezone": "America/New_York",
//     "timezone_offset": -14400,
//     "current": {
//         "dt": 1745601315,
//         "sunrise": 1745575746,
//         "sunset": 1745625000,
//         "temp": 296.58,
//         "feels_like": 296.47,
//         "pressure": 1022,
//         "humidity": 57,
//         "dew_point": 287.58,
//         "uvi": 5.91,
//         "clouds": 0,
//         "visibility": 10000,
//         "wind_speed": 5.66,
//         "wind_deg": 180,
//         "weather": [
//             {
//                 "id": 800,
//                 "main": "Clear",
//                 "description": "clear sky",
//                 "icon": "01d"
//             }
//         ]
//     },
//     "minutely": [
//         {
//             "dt": 1745601360,
//             "precipitation": 0
//         },
//         {
//             "dt": 1745601420,
//             "precipitation": 0
//         },
//         {
//             "dt": 1745601480,
//             "precipitation": 0
//         },
//         {
//             "dt": 1745601540,
//             "precipitation": 0
//         }
//     ],
//     "hourly": [
//         {
//             "dt": 1745600400,
//             "temp": 296.58,
//             "feels_like": 296.47,
//             "pressure": 1022,
//             "humidity": 57,
//             "dew_point": 287.58,
//             "uvi": 5.91,
//             "clouds": 0,
//             "visibility": 10000,
//             "wind_speed": 3.93,
//             "wind_deg": 181,
//             "wind_gust": 5.52,
//             "weather": [
//                 {
//                     "id": 800,
//                     "main": "Clear",
//                     "description": "clear sky",
//                     "icon": "01d"
//                 }
//             ],
//             "pop": 0
//         },
//         {
//             "dt": 1745604000,
//             "temp": 296.64,
//             "feels_like": 296.51,
//             "pressure": 1022,
//             "humidity": 56,
//             "dew_point": 287.37,
//             "uvi": 5.35,
//             "clouds": 20,
//             "visibility": 10000,
//             "wind_speed": 4.31,
//             "wind_deg": 179,
//             "wind_gust": 5.99,
//             "weather": [
//                 {
//                     "id": 801,
//                     "main": "Clouds",
//                     "description": "few clouds",
//                     "icon": "02d"
//                 }
//             ],
//             "pop": 0
//         },
//         {
//             "dt": 1745607600,
//             "temp": 296.8,
//             "feels_like": 296.63,
//             "pressure": 1021,
//             "humidity": 54,
//             "dew_point": 286.95,
//             "uvi": 4.49,
//             "clouds": 40,
//             "visibility": 10000,
//             "wind_speed": 4.38,
//             "wind_deg": 179,
//             "wind_gust": 6.24,
//             "weather": [
//                 {
//                     "id": 802,
//                     "main": "Clouds",
//                     "description": "scattered clouds",
//                     "icon": "03d"
//                 }
//             ],
//             "pop": 0
//         },
//         {
//             "dt": 1745611200,
//             "temp": 297.04,
//             "feels_like": 296.84,
//             "pressure": 1021,
//             "humidity": 52,
//             "dew_point": 286.6,
//             "uvi": 2.74,
//             "clouds": 60,
//             "visibility": 10000,
//             "wind_speed": 4.54,
//             "wind_deg": 175,
//             "wind_gust": 6.42,
//             "weather": [
//                 {
//                     "id": 803,
//                     "main": "Clouds",
//                     "description": "broken clouds",
//                     "icon": "04d"
//                 }
//             ],
//             "pop": 0
//         }
//     ],
//     "daily": [
//         {
//             "dt": 1745596800,
//             "sunrise": 1745575746,
//             "sunset": 1745625000,
//             "moonrise": 1745570880,
//             "moonset": 1745616600,
//             "moon_phase": 0.92,
//             "summary": "Expect a day of partly cloudy with rain",
//             "temp": {
//                 "day": 296.56,
//                 "min": 285.29,
//                 "max": 297.04,
//                 "night": 290.12,
//                 "eve": 296.23,
//                 "morn": 285.29
//             },
//             "feels_like": {
//                 "day": 296.44,
//                 "night": 290.09,
//                 "eve": 296.03,
//                 "morn": 284.65
//             },
//             "pressure": 1022,
//             "humidity": 57,
//             "dew_point": 287.56,
//             "wind_speed": 4.83,
//             "wind_deg": 164,
//             "wind_gust": 12.15,
//             "weather": [
//                 {
//                     "id": 500,
//                     "main": "Rain",
//                     "description": "light rain",
//                     "icon": "10d"
//                 }
//             ],
//             "clouds": 20,
//             "pop": 0.99,
//             "rain": 0.4,
//             "uvi": 5.99
//         },
//         {
//             "dt": 1745683200,
//             "sunrise": 1745662066,
//             "sunset": 1745711462,
//             "moonrise": 1745658840,
//             "moonset": 1745707680,
//             "moon_phase": 0.95,
//             "summary": "Expect a day of partly cloudy with rain",
//             "temp": {
//                 "day": 291.02,
//                 "min": 282.02,
//                 "max": 293.97,
//                 "night": 282.02,
//                 "eve": 291.56,
//                 "morn": 289.6
//             },
//             "feels_like": {
//                 "day": 291.24,
//                 "night": 278.9,
//                 "eve": 291.05,
//                 "morn": 289.86
//             },
//             "pressure": 1010,
//             "humidity": 91,
//             "dew_point": 289.52,
//             "wind_speed": 7.13,
//             "wind_deg": 288,
//             "wind_gust": 14.22,
//             "weather": [
//                 {
//                     "id": 501,
//                     "main": "Rain",
//                     "description": "moderate rain",
//                     "icon": "10d"
//                 }
//             ],
//             "clouds": 100,
//             "pop": 1,
//             "rain": 6.11,
//             "uvi": 3.71
//         }
//     ],
//     "alerts": [
//         {
//             "sender_name": "NWS Mount Holly NJ",
//             "event": "Special Weather Statement",
//             "start": 1745568540,
//             "end": 1745622000,
//             "description": "There is an increased risk for rapid fire spread this afternoon\nacross portions of New Jersey and eastern Pennsylvania. Minimum\nrelative humidity values will be around 25 to 35 percent combined\nwith southerly winds of 10 to 15 mph with gusts near 20 mph. High\ntemperatures will be in the mid 70s to near 80 degrees. These\nconditions, along with the continued drying of fine fuels, could\nsupport the rapid spread of any fires that ignite, which could\nquickly become difficult to control.\n\nOutdoor burning is strongly discouraged. Be sure to properly\nextinguish or dispose of any potential ignition sources, including\nsmoking materials such as cigarette butts.",
//             "tags": [
//                 "Other dangers"
//             ]
//         }
//     ]
// }

type Forecast struct {
	Alerts    []Alerts        `json:"alerts"`
	Currently CurrentWeather  `json:"currently"`
	Hourly    []HourlyWeather `json:"hourly"`
	Daily     DailyWeather    `json:"daily"`
	Latitude  float64         `json:"lat"`
	Longitude float64         `json:"lon"`
	Offset    int             `json:"timezone_offset"`
	Timezone  string          `json:"timezone"`
}

type CurrentWeather struct {
	Dt          int64       `json:"dt"`
	Sunrise     int64       `json:"sunrise"`
	Sunset      int64       `json:"sunset"`
	Temperature float64     `json:"temp"`
	FeelsLike   float64     `json:"feels_like"`
	Pressure    int         `json:"pressure"`
	Humidity    int         `json:"humidity"`
	DewPoint    float64     `json:"dew_point"`
	Uvi         float64     `json:"uvi"`
	Clouds      int         `json:"clouds"`
	Visibility  int         `json:"visibility"`
	WindSpeed   float64     `json:"wind_speed"`
	WindDegree  int         `json:"wind_deg"`
	Info        WeatherInfo `json:"weather"`
}

type DailyWeather struct {
	Dt          int64   `json:"dt"`
	Sunrise     int64   `json:"sunrise"`
	Sunset      int64   `json:"sunset"`
	Moonrise    int64   `json:"moonrise"`
	Moonset     int64   `json:"moonset"`
	MoonPhase   float64 `json:"moon_phase"`
	Summary     string  `json:"summary"`
	Temperature struct {
		Day   float64 `json:"day"`
		Min   float64 `json:"min"`
		Max   float64 `json:"max"`
		Night float64 `json:"night"`
		Eve   float64 `json:"eve"`
		Morn  float64 `json:"morn"`
	}
	FeelsLike struct {
		Day   float64 `json:"day"`
		Night float64 `json:"night"`
		Eve   float64 `json:"eve"`
		Morn  float64 `json:"morn"`
	}
	Pressure  int         `json:"pressure"`
	Humidity  int         `json:"humidity"`
	DewPoint  float64     `json:"dew_point"`
	WindSpeed float64     `json:"wind_speed"`
	WindDeg   int         `json:"wind_deg"`
	WindGust  float64     `json:"wind_gust"`
	Info      WeatherInfo `json:"weather"`
	Clouds    int         `json:"clouds"`
	Pop       float64     `json:"pop"`
	Rain      float64     `json:"rain"`
	Uvi       float64     `json:"uvi"`
}

type HourlyWeather struct {
	Dt          int64       `json:"dt"`
	Temperature float64     `json:"temp"`
	FeelsLike   float64     `json:"feels_like"`
	Pressure    int         `json:"pressure"`
	Humidity    int         `json:"humidity"`
	DewPoint    float64     `json:"dew_point"`
	Uvi         float64     `json:"uvi"`
	Clouds      int         `json:"clouds"`
	Visibility  int         `json:"visibility"`
	WindSpeed   float64     `json:"wind_speed"`
	WindDegree  int         `json:"wind_deg"`
	WindGust    float64     `json:"wind_gust"`
	Info        WeatherInfo `json:"weather"`
}

type Alerts struct {
	SenderName  string `json:"sender_name"`
	Event       string `json:"event"`
	Start       int64  `json:"start"`
	End         int64  `json:"end"`
	Description string `json:"description"`
}

type WeatherInfo struct {
	Id          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

func getForecast(data ForecastRequest) (forecast Forecast, err error) {
	client := &http.Client{}
	uri := "https://geocode.jessfraz.com/forecast"

	req, err := createRequest(uri, "POST", data)
	if err != nil {
		return forecast, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return forecast, fmt.Errorf("Http request to %s failed: %s", req.URL, err.Error())
	}
	defer resp.Body.Close()

	// decode the body
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&forecast)
	resp.Body.Close()

	if err != nil {
		return forecast, fmt.Errorf("Decoding the response from %s failed: %s", req.URL, err)
	}

	if forecast.Error != "" {
		return forecast, fmt.Errorf("The response returned: %s", forecast.Error)
	}

	return forecast, nil
}
