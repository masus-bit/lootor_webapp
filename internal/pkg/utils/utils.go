package utils

import (
	"github.com/google/uuid"
	"lootor/internal/core/models"
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
