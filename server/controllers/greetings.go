package controllers

import "fmt"

// Greetings represents actions when greeting
type Greetings struct{}

// HelloResponder represents response flow control for the Hello function
type HelloResponder struct {
	ErrC chan error
	ResC chan string
}

// Close will close all the responder channels
func (r *HelloResponder) Close() {
	close(r.ErrC)
	close(r.ResC)
}

// Hello is used to greet the server and get a response back
func (g *Greetings) Hello(name string) *HelloResponder {
	res := &HelloResponder{
		ErrC: make(chan error, 1),
		ResC: make(chan string, 1),
	}
	res.ResC <- fmt.Sprintf("Hello %s", name)

	return res
}
