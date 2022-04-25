package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zenthangplus/goccm"
)

var (
	regex_nitro  = regexp.MustCompile("(discord.com/gifts/|discordapp.com/gifts/|discord.gift/)([a-zA-Z0-9]+)")
	regex_invite = regexp.MustCompile("(discord.gg/)([0-9a-zA-Z]+)")
	blacklist    []string
)

func message_handler(_ *discordgo.Session, message *discordgo.MessageCreate) {
	go func() {
		for _, match := range regex_nitro.FindAllString(message.Content, -1) {
			go func(match string) {
				code := regex_nitro.FindStringSubmatch(match)[2]

				if len(code) >= 16 && len(code) <= 24 && !include(code, blacklist) && !include(message.Author.ID, blacklist_id) {
					log(fmt.Sprintf("Detected code: %s", code))
					blacklist = append(blacklist, code)

					if len(slot_list) > 0 {
						to_snipe := slot_list[len(slot_list)-1]
						res, resp := to_snipe.claim_nitro(code)

						if res == 1 {
							message := fmt.Sprintf("Claimed nitro code `%s` for `%s`", code, to_snipe.AccSlot.Username)
							to_snipe.webhook.send_hook(message)
							hit_hook.send_hook(message)
							to_snipe.AccSlot.Sniped++
							to_snipe.update_sniped()

							temp_slot := []*ClaimerAccount{}
							for _, slot := range slot_list {
								if slot.AccSlot.ClaimToken == to_snipe.AccSlot.ClaimToken {
									if slot.AccSlot.Sniped < slot.AccSlot.Slot {
										temp_slot = append(temp_slot, slot)
									} else {
										log("All nitro redeemed")
									}
								} else {
									temp_slot = append(temp_slot, slot)
								}
							}
							slot_list = temp_slot

						} else {
							public_hook.send_hook(fmt.Sprintf("Failed to claim nitro code `%s`. Code: *%s*", code, resp))
							blacklist_id = append(blacklist_id, message.Author.ID)
							log(fmt.Sprintf("Blacklised user: %s", message.Author.ID))
						}
					}
				}
			}(match)
		}
	}()

	go func() {
		for _, match := range regex_invite.FindAllString(message.Content, -1) {
			invite_code := strings.Split(match, "/")[1]

			if !include(invite_code, blacklist) {
				blacklist = append(blacklist, invite_code)

				f, _ := os.OpenFile("invites.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				f.WriteString(fmt.Sprintf("%s\n", invite_code))
				checkInvite(invite_code)
			}
		}
	}()
}

func run_sniper(token string) {
	client, err := discordgo.New(token)
	client.LogLevel = discordgo.LogError

	if err == nil {
		if client.Open() == nil {
			client_servers := len(client.State.Guilds)
			//log(fmt.Sprintf("Connected to %s", client.State.User.Username))

			if client_servers >= server_min {
				ttl_server += client_servers
				ttl_sniper++

				acc := NewSniperAccount(token, client)
				snp_list = append(snp_list, acc)

				for _, guild := range client.State.Guilds {
					joined_id = append(joined_id, guild.ID)
				}

				if client_servers >= 40 {
					blacklist_selfrep_tokens = append(blacklist_selfrep_tokens, token)
				}

				client.AddHandler(message_handler)
			} else {
				client.Close()
			}
		}
	}
}

func load_snipers() {
	c := goccm.New(load_threads)

	for _, token := range load_tokens("./tokens.txt") {
		c.Wait()

		go func(token string) {
			run_sniper(token)
			c.Done()
		}(token)
	}

	c.WaitAllDone()
	log("Finished to load snipers")
	public_hook.send_hook(fmt.Sprintf("Sniper load -> Sniping on *%d* guilds with *%d* accounts.", ttl_server, ttl_sniper))
}
