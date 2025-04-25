package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
)

// response from https://geocode.maps.co
// comes back in an array ordered by "importance". Probably safe to take the first one sorted by "importance":
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

// Response from ip-api.com
// {
//   "status": "success",
//   "country": "United States",
//   "countryCode": "US",
//   "region": "PA",
//   "regionName": "Pennsylvania",
//   "city": "Philadelphia",
//   "zip": "19143",
//   "lat": 39.9486,
//   "lon": -75.2339,
//   "timezone": "America/New_York",
//   "isp": "Verizon Communications",
//   "org": "MCI Communications Services, Inc. d/b/a Verizon Business",
//   "as": "AS701 Verizon Business",
//   "query": "71.185.185.227"
// }

type IpLocation struct {
	Status      string  `json:"success"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	TimeZone    string  `json:"timezone"`
	ISP         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
}

type GeoLocation struct {
	PlaceId     int      `json:"place_id"`
	License     string   `json:"license"`
	OsmType     string   `json:"osm_type"`
	OsmId       int      `json:"osm_id"`
	BoundingBox []string `json:"boundingbox"`
	Latitude    string   `json:"lat"`
	Longitude   string   `json:"lon"`
	DisplayName string   `json:"display_name"`
	Class       string   `json:"class"`
	Type        string   `json:"type"`
	Importance  float64  `json:"importance"`
}

type GeoLocations []GeoLocation

// If the user didn't enter a location, default to finding their current one by
// first determining their ip address, then doing a geo ip lookup. This function
// may be doing too much on its own.
//
// ToDo: Refactor?
func LocateByIp() (locationString string, err error) {
	ipLookupUri := "http://icanhazip.com"

	client := &http.Client{}
	req, err := http.NewRequest("GET", ipLookupUri, nil)
	response, err := client.Do(req)

	if err != nil {
		return locationString, fmt.Errorf("http request to %s failed: %s", req.URL, err.Error())
	}
	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error: reading the IP information from icanhazip.com", err)
		return
	}

	locationIp := strings.TrimSuffix(string(bodyBytes), "\n")

	geoLocationLookupUri := "http://ip-api.com/json/" + locationIp

	locationRequest, err := http.NewRequest("GET", geoLocationLookupUri, nil)
	locationResponse, err := client.Do(locationRequest)
	if err != nil {
		return locationString, fmt.Errorf("http request to %s failed: %s", locationRequest.URL, err.Error())
	}
	defer locationResponse.Body.Close()

	// Unmarshall the json response body so we can parse out the zip
	var ipLocation IpLocation
	dec := json.NewDecoder(locationResponse.Body)
	err = dec.Decode(&ipLocation)
	response.Body.Close()

	if err != nil {
		return "", fmt.Errorf("failed to decode and unmarshall response from ip-api.com %s", err.Error())
	}

	return ipLocation.Zip, err
}

// Using the location info given by the user, find thier lat and longs using the
// https://geocode.maps.co API
func locate(location string) (geolocation GeoLocation, err error) {

	if location == "" {
		location, err = LocateByIp()
	}

	// Concatenate together the url + location requested + API key (http GET)
	uri := "https://geocode.maps.co/search?q=" + location + "&api_key=" + os.Getenv("GEOCODING_API_KEY")

	client := &http.Client{}
	req, err := http.NewRequest("GET", uri, nil)
	response, err := client.Do(req)

	if err != nil {
		return geolocation, fmt.Errorf("http request to %s failed: %s", req.URL, err.Error())
	}
	defer response.Body.Close()

	// Decode the body, we should get back an array of Geolcations to unmarshall
	var locations []GeoLocation
	dec := json.NewDecoder(response.Body)
	err = dec.Decode(&locations)
	response.Body.Close()

	// Sort by the "importance" field in descending order. This should give us the _most_ relevant location
	sort.Slice(locations, func(i, j int) bool {
		return locations[i].Importance > locations[j].Importance
	})

	// Take the first, sorted location as it's the cloest match by "importance"
	if len(locations) > 0 {
		geolocation = locations[0]
	} else {
		return geolocation, fmt.Errorf("failed to return any locations %s", err.Error())
	}

	if err != nil {
		return geolocation, err
	}

	return geolocation, err
}
