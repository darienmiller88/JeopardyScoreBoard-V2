package controllers

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"JeopardyScoreBoardV2/repositories"
	"JeopardyScoreBoardV2/services"
	"JeopardyScoreBoardV2/encryption"
)

type Index struct{
	Router *chi.Mux
	viewsController ViewsController
	locationsController LocationsController
	playersController PlayersController
	savedGamesController SavedGamesController
}

func (i *Index) InitControllers(db *sqlx.DB, encryptionService *encryption.EncryptionService){
	i.Router = chi.NewRouter()

	//Initialize the views controller, and inject the following services
	i.viewsController.Init(
		&services.LocationServiceImpl{ Repository: repositories.GetSqlLocationRepository(db, encryptionService) },
		&services.PlayerServiceImpl{ 
			PlayerRepository: repositories.GetSqlPlayerRepository(db, encryptionService),
			LocationRepository: repositories.GetSqlLocationRepository(db, encryptionService),
		},
	)

	//Initialize the controllers, and inject the service and repo implementation
	i.locationsController.Init(&services.LocationServiceImpl{ Repository: repositories.GetSqlLocationRepository(db, encryptionService) })
	i.playersController.Init(&services.PlayerServiceImpl{ 
		PlayerRepository: repositories.GetSqlPlayerRepository(db, encryptionService),
		LocationRepository: repositories.GetSqlLocationRepository(db, encryptionService),
	})
	i.savedGamesController.Init(&services.SaveGameServiceImpl{ 
		SavedGameRepository: repositories.GetSqlSavedGameRepository(db, encryptionService),
		LocationRepository: repositories.GetSqlLocationRepository(db, encryptionService),
		TeamRepository: repositories.GetSqlTeamRepository(db, encryptionService),
		PlayerRepository: repositories.GetSqlPlayerRepository(db, encryptionService),
		EncryptionService: encryptionService,
	})

	//Afterwards, mount the views router onto this router, which wiil be mounted onto the main chi router
	//in main.go
	i.Router.Mount("/", i.viewsController.Router)
	i.Router.Mount("/locations", i.locationsController.Router)
	i.Router.Mount("/players", i.playersController.Router)
	i.Router.Mount("/savedgames", i.savedGamesController.Router)
}