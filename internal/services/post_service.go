package services

import (
	"context"
	"database/sql"
	"errors"
	"github.com/wanafiq/feed-api/internal/config"
	"github.com/wanafiq/feed-api/internal/middleware"
	"github.com/wanafiq/feed-api/internal/models"
	"github.com/wanafiq/feed-api/internal/repository"
	"github.com/wanafiq/feed-api/internal/types"
	"github.com/wanafiq/feed-api/internal/utils"
	"go.uber.org/zap"
	"time"
)

type PostService struct {
	config   *config.Config
	db       *sql.DB
	logger   *zap.SugaredLogger
	postRepo repository.PostRepository
	tagRepo  repository.TagRepository
	userRepo repository.UserRepository
}

func NewPostService(config *config.Config, db *sql.DB, logger *zap.SugaredLogger, postRepo repository.PostRepository, tagRepo repository.TagRepository, userRepo repository.UserRepository) *PostService {
	return &PostService{
		config:   config,
		db:       db,
		logger:   logger,
		postRepo: postRepo,
		tagRepo:  tagRepo,
		userRepo: userRepo,
	}
}

func (s *PostService) Save(ctx context.Context, userID string, req *types.PostRequest) (*models.Post, error) {
	author, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.logger.Errorw("failed to find user by id", "userID", userID, "error", err.Error())
		return nil, err
	}

	post := &models.Post{
		Title:     req.Title,
		Slug:      utils.GenerateSlug(req.Title),
		Content:   req.Content,
		CreatedAt: time.Now(),
		CreatedBy: author.Email,
		AuthorID:  author.ID,
	}

	if req.Publish {
		now := time.Now()
		post.IsPublished = true
		post.PublishedAt = &now
	}

	err = withTx(ctx, s.db, func(tx *sql.Tx) error {
		if err := s.postRepo.Save(ctx, tx, post); err != nil {
			s.logger.Errorw("failed to save post", "error", err.Error())
			return err
		}

		if err := s.processTags(ctx, tx, post, req.Tags); err != nil {
			return err
		}

		if err := s.postRepo.SavePostUser(ctx, tx, post.ID, post.AuthorID); err != nil {
			s.logger.Errorw("failed to save post user", "error", err.Error())
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	post.Author = *author

	return post, nil
}

func (s *PostService) GetAll(ctx context.Context, filter models.PostFilter) ([]*models.Post, int, error) {
	posts, count, err := s.postRepo.FindAll(ctx, filter)
	if err != nil {
		s.logger.Errorw("failed to find all posts", "error", err.Error())
		return nil, 0, err
	}

	return posts, count, nil
}

func (s *PostService) GetPostByID(ctx context.Context, postID string) (*models.Post, error) {
	post, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		s.logger.Errorw("failed to find post by id", "postID", postID, "error", err.Error())
		return nil, err
	}

	//tags, err := s.tagRepo.FindByPostID(ctx, postID)
	//if err != nil {
	//	return nil, err
	//}
	//post.Tags = tags

	return post, nil
}

func (s *PostService) Update(ctx context.Context, userCtx middleware.UserContext, postID string, req *types.PostRequest) (*models.Post, error) {
	post, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		s.logger.Errorw("failed to find post by id", "postID", postID, "error", err.Error())
		return nil, err
	}

	now := time.Now()

	post.Title = req.Title
	post.Content = req.Content
	post.IsPublished = req.Publish
	post.UpdatedAt = &now
	post.UpdatedBy = &userCtx.Username

	updatedPost, err := s.postRepo.Update(ctx, nil, post)
	if err != nil {
		s.logger.Errorw("failed to update post", "error", err.Error())
		return nil, err
	}

	return updatedPost, nil
}

func (s *PostService) Delete(ctx context.Context, postID string) error {
	return withTx(ctx, s.db, func(tx *sql.Tx) error {
		if err := s.postRepo.DeletePostTag(ctx, tx, postID); err != nil {
			s.logger.Errorw("failed to delete post tag", "error", err.Error())
			return err
		}

		if err := s.postRepo.DeletePostUser(ctx, tx, postID); err != nil {
			s.logger.Errorw("failed to delete post user", "error", err.Error())
			return err
		}

		if err := s.postRepo.Delete(ctx, tx, postID); err != nil {
			s.logger.Errorw("failed to delete post", "error", err.Error())
			return err
		}

		return nil
	})
}

func (s *PostService) processTags(ctx context.Context, tx *sql.Tx, post *models.Post, tagNames []string) error {
	for _, tagName := range tagNames {
		tag, err := s.tagRepo.FindByName(ctx, tagName)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				if err := s.createAndLinkTag(ctx, tx, post, tagName); err != nil {
					return err
				}
				continue
			}
			s.logger.Errorw("failed to find tag by name", "tagName", tagName, "error", err.Error())
			return err
		}

		if err := s.linkExistingTag(ctx, tx, post, tag); err != nil {
			return err
		}
	}
	return nil
}

func (s *PostService) createAndLinkTag(ctx context.Context, tx *sql.Tx, post *models.Post, tagName string) error {
	newTag := models.Tag{Name: tagName}
	if err := s.tagRepo.Save(ctx, tx, &newTag); err != nil {
		s.logger.Errorw("failed to save new tag", "tagName", tagName, "error", err.Error())
		return err
	}

	post.Tags = append(post.Tags, newTag)

	if err := s.postRepo.SavePostTag(ctx, tx, post.ID, newTag.ID); err != nil {
		s.logger.Errorw("failed to save post tag for new tag", "tagID", newTag.ID, "postID", post.ID, "error", err.Error())
		return err
	}

	return nil
}

func (s *PostService) linkExistingTag(ctx context.Context, tx *sql.Tx, post *models.Post, tag *models.Tag) error {
	//post.Tags = append(post.Tags, tag)
	//
	//if err := s.postRepo.SavePostTag(ctx, tx, post.ID, tag.ID); err != nil {
	//	s.logger.Errorw("failed to save post tag for existing tag", "tagID", tag.ID, "postID", post.ID, "error", err.Error())
	//	return err
	//}

	return nil
}
