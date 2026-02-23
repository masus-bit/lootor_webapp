package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"lootor/gen/go/microservices"
	"lootor/internal/core/dto"
	"lootor/internal/core/models"
	"lootor/internal/core/repositories"
	"lootor/internal/infrastructure/achievementsclient"
	"lootor/internal/infrastructure/postsclient"
	"lootor/internal/infrastructure/tagsclient"
	"lootor/internal/pkg/utils"
	"strconv"
	"time"
)

type PostsService struct {
	postsClient          *postsclient.GRPCPostsClient
	userRepo             *repositories.UsersRepository
	eventsService        *EventsService
	notificationsService *NotificationsService
	tagsClient           *tagsclient.GRPCTagsClient
	achService           *achievementsclient.GRPCAchievementsClient
}

func NewPostsService(
	postsClient *postsclient.GRPCPostsClient,
	userRepo *repositories.UsersRepository,
	eventsService *EventsService,
	notificationsService *NotificationsService,
	tagsClient *tagsclient.GRPCTagsClient,
	achService *achievementsclient.GRPCAchievementsClient,
) *PostsService {
	return &PostsService{
		postsClient:          postsClient,
		userRepo:             userRepo,
		eventsService:        eventsService,
		notificationsService: notificationsService,
		tagsClient:           tagsClient,
		achService:           achService,
	}
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
	strUint := strconv.FormatUint(post.Data.NumberId, 10)

	translit := utils.Slugify(post.Data.Title) + "_" + strUint
	if !request.IsDraft {
		go func() {
			postsCount, _ := s.postsClient.GetPostsCountByUserLogin(context.Background(), request.Author)

			xp, level := utils.GetAchievementPostsCreateData(postsCount.GetCount())
			err = utils.AddAchievement(
				s.achService,
				utils.AchievePostsCreated,
				request.Author,
				level,
				xp,
				postsCount.GetCount(),
			)

			err = s.userRepo.IncrementExperience(user.Login, utils.PostCreateExp+int(xp))
			if err != nil {
				return
			}
			err = s.userRepo.IncrementSocialScore(user.Login, 5)
			if err != nil {
				return
			}
		}()
		go func() {
			err = s.eventsService.AddEvent(
				user.Login,
				utils.EventActionCreate,
				utils.EventTargetPost,
				post.GetData().GetId(),
				&dto.EventsParams{TargetPostID: translit},
			)
			if err != nil {
				_ = fmt.Errorf("failed to add event: %v", err)
			}
		}()

	}

	_, _ = s.postsClient.UpdatePost(
		ctx, &dto.PostUpdateRequest{
			ID:       post.Data.Id,
			Translit: translit,
			Title:    post.Data.Title,
			IsDraft:  post.Data.IsDraft,
			Content:  content,
		},
	)
	var tags []dto.ShortTags
	if len(request.Tags) > 0 {
		go func() {
			tagsAdded, err := s.tagsClient.AddTagsToEntity(
				context.Background(), &microservices.AddFewTagsToEntityRequest{
					EntityType: "post",
					EntityId:   post.Data.Id,
					TagIds:     request.Tags,
					Author:     request.Author,
					ShowSearch: !request.IsDraft,
				},
			)
			if err != nil {
				fmt.Println(err)
			}
			if tagsAdded != nil {
				for _, tag := range tagsAdded.GetTags() {
					tagUUID, _ := uuid.Parse(tag.GetId())
					eventError := s.eventsService.AddEvent(
						post.GetData().GetAuthor(),
						utils.EventActionAddTag,
						utils.EventTargetTag,
						tag.Name,
						&dto.EventsParams{
							TargetTagID:          tagUUID,
							TargetPostID:         post.GetData().GetId(),
							TagRelatedEntityType: utils.EventTargetPost,
						},
					)
					if eventError != nil {
						log.Default().Print(eventError)
					}
				}
			}

		}()
		respTags, err := s.tagsClient.GetTagsByEntityId(
			context.Background(), &microservices.GetTagsByEntityIdRequest{
				EntityId: post.Data.Id,
			},
		)
		if err != nil {
			return nil, err
		}

		tags = utils.NormalizeTagsShort(respTags.GetTags())
	}

	result := dto.PostDataResponse{
		Data: dto.Posts{
			ID:   post.Data.Id,
			Date: post.Data.Date,
			Author: dto.SubUsers{
				Login:       user.Login,
				AvatarURL:   user.AvatarURL,
				ProfileName: user.ProfileName,
				IsPremium:   user.IsPremium,
			},
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
			IsDraft:        post.Data.IsDraft,
			Views:          0,
			CommentsCount:  0,
			Title:          post.Data.Title,
			Translit:       translit,
			Tags:           tags,
		},
	}

	return &result, nil
}

