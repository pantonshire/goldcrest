use bitvec::prelude::*;

pub type Indices = std::ops::Range<usize>;

fn indices_to_bits(ind: &Indices, len: usize) -> BitVec<Lsb0, u8> {
    let mut bits = bitvec![Lsb0, u8; 0; len];
    let mut i = ind.start;
    while i < ind.end && i < len {
        bits.set(i, true);
        i += 1
    }
    bits
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
            None        => self,
        }
    }

    /// Returns the Tweet text trimmed according to its display range. Inline entities such as
    /// URLs, mentions and hashtags can optionally be removed, depending on the provided
    /// TweetTextOptions.
    pub fn text(&self, options: tweet::TweetTextOptions) -> String {
        let l = self.raw_text.chars().count();
        let mut mask = indices_to_bits(&self.text_display_range, l);
        if !options.hashtags_included {
            for hashtag in self.hashtags.iter() {
                mask &= !indices_to_bits(&hashtag.indices, l);
            }
        }
        if !options.urls_included {
            for url in self.urls.iter() {
                mask &= !indices_to_bits(&url.indices, l);
            }
        }
        if !options.mentions_included {
            for mention in self.mentions.iter() {
                mask &= !indices_to_bits(&mention.indices, l);
            }
        }
        if !options.symbols_included {
            for symbol in self.symbols.iter() {
                mask &= !indices_to_bits(&symbol.indices, l);
            }
        }
        if !options.media_included {
            for media in self.media.iter() {
                mask &= !indices_to_bits(&media.url.indices, l);
            }
        }
        let mut text = String::new();
        for (i, c) in self.raw_text.chars().enumerate() {
            if mask[i] {
                text.push(c);
            }
        }
        text
    }
}

pub mod tweet {
    pub struct TweetTextOptions {
        pub(super) hashtags_included: bool,
        pub(super) urls_included: bool,
        pub(super) mentions_included: bool,
        pub(super) symbols_included: bool,
        pub(super) media_included: bool,
    }

    impl TweetTextOptions {
        pub fn all() -> TweetTextOptions {
            TweetTextOptions{
                hashtags_included: true,
                urls_included: true,
                mentions_included: true,
                symbols_included: true,
                media_included: true,
            }
        }

        pub fn none() -> TweetTextOptions {
            TweetTextOptions{
                hashtags_included: false,
                urls_included: false,
                mentions_included: false,
                symbols_included: false,
                media_included: false,
            }
        }

        pub fn hashtags(self, included: bool) -> Self {
            TweetTextOptions{
                hashtags_included: included,
                ..self
            }
        }

        pub fn urls(self, included: bool) -> Self {
            TweetTextOptions{
                urls_included: included,
                ..self
            }
        }

        pub fn mentions(self, included: bool) -> Self {
            TweetTextOptions{
                mentions_included: included,
                ..self
            }
        }

        pub fn symbols(self, included: bool) -> Self {
            TweetTextOptions{
                symbols_included: included,
                ..self
            }
        }

        pub fn media(self, included: bool) -> Self {
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
