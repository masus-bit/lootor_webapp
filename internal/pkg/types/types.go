package types

type CommonShortType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CommonShortTypeCollection struct {
	*CommonShortType
	Owner string `json:"owner"`
}

type CommonShortTypeItem struct {
	*CommonShortType
	Collection string `json:"collection"`
	Owner      string `json:"owner"`
}
