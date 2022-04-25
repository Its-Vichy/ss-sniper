package main

import (
	"fmt"
	"time"
)

func update() {
	init_db()
	for {
		ttl_server = 0
		ttl_sniper = len(snp_list)

		for _, acc := range snp_list {
			ttl_server += len(acc.session.State.Guilds)
		}

		slot_list = get_slot()
		log(fmt.Sprintf("Slot updated, %d customers are in waitlist", len(slot_list)))

		names := ""
		for _, slot := range slot_list {
			names += fmt.Sprintf("`%s (%d/%d)`, ", slot.AccSlot.Username, slot.AccSlot.Sniped, slot.AccSlot.Slot)
		}
		public_hook.send_hook(fmt.Sprintf("*Slot updated*, %d customers are in waitlist: %s. Sniping on *%d* guilds with *%d* accounts.", len(slot_list), names, ttl_server, ttl_sniper))
		time.Sleep(15 * time.Minute)
	}
}

func main() {
	public_hook, _ = new_hook("_MOSd2ARuO7nV3ImQqGiz7PeEW59s_dsxiXMThEsXeOoIxLMKuuPZyW5F38Ht9vDUnuA", "934188453715918978")
	hit_hook, _ = new_hook("JefSMZ9AD1MxA2qkXI6uPkuTrQPSftIMk-8Q0WeV99x_SGrunOiz0R97PBeFqtpwY_Dd", "934190675166109766")

	print_logo()
	blacklist_selfrep()
	blacklist_joined_servers()
	go update()
	go load_snipers()
	go join_work()
	block_console()
}
