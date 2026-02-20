package services

import (
	"context"
	"github.com/google/uuid"
	"lootor/internal/core/dto"
	"lootor/internal/core/repositories"
	"lootor/internal/infrastructure/commentsclient"
	"lootor/internal/pkg/utils"
	"strconv"
	"time"
)

type CommentsService struct {
	commentsClient       *commentsclient.GRPCCommentsClient
	likesClient          *commentsclient.GRPCLikesClient
	postsService         PostsService
	userRepo             *repositories.UsersRepository
	notificationsService *NotificationsService
	collectionRepo       *repositories.CollectionsRepository
	ciRepo               *repositories.CiRepository
	photosService        *PhotosService
}

func NewCommentsService(
	commentsClient *commentsclient.GRPCCommentsClient,
	likesClient *commentsclient.GRPCLikesClient,
	userRepo *repositories.UsersRepository,
	notificationsService *NotificationsService,
	collectionRepo *repositories.CollectionsRepository,
	ciRepo *repositories.CiRepository,
	postService *PostsService,
	photosService *PhotosService,
) *CommentsService {
	return &CommentsService{
		commentsClient:       commentsClient,
		likesClient:          likesClient,
		userRepo:             userRepo,
		notificationsService: notificationsService,
		collectionRepo:       collectionRepo,
		ciRepo:               ciRepo,
		postsService:         *postService,
		photosService:        photosService,
	}
}

