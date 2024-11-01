package main

import (
	"fmt"
)

type command struct {
  name string
  args []string 
}

type cliCommands struct {
  commands map[string]func(*state, command) error 
}

func (c *cliCommands) register(name string, f func(*state, command) error) {
  c.commands[name] = f 
}

func (c *cliCommands) run(s *state, cmd command) error {
  f, ok := c.commands[cmd.name]
  if !ok {
    return fmt.Errorf("command '%s' not found", cmd.name)
  }

  return f(s,cmd)
}


