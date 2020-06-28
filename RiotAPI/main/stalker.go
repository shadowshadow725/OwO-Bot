package main

import (
	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/sirupsen/logrus"
)

func stalk (username string, apikey string ) string {
	//username string
	//username := "Juke Jouster"
	//apikey := "RGAPI-cb9ee4d3-f44e-4aa4-abe2-3b1c38e046a7"
	client := golio.NewClient(apikey, golio.WithRegion(api.RegionNorthAmerica), golio.WithLogger(logrus.New().WithField("foo", "bar")))
	summoner, _ := client.Riot.Summoner.GetByName(username)
	gameinfo, err := client.Riot.Spectator.GetCurrent(summoner.ID)
	if err != nil{
		return ""
	}
	if err == nil{
		playerteam := 0
		for i := 0;i<len(gameinfo.Participants);i++{
			if username == gameinfo.Participants[i].SummonerName {
				playerteam = gameinfo.Participants[i].TeamID
			}
		}

		for i := 0;i<len(gameinfo.Participants);i++{
			if gameinfo.Participants[i].Spell1ID == 11 || gameinfo.Participants[i].Spell2ID == 11 {
				if playerteam != gameinfo.Participants[i].TeamID{
					result := SearchPlayer(apikey, gameinfo.Participants[i].SummonerName, 20)
					//fmt.Printf("%s\n%s",gameinfo.Participants[i].SummonerName,  result)
					return result
				}


			}
		}
	}
	return ""

}