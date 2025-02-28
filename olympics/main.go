//go:build !solution

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
)

var athletesData []Athlete

func getAthleteInfoHandler(w http.ResponseWriter, r *http.Request) {
	var athleteName string
	queryParams := r.URL.Query()

	nameParams, hasName := queryParams["name"]
	if !hasName || len(nameParams[0]) == 0 {
		http.Error(w, "missing 'name' parameter", http.StatusBadRequest)
		return
	}
	athleteName = nameParams[0]

	filteredAthletes := filter(athletesData, func(athlete Athlete) bool {
		return athlete.Athlete == athleteName
	})

	if len(filteredAthletes) == 0 {
		http.Error(w, "athlete not found", http.StatusNotFound)
		return
	}

	athleteInfo := athleteToInfo(&filteredAthletes[0])

	for _, athlete := range filteredAthletes {
		yearMedals, exists := athleteInfo.MedalsByYear[athlete.Year]
		if !exists {
			athleteInfo.MedalsByYear[athlete.Year] = &Medals{}
			yearMedals = athleteInfo.MedalsByYear[athlete.Year]
		}
		yearMedals.Gold += athlete.Gold
		yearMedals.Silver += athlete.Silver
		yearMedals.Bronze += athlete.Bronze
		yearMedals.Total += athlete.Gold + athlete.Silver + athlete.Bronze

		athleteInfo.Medals.Gold += athlete.Gold
		athleteInfo.Medals.Silver += athlete.Silver
		athleteInfo.Medals.Bronze += athlete.Bronze
		athleteInfo.Medals.Total += athlete.Gold + athlete.Silver + athlete.Bronze
	}

	responseData, err := json.Marshal(&athleteInfo)
	if err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, writeErr := w.Write(responseData)

	if writeErr != nil {
		http.Error(w, "error writing response", http.StatusInternalServerError)
	}
}

func getTopAthletesHandler(w http.ResponseWriter, r *http.Request) {
	var sportName string
	var resultLimit int
	queryParams := r.URL.Query()

	sportParams, hasSport := queryParams["sport"]
	if hasSport && sportParams[0] != "" {
		sportName = sportParams[0]
	} else {
		http.Error(w, "missing or empty 'sport' parameter", http.StatusBadRequest)
		return
	}

	limitParams, hasLimit := queryParams["limit"]
	if hasLimit && limitParams[0] != "" {
		var err error
		resultLimit, err = strconv.Atoi(limitParams[0])
		if err != nil {
			http.Error(w, "invalid 'limit' parameter", http.StatusBadRequest)
			return
		}
	} else {
		resultLimit = 3
	}

	filteredAthletes := filter(athletesData, func(athlete Athlete) bool {
		return athlete.Sport == sportName
	})

	if len(filteredAthletes) == 0 {
		http.Error(w, "no athletes found for the given sport", http.StatusNotFound)
		return
	}

	athletesInfo := getAthletes(filteredAthletes)

	infoList := make([]*AthleteInfo, 0, len(athletesInfo))
	for _, info := range athletesInfo {
		infoList = append(infoList, info)
	}

	sort.Slice(infoList, func(i, j int) bool {
		if infoList[i].Medals.Gold != infoList[j].Medals.Gold {
			return infoList[i].Medals.Gold > infoList[j].Medals.Gold
		}
		if infoList[i].Medals.Silver != infoList[j].Medals.Silver {
			return infoList[i].Medals.Silver > infoList[j].Medals.Silver
		}
		if infoList[i].Medals.Bronze != infoList[j].Medals.Bronze {
			return infoList[i].Medals.Bronze > infoList[j].Medals.Bronze
		}

		return infoList[i].Athlete < infoList[j].Athlete
	})

	limit := int(math.Min(float64(resultLimit), float64(len(infoList))))
	responseData, err := json.Marshal(infoList[:limit])

	if err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, writeErr := w.Write(responseData)

	if writeErr != nil {
		http.Error(w, "error writing response", http.StatusInternalServerError)
	}
}

func getTopCountriesHandler(w http.ResponseWriter, r *http.Request) {
	var yearFilter int
	var resultLimit int
	queryParams := r.URL.Query()

	yearParams, hasYear := queryParams["year"]
	if hasYear && yearParams[0] != "" {
		var err error
		yearFilter, err = strconv.Atoi(yearParams[0])
		if err != nil {
			http.Error(w, "invalid 'year' parameter", http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "missing or empty 'year' parameter", http.StatusBadRequest)
		return
	}

	limitParams, hasLimit := queryParams["limit"]
	if hasLimit && limitParams[0] != "" {
		var err error
		resultLimit, err = strconv.Atoi(limitParams[0])
		if err != nil {
			http.Error(w, "invalid 'limit' parameter", http.StatusBadRequest)
			return
		}
	} else {
		resultLimit = 3
	}

	filteredAthletes := filter(athletesData, func(athlete Athlete) bool {
		return athlete.Year == yearFilter
	})

	if len(filteredAthletes) == 0 {
		http.Error(w, "no data for the given year", http.StatusNotFound)
		return
	}

	countryInfos := getCountries(filteredAthletes)

	infoList := make([]*CountryInfo, 0, len(countryInfos))
	for _, info := range countryInfos {
		infoList = append(infoList, info)
	}

	sort.Slice(infoList, func(i, j int) bool {
		if infoList[i].Gold != infoList[j].Gold {
			return infoList[i].Gold > infoList[j].Gold
		}
		if infoList[i].Silver != infoList[j].Silver {
			return infoList[i].Silver > infoList[j].Silver
		}
		if infoList[i].Bronze != infoList[j].Bronze {
			return infoList[i].Bronze > infoList[j].Bronze
		}

		return infoList[i].Country < infoList[j].Country
	})

	limit := int(math.Min(float64(resultLimit), float64(len(infoList))))
	responseData, err := json.Marshal(infoList[:limit])

	if err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, writeErr := w.Write(responseData)

	if writeErr != nil {
		http.Error(w, "error writing response", http.StatusInternalServerError)
	}
}

func main() {
	portFlag := flag.String("port", "80", "server port")
	dataFilePath := flag.String("data", "./olympics/testdata/olympicWinners.json", "path to JSON data file")
	flag.Parse()

	dataFile, err := os.Open(*dataFilePath)
	if err != nil {
		log.Fatalf("error opening data file: %v", err)
	}
	defer dataFile.Close()

	dataContent, readErr := io.ReadAll(dataFile)
	if readErr != nil {
		log.Fatalf("error reading data file: %v", readErr)
	}

	err = json.Unmarshal(dataContent, &athletesData)
	if err != nil {
		log.Fatalf("error parsing JSON data: %v", err)
	}

	http.HandleFunc("/athlete-info", getAthleteInfoHandler)
	http.HandleFunc("/top-athletes-in-sport", getTopAthletesHandler)
	http.HandleFunc("/top-countries-in-year", getTopCountriesHandler)

	serverAddress := fmt.Sprintf(":%s", *portFlag)
	log.Fatal(http.ListenAndServe(serverAddress, nil))
}
