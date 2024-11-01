package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chaeanthony/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error { 
  if len(cmd.args) < 1 || len(cmd.args) > 2 {
		return fmt.Errorf("usage: %v <time_between_reqs>. example: %v 1s", cmd.name, cmd.name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	log.Printf("Collecting feeds every %s...", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	ctx := context.Background()

	// TODO: get next followed feed id to fetch in order to create post under correct current user
  feed, err := s.db.GetNextFeedToFetch(ctx) 
	if err != nil {
		log.Println("failed to get next feed to fetch: ", err)
		return
	}

	rssFeed, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", feed.Name, err)
		return
	}

	_, err = s.db.MarkFeedFetched(ctx, feed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		_, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
				ID: uuid.New(), 
				UpdatedAt: time.Now(), 
				Title: createNullString(item.Title), 
				Url: item.Link, 
				Description: createNullString(item.Description),
				PublishedAt: publishedAt,
				FeedID: feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}

func createNullString(str string) (sql.NullString) {
	if str == "" {
		return sql.NullString{Valid: false, String: ""}
	}

	return sql.NullString{Valid: true, String: str}
}