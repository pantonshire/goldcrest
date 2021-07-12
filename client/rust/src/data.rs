use std::ops::Range;

#[derive(Copy, Clone, Debug)]
pub struct Indices {
    pub start_inclusive: usize,
    pub end_exlcusive: usize,
}

impl Indices {
    #[inline]
    pub fn contains(&self, i: usize) -> bool {
        self.start_inclusive <= i && i < self.end_exlcusive
    }

    pub fn limit(&self, max: usize) -> Self {
        Self {
            start_inclusive: self.start_inclusive,
            end_exlcusive: self.end_exlcusive.min(max)
        }
    }

    pub fn range(&self) -> Range<usize> {
        self.start_inclusive..self.end_exlcusive
    }
}

/// Represents a single Tweet, otherwise known as a status in the Twitter API.
#[derive(Clone, Debug)]
pub struct Tweet {
    pub id: u64,
    pub created_at: chrono::DateTime<chrono::Utc>,
    pub raw_text: String,
    pub text_display_range: Indices,
    pub truncated: bool,
    pub source: String,
    pub user: Box<User>,
    pub replied_to: Option<ReplyData>,
    pub quoted: Option<Box<Tweet>>,
    pub retweeted: Option<Box<Tweet>>,
    pub quotes: u32,
    pub replies: u32,
    pub retweets: u32,
    pub likes: u32,
    pub you_liked: bool,
    pub you_retweeted: bool,
    pub your_retweet_id: Option<u64>,
    pub hashtags: Vec<Symbol>,
    pub urls: Vec<URL>,
    pub mentions: Vec<Mention>,
    pub symbols: Vec<Symbol>,
    pub media: Vec<Media>,
    pub polls: Vec<Poll>,
    pub sensitive: bool,
    pub filter_level: String,
    pub lang: String,
    pub withheld_copyright: bool,
    pub withheld_countries: Vec<String>,
    pub withheld_scope: String,
}

impl Tweet {
    /// If this Tweet is a retweet, returns the original Tweet. Otherwise, returns this Tweet.
    /// Moves the Tweet struct.
    pub fn original(self) -> Self {
        match self.retweeted {
            Some(tweet) => *tweet,
            None => self,
        }
    }

    /// Returns the Tweet text trimmed according to its display range. Inline entities such as
    /// URLs, mentions and hashtags can optionally be removed, depending on the provided
    /// TweetTextOptions.
    pub fn text(&self, options: &tweet::TweetTextOptions) -> String {
        let full_len = self.raw_text.chars().count();

        let exclude_indices = {
            let mut exclude_indices = vec![false; full_len];

            if !options.hashtags_included {
                for hashtag in self.hashtags.iter() {
                    for i in hashtag.indices.limit(full_len).range() {
                        exclude_indices[i] = true;
                    }
                }
            }

            if !options.urls_included {
                for url in self.urls.iter() {
                    for i in url.indices.limit(full_len).range() {
                        exclude_indices[i] = true;
                    }
                }
            }

            if !options.mentions_included {
                for mention in self.mentions.iter() {
                    for i in mention.indices.limit(full_len).range() {
                        exclude_indices[i] = true;
                    }
                }
            }

            if !options.symbols_included {
                for symbol in self.symbols.iter() {
                    for i in symbol.indices.limit(full_len).range() {
                        exclude_indices[i] = true;
                    }
                }
            }

            if !options.media_included {
                for media in self.media.iter() {
                    for i in media.url.indices.limit(full_len).range() {
                        exclude_indices[i] = true;
                    }
                }
            }

            exclude_indices
        };

        self.raw_text
            .chars()
            .enumerate()
            .filter_map(|(i, c)| if exclude_indices[i] {
                None
            } else {
                Some(c)
            })
            .collect()
    }
}

pub mod tweet {
    #[derive(Eq, PartialEq, Clone, Debug)]
    pub struct TweetTextOptions {
        pub(super) hashtags_included: bool,
        pub(super) urls_included: bool,
        pub(super) mentions_included: bool,
        pub(super) symbols_included: bool,
        pub(super) media_included: bool,
    }

