package main

import (
	"fmt"
	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/KnutZuidema/golio/riot"
	"github.com/sirupsen/logrus"
	"math"
	"strconv"
	"strings"
)

type Frame struct {
	invade bool
	location string
}
type Game struct {
	frames []Frame
	ever_invade bool
	team int
}
type Prediction struct{
	invade bool
	start_side string
	end_side string
	team int
}
func printmap(x int, y int){
	x = x/200
	y = y/200

	y = 80-y
	for i := 0; i<80;i++{
		for j := 0;j<80;j++{
			if (i == y && j == x){
				fmt.Printf("O ")
			}else {
				fmt.Printf("+ ")
			}
		}
		fmt.Printf("\n")
	}
}
func SearchPlayer(apikey string, playerign string, search_len int )string{

	return riotapi(apikey, playerign, search_len)

}
func locate_jungler(x int , y int , jgteam int )int {
	if jgteam == 100{
		x = x - 8500
		y = y - 8500
	}else {
		x = x - 7500
		y = y - 7500
	}
	new_x := math.Sin(math.Pi/4)*float64(y) + math.Cos(math.Pi/4)*float64(x)
	new_y := -math.Sin(math.Pi/4)*float64(x) + math.Cos(math.Pi/4)*float64(y)
	if new_x >= 0{
		if new_y >= 0 {
			return 1
		}else {
			return 4
		}
	}else {
		if new_y >= 0 {
			return 2
		}else {
			return 3
		}
	}
	return -1


}

func is_invade(location int, jgteam int)(bool, string){
	if jgteam == 100{
		if location == 2 {
			return false, "blue"
		} else if location == 3{
			return false,"red"
		}else if location == 1{
			return true, "red"
		}else if location == 4{
			return true, "blue"
		}
	}else if jgteam == 200{
		if location == 2 {
			return true, "blue"
		} else if location == 3{
			return true,"red"
		}else if location == 1{
			return false, "red"
		}else if location == 4{
			return false, "blue"
		}
	}
	return false, "nah"

}

func Predict(game Game)Prediction{
	predic := new(Prediction)
	predic.invade = game.ever_invade
	predic.team = game.team
	if game.frames[1].invade{
		predic.start_side = "invade " + game.frames[1].location
	}else {
		predic.start_side = "self " + game.frames[1].location
	}

	if game.frames[3].invade{
		predic.end_side= "invade " + game.frames[3].location
	}else {
		predic.end_side = "self " + game.frames[3].location
	}

	return *predic
}

func analyze(timeline riot.MatchTimeline, jgid int, jgteam int ) Prediction{
	x := 0
	y := 0
	ever_invade := false
	game := Game{}
	game.frames = []Frame{}
	game.team = jgteam
	for i := 1; i < 5; i++{
		for j := 1;j < len(timeline.Frames[0].ParticipantFrames); j++ {

			if timeline.Frames[i].ParticipantFrames[strconv.Itoa(j)].ParticipantID == jgid {
				x = timeline.Frames[i].ParticipantFrames[strconv.Itoa(j)].Position.X
				y = timeline.Frames[i].ParticipantFrames[strconv.Itoa(j)].Position.Y
				location := locate_jungler(x, y, jgteam)
				this_frame := new(Frame)
				this_frame.invade , this_frame.location = is_invade(location, jgteam)

				//fmt.Printf("cords (%d, %d) time = %d buff = %s invade? = %t\n", x, y, timeline.Frames[i].Timestamp, this_frame.location, this_frame.invade)
				//printmap(x,y)
				ever_invade = ever_invade || this_frame.invade
				game.frames = append(game.frames, *this_frame)

			}
		}
	}
	game.ever_invade = ever_invade
	return Predict(game)

}

func final_results(predictions []Prediction)string{
	invade_count := 0
	blue_start_count := 0
	blue_end_count := 0
	start_bot_side_count := 0
	end_bot_side_count := 0
	games_count := len(predictions)
	for i := 0;i<len(predictions);i++{
		if predictions[i].invade{
			invade_count++
		}
		if strings.Contains(predictions[i].start_side, "blue"){
			blue_start_count++
			if predictions[i].team == 200{
				start_bot_side_count++
			}
		}
		if strings.Contains(predictions[i].end_side, "blue"){
			blue_end_count++
			if predictions[i].team == 200{
				end_bot_side_count++
			}
		}
	}


	result := fmt.Sprintf("In the past %d games, the enemy jungler\nstarted on their blue buff first %d times\nstarted bot side %d times\n", games_count, blue_start_count, start_bot_side_count)
	result += fmt.Sprintf("ended on their blue %d times\nended bot side %d times\ninvaded the enemy jungle %d times\n", blue_end_count, end_bot_side_count, invade_count)
	result += fmt.Sprintf("Start Blue: %.2f%%\n", 100 * float64(blue_start_count)/float64(games_count))
	result += fmt.Sprintf("Start Bot: %.2f%%\n", 100 * float64(start_bot_side_count)/float64(games_count))
	result += fmt.Sprintf("End Bot: %.2f%%\n", 100 * float64(end_bot_side_count)/float64(games_count))
	result += fmt.Sprintf("End Blue: %.2f%%\n", 100 * float64(blue_end_count)/float64(games_count))
	result += fmt.Sprintf("Invade: %.2f%%\n", 100 * float64(invade_count)/float64(games_count))

	return result
}



func riotapi(apikey string, playerign string , search_len int )string {

	client := golio.NewClient(apikey, golio.WithRegion(api.RegionNorthAmerica), golio.WithLogger(logrus.New().WithField("foo", "bar")))
	summoner, _ := client.Riot.Summoner.GetByName(playerign)
	if summoner == nil{
		return "player not found"
	}
	//fmt.Printf("%s is a level %d summoner\n", summoner.Name, summoner.SummonerLevel)
	//fmt.Printf("champions id %d\n", games.Matches[0].Champion)
	games, _ := client.Riot.Match.List(summoner.AccountID, 0, search_len)
	predictions := []Prediction{}

	for imatch := 0; imatch < len(games.Matches);imatch++{

		//fmt.Printf("%d \n", games.Matches[imatch].GameID)

		match, _ := client.Riot.Match.Get(games.Matches[imatch].GameID)
		player_champion := games.Matches[imatch].Champion

		jgid := 0
		playerteam := 0
		for i := 0;i<len(match.Participants);i++{
			if match.Participants[i].ChampionID == player_champion {
				playerteam = match.Participants[i].TeamID
			}
		}
		jdindex := 0
		for i := 0;i<len(match.Participants);i++{
			if match.Participants[i].Spell1ID == 11 || match.Participants[i].Spell2ID == 11 {
				if playerteam == match.Participants[i].TeamID{
					jgid = match.Participants[i].ParticipantID
					jdindex = i
				}

			}
		}

		timeline, _ := client.Riot.Match.GetTimeline(games.Matches[imatch].GameID)
		jgteam := playerteam

		//fmt.Printf("%d \n", jgteam)
		if( match.QueueID == 420 || match.QueueID == 440 ) && match.Participants[jdindex].ChampionID == player_champion && match.GameDuration > 360{
			tmp := analyze(*timeline, jgid, jgteam)
			predictions = append(predictions, tmp)
		}


	}

	result := final_results(predictions)
	return result
	//fmt.Printf("%s\n", result)
}

