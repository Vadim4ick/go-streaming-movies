package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"myapp/database"
	"myapp/models"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms/openai"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var movieCollection *mongo.Collection = database.OpenCollection("movies")
var rankingCollection *mongo.Collection = database.OpenCollection("rankings")

var validate = validator.New()

func GetMovies()  gin.HandlerFunc{
	return func (c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var movies[] models.Movie

		cursor, err := movieCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, movies)
	}
}

func GetMovie()  gin.HandlerFunc{
	return func (c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		
		defer cancel()

		var movie models.Movie

		movieID := c.Param("movie_id")

		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
			return
		}

		fmt.Print(movieID)
 
		err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, movie)
	}
}

func AddMovie() gin.HandlerFunc {
	return func (c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var movie models.Movie

		if err := c.ShouldBindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err:=validate.Struct(movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error", "details": err.Error()})
			return
		}


		result, err := movieCollection.InsertOne(ctx, movie)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add movie"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func AdminReviewUpdate() gin.HandlerFunc {
 	return func (c *gin.Context) {
		movieId := c.Param("movie_id")

		if( movieId == "") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
			return
		}

        var req struct {
			AdminReview string `json:"admin_review"`
		}


		var resp struct {
			RankingName string `json:"ranking_name"`
			AdminReview string `json:"admin_review"`
		}

		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}


		sentiment, rankVal, err := GetReviewRankings(req.AdminReview)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		filter := bson.M{"imdb_id": movieId}

		update := bson.M{
			"$set": bson.M{
				"admin_review": req.AdminReview,
				"ranking": bson.M{
					"ranking_name": sentiment,
					"ranking_value": rankVal,
				},
			},
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := movieCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie review"})
			return
		}


		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		resp.RankingName = sentiment
		resp.AdminReview = req.AdminReview

		c.JSON(http.StatusOK, resp)
	}
}

func GetReviewRankings(admin_review string) (string, int, error){
	rankings, err := GetRankings()

	if err != nil {
		return "", 0, err
	}

	sentimentDelimited := ""

	for _, ranking := range rankings {
		if ranking.RankingValue != 999 {
			sentimentDelimited += ranking.RankingName + ", "
		}
	}

	sentimentDelimited = strings.Trim(sentimentDelimited, ",")

	err = godotenv.Load()

	if err != nil {
		log.Println("Warning: .env file not found")
	}

	OPENAI_API_KEY := os.Getenv("OPENAI_API_KEY")

	if OPENAI_API_KEY == "" {
		return "", 0, errors.New("could not read OPENAI_API_KEY")
	}

	llm, err := openai.New(openai.WithToken(OPENAI_API_KEY))

	if err != nil {
		return "", 0, err
	}

	promptTemplate := os.Getenv("BASE_PROMPT_TEMPLATE")

	if promptTemplate == "" {
		return "", 0, errors.New("could not read BASE_PROMPT_TEMPLATE")
	}

	basePrompt := strings.Replace(promptTemplate, "{rankings}", sentimentDelimited, 1)

	response, err := llm.Call(context.Background(), basePrompt + admin_review)

	if err != nil {
		return "", 0, err
	}

	rankVal := 0

	for _, ranking := range rankings {
		if ranking.RankingName == response {
			rankVal = ranking.RankingValue
			break
		}
	}

	return  response, rankVal, nil
}

func GetRankings() ([]models.Ranking, error) {
  	var rankings []models.Ranking

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	cursor, err := rankingCollection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &rankings); err != nil {
		return nil, err
	}

	return rankings, nil
}