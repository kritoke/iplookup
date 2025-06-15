package main

import (
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: iplookup <IP_ADDRESS>")
		os.Exit(1)
	}

	ipAddress := os.Args[1]
	ip := net.ParseIP(ipAddress)
	if ip == nil {
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
		country.Country.Names["en"],
		country.Country.IsoCode,
	)
}
