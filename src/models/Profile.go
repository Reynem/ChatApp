package models

import (
	"time"
)

type Profile struct {
	user_id      uint      `json:"id"`
	profile_name string    `json:"profile_name"`
	bio          string    `json:"bio"`
	avatar_url   string    `json:"avatar_url"`
	status       string    `json:"status"`
	last_seen    time.Time `json:"last_seen"`
}
