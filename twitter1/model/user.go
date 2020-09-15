package model

type User struct {
  ID                  uint64      `json:"id"`
  IDStr               string      `json:"id_str"`
  Name                string      `json:"name"`
  ScreenName          string      `json:"screen_name"`
  CreatedAt           TwitterTime `json:"created_at"`
  Location            *string     `json:"location"`
  URL                 *string     `json:"url"`
  Description         *string     `json:"description"`
  Protected           bool        `json:"protected"`
  Verified            bool        `json:"verified"`
  FollowersCount      uint        `json:"followers_count"`
  FriendsCount        uint        `json:"friends_count"`
  ListedCount         uint        `json:"listed_count"`
  FavoritesCount      uint        `json:"favourites_count"`
  StatusesCount       uint        `json:"statuses_count"`
  ProfileBanner       string      `json:"profile_banner_url"`
  ProfileImage        string      `json:"profile_image_url_https"`
  DefaultProfile      bool        `json:"default_profile"`
  DefaultProfileImage bool        `json:"default_profile_image"`
  WithheldCountries   []string    `json:"withheld_in_countries"`
  WithheldScope       *string     `json:"withheld_scope"`

  Entities struct {
    URL struct {
      URLs []URL `json:"urls"`
    } `json:"url"`
    Description struct {
      URLs []URL `json:"urls"`
    } `json:"description"`
  } `json:"entities"`

  //TODO: derived
}
