use crate::twitter1;

#[derive(Clone, Debug)]
pub struct Authentication {
    consumer_key: String,
    consumer_secret: String,
    token: String,
    token_secret: String,
}

impl Authentication {
    pub fn new(consumer_key: &str, consumer_secret: &str, token: &str, token_secret: &str) -> Self {
        Authentication{
            consumer_key: consumer_key.to_owned(),
            consumer_secret: consumer_secret.to_owned(),
            token: token.to_owned(),
            token_secret: token_secret.to_owned(),
        }
    }

    fn ser(self) -> twitter1::Authentication {
        twitter1::Authentication{
            consumer_key: self.consumer_key,
            secret_key: self.consumer_secret,
            access_token: self.token,
            secret_token: self.token_secret,
        }
    }
}

#[derive(Clone, Copy, Debug)]
pub enum TweetMode {
    Compatibility,
    Extended,
}

impl TweetMode {
    fn ser(self) -> i32 {
        use twitter1::tweet_options::Mode;
        (match self {
            TweetMode::Compatibility => Mode::Compat,
            TweetMode::Extended      => Mode::Extended,
        }) as i32
    }
}

#[derive(Clone, Debug)]
pub struct TweetOptions {
    _trim_user: bool,
    _include_your_retweet: bool,
    _include_entities: bool,
    _include_ext_alt_text: bool,
    _include_card_uri: bool,
    _mode: TweetMode,
}

impl TweetOptions {
    pub fn default() -> Self {
        TweetOptions{
            _trim_user: false,
            _include_your_retweet: true,
            _include_entities: true,
            _include_ext_alt_text: true,
            _include_card_uri: true,
            _mode: TweetMode::Extended,
        }
    }

    pub fn trim_user(self, trim: bool) -> Self {
        TweetOptions{
            _trim_user: trim,
            ..self
        }
    }

    pub fn include_your_retweet(self, include: bool) -> Self {
        TweetOptions{
            _include_your_retweet: include,
            ..self
        }
    }

    pub fn include_entities(self, include: bool) -> Self {
        TweetOptions{
            _include_entities: include,
            ..self
        }
    }

    pub fn include_alt(self, include: bool) -> Self {
        TweetOptions{
            _include_ext_alt_text: include,
            ..self
        }
    }

    pub fn include_card_uri(self, include: bool) -> Self {
        TweetOptions{
            _include_card_uri: include,
            ..self
        }
    }

    pub fn mode(self, mode: TweetMode) -> Self {
        TweetOptions{
            _mode: mode,
            ..self
        }
    }

    fn ser(self) -> twitter1::TweetOptions {
        twitter1::TweetOptions{
            trim_user: self._trim_user,
            include_my_retweet: self._include_your_retweet,
            include_entities: self._include_entities,
            include_ext_alt_text: self._include_ext_alt_text,
            include_card_uri: self._include_card_uri,
            mode: self._mode.ser(),
        }
    }
}

#[derive(Clone, Debug)]
pub struct TimelineOptions {
    _count: u32,
    _min: Option<u64>,
    _max: Option<u64>,
}

impl TimelineOptions {
    pub fn default() -> TimelineOptions {
        TimelineOptions{
            _count: 20,
            _min: None,
            _max: None,
        }
    }

    pub fn count(self, num: u32) -> Self {
        TimelineOptions{
            _count: num,
            ..self
        }
    }

    pub fn min_id(self, id: u64) -> Self {
        TimelineOptions{
            _min: Some(id),
            ..self
        }
    }

    pub fn max_id(self, id: u64) -> Self {
        TimelineOptions{
            _max: Some(id),
            ..self
        }
    }

    pub fn id_range(self, min: u64, max: u64) -> Self {
        TimelineOptions{
            _min: Some(min),
            _max: Some(max),
            ..self
        }
    }

    fn ser(self, twopts: TweetOptions) -> twitter1::TimelineOptions {
        twitter1::TimelineOptions{
            count: self._count,
            twopts: Some(twopts.ser()),
            min_id: self._min.map(u64::into),
            max_id: self._max.map(u64::into),
        }
    }
}

#[derive(Clone, Debug)]
pub enum UserIdentifier {
    Id(u64),
    Handle(String),
}

impl UserIdentifier {
    fn ser(self) -> twitter1::user_timeline_request::User {
        use twitter1::user_timeline_request::User;
        match self {
            UserIdentifier::Id(id)         => User::UserId(id),
            UserIdentifier::Handle(handle) => User::UserHandle(handle),
        }
    }
}

#[derive(Clone, Debug)]
pub struct TweetBuilder {
    _text: String,
    _reply_id: Option<u64>,
    _exclude_ids: Vec<u64>,
    _media_ids: Vec<u64>,
    _sensitive: bool,
    _enable_dm_commands: bool,
    _fail_dm_commands: bool,
    _attachment_url: Option<String>,
}

impl TweetBuilder {
    pub fn new(text: String) -> TweetBuilder {
        TweetBuilder{
            _text: text,
            _reply_id: None,
            _exclude_ids: Vec::new(),
            _media_ids: Vec::new(),
            _sensitive: false,
            _enable_dm_commands: false,
            _fail_dm_commands: false,
            _attachment_url: None,
        }
    }

