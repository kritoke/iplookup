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
    account-id "YOUR_ACCOUNT_ID"
    license-key "YOUR_LICENSE_KEY"
    db-path "."
}
```

- `account-id`: Your MaxMind account ID. Required for `iplookup update`.
- `license-key`: Your MaxMind license key. Required for `iplookup update`.
- `db-path`: Directory that contains `GeoLite2-Country.mmdb` and `GeoLite2-ASN.mmdb`.

## Usage

### Run Directly

```sh
go run . <IP_ADDRESS>
go run . update
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
./iplookup update
```

Or run the cross-compilation script:

```sh
./build.sh
```
