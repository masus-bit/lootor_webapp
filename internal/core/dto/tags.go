package dto

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type Tags struct {
	CreatedAt string         `json:"createdAt"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`

	ID          string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	PrimaryID   string `gorm:"type:uuid" json:"primaryId"`
	Primary     *Tags  `gorm:"foreignKey:PrimaryID" json:"primary"`
	IsPrimary   bool   `gorm:"default:true" json:"isPrimary"`
	Description string `gorm:"type:text" json:"description"`

	Slug   string   `gorm:"type:varchar(255);not null;unique" json:"slug"`
	Author SubUsers `json:"author"`

	TotalPosts           int64 `json:"totalPosts"`
	TotalCollections     int64 `json:"totalCollections"`
	TotalCollectionItems int64 `json:"totalCollectionItems"`
	TotalPhotos          int64 `json:"totalPhotos"`

	TagLinks []TagLinks `gorm:"foreignKey:TagID" json:"tagLinks"`
	Synonyms []Tags     `gorm:"foreignKey:PrimaryID" json:"synonyms"`

	Entities Entities `json:"entities"`

	CollectionItemsProps *CollectionItemsProps `json:"collectionItemsProps"`

	IsSeries      bool       `gorm:"default:true" json:"isSeries"`
	SeriesID      string     `gorm:"type:uuid" json:"seriesId"`
	SeriesEntries []Tags     `gorm:"foreignKey:PrimaryID" json:"seriesEntries"`
	Series        *ShortTags `gorm:"foreignKey:SeriesID" json:"series"`

	Images           []string       `gorm:"type:text[]" json:"images"`
	AdditionalFields datatypes.JSON `gorm:"type:jsonb" json:"additionalFields"`
}

type CollectionItemsProps struct {
	Books              int64 `json:"books"`
	VideoGames         int64 `json:"videoGames"`
	BoardGames         int64 `json:"boardGames"`
	Comics             int64 `json:"comics"`
	GamingHardware     int64 `json:"gamingHardware"`
	Vinyl              int64 `json:"vinyl"`
	Steelbooks         int64 `json:"steelbooks"`
	CollectibleCards   int64 `json:"collectibleCards"`
	CollectibleFigures int64 `json:"collectibleFigures"`
}

type Entities struct {
	Posts           []Posts                   `json:"posts"`
	CollectionItems []CollectionItemsResponse `json:"collectionItems"`
	Collections     []CollectionsResponse     `json:"collections"`
	Photos          []Photos                  `json:"photos"`
}

type TagLinks struct {
	TagID      string    `gorm:"type:uuid;primaryKey" json:"tagId"`
	EntityType string    `gorm:"type:varchar(50);primaryKey" json:"entityType"`
	EntityID   string    `gorm:"type:uuid;primaryKey" json:"entityId"`
	CreatedAt  time.Time `gorm:"default:now()" json:"createdAt"`

	Tag Tags `gorm:"foreignKey:TagID" json:"tag"`
}

type ShortTags struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	PrimaryID string `json:"primaryId"`
	SeriesID  string `json:"seriesId"`
	Type      string `json:"type"`
}

type TagsDataResponse struct {
	Data  []Tags `json:"data"`
	Total int64  `json:"total"`
}

type TagDataResponse struct {
	Data  Tags  `json:"data"`
	Total int64 `json:"total"`
}

type TagCreateRequest struct {
	Names            []string       `json:"names"`
	Author           string         `json:"author"`
	EntityID         string         `json:"entityId"`
	EntityType       string         `json:"entityType"`
	Images           []string       `json:"images"`
	AdditionalFields datatypes.JSON `json:"additionalFields"`
}

type TagsSearchRequest struct {
	Name string `json:"name"`
}

type MergeTagsRequest struct {
	FromTagIDs []string `json:"fromTagIds"`
	ToTagID    string   `json:"toTagId"`
}

type AddTagToEntityRequest struct {
	TagIDs     []string `json:"tagIds"`
	EntityType string   `json:"entityType"`
	EntityID   string   `json:"entityId"`
	Author     string   `json:"author"`
}

type RemoveTagsRequest struct {
	EntityID   string   `json:"entityId"`
	TagIDs     []string `json:"tagIds"`
	EntityType string   `json:"entityType"`
}

type TagUpdateRequest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

type TagIDsResponse struct {
	Data []string `json:"data"`
}

type TagsShortDataResponse struct {
	Data []ShortTags `json:"data"`
}

type ParsedSeries struct {
	OriginalName string  `json:"originalName"`
	CoreName     string  `json:"coreName"`
	Suffix       string  `json:"suffix"`
	Confidence   float64 `json:"confidence"`
	PatternUsed  string  `json:"patternUsed"`
	IsExactMatch bool    `json:"isExactMatch"`
	FoundTagID   string  `json:"foundTagId,omitempty"`
}

type ParsedTitle struct {
	Original   string  `json:"original"`
	SeriesCore string  `json:"seriesCore"`
	GamePart   string  `json:"gamePart"`
	Edition    string  `json:"edition"`
	Platform   string  `json:"platform"`
	Year       string  `json:"year"`
	Confidence float64 `json:"confidence"`
}

type Results struct {
	Tag           *Tags        `json:"tag"`
	SeriesTag     *Tags        `json:"seriesTag"`
	SeriesEntries []Tags       `json:"seriesEntries"`
	Confidence    float64      `json:"confidence"`
	MatchType     string       `json:"matchType"`
	Score         float64      `json:"score"`
	ParsedData    *ParsedTitle `json:"parsedData"`
}
type SearchResult struct {
	SearchTerm   string       `json:"searchTerm"`
	ParsedData   *ParsedTitle `json:"parsedData"`
	Results      []Results    `json:"results"`
	TotalCount   int32        `json:"totalCount"`
	SearchTimeNs int64        `json:"searchTimeNs"`
	DidLearn     bool         `json:"didLearn"`
}
type UserChoiceRequest struct {
	SearchQuery   string `json:"searchQuery"`
	SelectedTagID string `json:"selectedTagId"`
	SessionID     string `json:"sessionId"`
	UserID        string `json:"userId"`
	WasCorrect    bool   `json:"wasCorrect"`
	FeedbackScore int32  `json:"feedbackScore"`
}

type SearchSuggestions struct {
	Title      string  `json:"title"`
	Type       string  `json:"type"`
	Score      float64 `json:"score"`
	Confidence float64 `json:"confidence"`
	TagID      string  `json:"tagId"`
}
type SearchSuggestionResponse struct {
	Data []SearchSuggestions `json:"data"`
}
type SearchResponse struct {
	Data *SearchResult `json:"data"`
}

type DeleteTags struct {
	IDs []string `json:"ids"`
}