func (s *PostsService) GetAllPosts(ctx context.Context, order, limit, offset, authUser string) (
	*dto.PostsDataResponse,
	error,
) {
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

func (s *PostsService) GetPostsByUser(
	ctx context.Context,
	login, limit, offset, authUser string,
	isDraft bool,
) (*dto.PostsDataResponse, error) {
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

func (s *PostsService) DeletePost(ctx context.Context, id string, authUser string) (*dto.CommonResponse, error) {
	exists, err := s.postsClient.GetPostById(ctx, id, false)
	if err != nil {
		return nil, err
	}
	result, err := s.postsClient.DeletePost(ctx, id)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepo.GetUserByLogin(authUser)
	if err != nil {
		return nil, err
	}
	if !exists.Data.IsDraft {
		go func() {
			eventError := s.eventsService.AddEvent(
				user.Login,
				utils.EventActionDelete,
				utils.EventTargetPost,
				exists.GetData().GetId(),
				&dto.EventsParams{TargetPostID: exists.Data.Translit},
			)
			if eventError != nil {
				log.Default().Print(eventError)
			}
		}()
		go func() {
			_ = s.notificationsService.DeleteAllNotificationsByTargetID(
				context.Background(),
				id,
			)
		}()
		go func() {
			err = s.userRepo.DecrementExperience(user.Login, int(utils.PostCreateExp+result.ReactCount))
			if err != nil {
				return
			}
			err = s.userRepo.DecrementSocialScore(user.Login, 5)
			if err != nil {
				return
			}
			err = s.userRepo.DecrementPostCount(user.Login)
			if err != nil {
				return
			}
		}()
	}

	go func() {
		_, _ = s.tagsClient.RemoveEntityTags(
			context.Background(), &microservices.RemoveEntityTagsRequest{
				EntityId:   id,
				EntityType: utils.EventTargetPost,
			},
		)
	}()

	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}

func (s *PostsService) GetPostById(ctx context.Context, id string, authUser string) (*dto.PostDataResponse, error) {
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

	author, err := s.userRepo.GetUserByLogin(post.Data.GetAuthor())
	if err != nil {
		return nil, err
	}

	result := s.fillPost(post.Data, normalizedContent, reacts, author, nil)

	return &dto.PostDataResponse{Data: *result}, nil
}

func (s *PostsService) GetPostByTranslit(ctx context.Context, translit string, authUser string) (
	*dto.PostDataResponse,
	error,
) {
	user, err := s.userRepo.GetUserByLogin(authUser)
	var authUserIsPremium bool
	if err != nil {
		authUserIsPremium = false
	} else {
		authUserIsPremium = user.IsPremium
	}
	post, err := s.postsClient.GetPostByTranslit(ctx, translit, authUserIsPremium)
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

	author, err := s.userRepo.GetUserByLogin(post.Data.GetAuthor())
	if err != nil {
		return nil, err
	}

	result := s.fillPost(post.Data, normalizedContent, reacts, author, nil)

	return &dto.PostDataResponse{Data: *result}, nil
}

func (s *PostsService) React(ctx context.Context, req *dto.ReactRequest) (*dto.CommonResponse, error) {
	resp, err := s.postsClient.ReactPost(ctx, req)
	if err != nil {
		return nil, err
	}
	post, err := s.postsClient.GetPostById(ctx, req.PostID, false)
	if err != nil {
		return nil, err
	}
	go func() {
		postsReactCount, _ := s.postsClient.GetReactionsCountByUserLogin(
			context.Background(),
			post.GetData().GetAuthor(),
		)

		xp, level := utils.GetAchievementPostsReactionsData(postsReactCount.GetCount())
		err = utils.AddAchievement(
			s.achService,
			utils.AchievePostsReactions,
			post.GetData().GetAuthor(),
			level,
			xp,
			postsReactCount.GetCount(),
		)
		err = s.userRepo.IncrementExperience(resp.UserLogin, utils.PostReactExt+int(xp))
		if err != nil {
			return
		}
		err = s.userRepo.IncrementSocialScore(resp.UserLogin, 1)
		if err != nil {
			return
		}
	}()
	target := &dto.TargetItem{
		ID:              req.PostID,
		Name:            post.GetData().GetTitle(),
		Transliteration: post.GetData().GetTranslit(),
		TargetType:      "post",
	}
	go func() {
		err = s.notificationsService.SendNotification(
			context.Background(), &dto.NotificationsRequest{
				Login:       post.GetData().GetAuthor(),
				TargetID:    req.PostID,
				SenderLogin: req.UserLogin,
				Type:        utils.NotificationTypePost,
				Action:      utils.NotificationActionReact,
				Date:        time.Now().Format(time.RFC3339),
				OwnerLogin:  post.GetData().GetAuthor(),
			}, target,
		)

	}()
	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil

}

func (s *PostsService) Unreact(ctx context.Context, req *dto.ReactRequest) (*dto.CommonResponse, error) {
	resp, err := s.postsClient.ReactPostDecrement(ctx, req)
	if err != nil {
		return nil, err
	}
	backContext := context.Background()
	go func() {
		_ = s.userRepo.DecrementExperience(resp.UserLogin, utils.PostReactExt)

		_ = s.userRepo.DecrementSocialScore(resp.UserLogin, 1)

		_ = s.notificationsService.DeleteNotification(backContext, req.PostID, req.UserLogin)

	}()
	return &dto.CommonResponse{Data: dto.Resp{Success: true}}, nil
}

func (s *PostsService) UpdatePost(
	ctx context.Context,
	req *dto.PostUpdateRequest,
	authUser string,
) (*dto.PostDataResponse, error) {
	existPost, err := s.postsClient.GetPostById(ctx, req.ID, false)
	if err != nil {
		return nil, err
	}
	if req.Title != "" && existPost.Data.Title != req.Title {
		strUint := strconv.FormatUint(existPost.GetData().GetNumberId(), 10)
		translit := utils.Slugify(req.Title) + "_" + strUint
		req.Translit = translit
	} else if existPost.Data.Title == req.Title {
		req.Translit = existPost.Data.Translit
	}
	post, err := s.postsClient.UpdatePost(ctx, req)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepo.GetUserByLogin(post.Data.Author)
	if err != nil {
		return nil, err
	}

	if existPost.Data.IsDraft && !req.IsDraft {
		go func() {
			postsCount, _ := s.postsClient.GetPostsCountByUserLogin(context.Background(), authUser)

			xp, level := utils.GetAchievementPostsCreateData(postsCount.GetCount())
			err = utils.AddAchievement(
				s.achService,
				utils.AchievePostsCreated,
				authUser,
				level,
				xp,
				postsCount.GetCount(),
			)

			err = s.userRepo.IncrementExperience(user.Login, utils.PostCreateExp+int(xp))
			if err != nil {
				return
			}
			err = s.userRepo.IncrementSocialScore(user.Login, 5)
			if err != nil {
				return
			}
		}()
		go func() {
			err = s.eventsService.AddEvent(
				user.Login,
				utils.EventActionCreate,
				utils.EventTargetPost,
				existPost.GetData().GetId(),
				&dto.EventsParams{TargetPostID: existPost.Data.Translit},
			)
			if err != nil {
				_ = fmt.Errorf("failed to add event: %v", err)
			}
		}()
	} else if !existPost.Data.IsDraft && !req.IsDraft {
		go func() {
			err = s.eventsService.AddEvent(
				user.Login,
				utils.EventActionUpdate,
				utils.EventTargetPost,
				existPost.GetData().GetId(),
				&dto.EventsParams{TargetPostID: existPost.Data.Translit},
			)
			if err != nil {
				_ = fmt.Errorf("failed to add event: %v", err)
			}
		}()
	}

	var tags *microservices.TagsCreateResponse
	tags, err = s.tagsClient.UpdateTagsOfEntity(
		context.Background(), &microservices.UpdateTagsOfEntityRequest{
			EntityType: "post",
			EntityId:   existPost.GetData().GetId(),
			TagIds:     req.Tags,
			Author:     existPost.GetData().GetAuthor(),
		},
	)
	if err != nil {
		return nil, err
	}
	if len(req.Tags) > 0 {
		var IDs []string
		IDs = append(IDs, existPost.GetData().GetId())
		if !req.IsDraft {
			_, err = s.tagsClient.UpdateVisibleLinks(
				context.Background(), &microservices.UpdateVisibleLinksRequest{
					EntityIds: IDs,
					Visible:   true,
				},
			)
			if err != nil {
				return nil, err
			}
		} else {
			_, err = s.tagsClient.UpdateVisibleLinks(
				context.Background(), &microservices.UpdateVisibleLinksRequest{
					EntityIds: IDs,
					Visible:   false,
				},
			)
			if err != nil {
				return nil, err
			}
		}
		if tags != nil {
			for _, tag := range tags.GetTags() {
				tagUUID, _ := uuid.Parse(tag.GetId())
				eventError := s.eventsService.AddEvent(
					post.GetData().GetAuthor(),
					utils.EventActionAddTag,
					utils.EventTargetTag,
					tag.Name,
					&dto.EventsParams{
						TargetTagID:          tagUUID,
						TargetPostID:         req.ID,
						TagRelatedEntityType: utils.EventTargetPost,
					},
				)
				if eventError != nil {
					log.Default().Print(eventError)
				}
			}
		}
	}

	var resultTags []dto.ShortTags
	resultTags = make([]dto.ShortTags, 0)
	if tags != nil {
		for _, tag := range tags.GetTags() {
			resultTags = append(
				resultTags, dto.ShortTags{
					ID:        tag.GetId(),
					Name:      tag.GetName(),
					Slug:      tag.GetSlug(),
					PrimaryID: tag.GetPrimaryId(),
					SeriesID:  tag.GetSeriesId(),
				},
			)
		}
	}

	resultPost := s.fillPost(post.Data, utils.NormalizeContent(post.Data.Content), nil, user, nil)
	resultPost.Tags = resultTags

	return &dto.PostDataResponse{Data: *resultPost}, nil
}

func (s *PostsService) IncrementViews(ctx context.Context, req *dto.IncrementRequestInput) (
	*dto.CommonResponse,
	error,
) {
	resp, err := s.postsClient.IncrementViews(ctx, &dto.IncrementRequest{PostIDs: req.PostIDs})
	if err != nil {
		return nil, err
	}
	return &dto.CommonResponse{Data: dto.Resp{Success: resp.Success}}, nil
}

func (s *PostsService) IncrementCommentsCount(ctx context.Context, id string) (*dto.CommonResponse, error) {
	resp, err := s.postsClient.IncrementCommentsCount(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.CommonResponse{Data: dto.Resp{Success: resp.Success}}, nil
}

func (s *PostsService) DecrementCommentsCount(ctx context.Context, id string) (*dto.CommonResponse, error) {
	resp, err := s.postsClient.DecrementCommentsCount(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.CommonResponse{Data: dto.Resp{Success: resp.Success}}, nil
}

func (s *PostsService) GetCount(ctx context.Context, userLogin string) int64 {
	count, err := s.postsClient.GetCount(ctx, userLogin)
	if err != nil {
		return 0
	}
	return count.GetCount()
}

func (s *PostsService) GetPostsByIDs(
	ctx context.Context,
	ids []string,
	authUserLogin string,
	isPremium bool,
) (*dto.PostsDataResponse, error) {
	resp, _ := s.postsClient.GetPostsByIds(context.Background(), ids, authUserLogin, isPremium)
	var posts []dto.Posts
	var postIds []string
	postIds = append(postIds, ids...)
	tagsMap, err := s.tagsClient.GetTagsByEntityIdsMap(
		context.Background(), &microservices.GetTagsByEntityIdsMapRequest{
			EntityIds: postIds,
		},
	)
	if err != nil {
		return nil, err
	}

	if resp != nil && resp.Data != nil {
		for _, p := range resp.Data {
			var reactUsers []string
			itemTags := tagsMap.GetTags()[p.Id]

			var resultTags = utils.NormalizeTagsShort(itemTags.GetTags())
			for _, r := range p.Reactions {
				reactUsers = append(reactUsers, r.UserLogin)
			}
			reactsLen := len(p.GetReactions())
			var users []dto.SubUsers
			if reactsLen != 0 {
				users, _ = s.userRepo.GetForSubs(reactUsers)
			}
			author, _ := s.userRepo.GetUserByLogin(p.Author)
			postFormatted, err := utils.FormatPost(p, users, author, reactsLen)
			if err != nil {
				continue
			}

			posts = append(
				posts, dto.Posts{
					ID:             postFormatted.Data.ID,
					Title:          postFormatted.Data.Title,
					Author:         postFormatted.Data.Author,
					Translit:       postFormatted.Data.Translit,
					IsDraft:        postFormatted.Data.IsDraft,
					Content:        postFormatted.Data.Content,
					Views:          postFormatted.Data.Views,
					Date:           postFormatted.Data.Date,
					CommentsCount:  postFormatted.Data.CommentsCount,
					HeartCount:     postFormatted.Data.HeartCount,
					FireCount:      postFormatted.Data.FireCount,
					GlassesCount:   postFormatted.Data.GlassesCount,
					LaughCount:     postFormatted.Data.LaughCount,
					TearsCount:     postFormatted.Data.TearsCount,
					PokerFaceCount: postFormatted.Data.PokerFaceCount,
					EyesCount:      postFormatted.Data.EyesCount,
					AngryCount:     postFormatted.Data.AngryCount,
					ShitCount:      postFormatted.Data.ShitCount,
					ClownCount:     postFormatted.Data.ClownCount,
					TotalReactions: postFormatted.Data.TotalReactions,
					Reacted:        postFormatted.Data.Reacted,
					Tags:           resultTags,
				},
			)
		}
	}
	return &dto.PostsDataResponse{Data: posts}, nil
}

func (s *PostsService) formatPosts(posts *microservices.GetAllPostsResponse) (*dto.PostsDataResponse, error) {
	var result []dto.Posts
	var postIds []string
	for _, p := range posts.Data {
		postIds = append(postIds, p.Id)
	}
	tagsMap, err := s.tagsClient.GetTagsByEntityIdsMap(
		context.Background(), &microservices.GetTagsByEntityIdsMapRequest{
			EntityIds: postIds,
		},
	)
	if err != nil {
		return nil, err
	}
	for _, p := range posts.Data {
		content := utils.NormalizeContent(p.Content)
		itemTags := tagsMap.GetTags()[p.Id]

		var resultTags = utils.NormalizeTagsShort(itemTags.GetTags())
		user, err := s.userRepo.GetUserByLogin(p.Author)
		if err != nil {
			return nil, err
		}

		reacts, err := s.formatReacts(p)
		if err != nil {
			return nil, err
		}
		result = append(result, *s.fillPost(p, content, reacts, user, resultTags))
	}

	totalInt, _ := strconv.Atoi(posts.Total)

	return &dto.PostsDataResponse{Data: result, Total: int64(totalInt)}, nil
}

func (s *PostsService) formatReacts(post *microservices.PostItem) (*dto.ReactResponse, error) {
	var reactionsResult dto.ReactResponse
	reacts := map[string][]dto.React{}
	var reactUsers []string
	var fullReactUsers = map[string]dto.SubUsers{}
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

		reactionMapping := []struct {
			field *[]dto.React
			key   string
		}{
			{&reactionsResult.Fire, "fire"},
			{&reactionsResult.Heart, "heart"},
			{&reactionsResult.Glasses, "glasses"},
			{&reactionsResult.Laugh, "laugh"},
			{&reactionsResult.Tears, "tears"},
			{&reactionsResult.PokerFace, "pokerFace"},
			{&reactionsResult.Eyes, "eyes"},
			{&reactionsResult.Angry, "angry"},
			{&reactionsResult.Shit, "shit"},
			{&reactionsResult.Clown, "clown"},
		}

		for _, mapping := range reactionMapping {
			*mapping.field = []dto.React{}
		}

		for _, r := range post.GetReactions() {
			reactUUID, err := uuid.Parse(r.Id)
			if err != nil {
				return nil, err
			}
			reacts[r.GetReaction()] = append(
				reacts[r.GetReaction()], dto.React{
					ID:       reactUUID.String(),
					User:     fullReactUsers[r.UserLogin],
					Reaction: dto.ReactionType(r.Reaction),
				},
			)
		}

		for _, mapping := range reactionMapping {
			if slice, exists := reacts[mapping.key]; exists && len(slice) > 0 {
				*mapping.field = slice
			}
		}
	}

	return &reactionsResult, nil
}

func (s *PostsService) fillPost(
	p *microservices.PostItem,
	content []byte,
	reacts *dto.ReactResponse,
	user *models.Users,
	resTags []dto.ShortTags,
) *dto.Posts {
	var tags []dto.ShortTags
	if resTags != nil {
		tags = resTags
	} else {

		respTags, err := s.tagsClient.GetTagsByEntityId(
			context.Background(), &microservices.GetTagsByEntityIdRequest{
				EntityId: p.Id,
			},
		)
		if err != nil {
			return nil
		}

		tags = utils.NormalizeTagsShort(respTags.GetTags())
	}
	return &dto.Posts{
		ID:   p.GetId(),
		Date: p.GetDate(),
		Author: dto.SubUsers{
			Login:       user.Login,
			AvatarURL:   user.AvatarURL,
			ProfileName: user.ProfileName,
			IsPremium:   user.IsPremium,
		},
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
		Views:          int(p.Views),
		CommentsCount:  int(p.CommentsCount),
		IsDraft:        p.IsDraft,
		Title:          p.Title,
		Translit:       p.Translit,
		Tags:           tags,
		NumberID:       p.NumberId,
	}

}
