package controllers

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"JeopardyScoreBoardV2/repositories"
	"JeopardyScoreBoardV2/services"
	"JeopardyScoreBoardV2/encryption"
)

type IndexController struct{
	Router *chi.Mux
}

func NewIndexController(db *sqlx.DB, encryptionService *encryption.EncryptionService) *IndexController {
	i := &IndexController{
		Router: chi.NewRouter(),
	}

	i.registerRoutes(db, encryptionService)

	return i
}

func (i *IndexController) registerRoutes(db *sqlx.DB, es *encryption.EncryptionService) {
	sgs := services.NewSaveGameService(
			repositories.NewSavedGameRepository(db, es),
			repositories.NewLocationRepository(db, es),
			repositories.NewPlayerRepository(db, es),
			repositories.NewTeamRepository(db, es),
			es,
		)

	vc := NewViewsController(
		services.NewLocationService(repositories.NewLocationRepository(db, es)),
		services.NewPlayerService(repositories.NewPlayerRepository(db, es)),
		sgs,
		services.NewTeamService(repositories.NewTeamRepository(db, es)),
	)
	pc := NewPlayersController(services.NewPlayerService(repositories.NewPlayerRepository(db, es)))
	lc := NewLocationsController(services.NewLocationService(repositories.NewLocationRepository(db, es)))
	sgc := NewSavedGamesController(sgs)
	tc := NewTeamsController(services.NewTeamService(repositories.NewTeamRepository(db, es)))

	//Afterwards, mount the views router onto this router, which wiil be mounted onto the main chi router
	//in main.go
	i.Router.Mount("/", vc.Router)
	i.Router.Mount("/locations", lc.Router)
	i.Router.Mount("/players", pc.Router)
	i.Router.Mount("/savedgames", sgc.Router)
	i.Router.Mount("/teams", tc.Router)
}