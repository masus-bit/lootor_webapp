package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"lootor/gen/go/microservices"
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
	"lootor/internal/infrastructure/postsclient"
	"lootor/internal/pkg/dto"
	"lootor/internal/pkg/utils"
	"strconv"
)

type PostsService struct {
	postsClient *postsclient.GRPCPostsClient
	userRepo    *repositories.UsersRepository
}

func NewPostsService(postsClient *postsclient.GRPCPostsClient, userRepo *repositories.UsersRepository) *PostsService {
	return &PostsService{
		postsClient: postsClient, userRepo: userRepo}
}

func (s *PostsService) CreatePost(ctx context.Context, request *dto.PostRequest) (*dto.PostDataResponse, error) {
	post, err := s.postsClient.CreatePost(ctx, request)
	if err != nil {
		return nil, err
	}
	content := utils.NormalizeContent(post.Data.Content)
	user, err := s.userRepo.GetUserByLogin(request.Author)
	if err != nil {
		return nil, err
	}

	go func() {
		err = s.userRepo.IncrementExperience(user.Login, utils.PostCreateExp)
		if err != nil {
			return
		}
		err = s.userRepo.IncrementSocialScore(user.Login, 5)
		if err != nil {
			return
		}
	}()

	result := dto.PostDataResponse{
		Data: dto.Posts{
			Id:             post.Data.Id,
			Date:           post.Data.Date,
			Author:         user.Login,
			Content:        content,
			HeartCount:     0,
			FireCount:      0,
			GlassesCount:   0,
			LaughCount:     0,
			TearsCount:     0,
			PokerFaceCount: 0,
			EyesCount:      0,
			AngryCount:     0,
			ShitCount:      0,
			ClownCount:     0,
			TotalReactions: 0,
			Reactions:      &dto.ReactResponse{},
		},
	}

	return &result, nil
}

func (s *PostsService) GetAllPosts(ctx context.Context, order, limit, offset, authUser string) (*dto.PostsDataResponse, error) {
	user, err := s.userRepo.GetUserByLogin(authUser)
	var authUserIsPremium bool
	if err != nil {
		authUserIsPremium = false
	} else {
		authUserIsPremium = user.IsPremium
	}
	posts, err := s.postsClient.GetAllPosts(ctx, order, limit, offset, authUserIsPremium, authUser)
	if err != nil {
		return nil, err
	}

	return s.formatPosts(posts)

}

func (s *PostsService) GetPostsByUser(ctx context.Context, login, limit, offset, authUser string, isDraft bool) (*dto.PostsDataResponse, error) {
	user, err := s.userRepo.GetUserByLogin(authUser)
	var authUserIsPremium bool
	if err != nil {
		authUserIsPremium = false
	} else {
		authUserIsPremium = user.IsPremium
	}
	posts, err := s.postsClient.GetPostsByUser(ctx, limit, offset, login, authUserIsPremium, authUser, isDraft)
	if err != nil {
		return nil, err
	}

	return s.formatPosts(posts)

}

func (s *PostsService) DeletePost(ctx context.Context, id, authUser string) (*dto.CommonResponse, error) {
	_, err := s.postsClient.DeletePost(ctx, id)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepo.GetUserByLogin(authUser)
	if err != nil {
		return nil, err
	}
	go func() {
		err = s.userRepo.DecrementExperience(user.Login, utils.PostCreateExp)
		if err != nil {
			return
		}
		err = s.userRepo.DecrementSocialScore(user.Login, 5)
		if err != nil {
			return
		}
	}()

	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}

func (s *PostsService) GetPostById(ctx context.Context, id, authUser string) (*dto.PostDataResponse, error) {
	user, err := s.userRepo.GetUserByLogin(authUser)
	var authUserIsPremium bool
	if err != nil {
		authUserIsPremium = false
	} else {
		authUserIsPremium = user.IsPremium
	}
	post, err := s.postsClient.GetPostById(ctx, id, authUserIsPremium)
	if err != nil {
		return nil, err
	}

	normalizedContent := utils.NormalizeContent(post.Data.Content)
	var reacts *dto.ReactResponse
	if len(post.Data.GetReactions()) != 0 {
		reacts, err = s.formatReacts(post.Data)
		if err != nil {
			return nil, err
		}
	}

	result := dto.Posts{
		Id:             post.Data.Id,
		Date:           post.Data.Date,
		Author:         post.Data.Author,
		Content:        normalizedContent,
		HeartCount:     int(post.Data.HeartCount),
		FireCount:      int(post.Data.FireCount),
		GlassesCount:   int(post.Data.GlassesCount),
		LaughCount:     int(post.Data.LaughCount),
		TearsCount:     int(post.Data.TearsCount),
		PokerFaceCount: int(post.Data.PokerFaceCount),
		EyesCount:      int(post.Data.EyesCount),
		AngryCount:     int(post.Data.AngryCount),
		ShitCount:      int(post.Data.ShitCount),
		ClownCount:     int(post.Data.ClownCount),
		TotalReactions: int(post.Data.TotalReactions),
		Reactions:      reacts,
		Reacted:        post.Data.Reacted,
	}

	return &dto.PostDataResponse{Data: result}, nil
}

