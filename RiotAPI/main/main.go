package main

import "flag"

var key string
var discord_token string


func main (){
	flag.StringVar(&discord_token, "t", "", "Bot Token")
	flag.StringVar(&key, "k", "", "Riot API key")

	flag.Parse()
	StartBot(discord_token, key)

	//imgtest()
	//SearchPlayer(, "DISCORD TARZANED", 20)

}

