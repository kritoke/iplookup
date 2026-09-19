# iplookup

A command-line tool written in Go to look up country, city, and organization (ASN) information for an IP address or hostname using MaxMind GeoLite2 databases.

Hostnames are resolved with `LookupNetIP`. The first address is used, and the result is shown as `IP (hostname)`.

## Prerequisites

- Go 1.27+
- MaxMind GeoLite2 database files:
  - `GeoLite2-Country.mmdb`
  - `GeoLite2-City.mmdb`
  - `GeoLite2-ASN.mmdb`

If the databases are missing, run `iplookup update`.

## Configuration

Create or update `config.kdl` in the working directory. A default file is created on first run if it does not exist:

```kdl
maxmind {
    account-id "YOUR_ACCOUNT_ID"
    license-key "YOUR_LICENSE_KEY"
    db-path "."
}
```

- `account-id`: Your MaxMind account ID. Required for `iplookup update`.
- `license-key`: Your MaxMind license key. Required for `iplookup update`.
- `db-path`: Directory that contains `GeoLite2-Country.mmdb`, `GeoLite2-City.mmdb`, and `GeoLite2-ASN.mmdb`.

## Usage

```sh
iplookup <IP_OR_HOSTNAME>
iplookup update
iplookup help
```

### Run Directly

```sh
go run . 8.8.8.8
go run . dns.google
go run . update
```

IP lookup:

```text
IP Address: 8.8.8.8
Country: United States
City:
Organization: Google LLC
```

Hostname lookup:

```text
IP Address: 8.8.8.8 (dns.google)
Country: United States
City:
Organization: Google LLC
```

City is often empty. Anycast and similar IPs may report the ISP's registered country when no geolocated country is present.

### Build

```sh
go build -o iplookup
./iplookup 8.8.8.8
./iplookup update
```

Or run the cross-compilation script:

```sh
./build.sh
```
