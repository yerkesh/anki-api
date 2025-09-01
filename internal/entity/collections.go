package entity

type Language string

type Collection struct {
	ID               int32    `json:"id"`
	UserID           int32    `json:"user_id" validate:"required"`
	Name             string   `json:"name" validate:"required"`
	NativeLanguage   Language `json:"native_language" validate:"required"`
	LearningLanguage Language `json:"learning_language" validate:"required"`
}
