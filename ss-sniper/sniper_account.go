package main

// forked from vanshaj (rip<3) for bypass join

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/valyala/fasthttp"
)

type SniperAccount struct {
	session       *discordgo.Session
	base_request  fasthttp.Request
	token         string
	invite_joined int
	invite_max    int
	lasted_join   time.Time
}

func NewSniperAccount(token string, session *discordgo.Session) *SniperAccount {
	acc := SniperAccount{
		session:       session,
		base_request:  fasthttp.Request{},
		invite_max:    max_join_day,
		token:         token,
		invite_joined: 0,
		lasted_join:   time.Now(),
	}

	acc.initialize_header()
	acc.get_fingerprint()
	acc.get_cookies()
	go acc.reset_loop()

	return &acc
}

func (account *SniperAccount) reset_loop() {
	log(fmt.Sprintf("Start reset loop for %s", account.token))

	for {
		time.Sleep(24 * time.Hour)
		account.invite_joined = 0
		log(fmt.Sprintf("Reset invite delay of %s", account.token))
	}
}

func (account *SniperAccount) initialize_header() {
	account.base_request.Header.Set("accept", "*/*")
	account.base_request.Header.Set("Connection", "keep-alive")
	account.base_request.Header.Set("accept-encoding", "gzip, deflate, br")
	account.base_request.Header.Set("accept-language", "en-GB")
	account.base_request.Header.Set("content-type", "application/json")
	account.base_request.Header.Set("X-Debug-Options", "bugReporterEnabled")
	account.base_request.Header.Set("cache-control", "no-cache")
	account.base_request.Header.Set("sec-ch-ua", "'Chromium';v='92', ' Not A;Brand';v='99', 'Google Chrome';v='92'")
	account.base_request.Header.Set("sec-fetch-site", "same-origin")
	account.base_request.Header.Set("x-context-properties", "eyJsb2NhdGlvbiI6IkpvaW4gR3VpbGQiLCJsb2NhdGlvbl9ndWlsZF9pZCI6Ijg4NTkwNzE3MjMwNTgwOTUxOSIsImxvY2F0aW9uX2NoYW5uZWxfaWQiOiI4ODU5MDcxNzIzMDU4MDk1MjUiLCJsb2NhdGlvbl9jaGFubmVsX3R5cGUiOjB9")
	account.base_request.Header.Set("x-super-properties", "eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiRmlyZWZveCIsImRldmljZSI6IiIsInN5c3RlbV9sb2NhbGUiOiJlbi1VUyIsImJyb3dzZXJfdXNlcl9hZ2VudCI6Ik1vemlsbGEvNS4wIChXaW5kb3dzIE5UIDEwLjA7IFdpbjY0OyB4NjQ7IHJ2OjkzLjApIEdlY2tvLzIwMTAwMTAxIEZpcmVmb3gvOTMuMCIsImJyb3dzZXJfdmVyc2lvbiI6IjkzLjAiLCJvc192ZXJzaW9uIjoiMTAiLCJyZWZlcnJlciI6IiIsInJlZmVycmluZ19kb21haW4iOiIiLCJyZWZlcnJlcl9jdXJyZW50IjoiIiwicmVmZXJyaW5nX2RvbWFpbl9jdXJyZW50IjoiIiwicmVsZWFzZV9jaGFubmVsIjoic3RhYmxlIiwiY2xpZW50X2J1aWxkX251bWJlciI6MTAwODA0LCJjbGllbnRfZXZlbnRfc291cmNlIjpudWxsfQ==")
	account.base_request.Header.Set("sec-fetch-dest", "empty")
	account.base_request.Header.Set("sec-fetch-mode", "cors")
	account.base_request.Header.Set("sec-fetch-site", "same-origin")
	account.base_request.Header.Set("origin", "https://discord.com")
	account.base_request.Header.Set("referer", "https://discord.com/channels/@me")
	account.base_request.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) discord/0.0.16 Chrome/91.0.4472.164 Electron/13.4.0 Safari/537.36")
	account.base_request.Header.Set("te", "trailers")
	account.base_request.Header.Set("authorization", account.token)
}

