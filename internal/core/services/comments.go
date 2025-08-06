package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"lootor/gen/go/microservices"
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

	countInt, _ := strconv.Atoi(comment.Data.LikesCount)

	resultComment := dto.CommentDataResponse{
		Data: dto.CommentsResponse{
			Comments: &dto.Comments{
				Id:         id.String(),
				Date:       comment.Data.Date,
				TargetId:   comment.Data.TargetId,
				Author:     userStruct,
				ParentId:   comment.Data.ParentId,
				LikesCount: countInt,
				Content:    content,
				CreatedAt:  createdAtAsTime,
				UpdatedAt:  updatedAtAsTime,
				Likes:      make([]models.SubUsers, 0),
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
	fmt.Println(comments)
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

		countInt, _ := strconv.Atoi(f.LikesCount)

		user, err := s.userRepo.GetUserByLogin(f.Author)
		if err != nil {
			return nil, err
		}

		likesStrings := make([]string, 0)
		for _, l := range f.Likes {
			likesStrings = append(likesStrings, l.Author)
		}

		likes, _ := s.userRepo.GetForSubs(likesStrings)

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

			countIntCh, _ := strconv.Atoi(c.LikesCount)

			userCh, err := s.userRepo.GetUserByLogin(c.Author)
			if err != nil {
				return nil, err
			}

			likesStringsCh := make([]string, 0)
			for _, l := range c.Likes {
				likesStringsCh = append(likesStringsCh, l.Author)
			}

			likesCh, _ := s.userRepo.GetForSubs(likesStringsCh)

			userStructCh := models.SubUsers{
				Login:       userCh.Login,
				AvatarUrl:   userCh.AvatarUrl,
				ProfileName: userCh.ProfileName,
				IsPremium:   userCh.IsPremium,
			}
			resultChildrenComments = append(resultChildrenComments, dto.ChildrenComments{
				Comments: &dto.Comments{
					Id:         parsedId.String(),
					Date:       c.Date,
					TargetId:   c.TargetId,
					Author:     userStructCh,
					ParentId:   c.ParentId,
					LikesCount: countIntCh,
					Content:    contentChildren,
					CreatedAt:  createdAtAsTimeCh,
					UpdatedAt:  updatedAtAsTimeCh,
				},
				Likes: likesCh,
			})
		}

		resultComments = append(resultComments, dto.CommentsResponse{
			Comments: &dto.Comments{
				Id:         parsedId.String(),
				Date:       f.Date,
				TargetId:   f.TargetId,
				Author:     userStruct,
				ParentId:   f.ParentId,
				LikesCount: countInt,
				Content:    content,
				CreatedAt:  createdAtAsTime,
				UpdatedAt:  updatedAtAsTime,
				Likes:      likes,
			},
			ChildrenComments: resultChildrenComments,
			AnswersTotal:     int64(answerTotalInt),
		})
	}

	totalInt, _ := strconv.Atoi(comments.Total)

	return &dto.CommentsDataResponse{Data: resultComments, Total: int64(totalInt)}, nil
}

func collectUniqueUserLogins(comments []*microservices.CommentsItem) map[string]bool {
	uniqueLogins := make(map[string]bool)
	for _, f := range comments {
		uniqueLogins[f.Author] = true
		for _, l := range f.Likes {
			uniqueLogins[l.Author] = true
		}
		for _, c := range f.ChildrenComments {
			uniqueLogins[c.Author] = true
			for _, l := range c.Likes {
				uniqueLogins[l.Author] = true
			}
		}
	}
	return uniqueLogins
}

func (s *CommentsService) loadUsersData(logins map[string]bool) (map[string]models.SubUsers, error) {
	loginsSlice := make([]string, 0, len(logins))
	for login := range logins {
		loginsSlice = append(loginsSlice, login)
	}

	users, err := s.userRepo.GetUsersByLogins(loginsSlice)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	usersMap := make(map[string]models.SubUsers)
	for login, user := range users {
		usersMap[login] = models.SubUsers{
			Login:       user.Login,
			AvatarUrl:   user.AvatarUrl,
			ProfileName: user.ProfileName,
			IsPremium:   user.IsPremium,
		}
	}
	return usersMap, nil
}

func (s *CommentsService) convertComments(comments []*microservices.CommentsItem, usersMap map[string]models.SubUsers) ([]dto.CommentsResponse, error) {
	var result []dto.CommentsResponse

	for _, f := range comments {
		// Проверка и преобразование ID
		if f.Id == "" {
			return nil, fmt.Errorf("empty comment ID")
		}

		// Получение автора
		author, ok := usersMap[f.Author]
		if !ok {
			return nil, fmt.Errorf("author %s not found", f.Author)
		}

		// Обработка лайков
		commentLikes := make([]models.SubUsers, 0, len(f.Likes))
		for _, l := range f.Likes {
			if user, exists := usersMap[l.Author]; exists {
				commentLikes = append(commentLikes, user)
			}
		}

		// Обработка дочерних комментариев
		children, err := s.convertChildComments(f.ChildrenComments, usersMap, f.Id)
		if err != nil {
			return nil, err
		}

		// Формирование ответа
		createdAt, _ := time.Parse(time.RFC3339, f.CreatedAt)
		updatedAt, _ := time.Parse(time.RFC3339, f.UpdatedAt)
		result = append(result, dto.CommentsResponse{
			Comments: &dto.Comments{
				Id:         f.Id,
				Date:       f.Date,
				TargetId:   f.TargetId,
				Author:     author,
				ParentId:   f.ParentId,
				LikesCount: safeAtoi(f.LikesCount),
				Content:    utils.NormalizeContent(f.Content),
				CreatedAt:  createdAt,
				UpdatedAt:  updatedAt,
				Likes:      commentLikes,
			},
			ChildrenComments: children,
			AnswersTotal:     safeAtoi64(f.AnswersTotal),
		})
	}
	return result, nil
}

func (s *CommentsService) convertChildComments(children []*microservices.ChildrenComments, usersMap map[string]models.SubUsers, parentId string) ([]dto.ChildrenComments, error) {
	var result []dto.ChildrenComments

	for _, c := range children {
		// Проверка автора
		author, ok := usersMap[c.Author]
		if !ok {
			return nil, fmt.Errorf("child comment author %s not found", c.Author)
		}

		// Обработка лайков
		likes := make([]models.SubUsers, 0, len(c.Likes))
		for _, l := range c.Likes {
			if user, exists := usersMap[l.Author]; exists {
				likes = append(likes, user)
			}
		}
		createdAt, _ := time.Parse(time.RFC3339, c.CreatedAt)
		updatedAt, _ := time.Parse(time.RFC3339, c.UpdatedAt)
		result = append(result, dto.ChildrenComments{
			Comments: &dto.Comments{
				Id:         c.Id,
				Date:       c.Date,
				TargetId:   c.TargetId,
				Author:     author,
				ParentId:   parentId, // Используем parentId из параметра
				LikesCount: safeAtoi(c.LikesCount),
				Content:    utils.NormalizeContent(c.Content),
				CreatedAt:  createdAt,
				UpdatedAt:  updatedAt,
				Likes:      likes,
			},
			Likes:      likes,
			LikesCount: safeAtoi(c.LikesCount),
		})
	}
	return result, nil
}

// Вспомогательные функции
func safeAtoi(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}

func safeAtoi64(s string) int64 {
	val, _ := strconv.ParseInt(s, 10, 64)
	return val
}

func safeFormatTime(timeStr string) string {
	if timeStr == "" {
		return ""
	}
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return ""
	}
	return t.Format(time.RFC3339)
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
