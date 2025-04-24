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
