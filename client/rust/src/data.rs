use crate::twitter1;

use std::convert::TryFrom;
use chrono::{DateTime, NaiveDateTime, Duration, Utc};

type DesResult<T> = Result<T, DeserializationError>;

#[derive(Debug)]
pub enum DeserializationError {
    FieldMissing,
    FieldOutOfRange,
}

impl std::fmt::Display for DeserializationError {
    fn fmt(&self, fmt: &mut std::fmt::Formatter) -> std::fmt::Result {
        write!(fmt, "Deserialization error: {}", match self {
            DeserializationError::FieldMissing    => "missing field",
            DeserializationError::FieldOutOfRange => "field out of range",
        })
    }
}

impl std::error::Error for DeserializationError {}

trait Exists<T> {
    fn exists(self) -> Result<T, DeserializationError>;
}

impl<T> Exists<T> for Option<T> {
    fn exists(self) -> Result<T, DeserializationError> {
        self.ok_or(DeserializationError::FieldMissing)
    }
}

pub(super) trait Deserialize<T: Sized> {
    fn des(self) -> DesResult<T>;
}

impl<T,S> Deserialize<Option<S>> for Option<T> where T: Deserialize<S> {
    fn des(self) -> DesResult<Option<S>> {
        self.map(T::des).map_or(Ok(None), |x| x.map(Some))
    }
}

impl<T,S> Deserialize<Box<S>> for Box<T> where T: Deserialize<S> {
    fn des(self) -> DesResult<Box<S>> {
        Ok(Box::new((*self).des()?))
    }
}

impl<T,S> Deserialize<Vec<S>> for Vec<T> where T: Deserialize<S> {
    fn des(self) -> DesResult<Vec<S>> {
        self.into_iter().map(T::des).collect()
    }
}

impl Deserialize<DateTime<Utc>> for i64 {
    fn des(self) -> DesResult<DateTime<Utc>> {
        Ok(DateTime::from_utc(NaiveDateTime::from_timestamp(self, 0), Utc))
    }
}

impl Deserialize<DateTime<Utc>> for u64 {
    fn des(self) -> DesResult<DateTime<Utc>> {
        match i64::try_from(self) {
            Ok(x)  => x.des(),
            Err(_) => Err(DeserializationError::FieldOutOfRange)
        }
    }
}

pub type Indices = (usize, usize);

impl Deserialize<Indices> for twitter1::Indices {
    fn des(self) -> DesResult<Indices> {
        Ok((self.start as usize, self.end as usize))
    }
}

