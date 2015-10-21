package parser

import (
	"regexp"

	"github.com/alienchow/psi/region"
)

const (
	threeHourRegexString = `3-hr\s*PSI:\s*`
	overallRegexString   = `24-hr\s*PSI:\s*`
	regionRegexString    = `</span>\s*<span class="direction">`
	psiRegexString       = `\d+`
	psiRangeRegexString  = `\d+\s*-\s*\d+`
)

// Func is the parser function to be used for the  provided region
type Func func(string) string

// Func creates a closure to parse the page body based on the region specified
// If the region Invalid is provided, the default parserFunc will parse the 3 hour readings
func GetFunc(r region.Region) Func {
	switch r {
	case region.Invalid:
		return threeHourFunc
	case region.Overall:
		return overallFunc
	}
	return regionFunc(r)
}

// regionFunc is the closure function for parsing region 24-Hour PSI
func regionFunc(r region.Region) func(string) string {
	return func(pageHTML string) string {
		regex, _ := regexp.Compile(regionRegex(string(r)))
		matchString := regex.FindString(pageHTML)
		psiRegex, _ := regexp.Compile(psiRegexString)
		return psiRegex.FindString(matchString)
	}
}

// threeHourFunc is the closure function for parsing 3 hour PSI
func threeHourFunc(pageHTML string) string {
	regex, _ := regexp.Compile(threeHourRegexString + psiRegexString)
	matchString := regex.FindString(pageHTML)
	psiRegex, _ := regexp.Compile(psiRegexString)
	return psiRegex.FindString(matchString)
}

// overallFunc is the closure function for parsing the overall 24-hour PSI range
func overallFunc(pageHTML string) string {
	regex, _ := regexp.Compile(overallRegexString + psiRangeRegexString)
	matchString := regex.FindString(pageHTML)
	psiRegex, _ := regexp.Compile(psiRangeRegexString)
	return psiRegex.FindString(matchString)
}

func regionRegex(region string) string {
	return psiRegexString + regionRegexString + region
}
