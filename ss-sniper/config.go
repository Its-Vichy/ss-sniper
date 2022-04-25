package main

import "time"

var version = "0.0.5"

// Self replication configuration
var (
	max_join_day = 0
	delay_join   = 120 * time.Second
)

// Loader configuration
var (
	load_threads = 50
	server_min   = 0
)

// Temp storage
var (
	ttl_server = 0
	ttl_sniper = 0

	to_join_list             = []string{}
	joined                   = []string{}
	joined_id                = []string{}
	blacklist_selfrep_tokens = []string{}

	slot_list = []*ClaimerAccount{}
	snp_list  = []*SniperAccount{}

	blacklist_id = []string{}
	public_hook  = &hook{}
	hit_hook     = &hook{}
)

// Invite checker config
var (
	MinPercentage = 30
	MinMembers    = 1500
	ProxiesType   = "http"
	ProxiesPath   = "proxies.txt"

	proxies, err = readLines(ProxiesPath)
)
