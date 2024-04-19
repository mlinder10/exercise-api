package main

import (
	"exercise-api/src"
	"os"

	"github.com/gin-gonic/gin"
)

const PATH = "exercises.json"
const PORT = "8080"

func main() {
	exercises, err := src.ReadExercises(PATH)
	if err != nil {
		panic(err)
	}
	// Create a new Gin router
	router := gin.Default()

	// Middleware
	router.Use(src.AuthMiddleware())
	router.Use(src.DataMiddleware(exercises))

	// Routes
	router.GET("/images/:name", src.ImageHandler)
	router.GET("/ping", src.Ping)
	router.GET("/", src.GetAllExercises)
	router.GET("/exercise/:id", src.GetExerciseByID)
	router.GET("/filter", src.GetFilteredExercises)
	router.GET("/body_parts", src.GetBodyParts)
	router.GET("/exercise_types", src.GetExerciseTypes)
	router.GET("/muscle_groups", src.GetMuscleGroups)
	router.GET("/exercises_by_body_parts", src.GetExercisesByBodyParts)

	// Run the server on port 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = PORT
	}
	router.Run(":" + port)
}
