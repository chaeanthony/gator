package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/chaeanthony/blog-aggregator/internal/config"
	"github.com/chaeanthony/blog-aggregator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
  db   *database.Queries
  cfg  *config.Config
}

func main() {
  cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

  db, err := sql.Open("postgres", cfg.DBUrl)
  if err != nil {
    log.Fatalf("failed to open database. got: %v", err)
  }
  defer db.Close()
  dbQueries := database.New(db)

  programState := state{cfg: &cfg, db: dbQueries}
  
  cliCommands := cliCommands{commands: make(map[string]func(*state, command) error)}
  cliCommands.register("login", handlerLogin)
  cliCommands.register("register", handlerRegister)
  cliCommands.register("reset", handlerReset)
  cliCommands.register("users", handlerUsers)
  cliCommands.register("agg", handlerAgg)
  cliCommands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
  cliCommands.register("feeds", handlerFeeds)
  cliCommands.register("follow", middlewareLoggedIn(handlerFollow))
  cliCommands.register("following", middlewareLoggedIn(handlerListFollows))
  cliCommands.register("unfollow", middlewareLoggedIn(handlerUnfollow))
  cliCommands.register("browse", middlewareLoggedIn(handlerBrowse))

  args := os.Args
  if len(args) < 2 {
    log.Fatal("Usage: cli <command> [args...]")
    return 
  }

  cliCmd := command{name: args[1], args: args[2:]}
  if err := cliCommands.run(&programState, cliCmd); err != nil {
    log.Fatal(err)
  }
}
