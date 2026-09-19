package main

import (
	"log"
	"net/netip"

	"github.com/oschwald/geoip2-golang/v2"
)

type LookupResult struct {
	IP      string `json:"ip"`
	Country string `json:"country"`
	City    string `json:"city"`
	ASN     uint   `json:"asn"`
	Org     string `json:"org"`
}

type LookupService struct {
	countryDB *geoip2.Reader
	cityDB    *geoip2.Reader
	asnDB     *geoip2.Reader
}

func NewLookService(countryPath, cityPath, asnPath string) (*LookupService, error) {
	country, err := geoip2.Open(countryPath)
	if err != nil {
		return nil, err
	}

	city, err := geoip2.Open(cityPath)
	if err != nil {
		return nil, err
	}

	asn, err := geoip2.Open(asnPath)
	if err != nil {
		err := country.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	return &LookupService{
		countryDB: country,
		cityDB:    city,
		asnDB:     asn,
	}, nil
}

func (s *LookupService) Lookup(ipString string) (*LookupResult, error) {
	ip, _ := netip.ParseAddr(ipString)
	countryRecord, err := s.countryDB.Country(ip)
	if err != nil {
		log.Fatalf("Invalid IP address: %s", ip)
	}
	cityRecord, err := s.cityDB.City(ip)
	if err != nil {
		log.Fatalf("Invalid IP address: %s", ip)
	}

	asnRecord, _ := s.asnDB.ASN(ip) // Non-fatal if ASN not set

	result := &LookupResult{
		IP:      ipString,
		Country: countryRecord.Country.Names.English,
		City:    cityRecord.City.Names.English,
	}

	if asnRecord != nil {
		result.ASN = asnRecord.AutonomousSystemNumber
		result.Org = asnRecord.AutonomousSystemOrganization
	}

	return result, nil
}

// Close ensures all three databases release their memory maps and file handles
func (s *LookupService) Close() error {
	var firstErr error
	if s.countryDB != nil {
		if err := s.countryDB.Close(); err != nil {
			firstErr = err
		}
	}

	if s.cityDB != nil {
		if err := s.cityDB.Close(); err != nil {
			firstErr = err
		}
	}

	if s.asnDB != nil {
		if err := s.asnDB.Close(); err != nil {
			firstErr = err
		}
	}

	return firstErr
}
