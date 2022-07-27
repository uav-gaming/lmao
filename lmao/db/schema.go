package db

import "github.com/diamondburned/arikawa/v3/discord"

type UserInfo struct {
	ID            discord.UserID `json:"id"`
	Username      string         `json:"username"`
	Discriminator string         `json:"discriminator"`
}

// Leetcode Question slug. For example, "two-sum".
type QuestionSlug string

type UserLeetcodeInfo struct {
	// Discord username
	UserID discord.UserID `json:"id"`
	// Discord user discriminator
	Discriminator string
	// Leetcode username
	LeetcodeUsername string
	// Set of completed questions
	// We will encode them in custom bitmap since it's complication to use dynamodb's set in Golang
	// https://github.com/aws/aws-sdk-go/issues/1990#issuecomment-614087254
	CompletedQuestions []byte
}

// Each TODO list is bind to a discord user since
type TodoList struct {
	// Discord username
	UserID discord.UserID `json:"id"`
	// Set of TODO questions.
	Questions []byte
}
