package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/h24s_24/domain"
)

type UserRepository interface {
	GetUserPosts(userName string) ([]*domain.Post, error)
	GetUserReactionCount(userName string) (made int, get int, err error)
}

type UserHandler struct {
	rr ReactionRepository
	ur UserRepository
}

type GetUserResponse struct {
	UserName         string `json:"user_name"`
	PostCount        int    `json:"post_count"`
	ReactionCount    int    `json:"reaction_count"`
	GetReactionCount int    `json:"get_reaction_count"`
	Posts            []struct {
		ID               string          `json:"id"`
		UserName         string          `json:"user_name"`
		OriginalMessage  string          `json:"original_message"`
		ConvertedMessage string          `json:"converted_message"`
		ParentID         string          `json:"parent_id"`
		RootID           string          `json:"root_id"`
		Reactions        []reactionCount `json:"reactions"`
		MyReactions      []int           `json:"my_reactions"`
		CreatedAt        time.Time       `json:"created_at"`
	} `json:"posts"`
}

func (uh *UserHandler) GetUserHandler(c echo.Context) error {
	ctx := c.Request().Context()
	loginUser, err := getUserName(c)
	if err != nil {
		log.Println("failed to get user name: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user name")
	}

	userName := c.Param("userName")

	userPosts, err := uh.ur.GetUserPosts(userName)
	if err != nil {
		log.Println("failed to get user posts: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user posts")
	}

	userReactionCount, userGetReactionCount, err := uh.ur.GetUserReactionCount(userName)
	if err != nil {
		log.Println("failed to get user reactions: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user reactions")
	}

	userPostIDs := make([]uuid.UUID, 0, len(userPosts))
	for _, post := range userPosts {
		userPostIDs = append(userPostIDs, post.ID)
	}
	postIDReactionsMap, err := uh.rr.GetReactionsByPostIDs(ctx, userPostIDs)
	if err != nil {
		log.Println("failed to get post reactions: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get post reactions")
	}

	res := GetUserResponse{
		UserName:         userName,
		PostCount:        len(userPosts),
		ReactionCount:    userReactionCount,
		GetReactionCount: userGetReactionCount,
	}

	for _, post := range userPosts {
		reactions := postIDReactionsMap[post.ID]
		var reCount []reactionCount

		reCount = make([]reactionCount, 0, len(reactions))
		for _, reaction := range reactions {
			reCount = append(reCount, reactionCount{
				ID:    reaction.ReactionID,
				Count: reaction.Count,
			})
		}

		myReactions, err := uh.rr.GetReactionsByUserName(ctx, post.ID, loginUser)
		if err != nil {
			log.Println("failed to get my reactions: ", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get my reactions")
		}

		myReactionIDs := make([]int, 0, len(myReactions))
		for _, myReaction := range myReactions {
			myReactionIDs = append(myReactionIDs, myReaction.ReactionID)
		}

		res.Posts = append(res.Posts, struct {
			ID               string          `json:"id"`
			UserName         string          `json:"user_name"`
			OriginalMessage  string          `json:"original_message"`
			ConvertedMessage string          `json:"converted_message"`
			ParentID         string          `json:"parent_id"`
			RootID           string          `json:"root_id"`
			Reactions        []reactionCount `json:"reactions"`
			MyReactions      []int           `json:"my_reactions"`
			CreatedAt        time.Time       `json:"created_at"`
		}{
			ID:               post.ID.String(),
			UserName:         post.UserName,
			OriginalMessage:  post.OriginalMessage,
			ConvertedMessage: post.ConvertedMessage,
			ParentID:         post.ParentID.String(),
			RootID:           post.RootID.String(),
			Reactions:        reCount,
			MyReactions:      myReactionIDs,
			CreatedAt:        post.CreatedAt.Local(),
		})
	}

	return c.JSON(http.StatusOK, res)
}

type getMeResponse struct {
	Username string `json:"user_name"`
}

func (uh *UserHandler) GetMeHandler(c echo.Context) error {
	username, err := getUserName(c)
	if err != nil {
		log.Printf("failed to get username: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get username")
	}
	return c.JSON(http.StatusOK, getMeResponse{username})
}