func (s *PostsService) React(ctx context.Context, req *dto.ReactRequest) (*dto.CommonResponse, error) {
	resp, err := s.postsClient.ReactPost(ctx, req)
	if err != nil {
		return nil, err
	}
	go func() {
		err = s.userRepo.IncrementExperience(resp.UserLogin, utils.PostReactExt)
		if err != nil {
			return
		}
		err = s.userRepo.IncrementSocialScore(resp.UserLogin, 1)
		if err != nil {
			return
		}
	}()
	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil

}

func (s *PostsService) Unreact(ctx context.Context, req *dto.ReactRequest) (*dto.CommonResponse, error) {
	resp, err := s.postsClient.ReactPostDecrement(ctx, req)
	if err != nil {
		return nil, err
	}
	go func() {
		err = s.userRepo.DecrementExperience(resp.UserLogin, utils.PostReactExt)
		if err != nil {
			return
		}
		err = s.userRepo.DecrementSocialScore(resp.UserLogin, 1)
		if err != nil {
			return
		}
	}()
	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}

func (s *PostsService) UpdatePost(ctx context.Context, req *dto.PostUpdateRequest) (*dto.PostDataResponse, error) {
	post, err := s.postsClient.UpdatePost(ctx, req)
	if err != nil {
		return nil, err
	}
	resultPost := dto.Posts{
		Id:             post.Data.Id,
		Date:           post.Data.Date,
		Author:         post.Data.Author,
		Content:        utils.NormalizeContent(post.Data.Content),
		HeartCount:     int(post.Data.HeartCount),
		FireCount:      int(post.Data.FireCount),
		GlassesCount:   int(post.Data.GlassesCount),
		LaughCount:     int(post.Data.LaughCount),
		TearsCount:     int(post.Data.TearsCount),
		PokerFaceCount: int(post.Data.PokerFaceCount),
		EyesCount:      int(post.Data.EyesCount),
		AngryCount:     int(post.Data.AngryCount),
		ShitCount:      int(post.Data.ShitCount),
		ClownCount:     int(post.Data.ClownCount),
		TotalReactions: int(post.Data.TotalReactions),
		Reactions:      &dto.ReactResponse{},
		Reacted:        post.Data.Reacted,
	}
	return &dto.PostDataResponse{Data: resultPost}, nil
}

func (s *PostsService) formatPosts(posts *microservices.GetAllPostsResponse) (*dto.PostsDataResponse, error) {
	var result []dto.Posts
	for _, p := range posts.Data {
		content := utils.NormalizeContent(p.Content)

		user, err := s.userRepo.GetUserByLogin(p.Author)
		if err != nil {
			return nil, err
		}

		reacts, err := s.formatReacts(p)
		if err != nil {
			return nil, err
		}
		fmt.Println(p.GetTearsCount())
		result = append(result, dto.Posts{
			Id:             p.GetId(),
			Date:           p.GetDate(),
			Author:         user.Login,
			Content:        content,
			HeartCount:     int(p.GetHeartCount()),
			FireCount:      int(p.GetFireCount()),
			GlassesCount:   int(p.GetGlassesCount()),
			LaughCount:     int(p.GetLaughCount()),
			TearsCount:     int(p.GetTearsCount()),
			PokerFaceCount: int(p.GetPokerFaceCount()),
			EyesCount:      int(p.GetEyesCount()),
			AngryCount:     int(p.GetAngryCount()),
			ShitCount:      int(p.GetShitCount()),
			ClownCount:     int(p.GetClownCount()),
			TotalReactions: int(p.GetTotalReactions()),
			Reactions:      reacts,
			Reacted:        p.GetReacted(),
		})
	}

	totalInt, _ := strconv.Atoi(posts.Total)

	return &dto.PostsDataResponse{Data: result, Total: int64(totalInt)}, nil
}

func (s *PostsService) formatReacts(post *microservices.PostItem) (*dto.ReactResponse, error) {
	var reactionsResult dto.ReactResponse
	reacts := map[string][]dto.React{}
	var reactUsers []string
	var fullReactUsers = map[string]models.SubUsers{}
	if len(post.GetReactions()) != 0 {
		for _, r := range post.Reactions {
			reactUsers = append(reactUsers, r.UserLogin)
		}
		users, err := s.userRepo.GetForSubs(reactUsers)
		if err != nil {
			return nil, err
		}

		for _, u := range users {
			fullReactUsers[u.Login] = u
		}

		for _, r := range post.GetReactions() {
			reactUUID, err := uuid.Parse(r.Id)
			if err != nil {
				return nil, err
			}
			reacts[r.GetReaction()] = append(reacts[r.GetReaction()], dto.React{
				Id:       reactUUID.String(),
				User:     fullReactUsers[r.UserLogin],
				Reaction: dto.ReactionType(r.Reaction),
			})
		}

		reactionsResult.Fire = reacts["fire"]
		reactionsResult.Heart = reacts["heart"]
		reactionsResult.Glasses = reacts["glasses"]
		reactionsResult.Laugh = reacts["laugh"]
		reactionsResult.Tears = reacts["tears"]
		reactionsResult.PokerFace = reacts["pokerFace"]
		reactionsResult.Eyes = reacts["eyes"]
		reactionsResult.Angry = reacts["angry"]
		reactionsResult.Shit = reacts["shit"]
		reactionsResult.Clown = reacts["clown"]
	}

	return &reactionsResult, nil
}