    impl TweetTextOptions {
        pub const fn all() -> TweetTextOptions {
            TweetTextOptions{
                hashtags_included: true,
                urls_included: true,
                mentions_included: true,
                symbols_included: true,
                media_included: true,
            }
        }

        pub const fn none() -> TweetTextOptions {
            TweetTextOptions{
                hashtags_included: false,
                urls_included: false,
                mentions_included: false,
                symbols_included: false,
                media_included: false,
            }
        }

        pub const fn hashtags(self, included: bool) -> Self {
            TweetTextOptions{
                hashtags_included: included,
                ..self
            }
        }

        pub const fn urls(self, included: bool) -> Self {
            TweetTextOptions{
                urls_included: included,
                ..self
            }
        }

        pub const fn mentions(self, included: bool) -> Self {
            TweetTextOptions{
                mentions_included: included,
                ..self
            }
        }

        pub const fn symbols(self, included: bool) -> Self {
            TweetTextOptions{
                symbols_included: included,
                ..self
            }
        }

        pub const fn media(self, included: bool) -> Self {
            TweetTextOptions{
                media_included: included,
                ..self
            }
        }
    }
}

/// Represents a handle that uniquely identifies a user.
/// Referred to as "screen name" in the Twitter API.
/// Note that a user may change their handle, so user IDs should be preferred as a way of
/// identifying users rather than handles.
#[derive(Clone, Debug)]
pub struct Handle {
    /// The user's handle without the leading at symbol (@).
    pub name_only: String,
}

impl Handle {
    /// Returns the user's handle with an at symbol (@) prepended.
    pub fn at_name(&self) -> String {
        let mut at_name = "@".to_owned();
        at_name.push_str(&self.name_only);
        at_name
    }
}

#[derive(Clone, Debug)]
pub struct ReplyData {
    pub tweet_id: u64,
    pub user_id: u64,
    pub user_handle: Handle,
}

#[derive(Clone, Debug)]
pub struct User {
    pub id: u64,
    pub handle: Handle,
    pub name: String,
    pub created_at: chrono::DateTime<chrono::Utc>,
    pub bio: String,
    pub url: String,
    pub location: String,
    pub protected: bool,
    pub verified: bool,
    pub follower_count: u32,
    pub following_count: u32,
    pub listed_count: u32,
    pub favorites_count: u32,
    pub statuses_count: u32,
    pub profile_banner: String,
    pub profile_image: String,
    pub default_profile: bool,
    pub default_profile_image: bool,
    pub withheld_countries: Vec<String>,
    pub withheld_scope: String,
    pub urls: Vec<URL>,
    pub bio_urls: Vec<URL>,
}

#[derive(Clone, Debug)]
pub struct URL {
    pub indices: Indices,
    pub twitter_url: String,
    pub display_url: String,
    pub expanded_url: String,
}

#[derive(Clone, Debug)]
pub struct Symbol {
    pub indices: Indices,
    pub text: String,
}

#[derive(Clone, Debug)]
pub struct Mention {
    pub indices: Indices,
    pub user_id: u64,
    pub user_handle: Handle,
    pub user_name: String,
}

#[derive(Clone, Debug)]
pub struct Media {
    pub url: URL,
    pub id: u64,
    pub media_type: String,
    pub media_url: String,
    pub alt: String,
    pub source_tweet_id: Option<u64>,
    pub thumb: MediaSize,
    pub small: MediaSize,
    pub medium: MediaSize,
    pub large: MediaSize,
}

#[derive(Clone, Debug)]
pub struct MediaSize {
    pub width: u32,
    pub height: u32,
    pub resize: String,
}

#[derive(Clone, Debug)]
pub struct Poll {
    pub end_time: chrono::DateTime<chrono::Utc>,
    pub duration: chrono::Duration,
    pub options: Vec<PollOption>,
}

#[derive(Clone, Debug)]
pub struct PollOption {
    pub position: usize,
    pub text: String,
}
