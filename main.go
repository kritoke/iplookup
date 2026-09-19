package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	editionCountry = "GeoLite2-Country"
	editionCity    = "GeoLite2-City"
	editionASN     = "GeoLite2-ASN"
)

var geoLiteEditions = []string{editionCountry, editionCity, editionASN}

func mmdbPath(cfg *Config, edition string) string {
	return filepath.Join(cfg.MaxMind.DbPath, edition+".mmdb")
}

func main() {
	if len(os.Args) < 2 {
		prompt()
		os.Exit(1)
	}

	cfg, err := LoadOrCreateConfig("config.kdl")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	switch os.Args[1] {
	case "update":
		if err := updateDatabases(cfg); err != nil {
			log.Fatal(err)
		}
	case "-h", "--help", "help":
		prompt()
	default:
		if err := lookupIP(cfg, os.Args[1]); err != nil {
			log.Fatal(err)
		}
	}
}

func prompt() {
	fmt.Println("Usage: iplookup <IP_ADDRESS>")
	fmt.Println("       iplookup update")
}

func updateDatabases(cfg *Config) error {
	accountID := cfg.MaxMind.AccountID
	licenseKey := cfg.MaxMind.LicenseKey

	for _, edition := range geoLiteEditions {
		if err := UpdateMaxMindGeoLite2(accountID, licenseKey, edition, mmdbPath(cfg, edition)); err != nil {
			return fmt.Errorf("update %s database: %w", edition, err)
		}
	}

	fmt.Printf("Updated %s databases\n", strings.Join(geoLiteEditions, ", "))
	return nil
}

func lookupIP(cfg *Config, ipAddress string) error {
	svc, err := NewLookupService(
		mmdbPath(cfg, editionCountry),
		mmdbPath(cfg, editionCity),
		mmdbPath(cfg, editionASN),
	)
	if err != nil {
		return fmt.Errorf("database error: %w (ensure GeoLite2 databases exist; try: iplookup update)", err)
	}
	defer svc.Close()

	result := svc.Lookup(ipAddress)

	fmt.Printf("IP Address: %s\nCountry: %s\nCity: %s\nOrganization: %s\n",
		result.IP,
		result.Country,
		result.City,
		result.Org,
	)
	return nil
}