func (account *SniperAccount) get_fingerprint() {
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

func (account *SniperAccount) get_cookies() {
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

func (account *SniperAccount) bypass_screen(server_id string, invite_code string) {
	account.get_cookies()

	req := account.base_request
	req.SetRequestURIBytes([]byte(fmt.Sprintf("https://discord.com/api/v9/guilds/%s/requests/@me", server_id)))
	req.Header.SetMethodBytes([]byte("PUT"))
	req.SetBody([]byte("{\"response\":true}"))

	res := fasthttp.AcquireResponse()
	fasthttp.Do(&req, res)

	fasthttp.ReleaseRequest(&req)
	fasthttp.ReleaseResponse(res)

	if res.StatusCode() == 201 || res.StatusCode() == 204 || res.StatusCode() == 200 {
		log(fmt.Sprintf("Bypassed server %s", server_id))
	} else {
		log(fmt.Sprintf("Failed to bypass server %s --> %d", server_id, res.StatusCode()))
	}
}

func (account *SniperAccount) bypass_react_captcha(server_id string) {
	channels, err := account.session.GuildChannels(server_id)

	if err != nil {
		fmt.Println(err)
	}

	if len(channels) > 10 {
		log(fmt.Sprintf("To many channels, skipping react bypass %s", server_id))
		return
	}

	for _, channel := range channels {
		messages, _ := account.session.ChannelMessages(channel.ID, 5, "", "", "")
		for _, message := range messages {
			for _, reaction := range message.Reactions {
				account.session.MessageReactionAdd(message.ChannelID, message.ID, reaction.Emoji.Name)
				log(fmt.Sprintf("Reacted to message %s on %s", message.ID, server_id))
			}

			if message.Author.ID == "536991182035746816" {
				log(fmt.Sprintf("Found wick on %s, bypassing react captcha", server_id))

				file, err := os.OpenFile("wick.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					return
				}

				_, err = file.WriteString(fmt.Sprintf("%s\n", server_id))
				if err != nil {
					return
				}
			}
		}
	}
}

func (account *SniperAccount) join_server(invite_code string) {
	log(invite_code)
	account.get_cookies()

	req := account.base_request
	req.Header.Set("content-length", "2")
	req.SetRequestURIBytes([]byte(fmt.Sprintf("https://discord.com/api/v9/invites/%s", invite_code)))
	req.Header.SetMethodBytes([]byte("POST"))
	req.SetBody([]byte("{}"))

	res := fasthttp.AcquireResponse()
	fasthttp.Do(&req, res)

	fasthttp.ReleaseRequest(&req)
	body := res.Body()

	fasthttp.ReleaseResponse(res)

	p, _ := DecodeBr(body)

	type guild struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	type joinresponse struct {
		VerificationForm bool  `json:"show_verification_form"`
		GuildObj         guild `json:"guild"`
	}

	var ResponseBody joinresponse
	json.Unmarshal(p, &ResponseBody)

	if res.StatusCode() == 200 && ResponseBody.GuildObj.ID != "" {
		joined_id = append(joined_id, ResponseBody.GuildObj.ID)
		joined = append(joined, invite_code)
		account.lasted_join = time.Now()
		account.invite_joined++

		log(fmt.Sprintf("Joined guild %s", ResponseBody.GuildObj.Name))
		if ResponseBody.VerificationForm {
			if len(ResponseBody.GuildObj.ID) != 0 {
				log(fmt.Sprintf("Bypassing server rules for %s", ResponseBody.GuildObj.Name))
				account.bypass_screen(ResponseBody.GuildObj.ID, invite_code)
			}
		}

		account.bypass_react_captcha(ResponseBody.GuildObj.ID)
	} else {
		log(fmt.Sprintf("Error join server: %d - %s - %s", res.StatusCode(), res.Body(), account.token))
		joined = append(joined, invite_code)
	}
}