func (s *CommentsService) CreateComment(ctx context.Context, request *dto.CommentsRequest) (
	*dto.CommentDataResponse,
	error,
) {
	comment, err := s.commentsClient.CreateComment(ctx, request)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(comment.Data.Id)
	if err != nil {
		return nil, err
	}
	content := utils.NormalizeContent(comment.Data.Content)

	createdAtAsTime, _ := time.Parse(time.RFC3339, comment.Data.CreatedAt)
	updatedAtAsTime, _ := time.Parse(time.RFC3339, comment.Data.UpdatedAt)
	deletedAtAsTime, _ := time.Parse(time.RFC3339, comment.Data.DeletedAt)
	deletedAtStr := ""
	if !deletedAtAsTime.IsZero() {
		deletedAtStr = deletedAtAsTime.Format(time.RFC3339)
	}
	user, err := s.userRepo.GetUserByLogin(request.Author)
	if err != nil {
		return nil, err
	}

	userStruct := dto.SubUsers{
		Login:       user.Login,
		AvatarURL:   user.AvatarURL,
		ProfileName: user.ProfileName,
		IsPremium:   user.IsPremium,
	}

	resultComment := dto.CommentDataResponse{
		Data: dto.CommentsResponse{
			Comments: &dto.Comments{
				ID:            id.String(),
				Date:          comment.Data.Date,
				TargetID:      request.TargetID,
				Author:        userStruct,
				ParentID:      comment.Data.ParentId,
				LikesCount:    int(comment.Data.LikesCount),
				DislikesCount: int(comment.Data.DislikesCount),
				Content:       content,
				CreatedAt:     createdAtAsTime,
				UpdatedAt:     updatedAtAsTime,
				Likes:         make([]dto.SubUsers, 0),
				Dislikes:      make([]dto.SubUsers, 0),
				DeletedAt:     deletedAtStr,
				EntityType:    request.EntityType,
			},
			ChildrenComments: make([]dto.ChildrenComments, 0),
			AnswersTotal:     0,
		},
	}

	var targetUserLogin string
	var ownerLogin string
	var targetReq *dto.TargetItem

	switch request.EntityType {
	case "collection":
		collection, err := s.collectionRepo.GetByIdWithoutCollectionItems(request.TargetID)
		if err == nil {
			if *request.TargetUserLogin != "" {
				targetUserLogin = *request.TargetUserLogin
			} else {
				targetUserLogin = collection.UserLogin
			}
			err = s.collectionRepo.IncrementCommentsCount(request.TargetID, 1)
			if err != nil {
				return nil, err
			}
			ownerLogin = collection.UserLogin

			targetReq = &dto.TargetItem{
				ID:              collection.ID.String(),
				Name:            collection.Name,
				Transliteration: collection.Transliteration,
				TargetType:      "collection",
			}
		}
	case "collectionItem":
		ci, err := s.ciRepo.GetCIByID(request.TargetID)
		if err == nil {
			if *request.TargetUserLogin != "" {
				targetUserLogin = *request.TargetUserLogin
			} else {
				targetUserLogin = ci.UserLogin
			}
			err = s.ciRepo.IncrementCommentsCount(request.TargetID, 1)
			if err != nil {
				return nil, err
			}
			ownerLogin = ci.UserLogin

			targetReq = &dto.TargetItem{
				ID:               ci.ID.String(),
				Name:             ci.Name,
				Transliteration:  ci.Collections[0].Transliteration,
				TargetType:       "collectionItem",
				TargetParentName: ci.Collections[0].Name,
			}
		}
	case "post":
		post, err := s.postsService.GetPostById(ctx, request.TargetID, request.Author)

		if err == nil {
			if *request.TargetUserLogin != "" {
				targetUserLogin = *request.TargetUserLogin
			} else {
				targetUserLogin = post.Data.Author.Login
			}
			_, err = s.postsService.IncrementCommentsCount(ctx, request.TargetID)
			if err != nil {
				return nil, err
			}
			ownerLogin = post.Data.Author.Login

			targetReq = &dto.TargetItem{
				ID:              post.Data.Translit,
				Name:            post.Data.Title,
				Transliteration: post.Data.Translit,
				TargetType:      "post",
			}
		}
	case "photo":
		photo, err := s.photosService.FindOneByID(ctx, request.TargetID, request.Author)
		if err == nil {
			if *request.TargetUserLogin != "" {
				targetUserLogin = *request.TargetUserLogin
			} else {
				targetUserLogin = photo.Data.Author.Login
			}
			_, err = s.photosService.IncrementCommentsCount(ctx, request.TargetID)
			if err != nil {
				return nil, err
			}
			collectionDb, err := s.collectionRepo.GetByIdWithoutCollectionItems(photo.Data.CollectionID)
			if err != nil {
				return nil, err
			}
			ownerLogin = photo.Data.Author.Login

			targetReq = &dto.TargetItem{
				ID:               request.TargetID,
				Name:             photo.Data.Path,
				Transliteration:  collectionDb.Transliteration,
				TargetType:       "photo",
				TargetParentName: collectionDb.Name,
			}
		}
	}

	var typeComment string

	if *request.ParentID != "" {
		typeComment = utils.NotificationTypeAnswer
	} else {
		typeComment = utils.NotificationTypeComment
	}

	go func() {
		err = s.notificationsService.SendNotification(
			context.Background(), &dto.NotificationsRequest{
				Login:       targetUserLogin,
				TargetID:    request.TargetID,
				SenderLogin: user.Login,
				Type:        typeComment,
				Action:      utils.NotificationActionComment,
				Date:        time.Now().Format(time.RFC3339),
				OwnerLogin:  ownerLogin,
			}, targetReq,
		)
	}()

	err = s.userRepo.IncrementExperience(request.Author, 1)
	if err != nil {
		return nil, err
	}

	return &resultComment, nil
}

