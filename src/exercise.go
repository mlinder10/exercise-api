package src

type Exercise struct {
	ID                    string   `json:"id"`
	Name                  string   `json:"name"`
	Description           string   `json:"description"`
	Directions            []string `json:"directions"`
	Cues                  []string `json:"cues"`
	ImageURL              string   `json:"image_url"`
	VideoURL              string   `json:"video_url"`
	BodyPart              string   `json:"body_part"`
	PrimaryMuscleGroup    string   `json:"primary_muscle_group"`
	SecondaryMuscleGroups []string `json:"secondary_muscle_groups"`
	Compound              bool     `json:"compound"`
	ExerciseType          string   `json:"exercise_type"`
	GeneralBestRepRange   struct {
		Low  int `json:"low"`
		High int `json:"high"`
	} `json:"general_best_rep_range"`
}
