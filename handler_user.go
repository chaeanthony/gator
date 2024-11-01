package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chaeanthony/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
  if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}
  
  ctx := context.Background()
  _, err := s.db.GetUser(ctx, cmd.args[0])
  if err != nil {
    return fmt.Errorf("unable to get user: %w", err)
  }

  if err := s.cfg.SetUser(cmd.args[0]); err != nil {
    return fmt.Errorf("unable to set user: %w", err)
  }
  fmt.Printf("User set to: %s\n", cmd.args[0])

  return nil
}

func handlerRegister(s *state, cmd command) error {
  if len(cmd.args) != 1 {
    return fmt.Errorf("usage: %s <name>", cmd.name)
  }
  
  ctx := context.Background()
  name := cmd.args[0]

  usr, err := s.db.GetUser(ctx, name)
  if usr.ID != uuid.Nil {
    return fmt.Errorf("user %s already exists", name)
  }
  if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
      usr, err = s.db.CreateUser(ctx, database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: name})
      if err != nil {
        return fmt.Errorf("failed to create user. got: %w", err)
      }
      fmt.Printf("User created: %+v\n", usr)
    } else {
      return err
    }
  }

  if err = s.cfg.SetUser(name); err != nil {
    return fmt.Errorf("unable to set user: %w", err)
  }

  return nil
}

func handlerUsers(s *state, cmd command) error {
  users, err := s.db.GetUsers(context.Background())
  if err != nil {
    return fmt.Errorf("failed to get users: %v", err)
  } 
  
  fmt.Println("users:")
  for _, usr := range users {
    if usr.Name == s.cfg.CurrentUserName {
      fmt.Printf(" * %s (current) \n", usr.Name)
    } else {
      fmt.Printf(" * %s\n", usr.Name)
    }
  }
  return nil
}