package utils

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/structpb"
	"gorm.io/datatypes"
	"lootor/gen/go/microservices"
	"lootor/internal/core/dto"
	"time"

	"lootor/internal/core/models"
	"lootor/internal/infrastructure/achievementsclient"
	"slices"
	"strings"
	"unicode"
)

func RemoveByValue(slice []string, value string) []string {
	for i, v := range slice {
		if v == value {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func RemoveByValueStruct(slice []models.CollectionItems, id uuid.UUID) []models.CollectionItems {
	for i, v := range slice {
		if v.ID == id {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func DefineShareString(authorizedUser string, id string, collection *models.Collections) string {
	if authorizedUser == id && collection.IsPrivate {
		return collection.ShareString
	}
	return ""
}

func GetCollectionOrderBy(
	orderByInput string,
	collections []dto.CollectionsResponse,
	order string,
) []dto.CollectionsResponse {
	switch orderByInput {
	case "name":
		slices.SortFunc(
			collections, func(a, b dto.CollectionsResponse) int {
				if order == "asc" {
					return cmp.Compare(a.Name, b.Name)
				}
				return cmp.Compare(b.Name, a.Name)

			},
		)
	case "created":
		slices.SortFunc(
			collections, func(a, b dto.CollectionsResponse) int {
				if order == "asc" {
					return cmp.Compare(a.CreatedAt.UnixNano(), b.CreatedAt.UnixNano())
				}
				return cmp.Compare(b.CreatedAt.UnixNano(), a.CreatedAt.UnixNano())

			},
		)
	case "totalPrice":
		slices.SortFunc(
			collections, func(a, b dto.CollectionsResponse) int {
				if order == "asc" {
					return cmp.Compare(a.TotalPrice, b.TotalPrice)
				}
				return cmp.Compare(b.TotalPrice, a.TotalPrice)

			},
		)
	case "collectionItemsCount":
		slices.SortFunc(
			collections, func(a, b dto.CollectionsResponse) int {
				if order == "asc" {
					return cmp.Compare(a.CollectionItemsCount, b.CollectionItemsCount)
				}
				return cmp.Compare(b.CollectionItemsCount, a.CollectionItemsCount)

			},
		)
	case "likesCount":
		slices.SortFunc(
			collections, func(a, b dto.CollectionsResponse) int {
				if order == "asc" {
					return cmp.Compare(a.LikesCount, b.LikesCount)
				}
				return cmp.Compare(b.LikesCount, a.LikesCount)

			},
		)
	default:
		slices.SortFunc(
			collections, func(a, b dto.CollectionsResponse) int {
				return cmp.Compare(b.LikesCount, a.LikesCount)

			},
		)
	}
	return collections
}

type CiOrder struct {
	Joins string
	Order string
}

func GetCIOrderString(orderBy string, order string) *CiOrder {
	var result CiOrder

	switch orderBy {
	case "platform":
		result.Joins = "LEFT JOIN platforms ON platforms.id = collection_items.platform_id"
		if order == "desc" {
			result.Order = "platforms.name DESC"
		}
		result.Order = "platforms.name ASC"
	case "type":
		result.Joins = "LEFT JOIN item_types ON item_types.id = collection_items.item_type_id"
		if order == "desc" {
			result.Order = "item_types.ru_name DESC"
		}
		result.Order = "item_types.ru_name ASC"
	case "name":
		if order == "desc" {
			result.Order = "name DESC"
		}
		result.Order = "name ASC"
	case "rating":
		if order == "desc" {
			result.Order = "rating DESC"
		}
		result.Order = "rating ASC"
	case "purchaseDate":
		if order == "desc" {
			result.Order = "purchase_date DESC"
		}
		result.Order = "purchase_date ASC"
	case "purchasePrice":
		if order == "desc" {
			result.Order = "purchase_price DESC"
		}
		result.Order = "purchase_price ASC"
	default:
		result.Order = "purchase_date DESC"
	}
	return &result
}

func GetReformatedItemType(itemType string) string {
	switch itemType {
	case "videoGames":
		return "Video games"
	case "boardGames":
		return "Board games"
	case "comics":
		return "Comics"
	case "gamingHardware":
		return "Gaming hardware"
	case "collectibleFigures":
		return "Collectible figures"
	case "books":
		return "Books"
	case "vinyl":
		return "Vinyl"
	case "steelbooks":
		return "Steelbooks"
	case "collectibleCards":
		return "Collectible cards"
	default:
		return ""
	}
}

func FirstNonZero[T comparable](values ...T) T {
	var zero T
	for _, v := range values {
		if v != zero {
			return v
		}
	}
	return zero
}

func RemoveOrdered[T any](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}

func GetSubscriptionType(subscriptionType string) string {
	switch subscriptionType {
	case "monthly":
		return "149"
	case "yearly":
		return "1190"
	default:
		return "99"
	}
}

func NormalizeContent(content *structpb.Struct) []byte {
	contentMap := content.AsMap()
	contentJSON, err := json.Marshal(contentMap)

	if err != nil {
		return nil
	}
	return contentJSON
}

func GormJSONToProtoStruct(jsonData datatypes.JSON) (*structpb.Struct, error) {
	if len(jsonData) == 0 {
		return &structpb.Struct{Fields: make(map[string]*structpb.Value)}, nil
	}

	var contentMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &contentMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return structpb.NewStruct(contentMap)
}

func Slugify(input string) string {
	translitMap := map[rune]string{
		'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "yo",
		'ж': "zh", 'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m",
		'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u",
		'ф': "f", 'х': "kh", 'ц': "ts", 'ч': "ch", 'ш': "sh", 'щ': "shch",
		'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya",
		'А': "a", 'Б': "b", 'В': "v", 'Г': "g", 'Д': "d", 'Е': "e", 'Ё': "yo",
		'Ж': "zh", 'З': "z", 'И': "i", 'Й': "y", 'К': "k", 'Л': "l", 'М': "m",
		'Н': "n", 'О': "o", 'П': "p", 'Р': "r", 'С': "s", 'Т': "t", 'У': "u",
		'Ф': "f", 'Х': "kh", 'Ц': "ts", 'Ч': "ch", 'Ш': "sh", 'Щ': "shch",
		'Ъ': "", 'Ы': "y", 'Ь': "", 'Э': "e", 'Ю': "yu", 'Я': "ya",
	}

	var result strings.Builder
	hasRussian := false

	for _, char := range input {
		if unicode.Is(unicode.Cyrillic, char) {
			hasRussian = true
			break
		}
	}

	for _, char := range input {
		switch {
		case char == ' ':
			result.WriteString("-")
		case char == ':':
			result.WriteString("-")
		case hasRussian && unicode.Is(unicode.Cyrillic, char):
			if val, ok := translitMap[char]; ok {
				result.WriteString(val)
			}
		case unicode.IsUpper(char):
			result.WriteRune(unicode.ToLower(char))
		default:
			result.WriteRune(char)
		}
	}

	return strings.ToLower(result.String())
}

func UsersOrder(order string) string {
	var finalOrder string

	switch order {
	case "socialScore":
		finalOrder = "social_score"
	default:
		finalOrder = order
	}

	return finalOrder
}

func FormatPost(
	post *microservices.PostItem,
	users []dto.SubUsers,
	author *models.Users,
	reactLength int,
) (*dto.PostDataResponse, error) {
	content := NormalizeContent(post.Content)

	reacts, err := FormatReacts(post, users, reactLength)
	if err != nil {
		return nil, err
	}

	return &dto.PostDataResponse{Data: *fillPost(post, content, reacts, author)}, nil
}

func FormatReacts(post *microservices.PostItem, users []dto.SubUsers, reactLength int) (*dto.ReactResponse, error) {
	var reactionsResult dto.ReactResponse
	reacts := map[string][]dto.React{}
	var reactUsers []string
	var fullReactUsers = map[string]dto.SubUsers{}
	if reactLength != 0 {
		for _, r := range post.Reactions {
			reactUsers = append(reactUsers, r.UserLogin)
		}

		for _, u := range users {
			fullReactUsers[u.Login] = u
		}

		reactionMapping := []struct {
			field *[]dto.React
			key   string
		}{
			{&reactionsResult.Fire, "fire"},
			{&reactionsResult.Heart, "heart"},
			{&reactionsResult.Glasses, "glasses"},
			{&reactionsResult.Laugh, "laugh"},
			{&reactionsResult.Tears, "tears"},
			{&reactionsResult.PokerFace, "pokerFace"},
			{&reactionsResult.Eyes, "eyes"},
			{&reactionsResult.Angry, "angry"},
			{&reactionsResult.Shit, "shit"},
			{&reactionsResult.Clown, "clown"},
		}

		for _, mapping := range reactionMapping {
			*mapping.field = []dto.React{}
		}

		for _, r := range post.GetReactions() {
			reactUUID, err := uuid.Parse(r.Id)
			if err != nil {
				return nil, err
			}
			reacts[r.GetReaction()] = append(
				reacts[r.GetReaction()], dto.React{
					ID:       reactUUID.String(),
					User:     fullReactUsers[r.UserLogin],
					Reaction: dto.ReactionType(r.Reaction),
				},
			)
		}

		for _, mapping := range reactionMapping {
			if slice, exists := reacts[mapping.key]; exists && len(slice) > 0 {
				*mapping.field = slice
			}
		}
	}

	return &reactionsResult, nil
}

func fillPost(p *microservices.PostItem, content []byte, reacts *dto.ReactResponse, user *models.Users) *dto.Posts {
	return &dto.Posts{
		ID:   p.GetId(),
		Date: p.GetDate(),
		Author: dto.SubUsers{
			Login:       user.Login,
			AvatarURL:   user.AvatarURL,
			ProfileName: user.ProfileName,
			IsPremium:   user.IsPremium,
		},
		Content:        content,
		HeartCount:     int(p.GetHeartCount()),
		FireCount:      int(p.GetFireCount()),
		GlassesCount:   int(p.GetGlassesCount()),
		LaughCount:     int(p.GetLaughCount()),
		TearsCount:     int(p.GetTearsCount()),
		PokerFaceCount: int(p.GetPokerFaceCount()),
		EyesCount:      int(p.GetEyesCount()),
		AngryCount:     int(p.GetAngryCount()),
		ShitCount:      int(p.GetShitCount()),
		ClownCount:     int(p.GetClownCount()),
		TotalReactions: int(p.GetTotalReactions()),
		Reactions:      reacts,
		Reacted:        p.GetReacted(),
		Views:          int(p.Views),
		CommentsCount:  int(p.CommentsCount),
		IsDraft:        p.IsDraft,
		Title:          p.Title,
		Translit:       p.Translit,
	}

}

func getType(tagProto *microservices.TagItemShort) string {
	primaryID := tagProto.GetPrimaryId()
	seriesID := tagProto.GetSeriesId()

	if primaryID == "" {
		if seriesID == "" {
			return "series"
		}
		return "primarySeriesChild"
	}
	if seriesID == "" {
		return "synonymSeries"
	}
	return "synonymSeriesChild"
}

func NormalizeTagsShort(tagsProto []*microservices.TagItemShort) []dto.ShortTags {
	tempTags := make([]dto.ShortTags, 0, len(tagsProto))
	for _, tag := range tagsProto {
		addFields := NormalizeContent(tag.AdditionalFields)
		tempTags = append(
			tempTags, dto.ShortTags{
				ID:               tag.GetId(),
				Name:             tag.GetName(),
				Slug:             tag.GetSlug(),
				PrimaryID:        tag.GetPrimaryId(),
				SeriesID:         tag.GetSeriesId(),
				Type:             getType(tag),
				AdditionalFields: addFields,
			},
		)
	}

	return tempTags
}

func CanLike(likes []string, user, owner string) bool {
	if user == "" || user == owner {
		return true
	}
	return !slices.Contains(likes, user)
}

func GetAchievementSubscribersData(current int64) (int64, int64, bool) {
	var currentLevel int64
	switch {
	case current >= SubscribersLevel1 && current < SubscribersLevel2:
		currentLevel = 1
	case current >= SubscribersLevel2 && current < SubscribersLevel3:
		currentLevel = 2
	case current >= SubscribersLevel3:
		currentLevel = 3
	default:
		currentLevel = 0
	}

	var previousLevel int64
	previous := current - 1
	switch {
	case previous >= SubscribersLevel1 && previous < SubscribersLevel2:
		previousLevel = 1
	case previous >= SubscribersLevel2 && previous < SubscribersLevel3:
		previousLevel = 2
	case previous >= SubscribersLevel3:
		previousLevel = 3
	default:
		previousLevel = 0
	}

	var xp int64
	var level int64
	switch {
	case current >= SubscribersLevel1 && current < SubscribersLevel2:
		xp = XPSubscribersLevel1
		level = 1
	case current >= SubscribersLevel2 && current < SubscribersLevel3:
		xp = XPSubscribersLevel2
		level = 2
	case current >= SubscribersLevel3:
		xp = XPSubscribersLevel3
		level = 3
	default:
		xp = 0
		level = 0
	}

	levelUpdated := currentLevel > previousLevel

	return xp, level, levelUpdated
}

func GetAchievementBetaTesterData() (int64, int64) {
	return XPBetaTester, 1
}

func GetAchievementBetaTesterDonateData() (int64, int64) {
	return XPBetaTesterDonate, 1
}

func GetAchievementFirstCollectionCreateData() (int64, int64) {
	return XPFirstCollectionCreate, 1
}

func GetAchievementDonateData(monthsDonated int64) (int64, int64) {
	if monthsDonated >= 12 {
		return XPDonateLevel2, 1
	} else if monthsDonated >= 1 {
		return XPDonateLevel1, 1
	}
	return 0, 0
}

func GetAchievementYearlyRegisterData(yearsRegistered int64) (int64, int64) {
	if yearsRegistered >= 1 {
		return XPYearlyRegister * yearsRegistered, yearsRegistered
	}
	return 0, 0
}

func GetAchievementCollectionItemsAddData(current int64) (int64, int64, bool) {
	var currentLevel int64
	switch {
	case current >= CollectionItemsAddLevel1 && current < CollectionItemsAddLevel2:
		currentLevel = 1
	case current >= CollectionItemsAddLevel2 && current < CollectionItemsAddLevel3:
		currentLevel = 2
	case current >= CollectionItemsAddLevel3:
		currentLevel = 3
	default:
		currentLevel = 0
	}

	var previousLevel int64
	previous := current - 1
	switch {
	case previous >= CollectionItemsAddLevel1 && previous < CollectionItemsAddLevel2:
		previousLevel = 1
	case previous >= CollectionItemsAddLevel2 && previous < CollectionItemsAddLevel3:
		previousLevel = 2
	case previous >= CollectionItemsAddLevel3:
		previousLevel = 3
	default:
		previousLevel = 0
	}

	var xp int64
	var level int64
	switch {
	case current >= CollectionItemsAddLevel1 && current < CollectionItemsAddLevel2:
		xp = XPCollectionItemsAddLevel1
		level = 1
	case current >= CollectionItemsAddLevel2 && current < CollectionItemsAddLevel3:
		xp = XPCollectionItemsAddLevel2
		level = 2
	case current >= CollectionItemsAddLevel3:
		xp = XPCollectionItemsAddLevel3
		level = 3
	default:
		xp = 0
		level = 0
	}

	levelUpdated := currentLevel > previousLevel

	return xp, level, levelUpdated
}

func GetAchievementPhotosAddData(current int64) (int64, int64, bool) {
	var currentLevel int64
	switch {
	case current >= PhotosAddLevel1 && current < PhotosAddLevel2:
		currentLevel = 1
	case current >= PhotosAddLevel2 && current < PhotosAddLevel3:
		currentLevel = 2
	case current >= PhotosAddLevel3:
		currentLevel = 3
	default:
		currentLevel = 0
	}

	var previousLevel int64
	previous := current - 1
	switch {
	case previous >= PhotosAddLevel1 && previous < PhotosAddLevel2:
		previousLevel = 1
	case previous >= PhotosAddLevel2 && previous < PhotosAddLevel3:
		previousLevel = 2
	case previous >= PhotosAddLevel3:
		previousLevel = 3
	default:
		previousLevel = 0
	}

	var xp int64
	var level int64
	switch {
	case current >= PhotosAddLevel1 && current < PhotosAddLevel2:
		xp = XPPhotosAddLevel1
		level = 1
	case current >= PhotosAddLevel2 && current < PhotosAddLevel3:
		xp = XPPhotosAddLevel2
		level = 2
	case current >= PhotosAddLevel3:
		xp = XPPhotosAddLevel3
		level = 3
	default:
		xp = 0
		level = 0
	}

	levelUpdated := currentLevel > previousLevel
	return xp, level, levelUpdated
}

func GetAchievementTagsAddData(current int64) (int64, int64, bool) {
	var currentLevel int64
	switch {
	case current >= TagsAddLevel1 && current < TagsAddLevel2:
		currentLevel = 1
	case current >= TagsAddLevel2 && current < TagsAddLevel3:
		currentLevel = 2
	case current >= TagsAddLevel3:
		currentLevel = 3
	default:
		currentLevel = 0
	}

	var previousLevel int64
	previous := current - 1
	switch {
	case previous >= TagsAddLevel1 && previous < TagsAddLevel2:
		previousLevel = 1
	case previous >= TagsAddLevel2 && previous < TagsAddLevel3:
		previousLevel = 2
	case previous >= TagsAddLevel3:
		previousLevel = 3
	default:
		previousLevel = 0
	}

	var xp int64
	var level int64
	switch {
	case current >= TagsAddLevel1 && current < TagsAddLevel2:
		xp = XPTagsAddLevel1
		level = 1
	case current >= TagsAddLevel2 && current < TagsAddLevel3:
		xp = XPTagsAddLevel2
		level = 2
	case current >= TagsAddLevel3:
		xp = XPTagsAddLevel3
		level = 3
	default:
		xp = 0
		level = 0
	}

	levelUpdated := currentLevel > previousLevel
	return xp, level, levelUpdated
}

func GetAchievementPostsCreateData(current int64) (int64, int64, bool) {
	var currentLevel int64
	switch {
	case current >= PostsCreateLevel1 && current < PostsCreateLevel2:
		currentLevel = 1
	case current >= PostsCreateLevel2 && current < PostsCreateLevel3:
		currentLevel = 2
	case current >= PostsCreateLevel3:
		currentLevel = 3
	default:
		currentLevel = 0
	}

	var previousLevel int64
	previous := current - 1
	switch {
	case previous >= PostsCreateLevel1 && previous < PostsCreateLevel2:
		previousLevel = 1
	case previous >= PostsCreateLevel2 && previous < PostsCreateLevel3:
		previousLevel = 2
	case previous >= PostsCreateLevel3:
		previousLevel = 3
	default:
		previousLevel = 0
	}

	var xp int64
	var level int64
	switch {
	case current >= PostsCreateLevel1 && current < PostsCreateLevel2:
		xp = XPPostsCreateLevel1
		level = 1
	case current >= PostsCreateLevel2 && current < PostsCreateLevel3:
		xp = XPPostsCreateLevel2
		level = 2
	case current >= PostsCreateLevel3:
		xp = XPPostsCreateLevel3
		level = 3
	default:
		xp = 0
		level = 0
	}

	levelUpdated := currentLevel > previousLevel
	return xp, level, levelUpdated
}

func GetAchievementCollectionsLikesData(current int64) (int64, int64, bool) {
	var currentLevel int64
	switch {
	case current >= CollectionsLikesLevel1 && current < CollectionsLikesLevel2:
		currentLevel = 1
	case current >= CollectionsLikesLevel2 && current < CollectionsLikesLevel3:
		currentLevel = 2
	case current >= CollectionsLikesLevel3:
		currentLevel = 3
	default:
		currentLevel = 0
	}

	var previousLevel int64
	previous := current - 1
	switch {
	case previous >= CollectionsLikesLevel1 && previous < CollectionsLikesLevel2:
		previousLevel = 1
	case previous >= CollectionsLikesLevel2 && previous < CollectionsLikesLevel3:
		previousLevel = 2
	case previous >= CollectionsLikesLevel3:
		previousLevel = 3
	default:
		previousLevel = 0
	}

	var xp int64
	var level int64
	switch {
	case current >= CollectionsLikesLevel1 && current < CollectionsLikesLevel2:
		xp = XPCollectionsLikesLevel1
		level = 1
	case current >= CollectionsLikesLevel2 && current < CollectionsLikesLevel3:
		xp = XPCollectionsLikesLevel2
		level = 2
	case current >= CollectionsLikesLevel3:
		xp = XPCollectionsLikesLevel3
		level = 3
	default:
		xp = 0
		level = 0
	}

	levelUpdated := currentLevel > previousLevel
	return xp, level, levelUpdated
}

func GetAchievementCollectionItemsLikesData(current int64) (int64, int64, bool) {
	var currentLevel int64
	switch {
	case current >= CollectionItemsLikesLevel1 && current < CollectionItemsLikesLevel2:
		currentLevel = 1
	case current >= CollectionItemsLikesLevel2 && current < CollectionItemsLikesLevel3:
		currentLevel = 2
	case current >= CollectionItemsLikesLevel3:
		currentLevel = 3
	default:
		currentLevel = 0
	}

	var previousLevel int64
	previous := current - 1
	switch {
	case previous >= CollectionItemsLikesLevel1 && previous < CollectionItemsLikesLevel2:
		previousLevel = 1
	case previous >= CollectionItemsLikesLevel2 && previous < CollectionItemsLikesLevel3:
		previousLevel = 2
	case previous >= CollectionItemsLikesLevel3:
		previousLevel = 3
	default:
		previousLevel = 0
	}

	var xp int64
	var level int64
	switch {
	case current >= CollectionItemsLikesLevel1 && current < CollectionItemsLikesLevel2:
		xp = XPCollectionItemsLikesLevel1
		level = 1
	case current >= CollectionItemsLikesLevel2 && current < CollectionItemsLikesLevel3:
		xp = XPCollectionItemsLikesLevel2
		level = 2
	case current >= CollectionItemsLikesLevel3:
		xp = XPCollectionItemsLikesLevel3
		level = 3
	default:
		xp = 0
		level = 0
	}

	levelUpdated := currentLevel > previousLevel
	return xp, level, levelUpdated
}

func GetAchievementPhotosLikesData(current int64) (int64, int64, bool) {
	var currentLevel int64
	switch {
	case current >= PhotosLikesLevel1 && current < PhotosLikesLevel2:
		currentLevel = 1
	case current >= PhotosLikesLevel2 && current < PhotosLikesLevel3:
		currentLevel = 2
	case current >= PhotosLikesLevel3:
		currentLevel = 3
	default:
		currentLevel = 0
	}

	var previousLevel int64
	previous := current - 1
	switch {
	case previous >= PhotosLikesLevel1 && previous < PhotosLikesLevel2:
		previousLevel = 1
	case previous >= PhotosLikesLevel2 && previous < PhotosLikesLevel3:
		previousLevel = 2
	case previous >= PhotosLikesLevel3:
		previousLevel = 3
	default:
		previousLevel = 0
	}

	var xp int64
	var level int64
	switch {
	case current >= PhotosLikesLevel1 && current < PhotosLikesLevel2:
		xp = XPPhotosLikesLevel1
		level = 1
	case current >= PhotosLikesLevel2 && current < PhotosLikesLevel3:
		xp = XPPhotosLikesLevel2
		level = 2
	case current >= PhotosLikesLevel3:
		xp = XPPhotosLikesLevel3
		level = 3
	default:
		xp = 0
		level = 0
	}

	levelUpdated := currentLevel > previousLevel
	return xp, level, levelUpdated
}

func GetAchievementPostsReactionsData(current int64) (int64, int64, bool) {
	var currentLevel int64
	switch {
	case current >= PostsReactionsLevel1 && current < PostsReactionsLevel2:
		currentLevel = 1
	case current >= PostsReactionsLevel2 && current < PostsReactionsLevel3:
		currentLevel = 2
	case current >= PostsReactionsLevel3:
		currentLevel = 3
	default:
		currentLevel = 0
	}

	var previousLevel int64
	previous := current - 1
	switch {
	case previous >= PostsReactionsLevel1 && previous < PostsReactionsLevel2:
		previousLevel = 1
	case previous >= PostsReactionsLevel2 && previous < PostsReactionsLevel3:
		previousLevel = 2
	case previous >= PostsReactionsLevel3:
		previousLevel = 3
	default:
		previousLevel = 0
	}

	var xp int64
	var level int64
	switch {
	case current >= PostsReactionsLevel1 && current < PostsReactionsLevel2:
		xp = XPPostsReactionsLevel1
		level = 1
	case current >= PostsReactionsLevel2 && current < PostsReactionsLevel3:
		xp = XPPostsReactionsLevel2
		level = 2
	case current >= PostsReactionsLevel3:
		xp = XPPostsReactionsLevel3
		level = 3
	default:
		xp = 0
		level = 0
	}

	levelUpdated := currentLevel > previousLevel
	return xp, level, levelUpdated
}

func GetAchievementCollectionsSumData(current int64) (int64, int64, bool) {
	var currentLevel int64
	switch {
	case current >= CollectionsSumLevel1 && current < CollectionsSumLevel2:
		currentLevel = 1
	case current >= CollectionsSumLevel2 && current < CollectionsSumLevel3:
		currentLevel = 2
	case current >= CollectionsSumLevel3:
		currentLevel = 3
	default:
		currentLevel = 0
	}

	var previousLevel int64
	previous := current - 1
	switch {
	case previous >= CollectionsSumLevel1 && previous < CollectionsSumLevel2:
		previousLevel = 1
	case previous >= CollectionsSumLevel2 && previous < CollectionsSumLevel3:
		previousLevel = 2
	case previous >= CollectionsSumLevel3:
		previousLevel = 3
	default:
		previousLevel = 0
	}

	var xp int64
	var level int64
	switch {
	case current >= CollectionsSumLevel1 && current < CollectionsSumLevel2:
		xp = XPCollectionsSumLevel1
		level = 1
	case current >= CollectionsSumLevel2 && current < CollectionsSumLevel3:
		xp = XPCollectionsSumLevel2
		level = 2
	case current >= CollectionsSumLevel3:
		xp = XPCollectionsSumLevel3
		level = 3
	default:
		xp = 0
		level = 0
	}

	levelUpdated := currentLevel > previousLevel
	return xp, level, levelUpdated
}

func GetAchievementCollectionsShipSumData(current int64) (int64, int64, bool) {
	var currentLevel int64
	switch {
	case current >= CollectionsShipSumLevel1 && current < CollectionsShipSumLevel2:
		currentLevel = 1
	case current >= CollectionsShipSumLevel2 && current < CollectionsShipSumLevel3:
		currentLevel = 2
	case current >= CollectionsShipSumLevel3:
		currentLevel = 3
	default:
		currentLevel = 0
	}

	var previousLevel int64
	previous := current - 1
	switch {
	case previous >= CollectionsShipSumLevel1 && previous < CollectionsShipSumLevel2:
		previousLevel = 1
	case previous >= CollectionsShipSumLevel2 && previous < CollectionsShipSumLevel3:
		previousLevel = 2
	case previous >= CollectionsShipSumLevel3:
		previousLevel = 3
	default:
		previousLevel = 0
	}

	var xp int64
	var level int64
	switch {
	case current >= CollectionsShipSumLevel1 && current < CollectionsShipSumLevel2:
		xp = XPCollectionsShipSumLevel1
		level = 1
	case current >= CollectionsShipSumLevel2 && current < CollectionsShipSumLevel3:
		xp = XPCollectionsShipSumLevel2
		level = 2
	case current >= CollectionsShipSumLevel3:
		xp = XPCollectionsShipSumLevel3
		level = 3
	default:
		xp = 0
		level = 0
	}

	levelUpdated := currentLevel > previousLevel
	return xp, level, levelUpdated
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func GetThreshold(achievementCode string, currentLevel int) (nextThreshold int, maxLevel int, exists bool) {
	switch achievementCode {
	case AchieveSubscribers:
		thresholds := []int{SubscribersLevel1, SubscribersLevel2, SubscribersLevel3}
		return getNextThreshold(thresholds, currentLevel)

	case AchieveCollectionItemsAdded:
		thresholds := []int{CollectionItemsAddLevel1, CollectionItemsAddLevel2, CollectionItemsAddLevel3}
		return getNextThreshold(thresholds, currentLevel)

	case AchievePhotosAdded:
		thresholds := []int{PhotosAddLevel1, PhotosAddLevel2, PhotosAddLevel3}
		return getNextThreshold(thresholds, currentLevel)

	case AchieveTagsCreated:
		thresholds := []int{TagsAddLevel1, TagsAddLevel2, TagsAddLevel3}
		return getNextThreshold(thresholds, currentLevel)

	case AchievePostsCreated:
		thresholds := []int{PostsCreateLevel1, PostsCreateLevel2, PostsCreateLevel3}
		return getNextThreshold(thresholds, currentLevel)

	case AchieveCollectionsLikes:
		thresholds := []int{CollectionsLikesLevel1, CollectionsLikesLevel2, CollectionsLikesLevel3}
		return getNextThreshold(thresholds, currentLevel)

	case AchieveCollectionItemsLikes:
		thresholds := []int{CollectionItemsLikesLevel1, CollectionItemsLikesLevel2, CollectionItemsLikesLevel3}
		return getNextThreshold(thresholds, currentLevel)

	case AchievePhotosLikes:
		thresholds := []int{PhotosLikesLevel1, PhotosLikesLevel2, PhotosLikesLevel3}
		return getNextThreshold(thresholds, currentLevel)

	case AchievePostsReactions:
		thresholds := []int{PostsReactionsLevel1, PostsReactionsLevel2, PostsReactionsLevel3}
		return getNextThreshold(thresholds, currentLevel)

	case AchieveCollectionsSum:
		thresholds := []int{CollectionsSumLevel1, CollectionsSumLevel2, CollectionsSumLevel3}
		return getNextThreshold(thresholds, currentLevel)

	case AchieveCollectionShipSum:
		thresholds := []int{CollectionsShipSumLevel1, CollectionsShipSumLevel2, CollectionsShipSumLevel3}
		return getNextThreshold(thresholds, currentLevel)

	default:
		return 0, 0, false
	}
}

func getNextThreshold(thresholds []int, currentLevel int) (nextThreshold, maxLevel int, exists bool) {
	maxLevel = len(thresholds)

	if currentLevel >= maxLevel {
		return 0, maxLevel, false
	}

	return thresholds[currentLevel], maxLevel, true
}

func AddAchievement(
	client *achievementsclient.GRPCAchievementsClient,
	code, userLogin string,
	level, xp, value int64,
) error {
	achievement := &microservices.AddOrUpdateAchievementRequest{
		Code:      code,
		UserLogin: userLogin,
		Xp:        xp,
		Level:     level,
		Value:     value,
	}
	_, err := client.AddOrUpdateAchievement(context.Background(), achievement)
	if err != nil {
		return fmt.Errorf("AddAchievement: %w", err)
	}
	return nil
}

func DateZFormat(input string) string {
	parts := strings.Split(input, " ")
	cleanInput := strings.Join(parts[:len(parts)-1], " ")

	layoutInput := "2006-01-02 15:04:05.999 -0700"
	t, err := time.Parse(layoutInput, cleanInput)
	if err != nil {
		panic(err)
	}

	tUTC := t.UTC()
	output := tUTC.Format(time.RFC3339)

	return output
}
