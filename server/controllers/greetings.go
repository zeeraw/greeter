package controllers

import (
	"context"
	"fmt"
)

// Greetings represents actions when greeting
type Greetings struct{}

// Hello is used to greet the server and get a response back
func (g *Greetings) Hello(ctx context.Context, name string) (string, error) {
	return fmt.Sprintf("Hello %s", name), nil
}
