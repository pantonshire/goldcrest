use chrono::prelude::*;

#[derive(Clone, Debug)]
pub struct Authentication {
    pub(crate) consumer_key: String,
    pub(crate) consumer_secret: String,
    pub(crate) access_token: String,
    pub(crate) token_secret: String,
}

impl Authentication {
    pub fn new(consumer_key: String, consumer_secret: String, access_token: String, token_secret: String) -> Self {
        Authentication{
            consumer_key,
            consumer_secret,
            access_token,
            token_secret,
        }
    }
}

#[derive(Clone, Copy, Debug)]
pub enum TweetMode {
    Compatibility,
    Extended,
}

#[derive(Clone, Debug)]
pub struct TweetOptions {
    pub(crate) par_trim_user: bool,
    pub(crate) par_include_your_retweet: bool,
    pub(crate) par_include_entities: bool,
    pub(crate) par_include_ext_alt_text: bool,
    pub(crate) par_include_card_uri: bool,
    pub(crate) par_mode: TweetMode,
}

impl TweetOptions {
    pub fn default() -> Self {
        TweetOptions{
            par_trim_user: false,
            par_include_your_retweet: true,
            par_include_entities: true,
            par_include_ext_alt_text: true,
            par_include_card_uri: true,
            par_mode: TweetMode::Extended,
        }
    }

    pub fn trim_user(self, trim: bool) -> Self {
        TweetOptions{
            par_trim_user: trim,
            ..self
        }
    }

    pub fn include_your_retweet(self, include: bool) -> Self {
        TweetOptions{
            par_include_your_retweet: include,
            ..self
        }
    }

    pub fn include_entities(self, include: bool) -> Self {
        TweetOptions{
            par_include_entities: include,
            ..self
        }
    }

    pub fn include_alt(self, include: bool) -> Self {
        TweetOptions{
            par_include_ext_alt_text: include,
            ..self
        }
    }

    pub fn include_card_uri(self, include: bool) -> Self {
        TweetOptions{
            par_include_card_uri: include,
            ..self
        }
    }

    pub fn mode(self, mode: TweetMode) -> Self {
        TweetOptions{
            par_mode: mode,
            ..self
        }
    }
}

#[derive(Clone, Debug)]
pub struct TimelineOptions {
    pub(crate) par_count: u32,
    pub(crate) par_min: Option<u64>,
    pub(crate) par_max: Option<u64>,
}

impl TimelineOptions {
    pub fn default() -> TimelineOptions {
        TimelineOptions{
            par_count: 20,
            par_min: None,
            par_max: None,
        }
    }

    pub fn count(self, num: u32) -> Self {
        TimelineOptions{
            par_count: num,
            ..self
        }
    }

    pub fn min_id(self, id: u64) -> Self {
        TimelineOptions{
            par_min: Some(id),
            ..self
        }
    }

    pub fn max_id(self, id: u64) -> Self {
        TimelineOptions{
            par_max: Some(id),
            ..self
        }
    }

    pub fn id_range(self, min: u64, max: u64) -> Self {
        TimelineOptions{
            par_min: Some(min),
            par_max: Some(max),
            ..self
        }
    }
}

#[derive(Clone, Copy, Debug)]
pub enum SearchResultType {
    Mixed,
    Recent,
    Popular,
}

#[derive(Clone, Debug)]
pub struct SearchOptions {
    pub(crate) par_query: String,
    pub(crate) par_geocode: Option<String>,
    pub(crate) par_lang: Option<String>,
    pub(crate) par_locale: Option<String>,
    pub(crate) par_result_type: SearchResultType,
    pub(crate) par_until: Option<i64>,
}

impl SearchOptions {
    pub fn new(query: String) -> Self {
        SearchOptions{
            par_query: query,
            par_geocode: None,
            par_lang: None,
            par_locale: None,
            par_result_type: SearchResultType::Mixed,
            par_until: None,
        }
    }

