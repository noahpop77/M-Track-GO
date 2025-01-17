package endpoints

import (
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

type MatchDataEntry struct {
        ChampionName string `json:"championName"`
        Kills int `json:"kills"`
        Deaths int `json:"deaths"`
        Assists int `json:"assists"`
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

var timeFormat = "2006-01-02"

func GetGameTime(durationInSeconds int) string {
        // Convert seconds to minutes and calculate remaining seconds
        return fmt.Sprintf("02d:%02d", durationInSeconds/60, durationInSeconds%60)
}

func UnixToDateString(epoch int64) string {
        // Convert the epoch time to a time.Time object
        time.Unix(epoch/1000, 0).Format(timeFormat)
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

        var response Response
        if err := json.Unmarshal(body, &response); err != nil {
                http.Error(writer, "failed to parse Jank ass JSON", http.StatusBadRequest)
                return
        }

        if len(response.Attachments) == 0 {
                http.Error(writer, "No Attachments Found", http.StatusBadRequest)
                return
        }

        gameData := &Response.Attachments[0].Content

        riotID := fmt.Sprintf("%s:%s", gameData.Info.Participants[0].RiotIDGameName, gameData.Info.Participants[0].RiotIDTagline)

        var queueType string
        if gameData.Info.QueueID == 420 {
                queueType = "Ranked Solo/Duo"
        }

        // pre-allocation
        matchaData := make([]map[string]interface{}, len(gameData.Info.Participants))
        for i, p := range gameData.Info.Participants {
        // You will want to decide what data you actually want to store here.
        matchData[i] = map[string]interface{}{
            "championName": p.ChampionName,
            "kills":       p.Kills,
            "deaths":      p.Deaths,
            "assists":     p.Assists,
            "goldEarned":  p.GoldEarned,
            "items": []string{
                p.Item0, p.Item1, p.Item2, p.Item3,
                p.Item4, p.Item5, p.Item6,
            },
            "summonerName":  p.SummonerName,
            "riotId":        fmt.Sprintf("%s#%s", p.RiotIDGameName, p.RiotIDTagline),
            "totalDamage":   p.TotalDamageDealtToChampions,
            "visionScore":   p.VisionScore,
            "win":          p.Win,
        }
    }

    const query = `
                INSERT INTO matchHistroy" (
                        "gameID", "gameVer", "riotID", "gameDurationMinutes", "gameCreationTimestamp",
                        "gameEndTimestamp, "queueType", "gameData", "participants", "matchData"
                ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

    queueType := "Unknown"
    if gameData.Info.QueueID == 420 {
        queueType = "Ranked Solo/Duo"
    }

    riotID := fmt.Sprintf("%s%s,
        gameData.Info.Participants[0].RiotIDGameName",
        gameData.Info.Participants[0].riotIDTagline,
    )

        ctx := context.Background()
    err = dbpool.QueryRow(ctx, query,
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

    if err != nil {
        http.Error(writer, "Database error", http.StatusInternalServerError)
        return
    }

        writer.WriteHeader(http.StatusOK)
        writer.Write([]byte("Inserted into database successfully"))
}
