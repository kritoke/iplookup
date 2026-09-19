package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
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

func NewLookupService(countryPath, cityPath, asnPath string) (*LookupService, error) {
	s := &LookupService{}
	var err error

	if s.countryDB, err = geoip2.Open(countryPath); err != nil {
		return nil, err
	}
	if s.cityDB, err = geoip2.Open(cityPath); err != nil {
		_ = s.Close()
		return nil, err
	}
	if s.asnDB, err = geoip2.Open(asnPath); err != nil {
		_ = s.Close()
		return nil, err
	}
	return s, nil
}

func resolveAddr(s string) (netip.Addr, string, error) {
	if addr, err := netip.ParseAddr(s); err == nil {
		return addr, "", nil
	}
	ips, err := net.DefaultResolver.LookupNetIP(context.Background(), "ip", s)
	if err != nil {
		return netip.Addr{}, "", fmt.Errorf("could not resolve %q: %w", s, err)
	}
	if len(ips) == 0 {
		return netip.Addr{}, "", fmt.Errorf("lookup returned no addresses for %q", s)
	}
	return ips[0].Unmap(), s, nil
}

func (s *LookupService) Lookup(ipString string) *LookupResult {
	ip, host, err := resolveAddr(ipString)
	if err != nil {
		log.Fatalf("Invalid IP or domain address: %s", ipString)
	}

	countryRecord, err := s.countryDB.Country(ip)
	if err != nil {
		log.Fatalf("Country lookup failed for %s: %v", ipString, err)
	}
	cityRecord, err := s.cityDB.City(ip)
	if err != nil {
		log.Fatalf("City lookup failed for %s: %v", ipString, err)
	}

	asnRecord, _ := s.asnDB.ASN(ip) // Non-fatal if ASN not set

	ipName := ip.String()
	if host != "" {
		ipName = fmt.Sprintf("%s (%s)", ipName, host)
	}

	result := &LookupResult{
		IP:      ipName,
		Country: countryName(countryRecord),
		City:    cityRecord.City.Names.English,
	}

	if asnRecord != nil {
		result.ASN = asnRecord.AutonomousSystemNumber
		result.Org = asnRecord.AutonomousSystemOrganization
	}

	return result
}

// countryName prefers the geolocated country. Anycast and similar IPs often
// have no location country, only the ISP's registered country.
func countryName(rec *geoip2.Country) string {
	if rec.Country.Names.English != "" {
		return rec.Country.Names.English
	}
	return rec.RegisteredCountry.Names.English
}

// Close releases memory maps and file handles for any databases that were opened.
func (s *LookupService) Close() error {
	err := errors.Join(closeDB(s.countryDB), closeDB(s.cityDB), closeDB(s.asnDB))
	s.countryDB, s.cityDB, s.asnDB = nil, nil, nil
	return err
}

func closeDB(db *geoip2.Reader) error {
	if db == nil {
		return nil
	}
	return db.Close()
}
