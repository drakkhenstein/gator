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