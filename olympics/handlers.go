package main

func filter(list []Athlete, predicate func(Athlete) bool) []Athlete {
	var filtered []Athlete
	for _, athlete := range list {
		if predicate(athlete) {
			filtered = append(filtered, athlete)
		}
	}
	return filtered
}

func getAthletes(slc []Athlete) map[string]*AthleteInfo {
	athletesMap := make(map[string]*AthleteInfo)
	for _, athlete := range slc {
		addAthleteInfo(athletesMap, athlete)
	}
	return athletesMap
}

func addAthleteInfo(athletesMap map[string]*AthleteInfo, athlete Athlete) {
	info, exists := athletesMap[athlete.Athlete]
	if !exists {
		info = initializeAthleteInfo(athlete)
		athletesMap[athlete.Athlete] = info
	}
	updateAthleteInfo(info, athlete)
}

func initializeAthleteInfo(athlete Athlete) *AthleteInfo {
	return &AthleteInfo{
		Athlete:      athlete.Athlete,
		Country:      athlete.Country,
		Medals:       Medals{},
		MedalsByYear: make(map[int]*Medals),
	}
}

func updateAthleteInfo(info *AthleteInfo, athlete Athlete) {
	yearMedals := getOrCreateYearMedals(info.MedalsByYear, athlete.Year)
	updateMedals(yearMedals, athlete)
	updateMedals(&info.Medals, athlete)
}

func getOrCreateYearMedals(medalsByYear map[int]*Medals, year int) *Medals {
	if _, exists := medalsByYear[year]; !exists {
		medalsByYear[year] = &Medals{}
	}
	return medalsByYear[year]
}

func updateMedals(medals *Medals, athlete Athlete) {
	medals.Gold += athlete.Gold
	medals.Silver += athlete.Silver
	medals.Bronze += athlete.Bronze
	medals.Total += athlete.Gold + athlete.Silver + athlete.Bronze
}

func getCountries(slc []Athlete) map[string]*CountryInfo {
	countriesMap := make(map[string]*CountryInfo)
	for _, athlete := range slc {
		addCountryInfo(countriesMap, athlete)
	}
	return countriesMap
}

func addCountryInfo(countriesMap map[string]*CountryInfo, athlete Athlete) {
	countryInfo, exists := countriesMap[athlete.Country]
	if !exists {
		countryInfo = initializeCountryInfo(athlete.Country)
		countriesMap[athlete.Country] = countryInfo
	}
	updateCountryInfo(countryInfo, athlete)
}

func initializeCountryInfo(country string) *CountryInfo {
	return &CountryInfo{
		Country: country,
	}
}

func updateCountryInfo(countryInfo *CountryInfo, athlete Athlete) {
	countryInfo.Gold += athlete.Gold
	countryInfo.Silver += athlete.Silver
	countryInfo.Bronze += athlete.Bronze
	countryInfo.Total += athlete.Gold + athlete.Silver + athlete.Bronze
}
