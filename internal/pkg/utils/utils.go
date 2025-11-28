package utils

import (
	"cmp"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/structpb"
	"gorm.io/datatypes"
	"lootor/gen/go/microservices"
	"lootor/internal/core/models"
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

func GetCollectionOrderBy(orderByInput string, collections []models.CollectionsResponse, order string) []models.CollectionsResponse {
	switch orderByInput {
	case "name":
		slices.SortFunc(collections, func(a, b models.CollectionsResponse) int {
			if order == "asc" {
				return cmp.Compare(a.Name, b.Name)
			}
			return cmp.Compare(b.Name, a.Name)

		})
	case "created":
		slices.SortFunc(collections, func(a, b models.CollectionsResponse) int {
			if order == "asc" {
				return cmp.Compare(a.CreatedAt.UnixNano(), b.CreatedAt.UnixNano())
			}
			return cmp.Compare(b.CreatedAt.UnixNano(), a.CreatedAt.UnixNano())

		})
	case "totalPrice":
		slices.SortFunc(collections, func(a, b models.CollectionsResponse) int {
			if order == "asc" {
				return cmp.Compare(a.TotalPrice, b.TotalPrice)
			}
			return cmp.Compare(b.TotalPrice, a.TotalPrice)

		})
	case "collectionItemsCount":
		slices.SortFunc(collections, func(a, b models.CollectionsResponse) int {
			if order == "asc" {
				return cmp.Compare(a.CollectionItemsCount, b.CollectionItemsCount)
			}
			return cmp.Compare(b.CollectionItemsCount, a.CollectionItemsCount)

		})
	case "likesCount":
		slices.SortFunc(collections, func(a, b models.CollectionsResponse) int {
			if order == "asc" {
				return cmp.Compare(a.LikesCount, b.LikesCount)
			}
			return cmp.Compare(b.LikesCount, a.LikesCount)

		})
	default:
		slices.SortFunc(collections, func(a, b models.CollectionsResponse) int {
			return cmp.Compare(b.LikesCount, a.LikesCount)

		})
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

func FormatPost(post *microservices.PostItem, users []models.SubUsers, author *models.Users, reactLength int) (*models.PostDataResponse, error) {
	content := NormalizeContent(post.Content)

	reacts, err := FormatReacts(post, users, reactLength)
	if err != nil {
		return nil, err
	}

	return &models.PostDataResponse{Data: *fillPost(post, content, reacts, author)}, nil
}

func FormatReacts(post *microservices.PostItem, users []models.SubUsers, reactLength int) (*models.ReactResponse, error) {
	var reactionsResult models.ReactResponse
	reacts := map[string][]models.React{}
	var reactUsers []string
	var fullReactUsers = map[string]models.SubUsers{}
	if reactLength != 0 {
		for _, r := range post.Reactions {
			reactUsers = append(reactUsers, r.UserLogin)
		}

		for _, u := range users {
			fullReactUsers[u.Login] = u
		}

		reactionMapping := []struct {
			field *[]models.React
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
			*mapping.field = []models.React{}
		}

		for _, r := range post.GetReactions() {
			reactUUID, err := uuid.Parse(r.Id)
			if err != nil {
				return nil, err
			}
			reacts[r.GetReaction()] = append(reacts[r.GetReaction()], models.React{
				ID:       reactUUID.String(),
				User:     fullReactUsers[r.UserLogin],
				Reaction: models.ReactionType(r.Reaction),
			})
		}

		for _, mapping := range reactionMapping {
			if slice, exists := reacts[mapping.key]; exists && len(slice) > 0 {
				*mapping.field = slice
			}
		}
	}

	return &reactionsResult, nil
}

func fillPost(p *microservices.PostItem, content []byte, reacts *models.ReactResponse, user *models.Users) *models.Posts {
	return &models.Posts{
		ID:   p.GetId(),
		Date: p.GetDate(),
		Author: models.SubUsers{
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

func NormalizeTagsShort(tagsProto []*microservices.TagItemShort) []models.ShortTags {
	tempTags := make([]models.ShortTags, 0, len(tagsProto))
	for _, tag := range tagsProto {
		tempTags = append(tempTags, models.ShortTags{
			ID:        tag.GetId(),
			Name:      tag.GetName(),
			Slug:      tag.GetSlug(),
			PrimaryID: tag.GetPrimaryId(),
			SeriesID:  tag.GetSeriesId(),
			Type:      getType(tag),
		})
	}

	return tempTags
}

func CanLike(likes []string, user, owner string) bool {
	if user == "" || user == owner {
		return true
	}
	return !slices.Contains(likes, user)
}
