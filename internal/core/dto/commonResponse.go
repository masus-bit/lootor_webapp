package dto

type Resp struct {
	Success bool `json:"success"`
}

type CommonResponse struct {
	Data Resp `json:"data"`
}

type FeedbackDto struct {
}
