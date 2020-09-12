package twitter

type TweetIncludes struct {
  Tweets []Tweet `json:"tweets"`
  Users  []User  `json:"users"`
  Media  []Media `json:"media"`
}
