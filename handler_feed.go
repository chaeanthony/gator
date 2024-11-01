package main

import (
	"context"
	"fmt"
	"time"

	"github.com/chaeanthony/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
  if len(cmd.args) != 2 {
    return fmt.Errorf("usage: %s <name> <url>", cmd.name)
  }

  feed, err := s.db.CreateFeed(context.Background(), 
    database.CreateFeedParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.args[0], Url: cmd.args[1], UserID: user.ID})
  if err != nil {
    return fmt.Errorf("failed to create feed: %v", err)
  }
  
  fmt.Println("Feed created: ")
  printFeed(feed)
  fmt.Println()

  feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}
  fmt.Println("Feed followed.")
  fmt.Printf("%s - %s\n", feedFollow.FeedName, feedFollow.UserName)

  return nil
}

func handlerFeeds(s *state, cmd command) error {
  feeds, err := s.db.GetFeeds(context.Background())
  if err != nil {
    return fmt.Errorf("failed to get feeds: %v", err)
  }
  if len(feeds) == 0 {
		return fmt.Errorf("no feeds found")
	}
	
  for _, feed := range feeds {
    fmt.Printf(" - Name: %s\n - URL: %s\n - User: %s\n", feed.Name, feed.Url, feed.Name_2)
    fmt.Println()
  }

  return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", feed.UserID)
  fmt.Printf("* LastFetchedAt: %v\n", feed.LastFetchedAt.Time)
}
