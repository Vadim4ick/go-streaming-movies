package controllers

import (
	"context"
	"fmt"
	"myapp/database"
	"myapp/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var movieCollection *mongo.Collection = database.OpenCollection("movies")

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