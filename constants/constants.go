package constants

const(
	//CREATE
	insertLocation string = "INSERT INTO locations (location_name) VALUES (:location_name)"


	insertNewPlayer string = "INSERT INTO players (player_name, is_active, left_at, location_id, team_id)" +
	"VALUES(:player_name, :is_active, :left_at, :location_id, :team_id)"

	//READ


	//UPDATE


	//DESTROY
)