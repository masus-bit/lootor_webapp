package utils

import (
	"cmp"
	"fmt"
	"github.com/google/uuid"
	"lootor/internal/core/models"
	"slices"
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
	fmt.Println(orderByInput)
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
	default:
		slices.SortFunc(collections, func(a, b models.CollectionsResponse) int {
			return cmp.Compare(b.CreatedAt.UnixNano(), a.CreatedAt.UnixNano())

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
