package au

func (tweet Tweet) Unfold() Tweet {
  if tweet.Retweeted != nil {
    return *tweet.Retweeted
  }
  return tweet
}

func (tweet Tweet) TextOnly() string {
  var removeIndices []Indices
  if !tweet.TextDisplayRange.IsZero() {
    removeIndices = append(removeIndices, tweet.TextDisplayRange.Invert(uint(len(tweet.Text)))...)
  }
  for _, hashtag := range tweet.Hashtags {
    if !hashtag.Indices.IsZero() {
      removeIndices = append(removeIndices, hashtag.Indices)
    }
  }
  for _, url := range tweet.URLs {
    if !url.Indices.IsZero() {
      removeIndices = append(removeIndices, url.Indices)
    }
  }
  for _, mention := range tweet.Mentions {
    if !mention.Indices.IsZero() {
      removeIndices = append(removeIndices, mention.Indices)
    }
  }
  for _, symbol := range tweet.Symbols {
    if !symbol.Indices.IsZero() {
      removeIndices = append(removeIndices, symbol.Indices)
    }
  }
  for _, media := range tweet.Media {
    if !media.Indices.IsZero() {
      removeIndices = append(removeIndices, media.Indices)
    }
  }
  return removeFromString(tweet.Text, removeIndices...)
}

func (tweet Tweet) MakeReplyText(text string) string {
  return tweet.User.AtHandle() + " " + text
}

func (user User) AtHandle() string {
  return "@" + user.Handle
}

func (mention Mention) AtHandle() string {
  return "@" + mention.UserHandle
}

func (indices Indices) IsZero() bool {
  return indices.End <= indices.Start
}

func (indices Indices) Invert(l uint) []Indices {
  var inverted []Indices
  if indices.Start > 0 {
    inverted = append(inverted, Indices{Start: 0, End: indices.Start})
  }
  if indices.End < l {
    inverted = append(inverted, Indices{Start: indices.End, End: l})
  }
  return inverted
}
