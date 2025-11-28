package dto

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"lootor/internal/core/models"
	"time"
)

type Comments struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt string    `json:"deletedAt" gorm:"index"`

	ID            string          `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Date          string          `json:"date"`
	TargetID      string          `json:"targetId"`
	Author        models.SubUsers `json:"author"`
	ParentID      string          `json:"parentId"`
	LikesCount    int             `json:"likesCount"`
	DislikesCount int             `json:"dislikesCount"`
	EntityType    string          `json:"entityType"`

	Likes    []models.SubUsers `gorm:"foreignKey:CommentID;references:ID;constraint:OnDelete:CASCADE;" json:"likes"`
	Dislikes []models.SubUsers `gorm:"foreignKey:CommentID;references:ID;constraint:OnDelete:CASCADE;" json:"dislikes"`

	Content datatypes.JSON `gorm:"type:jsonb" json:"content"`
}

type CommentsDataResponse struct {
	Data  []CommentsResponse `json:"data"`
	Total int64              `json:"total"`
}

type ChildrenComments struct {
	*Comments
	Likes         []models.SubUsers `json:"likes"`
	Dislikes      []models.SubUsers `json:"dislikes"`
	LikesCount    int               `json:"likesCount"`
	DislikesCount int               `json:"dislikesCount"`
}

type CommentsResponse struct {
	ChildrenComments []ChildrenComments `json:"childrenComments"`
	*Comments
	AnswersTotal int64 `json:"answersTotal"`
}

type CommentsRequest struct {
	Date            string  `json:"date"`
	TargetID        string  `json:"targetId"`
	Author          string  `json:"author"`
	ParentID        *string `json:"parentId"`
	TargetUserLogin *string `json:"targetUserLogin"`
	EntityType      string  `json:"entityType"`

	Content datatypes.JSON `json:"content"`
}

type CommentsRequestSwag struct {
	Date       string    `json:"date"`
	TargetID   string    `json:"targetId"`
	Author     string    `json:"author"`
	ParentID   uuid.UUID `json:"parentId"`
	EntityType string    `json:"entityType"`

	Content map[string]interface{} `json:"content"`
}

type CommentDataResponse struct {
	Data CommentsResponse `json:"data"`
}

type Likes struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	ID        string          `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Author    models.SubUsers `json:"author"`
	CommentID string          `json:"commentID"`
	Comment   Comments        `gorm:"foreignKey:CommentID;references:ID;constraint:OnDelete:CASCADE;"`
}

type AnswersItem struct {
	*Comments
	Likes         []models.SubUsers `json:"likes"`
	Dislikes      []models.SubUsers `json:"dislikes"`
	LikesCount    int               `json:"likesCount"`
	DislikesCount int               `json:"dislikesCount"`
}

type AnswersRequest struct {
	ID     string `json:"id"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}
type AnswersDataResponse struct {
	Data  []AnswersItem `json:"data"`
	Total int64         `json:"total"`
}
