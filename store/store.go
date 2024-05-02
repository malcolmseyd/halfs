package store

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

const urlSite = "https://shorturl.haln.dev"

// const maxChunkSize = 3946
const maxChunkSize = 2958

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

func putChunk(ctx context.Context, chunk []byte) (string, error) {
	body := append([]byte("url="), encodeBytes(chunk)...)

	req, err := http.NewRequestWithContext(ctx, "POST", urlSite+"/generate", bytes.NewReader(body))
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

func ceil(n int, d int) int {
	return (n + d - 1) / d
}

func Put(blob []byte) (string, error) {
	numChunks := ceil(len(blob), maxChunkSize)

	group, ctx := errgroup.WithContext(context.Background())
	group.SetLimit(50)
	refs := make([]string, numChunks)
	for i := range numChunks {
		start, end := i*maxChunkSize, (i+1)*maxChunkSize
		chunk := blob[start:min(end, len(blob))]

		i := i
		group.Go(func() error {
			ref, err := putChunk(ctx, chunk)
			if err != nil {
				return err
			}
			// fmt.Print(">")
			refs[i] = ref
			return nil
		})
	}
	err := group.Wait()
	if err != nil {
		return "", err
	}

	ref := strings.Builder{}
	for i := range refs {
		if i != 0 {
			ref.WriteRune('.')
		}
		ref.WriteString(refs[i])
	}
	return ref.String(), nil
}

func getChunk(ctx context.Context, ref string) ([]byte, error) {
	httpClient := *httpClient
	// never follow redirects ;)
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	req, err := http.NewRequestWithContext(ctx, "GET", urlSite+"/"+ref, nil)
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
	chunk, err := decodeBytes([]byte(strings.TrimPrefix(url, "/")))
	if err != nil {
		return nil, fmt.Errorf("failed to decode chunk: %w", err)
	}

	return chunk, nil
}

func Get(ref string) ([]byte, error) {
	refs := strings.Split(ref, ".")

	group, ctx := errgroup.WithContext(context.Background())
	group.SetLimit(50)
	chunks := make([][]byte, len(refs))
	for i, ref := range refs {
		i := i
		ref := ref
		group.Go(func() error {
			chunk, err := getChunk(ctx, ref)
			if err != nil {
				return err
			}
			// fmt.Print("<")
			chunks[i] = chunk
			return nil
		})
	}
	err := group.Wait()
	if err != nil {
		return nil, err
	}
	blob := bytes.Join(chunks, nil)

	return blob, nil
}
