package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// UpdateMaxMindGeoLite2 downloads a MaxMind GeoLite2 edition and extracts the .mmdb file to destPath.
func UpdateMaxMindGeoLite2(accountID, licenseKey, edition, destPath string) error {
	url := fmt.Sprintf("https://download.maxmind.com/geoip/databases/%s/download?suffix=tar.gz", edition)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.SetBasicAuth(accountID, licenseKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download %s: %w", edition, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download %s: unexpected status %d", edition, resp.StatusCode)
	}

	gzReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar archive: %w", err)
		}

		if !strings.HasSuffix(header.Name, ".mmdb") {
			continue
		}

		tmpPath := destPath + ".tmp"
		outFile, err := os.Create(tmpPath)
		if err != nil {
			return fmt.Errorf("failed to create temp file: %w", err)
		}

		if _, err := io.Copy(outFile, tarReader); err != nil {
			_ = outFile.Close()
			os.Remove(tmpPath)
			return fmt.Errorf("failed to write database: %w", err)
		}
		if err := outFile.Close(); err != nil {
			os.Remove(tmpPath)
			return fmt.Errorf("failed to close temp file: %w", err)
		}

		return os.Rename(tmpPath, destPath)
	}

	return fmt.Errorf("no .mmdb file found in the tar archive")
}
