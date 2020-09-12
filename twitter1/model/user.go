package model

import "time"

type User struct {
  ID                  int64     `json:"id"`
  IDStr               string    `json:"id_str"`
  Name                string    `json:"name"`
  ScreenName          string    `json:"screen_name"`
  CreatedAt           time.Time `json:"created_at"`
  Location            *string   `json:"location"`
  URL                 *string   `json:"url"`
  Description         *string   `json:"description"`
  Protected           bool      `json:"protected"`
  Verified            bool      `json:"verified"`
  FollowersCount      int       `json:"followers_count"`
  FriendsCount        int       `json:"friends_count"`
  ListedCount         int       `json:"listed_count"`
  FavoritesCount      int       `json:"favourites_count"`
  StatusesCount       int       `json:"statuses_count"`
  ProfileBanner       string    `json:"profile_banner_url"`
  ProfileImage        string    `json:"profile_image_url_https"`
  DefaultProfile      bool      `json:"default_profile"`
  DefaultProfileImage bool      `json:"default_profile_image"`
  WithheldCountries   []string  `json:"withheld_in_countries"`
  WithheldScope       *string   `json:"withheld_scope"`

  //TODO: derived
}
