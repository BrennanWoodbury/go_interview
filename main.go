// input.txt test input {"location": "foo"} {"location": "foo.com"} {"location": "bar.com"} {"location": "httpbar.com"} {"location": "https://bar.com"}
// saved the input as its own file
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Location struct {
	URL string `json:"location"`
}

type Payload struct {
	Data string `json:"data"`
}

func main() {
	data, err := os.ReadFile("./input.txt")
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		data, err := getData(line)
		if err != nil {
			fmt.Printf("unable to get status: %v", err)
			continue
		}
		if data.(string) == "foo" {
			fmt.Printf("data found: %s", data)
			os.Exit(0)
		}
	}
}

func getData(line string) (any, error) {
	var location Location
	var payload Payload
	var buf *bytes.Buffer = &bytes.Buffer{}

	err := json.Unmarshal([]byte(line), &location)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal line: %w", err)
	}

	if !strings.HasPrefix(location.URL, "http://") && !strings.HasPrefix(location.URL, "https://") {
		location.URL = "http://" + location.URL
	}

	re := regexp.MustCompile(`^https?://`)
	if !re.MatchString(location.URL) {
		noColonSlash := regexp.MustCompile(`^https?[^://]`)
		if noColonSlash.MatchString(location.URL) {
			switch {
			case strings.HasPrefix(location.URL, "https://"):
				break
			case strings.HasPrefix(location.URL, "http://"):
				break
			case strings.HasPrefix(location.URL, "https:/") || strings.HasPrefix(location.URL, "http:/"):
				newStr := strings.Split(location.URL, ":/")
				location.URL = newStr[0] + "://" + newStr[1]
			case strings.HasPrefix(location.URL, "https:") || strings.HasPrefix(location.URL, "http:"):
				newStr := strings.Split(location.URL, ":")
				location.URL = newStr[0] + "://" + newStr[1]
			case strings.HasPrefix(location.URL, "https//") || strings.HasPrefix(location.URL, "http//"):
				newStr := strings.Split(location.URL, "//")
				location.URL = newStr[0] + "://" + newStr[1]
			case strings.HasPrefix(location.URL, "https/") || strings.HasPrefix(location.URL, "http/"):
				newStr := strings.Split(location.URL, "//")
				location.URL = newStr[0] + "://" + newStr[1]
			default:
				break
			}

		}
	}

	TLD := []string{
		`\.com$`,
		`\.org$`,
		`\.gov$`,
		`\.net$`,
		`\.edu$`,
		`\.mil$`,
	}

	var match bool = false
	for i := 0; i < len(TLD); i++ {
		tldMatch := regexp.MustCompile(TLD[i])
		if tldMatch.MatchString(location.URL) {
			match = true
			break
		}
	}
	if !match {
		location.URL = location.URL + ".com" // figured defaulting to .com made the most sense.
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, location.URL, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	_, err = io.Copy(buf, res.Body)
	if err != nil {
		return "", fmt.Errorf("error when reading the response body: %w", err)
	}

	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		return "", fmt.Errorf("error unmarshalling(totally a word) the response: %w", err)
	}

	return payload.Data, nil
}
