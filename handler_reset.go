package main

import (
	"context"
	"fmt"
)


func handlerReset(s *state, cmd command) error {
  ctx := context.Background()
  if err := s.db.DeleteAllUsers(ctx); err != nil {
    return fmt.Errorf("failed to reset: %v", err)
  }
  fmt.Println("reset successful")
  return nil
}



