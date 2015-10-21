/*
Package psi contains the exported functions for scraping and
parsing the 24-Hour and 3-Hour PSI readings in Singapore
*/
package psi

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/alienchow/psi/parser"
	"github.com/alienchow/psi/region"
)

const (
	psiSite = `http://www.haze.gov.sg/haze-updates/psi`
)

// Reading is the interface through which package users access PSI data
type Reading interface {
	Refresh() error
	Get(region.Region) string
}

// NewReading creates a new Reading instance for loading and parsing PSI values
func NewReading() Reading {
	return &readingImpl{}
}

// readingImpl is the logic implementation for the exported Reading interface
type readingImpl struct {
	pageHTML       string
	twentyFourHour map[region.Region]string
	threeHour      string
}

// Refresh loads the page HTML from NEA website
func (r *readingImpl) Refresh() error {
	resp, err := http.Get(psiSite)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	r.pageHTML = string(body)
	return r.parsePSIValues()
}

// Get retrieves the PSI value string for the provided region
// Invalid region returns the three hour reading instead
func (r *readingImpl) Get(refRegion region.Region) string {
	if refRegion == region.Invalid {
		return r.threeHour
	}
	return r.twentyFourHour[refRegion]
}

// parsePSIValues parses the loaded HTML content into the various values
func (r *readingImpl) parsePSIValues() error {
	for _, regionType := range region.All {
		PSI := parser.GetFunc(regionType)(r.pageHTML)
		if PSI == "" {
			return errors.New("NEA Website format has changed. Failed to parse HTML")
		}

		switch regionType {
		case region.Invalid:
			r.threeHour = PSI
		default:
			r.twentyFourHour[regionType] = PSI
		}
	}
	return nil
}
