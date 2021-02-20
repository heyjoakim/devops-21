package models

// FollowRequest defines a follow request
type FollowRequest struct {
	Follow   string `json:"follow"`
	Unfollow string `json:"unfollow"`
}
