package models

import (
	"time"
)

type Profile struct {
	User_id      uint      `json:"id"`
	Profile_name string    `json:"profile_name"`
	Bio          string    `json:"bio"`
	Avatar_url   string    `json:"avatar_url"`
	Status       string    `json:"status"`
	Last_seen    time.Time `json:"last_seen"`
}
