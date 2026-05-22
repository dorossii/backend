package controllers

import "os"

var inviteBaseURL string

func init() {
	inviteBaseURL = os.Getenv("INVITE_BASE_URL")
	if inviteBaseURL == "" {
		inviteBaseURL = "https://localhost:1111"
	}
}
