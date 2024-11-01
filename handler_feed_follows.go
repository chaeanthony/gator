package main

import (
	"context"
	"fmt"
	"time"

	"github.com/chaeanthony/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.name)
	}
  
	ctx := context.Background()

  feed, err := s.db.GetFeedByUrl(ctx, cmd.args[0])
	if err != nil {
		return fmt.Errorf("failed to get feed: %v", err)
	}
  
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), 
		database.CreateFeedFollowParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID: user.ID,
			FeedID: feed.ID,
		})
	if err != nil {
		return fmt.Errorf("failed to create feed follow: %v", err)
	}

	fmt.Println("Created feed follow:")
	fmt.Printf(" - Feed: %s\n - User: %s\n", feedFollow.FeedName, feedFollow.UserName)

	return nil
}

func handlerListFollows(s *state, cmd command, user database.User) error {
	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return fmt.Errorf("failed to get feed follows: %v", err)
	}
  
	fmt.Println("feed follows: ")
	for _, ff := range feedFollows {
		fmt.Printf(" - %s, %s\n", ff.FeedName, ff.UserName)
	}
	fmt.Println()

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.name)
	}

	url := cmd.args[0]
	if err := s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{UserID: user.ID, Url: url}); err != nil {
		return fmt.Errorf("failed to unfollow: %v", err)
	}

	fmt.Printf("Unfollowed feed '%s'\n", url)
	return nil
}