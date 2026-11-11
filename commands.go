package main

import (
	"fmt"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	command, ok := c.handlers[cmd.Name]
	if !ok {
		return fmt.Errorf("command not found")
	}

	err := command(s, cmd)
	if err != nil {
		return fmt.Errorf("failed to run the command: %s", err)
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	if c.handlers == nil {
		c.handlers = make(map[string]func(*state, command) error)
	}

	c.handlers[name] = f
}
