use crate::error::{DeserializationError, DeserializationResult, Exists};
use crate::data::*;
use crate::twitter1;

use std::convert::TryFrom;
use chrono::{DateTime, NaiveDateTime, Duration, Utc};

pub(super) trait Deserialize<T> where T: Sized {
    fn des(self) -> DeserializationResult<T>;
}

impl<T,S> Deserialize<Option<S>> for Option<T> where T: Deserialize<S> {
    fn des(self) -> DeserializationResult<Option<S>> {
        self.map(T::des).map_or(Ok(None), |x| x.map(Some))
    }
}

impl<T,S> Deserialize<Box<S>> for Box<T> where T: Deserialize<S> {
    fn des(self) -> DeserializationResult<Box<S>> {
        Ok(Box::new((*self).des()?))
    }
}

impl<T,S> Deserialize<Vec<S>> for Vec<T> where T: Deserialize<S> {
    fn des(self) -> DeserializationResult<Vec<S>> {
        self.into_iter().map(T::des).collect()
    }
}

impl Deserialize<DateTime<Utc>> for i64 {
    fn des(self) -> DeserializationResult<DateTime<Utc>> {
        Ok(DateTime::from_utc(NaiveDateTime::from_timestamp(self, 0), Utc))
    }
}

impl Deserialize<DateTime<Utc>> for u64 {
    fn des(self) -> DeserializationResult<DateTime<Utc>> {
        match i64::try_from(self) {
            Ok(x)  => x.des(),
            Err(_) => Err(DeserializationError::FieldOutOfRange)
        }
    }
}

impl Deserialize<Handle> for String {
    fn des(self) -> DeserializationResult<Handle> {
        Ok(Handle{
            name_only: self,
        })
    }
}

impl Deserialize<Indices> for twitter1::Indices {
    fn des(self) -> DeserializationResult<Indices> {
        Ok(self.start as usize .. self.end as usize)
    }
}

impl Deserialize<Tweet> for twitter1::Tweet {
    fn des(self) -> DeserializationResult<Tweet> {
        Ok(Tweet{
            id: self.id,
            created_at: self.created_at.des()?,
            raw_text: self.text,
            text_display_range: self.text_display_range.exists()?.des()?,
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
    fn des(self) -> DeserializationResult<Vec<Tweet>> {
        self.tweets.into_iter()
            .map(twitter1::Tweet::des)
            .collect()
    }
}

impl Deserialize<ReplyData> for twitter1::tweet::ReplyData {
    fn des(self) -> DeserializationResult<ReplyData> {
        Ok(ReplyData{
            tweet_id: self.reply_to_tweet_id,
            user_id: self.reply_to_user_id,
            user_handle: self.reply_to_user_handle.des()?,
        })
    }
}

impl Deserialize<User> for twitter1::User {
    fn des(self) -> DeserializationResult<User> {
        Ok(User{
            id: self.id,
            handle: self.handle.des()?,
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

impl Deserialize<URL> for twitter1::Url {
    fn des(self) -> DeserializationResult<URL> {
        Ok(URL{
            indices: self.indices.exists()?.des()?,
            twitter_url: self.twitter_url,
            display_url: self.display_url,
            expanded_url: self.expanded_url,
        })
    }
}

impl Deserialize<Symbol> for twitter1::Symbol {
    fn des(self) -> DeserializationResult<Symbol> {
        Ok(Symbol{
            indices: self.indices.exists()?.des()?,
            text: self.text,
        })
    }
}

impl Deserialize<Mention> for twitter1::Mention {
    fn des(self) -> DeserializationResult<Mention> {
        Ok(Mention{
            indices: self.indices.exists()?.des()?,
            user_id: self.user_id,
            user_handle: self.handle.des()?,
            user_name: self.display_name,
        })
    }
}

impl Deserialize<Media> for twitter1::Media {
    fn des(self) -> DeserializationResult<Media> {
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

impl Deserialize<MediaSize> for twitter1::media::Size {
    fn des(self) -> DeserializationResult<MediaSize> {
        Ok(MediaSize{
            width: self.width,
            height: self.height,
            resize: self.resize,
        })
    }
}

impl Deserialize<Poll> for twitter1::Poll {
    fn des(self) -> DeserializationResult<Poll> {
        Ok(Poll{
            end_time: self.end_time.des()?,
            duration: Duration::minutes(self.duration_minutes as i64),
            options: self.options.des()?,
        })
    }
}

impl Deserialize<PollOption> for twitter1::poll::Option {
    fn des(self) -> DeserializationResult<PollOption> {
        Ok(PollOption{
            position: self.position as usize,
            text: self.text,
        })
    }
}
