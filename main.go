package main

import (
	"fmt"
	"log"
	"net/netip"
	"os"

	"github.com/oschwald/geoip2-golang/v2"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: iplookup <IP_ADDRESS>")
		os.Exit(1)
	}

	ipAddress := os.Args[1]
	ip, err := netip.ParseAddr(ipAddress)
	if err != nil {
		log.Fatalf("Invalid IP address: %s", ipAddress)
	}

	db, err := geoip2.Open("GeoLite2-Country.mmdb")
	if err != nil {
		log.Fatalf("Database error: %v (Ensure GeoLite2-Country.mmdb exists)", err)
	}
	defer func(db *geoip2.Reader) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Database error: %v (Ensure GeoLite2-Country.mmdb exists)", err)
		}
	}(db)

	country, err := db.Country(ip)
	if err != nil {
		log.Fatalf("Geolocation error: %v", err)
	}

	fmt.Printf("Country: %s (ISO: %s)\n",
		country.Country.Names.English,
		country.Country.ISOCode,
	)
}