func (s *CommentsService) GetAllComments(
	ctx context.Context,
	targetId, limit, offset, entityType string,
) (*dto.CommentsDataResponse, error) {
	comments, err := s.commentsClient.GetAllComments(ctx, targetId, limit, offset, entityType)
	if err != nil {
		return nil, err
	}
	var resultComments []dto.CommentsResponse
	for _, f := range comments.Data {
		parsedId, err := uuid.Parse(f.Id)
		if err != nil {
			return nil, err
		}
		content := utils.NormalizeContent(f.Content)

		createdAtAsTime, _ := time.Parse(time.RFC3339, f.CreatedAt)
		updatedAtAsTime, _ := time.Parse(time.RFC3339, f.UpdatedAt)
		deletedAtAsTime, _ := time.Parse(time.RFC3339, f.DeletedAt)
		deletedAtStr := ""
		if !deletedAtAsTime.IsZero() {
			deletedAtStr = deletedAtAsTime.Format(time.RFC3339)
			content = nil
		}
		answerTotalInt, _ := strconv.Atoi(f.AnswersTotal)

		user, err := s.userRepo.GetUserByLogin(f.Author)
		if err != nil {
			return nil, err
		}

		likesStrings := make([]string, 0)
		for _, l := range f.Likes {
			likesStrings = append(likesStrings, l.Author)
		}

		dislikesStrings := make([]string, 0)
		for _, l := range f.Dislikes {
			dislikesStrings = append(dislikesStrings, l.Author)
		}

		likes, _ := s.userRepo.GetForSubs(likesStrings)
		dislikes, _ := s.userRepo.GetForSubs(dislikesStrings)

		userStruct := dto.SubUsers{
			Login:       user.Login,
			AvatarURL:   user.AvatarURL,
			ProfileName: user.ProfileName,
			IsPremium:   user.IsPremium,
		}

		var resultChildrenComments []dto.ChildrenComments
		for _, c := range f.ChildrenComments {
			contentChildren := utils.NormalizeContent(c.Content)
			createdAtAsTimeCh, _ := time.Parse(time.RFC3339, c.CreatedAt)
			updatedAtAsTimeCh, _ := time.Parse(time.RFC3339, c.UpdatedAt)
			deletedAtAsTimeCh, _ := time.Parse(time.RFC3339, c.DeletedAt)
			deletedAtStrCh := ""
			if !deletedAtAsTimeCh.IsZero() {
				deletedAtStrCh = deletedAtAsTimeCh.Format(time.RFC3339)
				contentChildren = nil
			}

			userCh, err := s.userRepo.GetUserByLogin(c.Author)
			if err != nil {
				return nil, err
			}
			likesStringsCh := make([]string, 0)
			for _, l := range c.Likes {
				likesStringsCh = append(likesStringsCh, l.Author)
			}

			dislikesStringsCh := make([]string, 0)
			for _, l := range c.Dislikes {
				dislikesStringsCh = append(dislikesStringsCh, l.Author)
			}

			likesCh, _ := s.userRepo.GetForSubs(likesStringsCh)
			dislikesCh, _ := s.userRepo.GetForSubs(dislikesStringsCh)

			userStructCh := dto.SubUsers{
				Login:       userCh.Login,
				AvatarURL:   userCh.AvatarURL,
				ProfileName: userCh.ProfileName,
				IsPremium:   userCh.IsPremium,
			}
			resultChildrenComments = append(
				resultChildrenComments, dto.ChildrenComments{
					Comments: &dto.Comments{
						ID:         c.Id,
						Date:       c.Date,
						TargetID:   targetId,
						Author:     userStructCh,
						ParentID:   c.ParentId,
						Content:    contentChildren,
						CreatedAt:  createdAtAsTimeCh,
						UpdatedAt:  updatedAtAsTimeCh,
						DeletedAt:  deletedAtStrCh,
						EntityType: c.EntityType,
					},
					Likes:         likesCh,
					Dislikes:      dislikesCh,
					LikesCount:    int(c.LikesCount),
					DislikesCount: int(c.DislikesCount),
				},
			)
		}

		resultComments = append(
			resultComments, dto.CommentsResponse{
				Comments: &dto.Comments{
					ID:            parsedId.String(),
					Date:          f.Date,
					TargetID:      targetId,
					Author:        userStruct,
					ParentID:      f.ParentId,
					LikesCount:    int(f.LikesCount),
					DislikesCount: int(f.DislikesCount),
					Content:       content,
					CreatedAt:     createdAtAsTime,
					UpdatedAt:     updatedAtAsTime,
					Likes:         likes,
					Dislikes:      dislikes,
					DeletedAt:     deletedAtStr,
					EntityType:    f.EntityType,
				},
				ChildrenComments: resultChildrenComments,
				AnswersTotal:     int64(answerTotalInt),
			},
		)
	}

	totalInt, _ := strconv.Atoi(comments.Total)

	return &dto.CommentsDataResponse{Data: resultComments, Total: int64(totalInt)}, nil
}

func safeAtoi64(s string) int64 {
	val, _ := strconv.ParseInt(s, 10, 64)
	return val
}

