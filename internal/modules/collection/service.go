package collection

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

//func (s *Service) Create(dto *models.CollectionCreateRequest) (*models.CreateCollectionResponse, error) {
//	var resultTags []*models.Tags
//	//for _, tag := range dto.Tags {
//	//	tags, _ := s.repo.GetTagByName(tag)
//	//	if tags == nil {
//	//		tags, _ = s.repo.AddTag(tag)
//	//	}
//	//	resultTags = append(resultTags, tags)
//	//}
//	dto.Tags = resultTags
//	collection := s.repo.Create(dto, resultTags)
//	return &models.CreateCollectionResponse{Data: models.ReturnedCollectionDto{Collection: collection}}, nil
//}
