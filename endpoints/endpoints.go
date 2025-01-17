package endpoints

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Modify PrintJsonHandler to accept the dbpool as a parameter
func PrintJsonHandler(writer http.ResponseWriter, requester *http.Request) {
	if requester.Method != http.MethodPost {
		http.Error(writer, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the body of the request
	body, err := io.ReadAll(requester.Body)
	if err != nil {
		http.Error(writer, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer requester.Body.Close()

	// Print the body to the console
	fmt.Printf("Received JSON body: %s\n", string(body))

	// Respond with a 200 OK
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Received successfully"))
}

// Main JSON object
type GameData struct {
	Info     Info     `json:"info"`
	Metadata Metadata `json:"metadata"`
}

// Metadata object
type Metadata struct {
	MatchID      string   `json:"matchId"`
	Participants []string `json:"participants"`
}

// Info object
type Info struct {
	GameCreation       int64         `json:"gameCreation"`
	GameDuration       int           `json:"gameDuration"`
	GameEndTimestamp   int64         `json:"gameEndTimestamp"`
	GameStartTimestamp int64         `json:"gameStartTimestamp"`
	GameVersion        string        `json:"gameVersion"`
	GameID             int64         `json:"gameId"`
	QueueID            int64         `json:"queueId"`
	Participants       []Participant `json:"participants"`
}

// Participant object
type Participant struct {
	Assists                       int    `json:"assists"`
	ChampExperience               int    `json:"champExperience"`
	ChampLevel                    int    `json:"champLevel"`
	ChampionID                    int    `json:"championId"`
	ChampionName                  string `json:"championName"`
	Deaths                        int    `json:"deaths"`
	GoldEarned                    int    `json:"goldEarned"`
	Item0                         string `json:"item0"`
	Item1                         string `json:"item1"`
	Item2                         string `json:"item2"`
	Item3                         string `json:"item3"`
	Item4                         string `json:"item4"`
	Item5                         string `json:"item5"`
	Item6                         string `json:"item6"`
	Kills                         int    `json:"kills"`
	NeutralMinionsKilled          int    `json:"neutralMinionsKilled"`
	Perks                         Perks  `json:"perks"`
	RiotIDGameName                string `json:"riotIdGameName"`
	RiotIDTagline                 string `json:"riotIdTagline"`
	Summoner1ID                   string `json:"summoner1Id"`
	Summoner2ID                   string `json:"summoner2Id"`
	SummonerName                  string `json:"summonerName"`
	TeamID                        int    `json:"teamId"`
	TotalAllyJungleMinionsKilled  int    `json:"totalAllyJungleMinionsKilled"`
	TotalDamageDealtToChampions   int    `json:"totalDamageDealtToChampions"`
	TotalEnemyJungleMinionsKilled int    `json:"totalEnemyJungleMinionsKilled"`
	TotalMinionsKilled            int    `json:"totalMinionsKilled"`
	VisionScore                   int    `json:"visionScore"`
	Win                           bool   `json:"win"`
}

// Perks object
type Perks struct {
	Styles []Style `json:"styles"`
}

// Style object
type Style struct {
	Selections []Selection `json:"selections,omitempty"`
	Style      string      `json:"style,omitempty"`
}

// Selection object
type Selection struct {
	Perk string `json:"perk"`
}

func GetGameTime(durationInSeconds int) string {
	// Convert seconds to minutes and calculate remaining seconds
	minutes := durationInSeconds / 60
	seconds := durationInSeconds % 60

	// Format the output as "mm:ss"
	formattedTime := fmt.Sprintf("%02d:%02d", minutes, seconds)
	return formattedTime
}

func UnixToDateString(epoch int64) string {
	// Convert the epoch time to a time.Time object
	t := time.Unix(epoch/1000, 0)
	// Format the time object as "year-month-day"
	return t.Format("2006-01-02")
}

func InsertIntoDatabase(writer http.ResponseWriter, requester *http.Request, dbpool *pgxpool.Pool) {
	if requester.Method != http.MethodPost {
		http.Error(writer, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the body of the request
	body, err := io.ReadAll(requester.Body)
	if err != nil {
		http.Error(writer, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer requester.Body.Close()

	// Print the body to the console
	// fmt.Printf("Received JSON body: %s\n", body)

	//var gameData GameData
	var rawJSON string
	err = json.Unmarshal([]byte(body), &rawJSON)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	var formattedJSON bytes.Buffer
	err = json.Indent(&formattedJSON, []byte(rawJSON), "", "  ")
	if err != nil {
		fmt.Println("Error formatting JSON:", err)
		return
	}

	//var gameData GameData
	var gameData GameData
	err = json.Unmarshal(formattedJSON.Bytes(), &gameData)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	riotID := fmt.Sprintf("%s:%s", gameData.Info.Participants[0].RiotIDGameName, gameData.Info.Participants[0].RiotIDTagline)

	var queueType string
	if gameData.Info.QueueID == 420 {
		queueType = "Ranked Solo/Duo"
	}

	// jsonData, err := json.Marshal(gameData.Info.Participants)
	// if err != nil {
	// 	fmt.Println("Error marshalling struct:", err)
	// 	return
	// }
	// fmt.Println(string(jsonData))

	matchData := []map[string]string{}
	for index, value := range gameData.Info.Participants {
		matchData = append(matchData, map[string]string{
			fmt.Sprintf("%d", index): fmt.Sprintf("%v", value),
		})
	}

	// Example of database interaction (read-only query for simplicity)
	var value string
	err = dbpool.QueryRow(context.Background(),
		`INSERT INTO "matchHistory" ("gameID", "gameVer", "riotID", "gameDurationMinutes", "gameCreationTimestamp", "gameEndTimestamp", "queueType", "gameDate", "participants", "matchData")
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		gameData.Info.GameID,
		gameData.Info.GameVersion,
		riotID,
		GetGameTime(gameData.Info.GameDuration),
		gameData.Info.GameCreation,
		gameData.Info.GameEndTimestamp,
		queueType,
		UnixToDateString(gameData.Info.GameCreation),
		gameData.Metadata.Participants,
		matchData,
	).Scan(&value)

	fmt.Println(gameData.Info.GameID)
	fmt.Println(gameData.Info.GameVersion)
	fmt.Println(riotID)
	fmt.Println(GetGameTime(gameData.Info.GameDuration))
	fmt.Println(gameData.Info.GameCreation)
	fmt.Println(gameData.Info.GameEndTimestamp)
	fmt.Println(queueType)
	fmt.Println(UnixToDateString(gameData.Info.GameCreation))

	for index, value := range gameData.Metadata.Participants {
		fmt.Println(index, value)
	}

	jsonData, err := json.Marshal(gameData.Info.Participants)
	if err != nil {
		fmt.Println("Error marshalling struct:", err)
		return
	}
	fmt.Println(string(jsonData))
	//fmt.Println(matchData)

	if err != nil {
		http.Error(writer, "Database error", http.StatusInternalServerError)
		return
	}

	// Respond with a 200 OK
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Inserted into database successfully"))
}
