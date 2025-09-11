package types

type CommonShortType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CommonShortTypeCollection struct {
	*CommonShortType
	Owner           string `json:"owner"`
	Transliteration string `json:"transliteration"`
}

type CommonShortTypeItem struct {
	*CommonShortType
	Collection     string `json:"collection"`
	Owner          string `json:"owner"`
	CollectionName string `json:"collectionName"`
}

type CommonShortTypePost struct {
	*CommonShortType
	Author string `json:"author"`
}
