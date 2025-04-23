package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// response from https://geocode.maps.co
// comes back in an array ordered by "importance". Probably safe to take the first one:
// [
//   {
//     "place_id": 318254891,
//     "licence": "Data © OpenStreetMap contributors, ODbL 1.0. https://osm.org/copyright",
//     "osm_type": "node",
//     "osm_id": 158623191,
//     "boundingbox": [
//       "40.0248306",
//       "40.0648306",
//       "-75.4588053",
//       "-75.4188053"
//     ],
//     "lat": "40.0448306",
//     "lon": "-75.4388053",
//     "display_name": "Berwyn, Easttown Township, Chester County, Pennsylvania, 19312, United States",
//     "class": "place",
//     "type": "village",
//     "importance": 0.469575729943735
//   },
//   {
//     "place_id": 318253942,
//     "licence": "Data © OpenStreetMap contributors, ODbL 1.0. https://osm.org/copyright",
//     "osm_type": "node",
//     "osm_id": 2058688121,
//     "boundingbox": [
//       "40.0430849",
//       "40.0530849",
//       "-75.4475098",
//       "-75.4375098"
//     ],
//     "lat": "40.0480849",
//     "lon": "-75.4425098",
//     "display_name": "Berwyn, East Lancaster Avenue, Berwyn, Easttown Township, Chester County, Pennsylvania, 19301, United States",
//     "class": "railway",
//     "type": "station",
//     "importance": 0.2419138440405379
//   }
// ]

type GeoLocation struct {
	PlaceId     string   `json:"place_id"`
	License     string   `json:"license"`
	OsmType     string   `json:"osm_type"`
	OsmId       string   `json:"osm_id"`
	BoundingBox []string `json:"boundingbox"`
	Latitude    float64  `json:"lat"`
	Longitude   float64  `json:"lon"`
	DisplayName string   `json:"display_name"`
	Class       string   `json:"class"`
	Type        string   `json:"type"`
	Importance  float64  `json:"importance"`
}

func requestLocation(req *http.Request) (geolocation GeoLocation, err error) {
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return geolocation, fmt.Errorf("http request to %s failed: %s", req.URL, err.Error())
	}
	defer resp.Body.Close()

	// decode the body
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&geolocation)
	resp.Body.Close()

	if err != nil {
		return geolocation, fmt.Errorf("decoding the response from %s failed: %s", req.URL, err)
	}

	return geolocation, nil
}

func autolocate() (geolocation GeoLocation, err error) {
	uri := "http://www.telize.com/geoip"

	req, err := createRequest(uri, "GET", map[string]string{})
	if err != nil {
		return geolocation, err
	}
	return requestLocation(req)
}

func locate(location string) (geolocation GeoLocation, err error) {
	if location == "" {
		return autolocate()
	}

	uri := "https://geocode.jessfraz.com/geocode"

	req, err := createRequest(uri, "POST", map[string]string{"location": location})
	if err != nil {
		return geolocation, err
	}

	return requestLocation(req)
}
