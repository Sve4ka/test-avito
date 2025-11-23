package entity

import "time"

type PullRequestStatus string

const (
	PRStatusOPEN   PullRequestStatus = "OPEN"
	PRStatusMERGED PullRequestStatus = "MERGED"
)

type PullRequestCreate struct {
	AuthorId        string `json:"author_id"`
	PullRequestId   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
}

type PullRequest struct {
	AssignedReviewers []string          `json:"assigned_reviewers"`
	AuthorId          string            `json:"author_id"`
	CreatedAt         *time.Time        `json:"createdAt"`
	MergedAt          *time.Time        `json:"mergedAt"`
	PullRequestId     string            `json:"pull_request_id"`
	PullRequestName   string            `json:"pull_request_name"`
	Status            PullRequestStatus `json:"status"`
}

type PullRequestShort struct {
	AuthorId        string            `json:"author_id"`
	PullRequestId   string            `json:"pull_request_id"`
	PullRequestName string            `json:"pull_request_name"`
	Status          PullRequestStatus `json:"status"`
}

type Team struct {
	Members  []TeamMember `json:"members"`
	TeamName string       `json:"team_name"`
}

type TeamMember struct {
	IsActive bool   `json:"is_active"`
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}

type User struct {
	IsActive bool   `json:"is_active"`
	TeamName string `json:"team_name"`
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}

type UserStat struct {
	AvgDuration *float64 `json:"avg_duration"`
	CountPr     int      `json:"count_pr"`
	IsActive    bool     `json:"is_active"`
	UserId      string   `json:"user_id"`
}

type TeamStat struct {
	TeamName    string     `json:"team_name"`
	UsersStat   []UserStat `json:"users_stat"`
	AvgDuration float64    `json:"avg_duration"`
}
