package store

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const urlSite = "https://shorturl.haln.dev"

var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

func encodeBytes(decoded []byte) []byte {
	encoded := base64.URLEncoding.AppendEncode(nil, decoded)
	for i := len(encoded) - 1; i >= 0; i-- {
		if encoded[i] != '=' {
			break
		}
		encoded[i] = '.'
	}
	return encoded
}

// mutates `encoded` in place!
func decodeBytes(encoded []byte) ([]byte, error) {
	for i := len(encoded) - 1; i >= 0; i-- {
		if encoded[i] != '.' {
			break
		}
		encoded[i] = '='
	}
	return base64.URLEncoding.AppendDecode(nil, encoded)
}

var generatedURLPattern = regexp.MustCompile(`id="generated-url" href="([^"]*)"`)

func Put(blob []byte) (string, error) {
	body := append([]byte("url="), encodeBytes(blob)...)

	req, err := http.NewRequest("POST", urlSite+"/generate", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	pageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err

	}
	matches := generatedURLPattern.FindSubmatch(pageBytes)
	if len(matches) < 2 {
		return "", fmt.Errorf("failed to find URL on page")
	}
	ref := strings.TrimPrefix(string(matches[1]), "/")

	return ref, nil
}

func Get(ref string) ([]byte, error) {
	httpClient := *httpClient
	// never follow redirects ;)
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	req, err := http.NewRequest("GET", urlSite+"/"+ref, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusPermanentRedirect {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("response was not a redirect: %v\n%v", resp.Status, string(body))
	}

	url := resp.Header.Get("Location")
	blob, err := decodeBytes([]byte(strings.TrimPrefix(url, "/")))
	if err != nil {
		return nil, fmt.Errorf("failed to decode blob: %w", err)
	}

	return blob, nil
}
