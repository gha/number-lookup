package main

import (
	"bufio"
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	lookupURL = "http://portal.aql.com/telecoms/network_lookup.php?number=%s&nlSubmit=submit"
)

var (
	numberFile = flag.String("input", "", "Location of number list")
	concurrent = flag.Int("concurrent", 0, "Enable concurrent mode")
	proxyAddr  = flag.String("proxy", "", "Proxy address to use")
)

func main() {
	flag.Parse()

	if len(*proxyAddr) != 0 {
		os.Setenv("HTTP_PROXY", *proxyAddr)
	}

	if len(*numberFile) == 0 {
		fmt.Println("Usage: numnberlookup -input /path/to/number/list.txt")
		return
	}

	if err := checkNumberFile(*numberFile); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Starting number lookup")

	file, err := os.Open(*numberFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	s := bufio.NewScanner(file)
	for s.Scan() {
		number := s.Text()
		number = strings.Replace(number, " ", "", -1)
		number = strings.Replace(number, "-", "", -1)
		if *concurrent > 0 {
			fmt.Println("Concurrent lookup not supported")
		} else {
			net, err := lookupNetwork(number)
			if err != nil && err != io.EOF {
				fmt.Println("Unable to lookup network for", number, ":", err)
				continue
			}

			fmt.Println(number, "-", net)
		}
	}

	if err := s.Err(); err != nil {
		fmt.Println(err)
	}
}

func checkNumberFile(file string) error {
	_, err := os.Stat(file)
	return err
}

func lookupNetwork(number string) (string, error) {
	URL := fmt.Sprintf(lookupURL, number)
	resp, err := http.Get(URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	z := html.NewTokenizer(resp.Body)
	var extractProp bool
	var extractData bool
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return "not found", z.Err()
		case html.StartTagToken:
			t := z.Token()
			if t.Data == "td" {
				extractProp = true
			}
		case html.EndTagToken:
			extractProp = false
		case html.TextToken:
			if extractProp {
				data := string(z.Text())

				if extractData {
					return data, nil
				}

				if data == "Network" {
					extractData = true
				}
			}
		}
	}

	return "", nil
}
