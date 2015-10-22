/*
Package psi contains the exported functions for scraping and
parsing the 24-Hour and 3-Hour PSI readings in Singapore
*/
package psi

import (
	"errors"
	"io/ioutil"
	"net/http"
	"sync"

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
	return &readingImpl{
		twentyFourHour: map[region.Region]string{},
	}
}

// readingImpl is the logic implementation for the exported Reading interface
type readingImpl struct {
	sync.Mutex
	pageHTML       string
	twentyFourHour map[region.Region]string
	threeHour      string
	err            error
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

// error checks the values of the reading and returns and error if fields are empty
func (r readingImpl) error() error {
	for region := range r.twentyFourHour {
		if r.twentyFourHour[region] == "" {
			return errors.New("Failed to parse PSI data for region: " + string(region))
		}
	}

	if r.threeHour == "" {
		return errors.New("Failed to parse PSI data for three hour reading")
	}
	return nil
}

// setPSI takes in the region for the respective PSI value and stores it to the appropriate field
func (r *readingImpl) setPSI(refRegion region.Region, PSI string) {
	r.Lock()
	defer r.Unlock()

	if refRegion == region.Invalid {
		r.threeHour = PSI
		return
	}
	r.twentyFourHour[refRegion] = PSI
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
	wg := &sync.WaitGroup{}

	for _, regionType := range region.All {
		wg.Add(1)

		go func(currentRegion region.Region) {
			defer wg.Done()
			r.setPSI(currentRegion, parser.GetFunc(currentRegion)(r.pageHTML))
		}(regionType)
	}

	wg.Wait()
	return r.error()
}
