package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/traP-jp/h24s_24/domain"
)

type UserRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (ur *UserRepository) GetUserPosts(userName string) ([]*domain.Post, error) {
	var userPosts []post
	err := ur.DB.Select(&userPosts, "SELECT * FROM posts WHERE user_name = ? ORDER BY created_at DESC", userName)
	if err != nil {
		return nil, err
	}

	posts := make([]*domain.Post, 0, len(userPosts))
	for _, post := range userPosts {
		posts = append(posts, &domain.Post{
			ID:               post.ID,
			UserName:         post.UserName,
			OriginalMessage:  post.OriginalMessage,
			ConvertedMessage: post.ConvertedMessage,
			ParentID:         post.ParentID,
			RootID:           post.RootID,
			CreatedAt:        post.CreatedAt,
		})
	}

	return posts, nil
}

func (ur *UserRepository) GetUserReactionCount(userName string) (int, int, error) {
	var madeCount int
	err := ur.DB.Get(&madeCount, "SELECT COUNT(*) FROM posts_reactions WHERE user_name = ?", userName)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get reaction count by user: %w", err)
	}

	var getCount int
	err = ur.DB.Get(&getCount, "SELECT COUNT(*) FROM posts_reactions AS pr JOIN posts AS p ON pr.post_id = p.id WHERE p.user_name = ?", userName)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get reaction count earned by user: %w", err)
	}

	return madeCount, getCount, nil
}
