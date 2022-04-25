package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/andybalholm/brotli"
)

func load_tokens(file_path string) []string {
	file, _ := os.Open(file_path)
	scanner := bufio.NewScanner(file)
	var found_token []string

	for scanner.Scan() {
		for _, match := range regexp.MustCompile(`[\w-]{24}\.[\w-]{6}\.[\w-]{27}|mfa\.[\w-]{84}`).FindAllString(scanner.Text(), -1) {
			found_token = append(found_token, match)
		}
	}

	return found_token
}

func include(match string, list []string) bool {
	for _, item := range list {
		if match == item {
			return true
		}
	}

	return false
}

func DecodeBr(data []byte) ([]byte, error) {
	r := bytes.NewReader(data)
	br := brotli.NewReader(r)
	return ioutil.ReadAll(br)
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	defer file.Close()
	return lines, scanner.Err()
}

func blacklist_joined_servers() {
	invites, _ := readLines("code.txt")

	for _, invite := range invites {
		invite_code := strings.Split(invite, ":")[0]
		server_id := strings.Split(invite, ":")[1]

		joined_id = append(joined_id, server_id)
		joined = append(joined, invite_code)
	}

	log(fmt.Sprintf("Blacklisted %d servers id & invite", len(joined_id)))
}

func blacklist_selfrep() {
	tokens, _ := readLines("./blacklist_token_selfrep.txt")

	for _, token := range tokens {
		blacklist_selfrep_tokens = append(blacklist_selfrep_tokens, token)
	}
}

func get_proxy_http_client() *http.Client {
	url_i := url.URL{}
	url_proxy, _ := url_i.Parse(ProxiesType + "://" + proxies[rand.Intn(len(proxies))])

	transport := http.Transport{}
	transport.Proxy = http.ProxyURL(url_proxy)
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := http.Client{Transport: &transport}

	return &client
}