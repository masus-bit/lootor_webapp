package dto

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"lootor/internal/core/models"
	"time"
)

type ReactRequest struct {
	Reaction  string `json:"reaction"`
	UserLogin string `json:"userLogin"`
	PostId    string `json:"postId"`
}

type PostsDataResponse struct {
	Data  []Posts `json:"data"`
	Total int64   `json:"total"`
}

type PostDataResponse struct {
	Data Posts `json:"data"`
}

type PostUpdateRequest struct {
	Id      string
	Content datatypes.JSON `json:"content"`
	IsDraft bool           `json:"isDraft"`
}

type PostRequest struct {
	Date    string         `json:"date"`
	Author  string         `json:"author"`
	Content datatypes.JSON `json:"content"`
	IsDraft bool           `json:"isDraft"`
}

type Posts struct {
	Id             string          `json:"id"`
	Date           string          `json:"date"`
	Author         models.SubUsers `json:"author"`
	HeartCount     int             `json:"heartCount"`
	FireCount      int             `json:"fireCount"`
	GlassesCount   int             `json:"glassesCount"`
	LaughCount     int             `json:"laughCount"`
	TearsCount     int             `json:"tearsCount"`
	PokerFaceCount int             `json:"pokerFaceCount"`
	EyesCount      int             `json:"eyesCount"`
	AngryCount     int             `json:"angryCount"`
	ShitCount      int             `json:"shitCount"`
	ClownCount     int             `json:"clownCount"`
	Reacted        string          `json:"reacted"`
	TotalReactions int             `gorm:"->;type:int generated always as (heart_count + fire_count + glasses_count + laugh_count + tears_count + poker_face_count + eyes_count + angry_count + shit_count + clown_count) stored" json:"totalReactions"`
	IsDraft        bool            `json:"isDraft"`

	Content datatypes.JSON `gorm:"type:jsonb" json:"content"`

	Reactions *ReactResponse `gorm:"foreignKey:PostId" json:"reactions"`
}

type ReactionType string

const (
	Fire      ReactionType = "fire"
	Heart     ReactionType = "heart"
	Glasses   ReactionType = "glasses"
	Laugh     ReactionType = "laugh"
	Tears     ReactionType = "tears"
	PokerFace ReactionType = "pokerFace"
	Eyes      ReactionType = "eyes"
	Angry     ReactionType = "angry"
	Shit      ReactionType = "shit"
	Clown     ReactionType = "clown"
)

type PostReactions struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`

	Id        uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PostId    string       `gorm:"not null;uniqueIndex:idx_user_post" json:"postId"`
	UserLogin string       `gorm:"not null;uniqueIndex:idx_user_post" json:"userLogin"`
	Reaction  ReactionType `json:"reaction"`

	Post Posts `gorm:"foreignKey:PostId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type React struct {
	Id       string          `json:"id"`
	User     models.SubUsers `json:"user"`
	Reaction ReactionType    `json:"reaction"`
}

type ReactResponse struct {
	Fire      []React `json:"fire"`
	Heart     []React `json:"heart"`
	Glasses   []React `json:"glasses"`
	Laugh     []React `json:"laugh"`
	Tears     []React `json:"tears"`
	PokerFace []React `json:"pokerFace"`
	Eyes      []React `json:"eyes"`
	Angry     []React `json:"angry"`
	Shit      []React `json:"shit"`
	Clown     []React `json:"clown"`
}
