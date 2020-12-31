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
    pub fn original(self) -> Self {
        match self.retweeted {
            Some(tweet) => *tweet,
            None        => self,
        }
    }

    pub fn text(&self) -> String {
        let l = self.raw_text.chars().count();
        let mut mask = indices_to_bits(&self.text_display_range, l);
        for hashtag in self.hashtags.iter() {
            mask &= !indices_to_bits(&hashtag.indices, l);
        }
        for url in self.urls.iter() {
            mask &= !indices_to_bits(&url.indices, l);
        }
        for mention in self.mentions.iter() {
            mask &= !indices_to_bits(&mention.indices, l);
        }
        for symbol in self.symbols.iter() {
            mask &= !indices_to_bits(&symbol.indices, l);
        }
        for media in self.media.iter() {
            mask &= !indices_to_bits(&media.url.indices, l);
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

#[derive(Clone, Debug)]
pub struct Handle {
    pub name_only: String,
}

impl Handle {
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
