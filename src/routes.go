package src

import (
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type filter struct {
	Name                  string   `form:"name"`
	BodyPart              string   `form:"body_part"`
	PrimaryMuscleGroup    string   `form:"primary_muscle_group"`
	SecondaryMuscleGroups []string `form:"secondary_muscle_groups"`
	Compound              string   `form:"compound"`
	ExerciseType          string   `form:"exercise_type"`
}

func filterExercises(exercises []Exercise, filter filter) []Exercise {
	filtered := []Exercise{}

	for _, exercise := range exercises {
		nameFilterPassed := filter.Name == "" || strings.Contains(strings.ToLower(exercise.Name), strings.ToLower(filter.Name))
		bodyPartFilterPassed := filter.BodyPart == "" || exercise.BodyPart == filter.BodyPart
		primaryMuscleGroupFilterPassed := filter.PrimaryMuscleGroup == "" || exercise.PrimaryMuscleGroup == filter.PrimaryMuscleGroup
		secondaryMuscleGroupFilterPassed := true
		for _, group := range filter.SecondaryMuscleGroups {
			if !slices.Contains(exercise.SecondaryMuscleGroups, group) {
				secondaryMuscleGroupFilterPassed = false
				break
			}
		}
		compoundFilterPassed := filter.Compound == "" || (exercise.Compound && filter.Compound == "true") || (!exercise.Compound && filter.Compound == "false")
		exerciseTypeFilterPassed := filter.ExerciseType == "" || exercise.ExerciseType == filter.ExerciseType

		if nameFilterPassed && bodyPartFilterPassed && primaryMuscleGroupFilterPassed && secondaryMuscleGroupFilterPassed && compoundFilterPassed && exerciseTypeFilterPassed {
			filtered = append(filtered, exercise)
		}
	}

	return filtered
}

func imposeLimitAndOffset(exercises []Exercise, c *gin.Context) []Exercise {
	limit := c.MustGet("limit").(int)
	offset := c.MustGet("offset").(int) % len(exercises)

	target := []Exercise{}

	i := 0
	for i < limit && i < len(exercises) {
		target = append(target, exercises[offset+i%len(exercises)])
		i++
	}

	return target
}

func AuthMiddleware() gin.HandlerFunc {
	rapidApiKey := os.Getenv("X_RAPIDAPI_PROXY_SECRET")
	return func(c *gin.Context) {
		rapidApiKeyHeader := c.GetHeader("X-RapidAPI-Proxy-Secret")
		if rapidApiKeyHeader != rapidApiKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func DataMiddleware(exercises []Exercise) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := c.Query("limit")
		offset := c.Query("offset")

		if limit == "" {
			c.Set("limit", 10)
		} else {
			limitInt, err := strconv.Atoi(limit)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
				c.Abort()
				return
			}
			if limitInt < 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
				c.Abort()
				return
			}
			c.Set("limit", limitInt)
		}

		if offset == "" {
			c.Set("offset", 0)
		} else {
			offsetInt, err := strconv.Atoi(offset)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
				c.Abort()
				return
			}
			if offsetInt < 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
				c.Abort()
				return
			}
			c.Set("offset", offsetInt)
		}

		c.Set("exercises", exercises)
		c.Next()
	}
}

func ImageHandler(c *gin.Context) {
	c.File("./assets/" + c.Param("name") + ".png")
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, "success")
}

func GetAllExercises(c *gin.Context) {
	exercises := c.MustGet("exercises").([]Exercise)
	target := imposeLimitAndOffset(exercises, c)

	c.JSON(http.StatusOK, target)
}

func GetExerciseByID(c *gin.Context) {
	id := c.Param("id")
	exercises := c.MustGet("exercises").([]Exercise)
	for _, exercise := range exercises {
		if exercise.ID == id {
			c.JSON(http.StatusOK, exercise)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "exercise not found"})
}

func GetMuscleGroups(c *gin.Context) {
	exercises := c.MustGet("exercises").([]Exercise)
	encountered := map[string]bool{}
	muscleGroups := []string{}

	for _, exercise := range exercises {
		if !encountered[exercise.PrimaryMuscleGroup] {
			encountered[exercise.PrimaryMuscleGroup] = true
			muscleGroups = append(muscleGroups, exercise.PrimaryMuscleGroup)
		}
	}

	c.JSON(http.StatusOK, muscleGroups)
}

func GetExerciseTypes(c *gin.Context) {
	exercises := c.MustGet("exercises").([]Exercise)
	encountered := map[string]bool{}
	exerciseTypes := []string{}

	for _, exercise := range exercises {
		if !encountered[exercise.ExerciseType] {
			encountered[exercise.ExerciseType] = true
			exerciseTypes = append(exerciseTypes, exercise.ExerciseType)
		}
	}

	c.JSON(http.StatusOK, exerciseTypes)
}

func GetBodyParts(c *gin.Context) {
	exercises := c.MustGet("exercises").([]Exercise)
	encountered := map[string]bool{}
	bodyParts := []string{}

	for _, exercise := range exercises {
		if !encountered[exercise.BodyPart] {
			encountered[exercise.BodyPart] = true
			bodyParts = append(bodyParts, exercise.BodyPart)
		}
	}

	c.JSON(http.StatusOK, bodyParts)
}

func GetFilteredExercises(c *gin.Context) {
	var filter filter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	exercises := c.MustGet("exercises").([]Exercise)
	filtered := filterExercises(exercises, filter)
	target := imposeLimitAndOffset(filtered, c)

	c.JSON(http.StatusOK, target)
}

func GetExercisesByBodyParts(c *gin.Context) {
	exercises := c.MustGet("exercises").([]Exercise)

	payload := map[string][]Exercise{}
	for _, exercise := range exercises {
		payload[exercise.BodyPart] = append(payload[exercise.BodyPart], exercise)
	}

	c.JSON(http.StatusOK, payload)
}