    pub fn geocode(self, geocode: String) -> Self {
        SearchOptions{
            par_geocode: Some(geocode),
            ..self
        }
    }

    pub fn lang(self, lang: String) -> Self {
        SearchOptions{
            par_lang: Some(lang),
            ..self
        }
    }

    pub fn locale(self, locale: String) -> Self {
        SearchOptions{
            par_locale: Some(locale),
            ..self
        }
    }

    pub fn result_type(self, result_type: SearchResultType) -> Self {
        SearchOptions{
            par_result_type: result_type,
            ..self
        }
    }

    pub fn until<Tz: TimeZone>(self, time: DateTime<Tz>) -> Self {
        self.until_unix(time.timestamp())
    }

    pub fn until_unix(self, timestamp: i64) -> Self {
        SearchOptions{
            par_until: Some(timestamp),
            ..self
        }
    }
}

#[derive(Clone, Debug)]
pub enum UserIdentifier {
    Id(u64),
    Handle(String),
}

#[derive(Clone, Debug)]
pub struct TweetBuilder {
    pub(crate) par_text: String,
    pub(crate) par_reply_id: Option<u64>,
    pub(crate) par_exclude_ids: Vec<u64>,
    pub(crate) par_media_ids: Vec<u64>,
    pub(crate) par_sensitive: bool,
    pub(crate) par_enable_dm_commands: bool,
    pub(crate) par_fail_dm_commands: bool,
    pub(crate) par_attachment_url: Option<String>,
}

impl TweetBuilder {
    pub fn new(text: String) -> TweetBuilder {
        TweetBuilder{
            par_text: text,
            par_reply_id: None,
            par_exclude_ids: Vec::new(),
            par_media_ids: Vec::new(),
            par_sensitive: false,
            par_enable_dm_commands: false,
            par_fail_dm_commands: false,
            par_attachment_url: None,
        }
    }

    pub fn reply_to(&mut self, id: u64) -> &mut Self {
        self.par_reply_id = Some(id);
        self
    }

    pub fn exclude_id(&mut self, id: u64) -> &mut Self {
        self.par_exclude_ids.push(id);
        self
    }

    pub fn exclude_ids<I: Iterator<Item = u64>>(&mut self, ids: I) -> &mut Self {
        self.par_exclude_ids.extend(ids);
        self
    }

    pub fn mark_sensitive(&mut self) -> &mut Self {
        self.par_sensitive = true;
        self
    }

    pub fn enable_dm_commands(&mut self) -> &mut Self {
        self.par_enable_dm_commands = true;
        self
    }

    pub fn fail_dm_commands(&mut self) -> &mut Self {
        self.par_fail_dm_commands = true;
        self
    }

    pub fn attach(&mut self, url: String) -> &mut Self {
        self.par_attachment_url = Some(url);
        self
    }
}

#[derive(Clone, Debug)]
pub struct ProfileBuilder {
    pub(crate) par_name: Option<String>,
    pub(crate) par_url: Option<String>,
    pub(crate) par_location: Option<String>,
    pub(crate) par_bio: Option<String>,
    pub(crate) par_link_color: Option<String>,
}

impl ProfileBuilder {
    pub fn new() -> ProfileBuilder {
        ProfileBuilder{
            par_name: None,
            par_url: None,
            par_location: None,
            par_bio: None,
            par_link_color: None,
        }
    }

    pub fn name(&mut self, name: String) -> &mut Self {
        self.par_name = Some(name);
        self
    }

    pub fn url(&mut self, url: String) -> &mut Self {
        self.par_url = Some(url);
        self
    }

    pub fn location(&mut self, location: String) -> &mut Self {
        self.par_location = Some(location);
        self
    }

    pub fn bio(&mut self, bio: String) -> &mut Self {
        self.par_bio = Some(bio);
        self
    }

    pub fn link_color(&mut self, color: String) -> &mut Self {
        self.par_link_color = Some(color);
        self
    }
}
