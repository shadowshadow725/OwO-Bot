package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)
var token string
var apikey string
var helpmsg = "!owo search <jungler username> - try to predict enemy jungle\n!owo help - display this help message\n!owo update - update the riot api key (only works for user set )\n"
var apikey_updater string // discord user that's allowed to update the api key

var buffer = make([][]byte, 0)

func startStalker(username string, apikey string, s *discordgo.Session){
	for true{
		fmt.Printf("ran in parallel\n")
		result := stalk(username, apikey)
		pastresult := ""
		if result != ""{
			channel, err := s.UserChannelCreate(updater_id)
			if err != nil{
				fmt.Printf("Error creating message\n")
				os.Exit(2)
			}
			if result != pastresult{
				s.ChannelMessageSend(channel.ID, result)
				pastresult = result
				time.Sleep(5*time.Minute)

			}
		}else{
			time.Sleep(time.Minute)
		}



	}

}
func StartBot(discord_token string, key string, updater string, stalk_ign string){

	apikey_updater = updater
	token = discord_token
	apikey = key
	if token == "" {
		fmt.Println("No token provided. Exiting")
		return
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}
	// go startStalker(stalk_ign, key, dg) // disabled due to issues 
	// Register ready as a callback for the ready events.
	dg.AddHandler(ready)

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(messageCreate)


	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("OwO is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()

}
// This function will be called (due to AddHandler above) when the bot receives
// the "ready" event from Discord.
func ready(s *discordgo.Session, event *discordgo.Ready) {
	// Set the playing status.
	s.UpdateStatus(0, "!owo help")

}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.HasPrefix(m.Content, "!owo") {
		fmt.Printf("%s called the bot\n", m.Author.ID)
		// Find the channel that the message came from.
		_, err := s.State.Channel(m.ChannelID)
		if err != nil {
			// Could not find channel.
			return
		}
		if strings.HasPrefix(m.Content, "!owo search"){
			username := strings.Replace(m.Content, "!owo search ", "", 1)
			result := SearchPlayer(apikey, username, 20)
			s.ChannelMessageSend(m.ChannelID, result)
		}else if strings.HasPrefix(m.Content, "!owo update"){
			if m.Author.ID == apikey_updater{
				apikey = strings.Replace(m.Content, "!owo update ", "", 1)
				for strings.Contains(apikey, " "){
					apikey = strings.Replace(apikey, " ", "", 1)
				}
				s.ChannelMessageSend(m.ChannelID, "api key update success\n")
			}
		}else if strings.HasPrefix(m.Content, "!owo help"){
			s.ChannelMessageSend(m.ChannelID, helpmsg)
		}else{
			s.ChannelMessageSend(m.ChannelID, "IDK what you want me to do\n")
		}
	}

}
