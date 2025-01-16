package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
	fmt.Printf("Received JSON body: %s\n", string(body))

	// Parse the JSON string into a map
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("Error parsing JSON:", err)
		http.Error(writer, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	// Extract values from the parsed data map
	info, ok := data["info"].(map[string]interface{})
	if !ok {
		http.Error(writer, "Missing 'info' field", http.StatusBadRequest)
		return
	}

	// TODO: Create a full structure body for what is in the payloads
	// The following is an EXAMPLE and not to be used for realsies
	// Example:
	/*
			// Participant struct for each player in the participants array
		type Participant struct {
		    PlayerID int    `json:"playerID"`
		    Name     string `json:"name"`
		    Kills    int    `json:"kills"`
		    Deaths   int    `json:"deaths"`
		}

		// Info struct for the main game info
		type Info struct {
		    GameID        int            `json:"gameID"`
		    GameVersion   string         `json:"gameVersion"`
		    GameDuration  int            `json:"gameDuration"`
		    Participants  []Participant  `json:"participants"`
		}

		// Metadata struct for the metadata section
		type Metadata struct {
		    MatchID   string `json:"matchId"`
		    Timestamp int    `json:"timestamp"`
		}

		// GameData struct for the full JSON structure
		type GameData struct {
		    Info     Info     `json:"info"`
		    Metadata Metadata `json:"metadata"`
		}
	*/
	// Temporary dog water variables
	gameID := info["gameId"].(string)
	gameVer := info["gameVersion"].(string)
	riotID := fmt.Sprintf("%s:%s", info["riotIdGameName"].(string), info["riotIdTagline"].(string))
	gameDurationMinutes := fmt.Sprintf("%.0f", info["gameDuration"].(float64))
	gameCreationTimestamp := info["gameCreation"].(string)
	gameEndTimestamp := info["gameEndTimestamp"].(string)
	queueType := "bob"
	gameDate := "bob"
	participants := info["gameEndTimestamp"].(string)
	matchData := "bob"

	// Example of database interaction (read-only query for simplicity)
	var value string
	err = dbpool.QueryRow(context.Background(),
		`INSERT INTO "matchHistory" ("gameID", "gameVer", "riotID", "gameDurationMinutes", "gameCreationTimestamp", "gameEndTimestamp", "queueType", "gameDate", "participants", "matchData")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		gameID, gameVer, riotID, gameDurationMinutes, gameCreationTimestamp, gameEndTimestamp, queueType, gameDate, participants, matchData,
	).Scan(&value)

	if err != nil {
		http.Error(writer, "Database error", http.StatusInternalServerError)
		return
	}

	// Respond with a 200 OK
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Inserted into database successfully"))
}
