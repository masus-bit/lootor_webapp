package services

import (
	"context"
	"github.com/google/uuid"
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
	"lootor/internal/infrastructure/commentsclient"
	"lootor/internal/pkg/dto"
	"lootor/internal/pkg/utils"
	"strconv"
	"time"
)

type CommentsService struct {
	commentsClient *commentsclient.GRPCCommentsClient
	likesClient    *commentsclient.GRPCLikesClient
	userRepo       *repositories.UsersRepository
}

func NewCommentsService(commentsClient *commentsclient.GRPCCommentsClient, likesClient *commentsclient.GRPCLikesClient, userRepo *repositories.UsersRepository) *CommentsService {
	return &CommentsService{
		commentsClient: commentsClient, likesClient: likesClient, userRepo: userRepo}
}

func (s *CommentsService) CreateComment(ctx context.Context, request *dto.CommentsRequest) (*dto.CommentDataResponse, error) {
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

	user, err := s.userRepo.GetUserByLogin(comment.Data.Author)
	if err != nil {
		return nil, err
	}

	userStruct := models.SubUsers{
		Login:       user.Login,
		AvatarUrl:   user.AvatarUrl,
		ProfileName: user.ProfileName,
		IsPremium:   user.IsPremium,
	}

	resultComment := dto.CommentDataResponse{
		Data: dto.CommentsResponse{
			Comments: &dto.Comments{
				Id:            id.String(),
				Date:          comment.Data.Date,
				TargetId:      comment.Data.TargetId,
				Author:        userStruct,
				ParentId:      comment.Data.ParentId,
				LikesCount:    int(comment.Data.LikesCount),
				DislikesCount: int(comment.Data.DislikesCount),
				Content:       content,
				CreatedAt:     createdAtAsTime,
				UpdatedAt:     updatedAtAsTime,
				Likes:         make([]models.SubUsers, 0),
				Dislikes:      make([]models.SubUsers, 0),
			},
			ChildrenComments: make([]dto.ChildrenComments, 0),
			AnswersTotal:     0,
		},
	}

	return &resultComment, nil
}

func (s *CommentsService) GetAllComments(ctx context.Context, targetId, limit, offset string) (*dto.CommentsDataResponse, error) {
	comments, err := s.commentsClient.GetAllComments(ctx, targetId, limit, offset)
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

		userStruct := models.SubUsers{
			Login:       user.Login,
			AvatarUrl:   user.AvatarUrl,
			ProfileName: user.ProfileName,
			IsPremium:   user.IsPremium,
		}

		var resultChildrenComments []dto.ChildrenComments
		for _, c := range f.ChildrenComments {
			contentChildren := utils.NormalizeContent(c.Content)

			createdAtAsTimeCh, _ := time.Parse(time.RFC3339, c.CreatedAt)
			updatedAtAsTimeCh, _ := time.Parse(time.RFC3339, c.UpdatedAt)

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

			userStructCh := models.SubUsers{
				Login:       userCh.Login,
				AvatarUrl:   userCh.AvatarUrl,
				ProfileName: userCh.ProfileName,
				IsPremium:   userCh.IsPremium,
			}
			resultChildrenComments = append(resultChildrenComments, dto.ChildrenComments{
				Comments: &dto.Comments{
					Id:            parsedId.String(),
					Date:          c.Date,
					TargetId:      c.TargetId,
					Author:        userStructCh,
					ParentId:      c.ParentId,
					LikesCount:    int(c.LikesCount),
					DislikesCount: int(c.DislikesCount),
					Content:       contentChildren,
					CreatedAt:     createdAtAsTimeCh,
					UpdatedAt:     updatedAtAsTimeCh,
				},
				Likes:    likesCh,
				Dislikes: dislikesCh,
			})
		}

		resultComments = append(resultComments, dto.CommentsResponse{
			Comments: &dto.Comments{
				Id:            parsedId.String(),
				Date:          f.Date,
				TargetId:      f.TargetId,
				Author:        userStruct,
				ParentId:      f.ParentId,
				LikesCount:    int(f.LikesCount),
				DislikesCount: int(f.DislikesCount),
				Content:       content,
				CreatedAt:     createdAtAsTime,
				UpdatedAt:     updatedAtAsTime,
				Likes:         likes,
				Dislikes:      dislikes,
			},
			ChildrenComments: resultChildrenComments,
			AnswersTotal:     int64(answerTotalInt),
		})
	}

	totalInt, _ := strconv.Atoi(comments.Total)

	return &dto.CommentsDataResponse{Data: resultComments, Total: int64(totalInt)}, nil
}

func safeAtoi64(s string) int64 {
	val, _ := strconv.ParseInt(s, 10, 64)
	return val
}

func (s *CommentsService) DeleteComment(ctx context.Context, id string) (*dto.CommonResponse, error) {
	_, err := s.commentsClient.DeleteComments(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}

func (s *CommentsService) LikeComment(ctx context.Context, id, login string, isLike bool) (*dto.CommonResponse, error) {
	_, err := s.likesClient.Like(ctx, id, login, isLike)
	if err != nil {
		return nil, err
	}
	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}

func (s *CommentsService) DislikeComment(ctx context.Context, id, login string, isLike bool) (*dto.CommonResponse, error) {
	_, err := s.likesClient.Dislike(ctx, id, login, isLike)
	if err != nil {
		return nil, err
	}
	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}

func (s *CommentsService) LoadAnswers(ctx context.Context, id, limit, offset string) (*dto.AnswersDataResponse, error) {

	limitInt, _ := strconv.Atoi(limit)
	offsetInt, _ := strconv.Atoi(offset)

	request := &dto.AnswersRequest{
		Id:     id,
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

		userStruct := models.SubUsers{
			Login:       user.Login,
			AvatarUrl:   user.AvatarUrl,
			ProfileName: user.ProfileName,
			IsPremium:   user.IsPremium,
		}

		result = append(result, dto.AnswersItem{
			Comments: &dto.Comments{
				Id:            f.Id,
				Date:          f.Date,
				TargetId:      f.TargetId,
				Author:        userStruct,
				ParentId:      f.ParentId,
				LikesCount:    int(f.LikesCount),
				DislikesCount: int(f.DislikesCount),
				Content:       utils.NormalizeContent(f.Content),
				CreatedAt:     createdAt,
				UpdatedAt:     updatedAt,
				Likes:         likes,
				Dislikes:      dislikes,
			},
		})
	}
	return &dto.AnswersDataResponse{Data: result, Total: safeAtoi64(response.Total)}, nil

}
