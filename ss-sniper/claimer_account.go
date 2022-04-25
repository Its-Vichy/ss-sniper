package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/valyala/fasthttp"
)

func (account *ClaimerAccount) initialize_header() {
	account.base_request.Header.Set("x-super-properties", "eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiRmlyZWZveCIsImRldmljZSI6IiIsInN5c3RlbV9sb2NhbGUiOiJlbi1VUyIsImJyb3dzZXJfdXNlcl9hZ2VudCI6Ik1vemlsbGEvNS4wIChXaW5kb3dzIE5UIDEwLjA7IFdpbjY0OyB4NjQ7IHJ2OjkzLjApIEdlY2tvLzIwMTAwMTAxIEZpcmVmb3gvOTMuMCIsImJyb3dzZXJfdmVyc2lvbiI6IjkzLjAiLCJvc192ZXJzaW9uIjoiMTAiLCJyZWZlcnJlciI6IiIsInJlZmVycmluZ19kb21haW4iOiIiLCJyZWZlcnJlcl9jdXJyZW50IjoiIiwicmVmZXJyaW5nX2RvbWFpbl9jdXJyZW50IjoiIiwicmVsZWFzZV9jaGFubmVsIjoic3RhYmxlIiwiY2xpZW50X2J1aWxkX251bWJlciI6MTAwODA0LCJjbGllbnRfZXZlbnRfc291cmNlIjpudWxsfQ==")
	account.base_request.Header.Set("sec-fetch-dest", "empty")
	account.base_request.Header.Set("sec-fetch-mode", "cors")
	account.base_request.Header.Set("sec-fetch-site", "same-origin")
	account.base_request.Header.Set("x-context-properties", "eyJsb2NhdGlvbiI6IkpvaW4gR3VpbGQiLCJsb2NhdGlvbl9ndWlsZF9pZCI6Ijg4NTkwNzE3MjMwNTgwOTUxOSIsImxvY2F0aW9uX2NoYW5uZWxfaWQiOiI4ODU5MDcxNzIzMDU4MDk1MjUiLCJsb2NhdGlvbl9jaGFubmVsX3R5cGUiOjB9")
	account.base_request.Header.Set("sec-ch-ua", "'Chromium';v='92', ' Not A;Brand';v='99', 'Google Chrome';v='92'")
	account.base_request.Header.Set("accept", "*/*")
	account.base_request.Header.Set("accept-language", "en-GB")
	account.base_request.Header.Set("content-type", "application/json")
	account.base_request.Header.Set("sec-ch-ua-mobile", "?0")
	account.base_request.Header.Set("sec-ch-ua-platform", "\"Windows\"")
	account.base_request.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) discord/0.0.16 Chrome/91.0.4472.164 Electron/13.4.0 Safari/537.36")
	account.base_request.Header.Set("authorization", account.AccSlot.ClaimToken)
}

func (account *ClaimerAccount) get_fingerprint() {
	client := get_proxy_http_client()
	resp, _ := client.Get("https://discordapp.com/api/v9/experiments")
	body, _ := ioutil.ReadAll(resp.Body)

	type Fingerprintx struct {
		Fingerprint string `json:"fingerprint"`
	}

	var fingerprinty Fingerprintx
	json.Unmarshal(body, &fingerprinty)

	account.base_request.Header.Set("x-fingerprint", fingerprinty.Fingerprint)
}

func (account *ClaimerAccount) get_cookies() {
	client := get_proxy_http_client()
	resp, _ := client.Get("https://discord.com")

	type cookie struct {
		dcfduid  string
		sdcfduid string
	}

	c := cookie{}
	if resp.Cookies() != nil {
		for _, cookie := range resp.Cookies() {
			if cookie.Name == "__dcfduid" {
				c.dcfduid = cookie.Value
			}
			if cookie.Name == "__sdcfduid" {
				c.sdcfduid = cookie.Value
			}
		}
	}

	account.base_request.Header.Set("cookie", fmt.Sprintf("__dcfduid=%s; __sdcfduid=%s; locale=us", c.dcfduid, c.sdcfduid))
}

func (account *ClaimerAccount) get_payement_source_id() {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethodBytes([]byte("GET"))
	req.SetRequestURIBytes([]byte("https://discord.com/api/v8/users/@me/billing/payment-sources"))
	res := fasthttp.AcquireResponse()

	fasthttp.Do(req, res)
	fasthttp.ReleaseRequest(req)
	id := regexp.MustCompile(`("id": ")([0-9]+)"`).FindStringSubmatch(string(res.Body()))

	if id == nil {
		account.payement_method = "null"
	} else if len(id) > 1 {
		account.payement_method = id[2]
	}
}

func (account *ClaimerAccount) claim_nitro(nitro_code string) (int, string) {
	req := account.base_request
	req.SetRequestURIBytes([]byte("https://discordapp.com/api/v8/entitlements/gift-codes/" + nitro_code + "/redeem"))
	req.Header.SetMethodBytes([]byte("POST"))
	log(account.payement_method)
	req.SetBody([]byte(`{"channel_id":` + "null" + `,"payment_source_id": ` + account.payement_method + `}`))

	res := fasthttp.AcquireResponse()
	fasthttp.Do(&req, res)

	fasthttp.ReleaseRequest(&req)
	body := res.Body()

	bodyString := string(body)
	fasthttp.ReleaseResponse(res)

	log(string(req.Body()))
	
	if strings.Contains(bodyString, "Unknown Gift Code") {
		return 0, "Unknown Gift Code"
	} else if strings.Contains(bodyString, "The resource is being rate limited.") {
		return 0, "The resource is being rate limited."
	} else if strings.Contains(bodyString, "This gift has been redeemed already.") {
		return 0, "This gift has been redeemed already."
	} else {
		if strings.Contains(bodyString, "gifter_user_id") {
			return 1, "Ok"
		} else {
			go log(bodyString)
			return 0, bodyString
		}
	}
}