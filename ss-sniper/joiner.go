package main

import (
	"strings"
	"time"
)

func join_work() {
	log("Start join work")
	time.Sleep(1 * time.Second)

	for {
		for _, acc := range snp_list {
			for _, invite := range to_join_list {
				code := strings.Split(invite, ":")[0]
				serv_id := strings.Split(invite, ":")[1]

				if !include(acc.token, blacklist_selfrep_tokens) && !include(code, joined) && !include(serv_id, joined_id) && time.Since(acc.lasted_join).Seconds() >= time.Duration(delay_join).Seconds() {
					if acc.invite_joined <= acc.invite_max && len(acc.session.State.Guilds) < 100 {
						acc.join_server(code)
					}
				}
			}
		}

		time.Sleep(5 * time.Second)
	}
}
