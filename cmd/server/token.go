package main

type manager struct {
	tokens map[string]token
}

type token struct {
	text     string `json:"token"`
	username string `json:"username"`
}

func newToken(text, username string) *token {
	return &token{text: text, username: username}
}
