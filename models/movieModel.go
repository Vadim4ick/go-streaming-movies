package model

import "go.mongodb.org/mongo-driver/v2/bson"

type Genre struct { 
	GenreID int `json:"genre_id" bson:"genre_id" validate:"required"`
	GenreName string `json:"genre_name" bson:"genre_name" validate:"required, min=2, max=100"`
}

type Ranking struct {
	RankingValue int `json:"ranking_value" bson:"ranking_value" validate:"required"`
	RankingName string `json:"ranking_name" bson:"ranking_name" validate:"required"`
}

type Movie struct {
    ID bson.ObjectID `json:"_id" bson:"_id"`
	ImdbID string `json:"imdb_id" bson:"imdb_id" validate:"required"`
	Title string `json:"title" bson:"title" validate:"required, min=2, max=500"`
	PosterPath string `json:"poster_path" bson:"poster_path" validate:"required, url"`
	YoutubeID string `json:"youtube_id" bson:"youtube_id" validate:"required"`
	Genre []Genre `json:"genre" bson:"genre" validate:"required, dive"`
	AdminReview string `json:"admin_review" bson:"admin_review" validate:"required"`
	Ranking Ranking `json:"ranking" bson:"ranking" validate:"required"`
}