    pub fn reply_to(&mut self, id: u64) -> &mut Self {
        self._reply_id = Some(id);
        self
    }

    pub fn exclude_id(&mut self, id: u64) -> &mut Self {
        self._exclude_ids.push(id);
        self
    }

    pub fn exclude_ids<I: Iterator<Item = u64>>(&mut self, ids: I) -> &mut Self {
        self._exclude_ids.extend(ids);
        self
    }

    pub fn mark_sensitive(&mut self) -> &mut Self {
        self._sensitive = true;
        self
    }

    pub fn enable_dm_commands(&mut self) -> &mut Self {
        self._enable_dm_commands = true;
        self
    }

    pub fn fail_dm_commands(&mut self) -> &mut Self {
        self._fail_dm_commands = true;
        self
    }

    pub fn attach(&mut self, url: String) -> &mut Self {
        self._attachment_url = Some(url);
        self
    }
}

#[derive(Clone, Debug)]
pub struct ProfileBuilder {
    _name: Option<String>,
    _url: Option<String>,
    _location: Option<String>,
    _bio: Option<String>,
    _link_color: Option<String>,
}

impl ProfileBuilder {
    pub fn new() -> ProfileBuilder {
        ProfileBuilder{
            _name: None,
            _url: None,
            _location: None,
            _bio: None,
            _link_color: None,
        }
    }

    pub fn name(&mut self, name: String) -> &mut Self {
        self._name = Some(name);
        self
    }

    pub fn url(&mut self, url: String) -> &mut Self {
        self._url = Some(url);
        self
    }

    pub fn location(&mut self, location: String) -> &mut Self {
        self._location = Some(location);
        self
    }

    pub fn bio(&mut self, bio: String) -> &mut Self {
        self._bio = Some(bio);
        self
    }

    pub fn link_color(&mut self, color: String) -> &mut Self {
        self._link_color = Some(color);
        self
    }
}

pub(super) fn new_tweet_request(auth: Authentication, id: u64, twopts: TweetOptions) -> twitter1::TweetRequest {
    twitter1::TweetRequest{
        auth: Some(auth.ser()),
        id: id,
        twopts: Some(twopts.ser()),
    }
}

pub(super) fn new_tweets_request(auth: Authentication, ids: Vec<u64>, twopts: TweetOptions) -> twitter1::TweetsRequest {
    twitter1::TweetsRequest{
        auth: Some(auth.ser()),
        ids: ids,
        twopts: Some(twopts.ser()),
    }
}

pub(super) fn new_home_timeline_request(auth: Authentication, tlopts: TimelineOptions, twopts: TweetOptions, replies: bool) -> twitter1::HomeTimelineRequest {
    twitter1::HomeTimelineRequest{
        auth: Some(auth.ser()),
        timeline_options: Some(tlopts.ser(twopts)),
        include_replies: replies,
    }
}

pub(super) fn new_mention_timeline_request(auth: Authentication, tlopts: TimelineOptions, twopts: TweetOptions) -> twitter1::MentionTimelineRequest {
    twitter1::MentionTimelineRequest{
        auth: Some(auth.ser()),
        timeline_options: Some(tlopts.ser(twopts)),
    }
}

pub(super) fn new_user_timeline_request(auth: Authentication, user: UserIdentifier, tlopts: TimelineOptions, twopts: TweetOptions, replies: bool, retweets: bool) -> twitter1::UserTimelineRequest {
    twitter1::UserTimelineRequest{
        auth: Some(auth.ser()),
        timeline_options: Some(tlopts.ser(twopts)),
        include_replies: replies,
        include_retweets: retweets,
        user: Some(user.ser()),
    }
}

pub(super) fn new_publish_tweet_request(auth: Authentication, builder: TweetBuilder, twopts: TweetOptions) -> twitter1::PublishTweetRequest {
    twitter1::PublishTweetRequest{
        auth: Some(auth.ser()),
        text: builder._text,
        auto_populate_reply_metadata: builder._reply_id.is_some(),
        exclude_reply_user_ids: builder._exclude_ids,
        media_ids: builder._media_ids,
        possibly_sensitive: builder._sensitive,
        enable_dm_commands: builder._enable_dm_commands,
        fail_dm_commands: builder._fail_dm_commands,
        twopts: Some(twopts.ser()),
        reply_id: builder._reply_id.map(u64::into),
        attachment_url: builder._attachment_url.map(String::into),
    }
}

pub(super) fn new_update_profile_request(auth: Authentication, builder: ProfileBuilder, entities: bool, statuses: bool) -> twitter1::UpdateProfileRequest {
    twitter1::UpdateProfileRequest{
        auth: Some(auth.ser()),
        include_entities: entities,
        include_statuses: statuses,
        name: builder._name.map(String::into),
        url: builder._url.map(String::into),
        location: builder._location.map(String::into),
        bio: builder._bio.map(String::into),
        link_color: builder._link_color.map(String::into),
    }
}
