package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	"JeopardyScoreBoardV2/controllers"
	"JeopardyScoreBoardV2/database"
	"JeopardyScoreBoardV2/encryption"
)

//push from work laptop. changed my author email
//push from old hp, oligatory push of thr day pt 3 lol
func main() {
	//Load env file immediately at the start of the program
	godotenv.Load()

	//Create new chi router instance to push handlers to.
	router := chi.NewRouter()

	//Middleware stack, keeping it basic for now.
	router.Use(middleware.Logger, middleware.Recoverer)

	//Initiate the database connection to SQL, and defer its disconnection.
	database.Init()
	defer database.CloseSQLDB()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password
		DB:       0,  // use default DB
		Protocol: 2,
	})

	ctx := context.Background()



	key := os.Getenv("ENCRYPTION_KEY")
	keyB64, err := base64.StdEncoding.DecodeString(key)


	if err != nil {
		log.Fatal(err)
	}

	//Initialize encryption service
	encryptionService := encryption.NewService(keyB64)

	//Initialize the parent controller router, and its children
	index := controllers.Index{}
	index.InitControllers(database.GetDB(), encryptionService)

	// encryptNames(database.GetDB(), encryptionService)
	// addPlayers(database.GetDB(), encryptionService)
	// addSavedGames(database.GetDB(), encryptionService)

	//Afterwards, mount that router onto this one.
	router.Mount("/", index.Router)

	//Serve static files along the "/static" route
	fs := http.FileServer(http.Dir("static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	//Finally, listen and serve on the port in the env, which is 8080 on local machine.
	fmt.Println("Listening on Port:", os.Getenv("PORT"))
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router)
}

func addSavedGames(db *sqlx.DB, es *encryption.EncryptionService){
	data, err := os.ReadFile("savedgames.json")

	if err != nil { 
		panic(err)
	}

	games := []struct{
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
		// PlayerName string `json:"player"`
		Location   string `json:"location_name"`
		TotalPoints int   `json:"total_points"`
		AveragePoints float64 `json:"average_points"`
		Winner      struct{
			Score int `json:"score"`
			Username string `json:"username"`
		}   `json:"winner"`
		Players    []struct{
			PlayerName string `json:"username"`
			Score      int    `json:"score"`
		} `json:"players"`
	}{}

	if err := json.Unmarshal(data, &games); err != nil{
		panic(err)
	}

	for _, game := range games {
		encrypted, _ := es.Encrypt(game.Winner.Username)	
		hash := encryption.NameHash(game.Winner.Username)

		var locationID int
		err := db.QueryRow(
			"SELECT id FROM locations WHERE location_name=$1",
			game.Location,
		).Scan(&locationID)

		if err != nil {
			panic(fmt.Errorf("location not found: %s", game.Location))
		}

		var playerID int
		err = db.QueryRow(
			"SELECT id FROM players WHERE player_name_hash=$1 AND location_id=$2",
			hash,
			locationID,
		).Scan(&playerID)

		if err != nil {
			panic(fmt.Errorf("player not found: %s at location %s", game.Winner.Username, game.Location))
		}

		savedGameId := 0
		err = db.QueryRow(`INSERT INTO savedgames 
			(
			created_at,
			updated_at,
			location_id, 
			winning_player_id, 
			winning_player_name_encrypted, 
			winning_player_name_hash, 
			total_score, 
			average_score
			)
			VALUES(
				$1,
				$2,
				$3,
				$4,
				$5,
				$6,
				$7,
				$8
			) RETURNING id
			`,
			game.CreatedAt,
			game.UpdatedAt,
			locationID,
			playerID,
			encrypted,
			hash,
			game.TotalPoints,
			game.AveragePoints,
		).Scan(&savedGameId)

		if err != nil {
			panic(fmt.Errorf("saved game err:%v", err))
		}

		for _, player := range game.Players {
			encrypted, _ = es.Encrypt(player.PlayerName)	
			hash = encryption.NameHash(player.PlayerName)

			var playerID int
			err := db.QueryRow(
				"SELECT id FROM players WHERE player_name_hash=$1",
				hash,
			).Scan(&playerID)

			if err != nil {
				fmt.Println("player:", player.PlayerName, "doesnt exist for some reason at location", locationID)
				panic(fmt.Errorf("player err:%v", err))
			}

			_, err = db.Exec(
				`
					INSERT INTO savedgamesplayers
					(
						created_at,
						updated_at,
						player_id,
						saved_game_id,
						player_score,
						player_name_encrypted,
						player_name_hash
					)
					VALUES
					(
						$1,
						$2,
						$3,
						$4,
						$5,
						$6,
						$7
					)
				`,
				game.CreatedAt,
				game.UpdatedAt,
				playerID,
				savedGameId,
				player.Score,
				encrypted,
				hash,
			)

			if err != nil {
				panic(fmt.Errorf("player insertion err:%v", err))
			}
		}

	}

	fmt.Println("All games added successfully")
}

func addPlayers(db *sqlx.DB, es *encryption.EncryptionService){
	data, err := os.ReadFile("players.json")

	if err != nil { 
		panic(err)
	}

	players := []struct{
		PlayerName string `json:"player"`
		Location   string `json:"location"`
	}{}

	if err := json.Unmarshal(data, &players); err != nil{
		panic(err)
	}

	fmt.Println("players:", len(players))

	for _, player := range players {
		encrypted, _ := es.Encrypt(player.PlayerName)	
		hash := encryption.NameHash(player.PlayerName)


		_, err := db.Exec(`INSERT INTO players 
			(location_id, player_name_encrypted, player_name_hash)
			VALUES(
				(SELECT id FROM locations WHERE location_name=$1),
				$2,
				$3
			)`,
			player.Location,
			encrypted,
			hash,
		)

		if err != nil {
			panic(err)
		}
	}

}

func encryptNames(db *sqlx.DB, es *encryption.EncryptionService){
	rows, _ := db.Query(`
		SELECT id, player_name FROM players
	`)

	for rows.Next() {
		var id int
		var name string

		rows.Scan(&id, &name)

		if id == 10 || id == 12 {
			fmt.Println("encrypting", name)

			encrypted, _ := es.Encrypt(name)
			hash := encryption.NameHash(name)
	
			db.Exec(`
				UPDATE savedgames
				SET winning_player_name_encrypted = $1,
					winning_player_name_hash = $2
				WHERE winning_player_id = $3
			`, encrypted, hash, id)
		}
	}

}