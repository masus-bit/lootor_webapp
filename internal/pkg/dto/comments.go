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

	Id            string          `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Date          string          `json:"date"`
	TargetId      string          `json:"targetId"`
	Author        models.SubUsers `json:"author"`
	ParentId      string          `json:"parentId"`
	LikesCount    int             `json:"likesCount"`
	DislikesCount int             `json:"dislikesCount"`

	Likes    []models.SubUsers `gorm:"foreignKey:CommentId;references:Id;constraint:OnDelete:CASCADE;" json:"likes"`
	Dislikes []models.SubUsers `gorm:"foreignKey:CommentId;references:Id;constraint:OnDelete:CASCADE;" json:"dislikes"`

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
	TargetId        string  `json:"targetId"`
	Author          string  `json:"author"`
	ParentId        *string `json:"parentId"`
	TargetUserLogin *string `json:"targetUserLogin"`

	Content datatypes.JSON `json:"content"`
}

type CommentsRequestSwag struct {
	Date     string    `json:"date"`
	TargetId string    `json:"targetId"`
	Author   string    `json:"author"`
	ParentId uuid.UUID `json:"parentId"`

	Content map[string]interface{} `json:"content"`
}

type CommentDataResponse struct {
	Data CommentsResponse `json:"data"`
}

type Likes struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Id        string          `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Author    models.SubUsers `json:"author"`
	CommentId string          `json:"commentId"`
	Comment   Comments        `gorm:"foreignKey:CommentId;references:Id;constraint:OnDelete:CASCADE;"`
}

type AnswersItem struct {
	*Comments
	Likes         []models.SubUsers `json:"likes"`
	Dislikes      []models.SubUsers `json:"dislikes"`
	LikesCount    int               `json:"likesCount"`
	DislikesCount int               `json:"dislikesCount"`
}

type AnswersRequest struct {
	Id     string `json:"id"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}
type AnswersDataResponse struct {
	Data  []AnswersItem `json:"data"`
	Total int64         `json:"total"`
}
