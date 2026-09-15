package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: iplookup <IP_ADDRESS>")
		os.Exit(1)
	}

	cfg, err := LoadConfig("config.kdl")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	countryPath := filepath.Join(cfg.MaxMind.DbPath, "GeoLite2-Country.mmdb")
	asnPath := filepath.Join(cfg.MaxMind.DbPath, "GeoLite2-ASN.mmdb")

	svc, err := NewLookService(countryPath, asnPath)
	if err != nil {
		log.Fatalf("Database error: %v (Ensure GeoLite2-ASN.mmdb exists)", err)
	}
	defer func(svc *LookupService) {
		err := svc.Close()
		if err != nil {
			log.Fatalf("Database closing failed: %v", err)
		}
	}(svc) // Closes both countryDb and asnDb on exit

	ipAddress := os.Args[1]
	result, err := svc.Lookup(ipAddress)
	if err != nil {
		log.Fatalf("Invalid IP address: %s", ipAddress)
	}

	fmt.Printf("IP Address: %s\nCountry: %s\nOrganization: %s\n",
		result.IP,
		result.Country,
		result.Org,
	)
}
