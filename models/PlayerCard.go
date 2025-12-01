package models

type PlayerCard struct{
	Score int    `bson:"score"` 
	Name  string `bson:"name"`
}