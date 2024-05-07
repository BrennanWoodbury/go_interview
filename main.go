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
	"strings"
)

var buf *bytes.Buffer

func main() {
	data, err := os.ReadFile("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		data, err := getData(line)
		if err != nil {
			log.Printf("unable to get status: %v", err)
		}
		if data.(string) == "foo" {
			fmt.Printf("data found: %s", data)
			os.Exit(0)
		}
	}
}

func getData(line string) (any, error) {
	var location struct{ URL string }
	err := json.Unmarshal([]byte(line), &location)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, location.URL, nil)
	if err != nil {
		return nil, err
	}
	c := &http.Client{}
	res, err := c.Do(req)
	fmt.Println(res)
	if err != nil {
		return nil, err
	}
	data, err := io.Copy(buf, res.Body)
	fmt.Println(data)
	if err != nil {
		return nil, err
	}
	var payload struct {
		Data string `bson:"data"`
	}
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		return nil, err
	}
	return payload.Data, nil
}