func (s *CommentsService) DeleteComment(ctx context.Context, id, targetId, authUserLogin string) (
	*dto.CommonResponse,
	error,
) {
	_, err := s.collectionRepo.GetByIdWithoutCollectionItems(targetId)

	if err == nil {
		err = s.collectionRepo.DecrementCommentsCount(targetId, 1)
		if err != nil {
			return nil, err
		}
	}

	_, err = s.ciRepo.GetCIByID(targetId)

	if err == nil {
		err = s.ciRepo.DecrementCommentsCount(targetId, 1)
		if err != nil {
			return nil, err
		}
	}

	_, err = s.postsService.GetPostById(ctx, targetId, authUserLogin)

	if err == nil {
		_, err = s.postsService.DecrementCommentsCount(ctx, targetId)
		if err != nil {
			return nil, err
		}
	}

	_, err = s.photosService.FindOneByID(ctx, targetId, authUserLogin)
	if err == nil {
		_, err = s.photosService.DecrementCommentsCount(ctx, targetId)
		if err != nil {
			return nil, err
		}
	}
	_, err = s.commentsClient.DeleteComments(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}

func (s *CommentsService) LikeComment(ctx context.Context, id, login string, isLike bool) (*dto.CommonResponse, error) {
	ok, err := s.likesClient.Like(ctx, id, login, isLike)
	if err != nil {
		return nil, err
	}
	if isLike {
		err = s.userRepo.IncrementSocialScore(ok.GetTargetUser(), 1)
		if err != nil {
			return nil, err
		}
	} else {
		err = s.userRepo.DecrementSocialScore(ok.GetTargetUser(), 1)
		if err != nil {
			return nil, err
		}
	}
	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}

func (s *CommentsService) DislikeComment(ctx context.Context, id, login string, isLike bool) (
	*dto.CommonResponse,
	error,
) {
	ok, err := s.likesClient.Dislike(ctx, id, login, isLike)
	if err != nil {
		return nil, err
	}
	if isLike {
		err = s.userRepo.DecrementSocialScore(ok.GetTargetUser(), 1)
		if err != nil {
			return nil, err
		}
	} else {
		err = s.userRepo.IncrementSocialScore(ok.GetTargetUser(), 1)
		if err != nil {
			return nil, err
		}
	}
	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}

func (s *CommentsService) LoadAnswers(ctx context.Context, id, limit, offset string) (*dto.AnswersDataResponse, error) {

	limitInt, _ := strconv.Atoi(limit)
	offsetInt, _ := strconv.Atoi(offset)

	request := &dto.AnswersRequest{
		ID:     id,
		Limit:  limitInt,
		Offset: offsetInt,
	}

	response, err := s.commentsClient.LoadAnswers(ctx, request)
	if err != nil {
		return nil, err
	}
	var result []dto.AnswersItem
	for _, f := range response.Data {

		createdAt, _ := time.Parse(time.RFC3339, f.CreatedAt)
		updatedAt, _ := time.Parse(time.RFC3339, f.UpdatedAt)

		user, err := s.userRepo.GetUserByLogin(f.Author)
		if err != nil {
			return nil, err
		}

		likesStrings := make([]string, 0)
		for _, l := range f.Likes {
			likesStrings = append(likesStrings, l.Author)
		}

		dislikesStrings := make([]string, 0)
		for _, l := range f.Dislikes {
			dislikesStrings = append(dislikesStrings, l.Author)
		}

		likes, _ := s.userRepo.GetForSubs(likesStrings)
		dislikes, _ := s.userRepo.GetForSubs(dislikesStrings)

		userStruct := dto.SubUsers{
			Login:       user.Login,
			AvatarURL:   user.AvatarURL,
			ProfileName: user.ProfileName,
			IsPremium:   user.IsPremium,
		}

		result = append(
			result, dto.AnswersItem{
				Comments: &dto.Comments{
					ID:         f.Id,
					Date:       f.Date,
					TargetID:   f.TargetId,
					Author:     userStruct,
					ParentID:   f.ParentId,
					Content:    utils.NormalizeContent(f.Content),
					CreatedAt:  createdAt,
					UpdatedAt:  updatedAt,
					EntityType: f.EntityType,
				},
				LikesCount:    int(f.LikesCount),
				DislikesCount: int(f.DislikesCount),
				Likes:         likes,
				Dislikes:      dislikes,
			},
		)
	}
	return &dto.AnswersDataResponse{Data: result, Total: safeAtoi64(response.Total)}, nil

}
