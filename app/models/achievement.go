package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateAchievementRequest struct {
	AchievementType string                 `json:"achievementType" validate:"required"`
	Title           string                 `json:"title" validate:"required"`
	Description     string                 `json:"description"`
	Details         map[string]interface{} `json:"details"`
	Tags            []string               `json:"tags"`
	EventDate       time.Time              `json:"eventDate"`
}

type AchievementMongo struct {
	ID              primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	StudentID       string                 `bson:"studentId" json:"student_id"`
	AchievementType string                 `bson:"achievementType" json:"achievement_type"`
	Title           string                 `bson:"title" json:"title"`
	Description     string                 `bson:"description" json:"description"`
	Details         map[string]interface{} `bson:"details" json:"details"`
	Tags            []string               `bson:"tags" json:"tags"`
	Points          int                    `bson:"points" json:"points"`
	CreatedAt       time.Time              `bson:"createdAt" json:"created_at"`
	UpdatedAt       time.Time              `bson:"updatedAt" json:"updated_at"`
}

type AchievementReference struct {
	ID                 string    `json:"id"`
	StudentID          string    `json:"student_id"`
	MongoAchievementID string    `json:"mongo_achievement_id"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"created_at"`
}

// Struct Response untuk List (Admin View)
type AchievementResponse struct {
	ID              string                 `json:"id"`
	MongoID         string                 `json:"mongo_id"`
	StudentID       string                 `json:"student_id"`
	StudentName     string                 `json:"student_name"`
	StudentNIM      string                 `json:"student_nim"`
	AchievementType string                 `json:"achievement_type"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	Status          string                 `json:"status"`
	Points          int                    `json:"points"`
	Tags            []string               `json:"tags"`
	Details         map[string]interface{} `json:"details"`

	// Field Tambahan untuk Detail
	SubmittedAt     *time.Time             `json:"submitted_at,omitempty"`
	VerifiedAt      *time.Time             `json:"verified_at,omitempty"`
	RejectionNote   string                 `json:"rejection_note,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
}

// Struct untuk item history (timeline)
type HistoryItem struct {
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
	Actor     string    `json:"actor"`
	Note      string    `json:"note,omitempty"`
}