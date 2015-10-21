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

			PSI := parser.GetFunc(currentRegion)(r.pageHTML)

			r.Lock()
			defer r.Unlock()

			if PSI == "" {
				r.err = errors.New("NEA Website format has changed, failed to parse HTML")
				return
			}

			switch currentRegion {
			case region.Invalid:
				r.threeHour = PSI
			default:
				r.twentyFourHour[currentRegion] = PSI
			}
		}(regionType)
	}

	wg.Wait()

	if r.err != nil {
		return r.err
	}
	return nil
}
