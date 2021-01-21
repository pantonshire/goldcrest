use crate::request::*;
use crate::twitter1;

pub(crate) fn ser_tweet_request(auth: Authentication, id: u64, tweet_opts: TweetOptions) -> twitter1::TweetRequest {
    twitter1::TweetRequest{
        auth: Some(ser_authentication(auth)),
        id: id,
        twopts: Some(ser_tweet_options(tweet_opts)),
    }
}

pub(crate) fn ser_tweets_request(auth: Authentication, ids: Vec<u64>, tweet_opts: TweetOptions) -> twitter1::TweetsRequest {
    twitter1::TweetsRequest{
        auth: Some(ser_authentication(auth)),
        ids: ids,
        twopts: Some(ser_tweet_options(tweet_opts)),
    }
}

pub(crate) fn ser_search_request(auth: Authentication, search_opts: SearchOptions, tweet_opts: TweetOptions, timeline_opts: TimelineOptions) -> twitter1::SearchRequest {
    twitter1::SearchRequest{
        auth: Some(ser_authentication(auth)),
        query: search_opts.par_query,
        geocode: search_opts.par_geocode.map(String::into),
        lang: search_opts.par_lang.map(String::into),
        locale: search_opts.par_locale.map(String::into),
        result_type: ser_search_result_type(search_opts.par_result_type),
        until_timestamp: search_opts.par_until.map(|x| (x as u64).into()),
        timeline_options: Some(ser_timeline_options(timeline_opts, tweet_opts)),
    }
}

pub(crate) fn ser_home_timeline_request(auth: Authentication, timeline_opts: TimelineOptions, tweet_opts: TweetOptions, replies: bool) -> twitter1::HomeTimelineRequest {
    twitter1::HomeTimelineRequest{
        auth: Some(ser_authentication(auth)),
        timeline_options: Some(ser_timeline_options(timeline_opts, tweet_opts)),
        include_replies: replies,
    }
}

pub(crate) fn ser_mention_timeline_request(auth: Authentication, timeline_opts: TimelineOptions, tweet_opts: TweetOptions) -> twitter1::MentionTimelineRequest {
    twitter1::MentionTimelineRequest{
        auth: Some(ser_authentication(auth)),
        timeline_options: Some(ser_timeline_options(timeline_opts, tweet_opts)),
    }
}

pub(crate) fn ser_user_timeline_request(auth: Authentication, user: UserIdentifier, timeline_opts: TimelineOptions, tweet_opts: TweetOptions, replies: bool, retweets: bool) -> twitter1::UserTimelineRequest {
    twitter1::UserTimelineRequest{
        auth: Some(ser_authentication(auth)),
        timeline_options: Some(ser_timeline_options(timeline_opts, tweet_opts)),
        include_replies: replies,
        include_retweets: retweets,
        user: Some(ser_user_identifier(user)),
    }
}

pub(crate) fn ser_publish_tweet_request(auth: Authentication, builder: TweetBuilder, tweet_opts: TweetOptions) -> twitter1::PublishTweetRequest {
    twitter1::PublishTweetRequest{
        auth: Some(ser_authentication(auth)),
        text: builder.par_text,
        auto_populate_reply_metadata: builder.par_reply_id.is_some(),
        exclude_reply_user_ids: builder.par_exclude_ids,
        media_ids: builder.par_media_ids,
        possibly_sensitive: builder.par_sensitive,
        enable_dm_commands: builder.par_enable_dm_commands,
        fail_dm_commands: builder.par_fail_dm_commands,
        twopts: Some(ser_tweet_options(tweet_opts)),
        reply_id: builder.par_reply_id.map(u64::into),
        attachment_url: builder.par_attachment_url.map(String::into),
    }
}

pub(crate) fn ser_update_profile_request(auth: Authentication, builder: ProfileBuilder, entities: bool, statuses: bool) -> twitter1::UpdateProfileRequest {
    twitter1::UpdateProfileRequest{
        auth: Some(ser_authentication(auth)),
        include_entities: entities,
        include_statuses: statuses,
        name: builder.par_name.map(String::into),
        url: builder.par_url.map(String::into),
        location: builder.par_location.map(String::into),
        bio: builder.par_bio.map(String::into),
        link_color: builder.par_link_color.map(String::into),
    }
}

fn ser_authentication(auth: Authentication) -> twitter1::Authentication {
    twitter1::Authentication{
        consumer_key: auth.consumer_key,
        secret_key: auth.consumer_secret,
        access_token: auth.access_token,
        secret_token: auth.token_secret,
    }
}

fn ser_tweet_options(tweet_opts: TweetOptions) -> twitter1::TweetOptions {
    twitter1::TweetOptions{
        trim_user: tweet_opts.par_trim_user,
        include_my_retweet: tweet_opts.par_include_your_retweet,
        include_entities: tweet_opts.par_include_entities,
        include_ext_alt_text: tweet_opts.par_include_ext_alt_text,
        include_card_uri: tweet_opts.par_include_card_uri,
        mode: ser_tweet_mode(tweet_opts.par_mode),
    }
}

fn ser_tweet_mode(mode: TweetMode) -> i32 {
    use twitter1::tweet_options::Mode;
    (match mode {
        TweetMode::Compatibility => Mode::Compat,
        TweetMode::Extended      => Mode::Extended,
    }) as i32
}

fn ser_timeline_options(timeline_options: TimelineOptions, tweet_options: TweetOptions) -> twitter1::TimelineOptions {
    twitter1::TimelineOptions{
        count: timeline_options.par_count,
        min_id: timeline_options.par_min.map(u64::into),
        max_id: timeline_options.par_max.map(u64::into),
        twopts: Some(ser_tweet_options(tweet_options)),
    }
}

fn ser_search_result_type(result_type: SearchResultType) -> i32 {
    use twitter1::search_request::ResultType;
    (match result_type {
        SearchResultType::Mixed   => ResultType::Mixed,
        SearchResultType::Recent  => ResultType::Recent,
        SearchResultType::Popular => ResultType::Popular,
    }) as i32
}

fn ser_user_identifier(user: UserIdentifier) -> twitter1::user_timeline_request::User {
    use twitter1::user_timeline_request::User;
    match user {
        UserIdentifier::Id(id)         => User::UserId(id),
        UserIdentifier::Handle(handle) => User::UserHandle(handle),
    }
}
