package main

import (
	"flag"
)

var key string
var discord_token string
var updater_id string
var config_filename string
var stalk_ign string

func main (){
	flag.StringVar(&discord_token, "t", "", "Bot Token")
	flag.StringVar(&key, "k", "", "Riot API key")
	flag.StringVar(&updater_id, "u", "", "updater_id")
	flag.StringVar(&config_filename, "c", "", "config_filename")
	flag.Parse()
	var cfg Config

	if &config_filename != nil{
		ParseConfig(config_filename, &cfg)
		discord_token = cfg.Keys.Discord
		key = cfg.Keys.Riot
		updater_id = cfg.Users.DiscordID
		stalk_ign = cfg.Users.IGN
	}
	//fmt.Printf("%s\n%s\n%s\n%s\n", discord_token, key, updater_id, stalk_ign)

	StartBot(discord_token, key, updater_id, stalk_ign)

	//imgtest()
	//SearchPlayer(, "DISCORD TARZANED", 20)




	/*
	Keys:
	  Discord_key: "a"
	  Riot_key: "v"
	Users:
	  disocrd_user_id: "admin"
	  lol_ign: "super-pedro-1980"
	 */
}

