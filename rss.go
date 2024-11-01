package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
  rssFeed := RSSFeed{}

  req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
  if err != nil {
    return &rssFeed, err
  }
  req.Header.Add("User-Agent", "gator")

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    return &rssFeed, fmt.Errorf("failed to fetch feed: %v", err)
  }
  defer resp.Body.Close()

  data, err := io.ReadAll(resp.Body)
  if err != nil {
    return &rssFeed, err
  }

  if err := xml.Unmarshal(data, &rssFeed); err != nil {
    return &rssFeed, err
  }

  rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
  rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)

  for i, item := range rssFeed.Channel.Item {
    item.Title = html.UnescapeString(item.Title)
    item.Description = html.UnescapeString(item.Description)
    rssFeed.Channel.Item[i] = item
  }

  return &rssFeed, nil
}

