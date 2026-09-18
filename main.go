package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

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
	countryPath := filepath.Join(cfg.MaxMind.DbPath, "GeoLite2-Country.mmdb")
	asnPath := filepath.Join(cfg.MaxMind.DbPath, "GeoLite2-ASN.mmdb")

	if err := UpdateMaxMindGeoLite2(accountID, licenseKey, "GeoLite2-Country", countryPath); err != nil {
		return fmt.Errorf("update country database: %w", err)
	}
	if err := UpdateMaxMindGeoLite2(accountID, licenseKey, "GeoLite2-ASN", asnPath); err != nil {
		return fmt.Errorf("update ASN database: %w", err)
	}

	fmt.Println("Updated GeoLite2-Country and GeoLite2-ASN databases")
	return nil
}

func lookupIP(cfg *Config, ipAddress string) error {
	countryPath := filepath.Join(cfg.MaxMind.DbPath, "GeoLite2-Country.mmdb")
	asnPath := filepath.Join(cfg.MaxMind.DbPath, "GeoLite2-ASN.mmdb")

	svc, err := NewLookService(countryPath, asnPath)
	if err != nil {
		return fmt.Errorf("database error: %w (ensure GeoLite2 databases exist; try: iplookup update)", err)
	}
	defer svc.Close()

	result, err := svc.Lookup(ipAddress)
	if err != nil {
		return fmt.Errorf("invalid IP address: %s", ipAddress)
	}

	fmt.Printf("IP Address: %s\nCountry: %s\nOrganization: %s\n",
		result.IP,
		result.Country,
		result.Org,
	)
	return nil
}
