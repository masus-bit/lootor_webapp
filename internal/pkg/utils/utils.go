package utils

import (
	"cmp"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/structpb"
	"gorm.io/datatypes"
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
		if v.Id == id {
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
			} else {
				return cmp.Compare(b.Name, a.Name)
			}
		})
		break
	case "created":
		slices.SortFunc(collections, func(a, b models.CollectionsResponse) int {
			if order == "asc" {
				return cmp.Compare(a.CreatedAt.UnixNano(), b.CreatedAt.UnixNano())
			} else {
				return cmp.Compare(b.CreatedAt.UnixNano(), a.CreatedAt.UnixNano())
			}
		})
		break
	case "totalPrice":
		slices.SortFunc(collections, func(a, b models.CollectionsResponse) int {
			if order == "asc" {
				return cmp.Compare(a.TotalPrice, b.TotalPrice)
			} else {
				return cmp.Compare(b.TotalPrice, a.TotalPrice)
			}
		})
		break
	case "collectionItemsCount":
		slices.SortFunc(collections, func(a, b models.CollectionsResponse) int {
			if order == "asc" {
				return cmp.Compare(a.CollectionItemsCount, b.CollectionItemsCount)
			} else {
				return cmp.Compare(b.CollectionItemsCount, a.CollectionItemsCount)
			}
		})
		break
	case "likesCount":
		slices.SortFunc(collections, func(a, b models.CollectionsResponse) int {
			if order == "asc" {
				return cmp.Compare(a.LikesCount, b.LikesCount)
			} else {
				return cmp.Compare(b.LikesCount, a.LikesCount)
			}
		})
		break
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
		break
	case "type":
		result.Joins = "LEFT JOIN item_types ON item_types.id = collection_items.item_type_id"
		if order == "desc" {
			result.Order = "item_types.ru_name DESC"
		}
		result.Order = "item_types.ru_name ASC"
		break
	case "name":
		if order == "desc" {
			result.Order = "name DESC"
		}
		result.Order = "name ASC"
		break
	case "rating":
		if order == "desc" {
			result.Order = "rating DESC"
		}
		result.Order = "rating ASC"
		break
	case "purchaseDate":
		if order == "desc" {
			result.Order = "purchase_date DESC"
		}
		result.Order = "purchase_date ASC"
		break
	case "purchasePrice":
		if order == "desc" {
			result.Order = "purchase_price DESC"
		}
		result.Order = "purchase_price ASC"
		break
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
	case "creatorRating":
		finalOrder = "creator_rating"
	default:
		finalOrder = order
	}

	return finalOrder
}
