package main

import (
	"context"
	"fmt"
	"time"

	"github.com/drakkhenstein/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <username>", cmd.Name)
	}

	name := cmd.Args[0]
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
    return fmt.Errorf("couldn't find user: %w", err)
	}

	
	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("error setting user: %w", err)
	}
	fmt.Println("User switched successfully")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <username>", cmd.Name)
	}
	name := cmd.Args[0]
	userID := uuid.New()
	_, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:   userID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: name,
	})
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	err = s.cfg.SetUser(name)
	if err != nil {
    	return fmt.Errorf("couldn't set current user: %w", err)
	}
	fmt.Println("User created successfully")
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error deleting users: %w", err)
	}
	err = s.cfg.SetUser("")
	if err != nil {
		return fmt.Errorf("couldn't reset current user: %w", err)
	}
	fmt.Println("Users reset successfully")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error getting users: %w", err)
	}
	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("%s (current)\n", user.Name)
		} else {
			fmt.Printf("%s\n", user.Name)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
    return fmt.Errorf("usage: agg <time_between_reqs>")
	}
	time_between_reqs := cmd.Args[0]
	duration, err := time.ParseDuration(time_between_reqs)
	if err != nil {
		return fmt.Errorf("error parsing duration: %w", err)
	}
	ticker := time.NewTicker(duration)
	fmt.Printf("Collecting feeds every %s\n", duration)
	for ; ; <-ticker.C {
    scrapeFeeds(s)
	}
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <feed name> <feed url>", cmd.Name)
	}
	name := cmd.Args[0]
	url := cmd.Args[1]
	feed, err := s.db.AddFeed(context.Background(), database.AddFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: name,
		Url: url,
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("error adding feed: %w", err)
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error following feed: %w", err)
	}
	fmt.Println("Feed added successfully")

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds: %w", err)
	}
	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("error getting user: %w", err)
		}
		fmt.Printf("All Feeds: %s %s (User: %s)\n", feed.Name, feed.Url, user.Name)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feed url>", cmd.Name)
	}

	url := cmd.Args[0]
	feedToFollow, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error getting feed to follow: %w", err)
	}
	
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: feedToFollow.ID,
	})
	if err != nil {
		return fmt.Errorf("error following feed: %w", err)
	}
	fmt.Printf("Now following feed: %s\n", feedToFollow.Name)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	followingFeeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error getting following feeds: %w", err)
	}
	if len(followingFeeds) == 0 {
		fmt.Println("You are not following any feeds.")
		return nil
	}
	fmt.Printf("Feed follows for user %s:\n", user.Name)
	for _, feed := range followingFeeds {
		fmt.Printf("Following Feed: %s\n", feed.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feed url>", cmd.Name)
	}

	url := cmd.Args[0]
	feedToUnfollow, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error getting feed to unfollow: %w", err)
	}
	
	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feedToUnfollow.ID,
	})
	if err != nil {
		return fmt.Errorf("error unfollowing feed: %w", err)
	}
	fmt.Printf("Unfollowed feed: %s\n", feedToUnfollow.Name)
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.Args) > 0 {
		if n, err := fmt.Sscanf(cmd.Args[0], "%d", &limit); err != nil || n != 1 {
			return fmt.Errorf("invalid limit: %s", cmd.Args[0])
		}
	}
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("error getting posts for user: %w", err)
	}
	if len(posts) == 0 {
		fmt.Println("No posts to display. Follow some feeds to see posts here.")
		return nil
	}
	fmt.Printf("Posts for user %s:\n", user.Name)
	for _, post := range posts {
		fmt.Printf("Title: %s\nURL: %s\nDescription: %s\nPublished At: %s\nFeed Name: %s\n\n", post.Title, post.Url, post.Description.String, post.PublishedAt.Time.Format(time.RFC1123), post.FeedName)
	}
	return nil
}