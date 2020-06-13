package main

import "flag"

var key string
var discord_token string
var updater_id string

func main (){
	flag.StringVar(&discord_token, "t", "", "Bot Token")
	flag.StringVar(&key, "k", "", "Riot API key")
	flag.StringVar(&updater_id, "u", "", "updater_id")
	flag.Parse()
	StartBot(discord_token, key, updater_id)

	//imgtest()
	//SearchPlayer(, "DISCORD TARZANED", 20)

}

