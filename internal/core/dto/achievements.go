package dto

type Achievements struct {
	ID   string `json:"id"`
	Code string `json:"code"`
}

type AchievedAchievement struct {
	AchievedDate string `json:"achievedDate"`

	ID                      string       `json:"id"`
	AchievementID           string       `gorm:"type:uuid;primaryKey" json:"achievementId"`
	AchievementCode         string       `gorm:"type:varchar(50);primaryKey" json:"achievementCode"`
	UserLogin               string       `gorm:"type:varchar(50);primaryKey" json:"userLogin"`
	Achievement             Achievements `gorm:"foreignKey:AchievementID" json:"achievement"`
	Exp                     int64        `json:"exp"`
	Level                   int64        `json:"level"`
	TotalAchieved           int64        `json:"totalAchieved"`
	TotalAchievedPercentage float64      `json:"totalAchievedPercentage"`
	CurrentProgress         int64        `json:"currentProgress"`
	Threshold               int64        `json:"threshold"`
}

type AchievementsResponse struct {
	Data []AchievedAchievement `json:"data"`
}

type AchievementResponse struct {
	Data AchievedAchievement `json:"data"`
}

type AchievementsItemsResponse struct {
	Data []Achievements `json:"data"`
}
