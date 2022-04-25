package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type InviteCode struct {
	Code      string      `json:"code"`
	Type      int         `json:"type"`
	ExpiresAt interface{} `json:"expires_at"`
	Guild     struct {
		ID                string      `json:"id"`
		Name              string      `json:"name"`
		Splash            string      `json:"splash"`
		Banner            string      `json:"banner"`
		Description       interface{} `json:"description"`
		Icon              string      `json:"icon"`
		Features          []string    `json:"features"`
		VerificationLevel int         `json:"verification_level"`
		VanityURLCode     string      `json:"vanity_url_code"`
		Nsfw              bool        `json:"nsfw"`
		NsfwLevel         int         `json:"nsfw_level"`
	} `json:"guild"`
	Channel struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Type int    `json:"type"`
	} `json:"channel"`
	ApproximateMemberCount   int `json:"approximate_member_count"`
	ApproximatePresenceCount int `json:"approximate_presence_count"`
}

func checkInvite(invite string) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://canary.discord.com/api/v6/invite/%s?with_counts=true", invite), nil)
	if err != nil {
		log(fmt.Sprintf("[ERROR] %s", err))
		return
	}

	client := get_proxy_http_client()
	resp, err := client.Do(req)
	if err != nil {
		log(fmt.Sprintf("[ERROR] %s", err))
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log(fmt.Sprintf("[ERROR] %s", err))
		return
	}

	var code = InviteCode{}
	err = json.Unmarshal(body, &code)
	if err != nil {
		log(fmt.Sprintf("Invalid invite: %s", invite))
		return
	}

	if include(code.Guild.ID, joined_id) || include(code.Code, joined) {
		log(fmt.Sprintf("Already joined %s", code.Code))
		return
	}

	percentage := float64(code.ApproximatePresenceCount) / float64(code.ApproximateMemberCount) * 100
	if code.ApproximateMemberCount < MinMembers {
		log(fmt.Sprintf("%s - Not enough members: %d/%d", code.Code, code.ApproximateMemberCount, MinMembers))
		return
	}
	if percentage < float64(MinPercentage) {
		log(fmt.Sprintf("%s - Not enough online members: %d/%d", code.Code, code.ApproximatePresenceCount, code.ApproximateMemberCount))
		return
	}

	log(fmt.Sprintf("%s - Online: %d/%d (%.2f%%)", code.Code, code.ApproximatePresenceCount, code.ApproximateMemberCount, percentage))
	to_join_list = append(to_join_list, fmt.Sprintf("%s:%s", code.Code, code.Guild.ID))

	file, err := os.OpenFile("code.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}

	_, err = file.WriteString(fmt.Sprintf("%s:%s\n", code.Code, code.Guild.ID))
	if err != nil {
		return
	} else {
	}

	file.Close()
}