#[derive(Clone, Debug)]
pub struct Tweet {
    pub id: u64,
    pub created_at: chrono::DateTime<chrono::Utc>,
    pub text: String,
    pub text_display_range: Option<Indices>,
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

impl Deserialize<Tweet> for twitter1::Tweet {
    fn des(self) -> DesResult<Tweet> {
        Ok(Tweet{
            id: self.id,
            created_at: self.created_at.des()?,
            text: self.text,
            text_display_range: self.text_display_range.des()?,
            truncated: self.truncated,
            source: self.source,
            user: Box::new(self.user.exists()?.des()?),
            replied_to: self.replied_tweet.des()?,
            quoted: self.quoted_tweet.des()?,
            retweeted: self.retweeted_tweet.des()?,
            quotes: self.quote_count,
            replies: self.reply_count,
            retweets: self.retweet_count,
            likes: self.favorite_count,
            you_liked: self.favorited,
            you_retweeted: self.retweeted,
            your_retweet_id: self.current_user_retweet_id.map(u64::from),
            hashtags: self.hashtags.des()?,
            urls: self.urls.des()?,
            mentions: self.mentions.des()?,
            symbols: self.symbols.des()?,
            media: self.media.des()?,
            polls: self.polls.des()?,
            sensitive: self.possibly_sensitive,
            filter_level: self.filter_level,
            lang: self.lang,
            withheld_copyright: self.withheld_copyright,
            withheld_countries: self.withheld_countries,
            withheld_scope: self.withheld_scope,
        })
    }
}

impl Deserialize<Vec<Tweet>> for twitter1::Tweets {
    fn des(self) -> DesResult<Vec<Tweet>> {
        self.tweets.into_iter()
            .map(twitter1::Tweet::des)
            .collect()
    }
}

#[derive(Clone, Debug)]
pub struct ReplyData {
    pub tweet_id: u64,
    pub user_id: u64,
    pub user_handle: String,
}

impl Deserialize<ReplyData> for twitter1::tweet::ReplyData {
    fn des(self) -> DesResult<ReplyData> {
        Ok(ReplyData{
            tweet_id: self.reply_to_tweet_id,
            user_id: self.reply_to_user_id,
            user_handle: self.reply_to_user_handle,
        })
    }
}

#[derive(Clone, Debug)]
pub struct User {
    pub id: u64,
    pub handle: String,
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

impl Deserialize<User> for twitter1::User {
    fn des(self) -> DesResult<User> {
        Ok(User{
            id: self.id,
            handle: self.handle,
            name: self.display_name,
            created_at: self.created_at.des()?,
            bio: self.bio,
            url: self.url,
            location: self.location,
            protected: self.protected,
            verified: self.verified,
            follower_count: self.follower_count,
            following_count: self.following_count,
            listed_count: self.listed_count,
            favorites_count: self.favorites_count,
            statuses_count: self.statuses_count,
            profile_banner: self.profile_banner,
            profile_image: self.profile_image,
            default_profile: self.default_profile,
            default_profile_image: self.default_profile_image,
            withheld_countries: self.withheld_countries,
            withheld_scope: self.withheld_scope,
            urls: self.url_urls.des()?,
            bio_urls: self.bio_urls.des()?,
        })
    }
}

#[derive(Clone, Debug)]
pub struct URL {
    pub indices: Indices,
    pub twitter_url: String,
    pub display_url: String,
    pub expanded_url: String,
}

impl Deserialize<URL> for twitter1::Url {
    fn des(self) -> DesResult<URL> {
        Ok(URL{
            indices: self.indices.exists()?.des()?,
            twitter_url: self.twitter_url,
            display_url: self.display_url,
            expanded_url: self.expanded_url,
        })
    }
}

#[derive(Clone, Debug)]
pub struct Symbol {
    pub indices: Indices,
    pub text: String,
}

impl Deserialize<Symbol> for twitter1::Symbol {
    fn des(self) -> DesResult<Symbol> {
        Ok(Symbol{
            indices: self.indices.exists()?.des()?,
            text: self.text,
        })
    }
}

#[derive(Clone, Debug)]
pub struct Mention {
    pub indices: Indices,
    pub user_id: u64,
    pub user_handle: String,
    pub user_name: String,
}

impl Deserialize<Mention> for twitter1::Mention {
    fn des(self) -> DesResult<Mention> {
        Ok(Mention{
            indices: self.indices.exists()?.des()?,
            user_id: self.user_id,
            user_handle: self.handle,
            user_name: self.display_name,
        })
    }
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

impl Deserialize<Media> for twitter1::Media {
    fn des(self) -> DesResult<Media> {
        Ok(Media{
            url: self.url.exists()?.des()?,
            id: self.id,
            media_type: self.r#type,
            media_url: self.media_url,
            alt: self.alt,
            source_tweet_id: self.source_tweet_id.map(u64::from),
            thumb: self.thumb.exists()?.des()?,
            small: self.small.exists()?.des()?,
            medium: self.medium.exists()?.des()?,
            large: self.large.exists()?.des()?,
        })
    }
}

#[derive(Clone, Debug)]
pub struct MediaSize {
    pub width: u32,
    pub height: u32,
    pub resize: String,
}

impl Deserialize<MediaSize> for twitter1::media::Size {
    fn des(self) -> DesResult<MediaSize> {
        Ok(MediaSize{
            width: self.width,
            height: self.height,
            resize: self.resize,
        })
    }
}

#[derive(Clone, Debug)]
pub struct Poll {
    pub end_time: chrono::DateTime<chrono::Utc>,
    pub duration: chrono::Duration,
    pub options: Vec<PollOption>,
}

impl Deserialize<Poll> for twitter1::Poll {
    fn des(self) -> DesResult<Poll> {
        Ok(Poll{
            end_time: self.end_time.des()?,
            duration: Duration::minutes(self.duration_minutes as i64),
            options: self.options.des()?,
        })
    }
}

#[derive(Clone, Debug)]
pub struct PollOption {
    pub position: usize,
    pub text: String,
}

impl Deserialize<PollOption> for twitter1::poll::Option {
    fn des(self) -> DesResult<PollOption> {
        Ok(PollOption{
            position: self.position as usize,
            text: self.text,
        })
    }
}
