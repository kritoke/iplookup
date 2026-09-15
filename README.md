# iplookup

A command-line tool written in Go to look up country and organization (ASN) information for an IP address using MaxMind GeoLite2 databases.

## Prerequisites

- Go 1.20+ (or compatible version)
- MaxMind GeoLite2 database files:
  - `GeoLite2-Country.mmdb`
  - `GeoLite2-ASN.mmdb`

## Configuration

Create or update `config.kdl` in the root directory:

```kdl
maxmind {
    db-path "."
}
```

- `db-path`: The directory path where the `GeoLite2-Country.mmdb` and `GeoLite2-ASN.mmdb` database files are located.

## Usage

### Run Directly

```sh
go run . <IP_ADDRESS>
```

Example:

```sh
go run . 8.8.8.8
```

Output:

```text
IP Address: 8.8.8.8
Country: United States
Organization: GOOGLE
```

### Build

To build binaries locally:

```sh
go build -o iplookup
./iplookup <IP_ADDRESS>
```

Or run the cross-compilation script:

```sh
./build.sh
```
