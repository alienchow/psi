package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const (
	psiSite = `http://www.haze.gov.sg/haze-updates/psi`
)

func main() {
	var arg string
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}

	var PSI string
	switch strings.ToUpper(arg) {
	case "N", "NORTH":
		PSI = getRegionPSI("North")
	case "S", "SOUTH":
		PSI = getRegionPSI("South")
	case "E", "EAST":
		PSI = getRegionPSI("East")
	case "W", "WEST":
		PSI = getRegionPSI("West")
	case "C", "CENTRAL":
		PSI = getRegionPSI("Central")
	default:
		PSI = get3HourPSI()
	}
	fmt.Println(PSI)
}

func getRegionPSI(region string) string {
	body, err := getPage()
	if err != nil {
		return "Error loading page at " + psiSite + " while loading PSI readings for " + region
	}
	PSI := getPSIByRegex(regionRegex(region), body)
	if PSI == "" {
		return psiSite + " HTML format changed. Failed to load PSI."
	}
	return "3-hour PSI Reading for " + region + ": " + PSI
}

func get3HourPSI() string {
	body, err := getPage()
	if err != nil {
		return "Error loading page at " + psiSite + " while loading overall readings."
	}
	PSI := getPSIByRegex("3-hr PSI: \\d+", body)
	if PSI == "" {
		return psiSite + " HTML format changed. Failed to load PSI."
	}
	return "Over 3-hour PSI reading is: " + PSI
}

func getPSIByRegex(regex, page string) string {
	r, _ := regexp.Compile(regex)
	matchString := r.FindString(page)
	matchString = strings.Replace(matchString, "3-hr PSI:", "", -1)
	psiRegex, _ := regexp.Compile("\\d{1,}")
	return psiRegex.FindString(matchString)
}

func regionRegex(region string) string {
	return "\\d{1,}</span>\\s*<span class=\"direction\">" + region
}

func getPage() (string, error) {
	resp, err := http.Get(psiSite)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
