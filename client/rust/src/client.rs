use crate::response;
use crate::request;
use crate::serialize;
use crate::deserialize::Deserialize;
use crate::data;
use crate::error::{RequestResult, ClientError};
use crate::twitter1::{self, twitter_client::TwitterClient};

use std::future::Future;
use chrono::prelude::*;
use tonic::transport::{Endpoint, Channel};
use tokio::time;

pub struct ClientBuilder {
    par_scheme: String,
    par_host: String,
    par_port: u32,
    par_request_timeout: chrono::Duration,
    par_wait_timeout: chrono::Duration,
    par_concurrency_limit: Option<usize>,
    par_authentication: Option<request::Authentication>,
}

impl ClientBuilder {
    pub fn new() -> ClientBuilder {
        ClientBuilder{
            par_scheme: "http".to_owned(),
            par_host: "localhost".to_owned(),
            par_port: 8000,
            par_request_timeout: chrono::Duration::zero(),
            par_wait_timeout: chrono::Duration::zero(),
            par_concurrency_limit: None,
            par_authentication: None,
        }
    }

    pub async fn connect(self) -> Result<Client, Box<dyn std::error::Error>> {
        if self.par_authentication.is_none() {
            return Err(Box::new(ClientError::Unauthenticated));
        }
        let uri = format!("{}://{}:{}", self.par_scheme, self.par_host, self.par_port);
        let mut ep = Endpoint::from_shared(uri)?;
        if !self.par_request_timeout.is_zero() {
            ep = ep.timeout(self.par_request_timeout.to_std()?);
        }
        if self.par_concurrency_limit.is_some() {
            ep = ep.concurrency_limit(self.par_concurrency_limit.unwrap());
        }
        let channel = ep.connect().await?;
        Ok(Client{
            au_client: TwitterClient::new(channel),
            wait_timeout: if self.par_wait_timeout.is_zero() {
                None
            } else {
                Some(self.par_wait_timeout)
            },
            authentication: self.par_authentication.unwrap(),
        })
    }

    pub fn scheme(&mut self, scheme: &str) -> &mut Self {
        self.par_scheme = scheme.to_owned();
        self
    }

    pub fn host(&mut self, host: &str) -> &mut Self {
        self.par_host = host.to_owned();
        self
    }

    pub fn port(&mut self, port: u32) -> &mut Self {
        self.par_port = port;
        self
    }

    pub fn socket(&mut self, host: &str, port: u32) -> &mut Self {
        self.host(host)
            .port(port)
    }

    pub fn authenticate(&mut self, auth: request::Authentication) -> &mut Self {
        self.par_authentication = Some(auth);
        self
    }

    pub fn request_timeout(&mut self, timeout: chrono::Duration) -> &mut Self {
        self.par_request_timeout = timeout;
        self
    }

    pub fn wait_timeout(&mut self, timeout: chrono::Duration) -> &mut Self {
        self.par_wait_timeout = timeout;
        self
    }

    pub fn concurrency_limit(&mut self, limit: usize) -> &mut Self {
        self.par_concurrency_limit = Some(limit);
        self
    }
}

#[derive(Clone)]
pub struct Client {
    au_client: TwitterClient<Channel>,
    wait_timeout: Option<chrono::Duration>,
    authentication: request::Authentication,
}

macro_rules! request_cls {
    ($f:ident) => {
        |mut t, r| async move { t.$f(r).await }
    }
}

macro_rules! request {
    ($client:expr, $req:expr, $f:ident) => {
        $client.request($req, request_cls!($f))
            .await
            .and_then(|t| Ok(t.des()?))
    }
}

macro_rules! tweet_request {
    ($client:expr, $id:expr, $twopts: expr, $f:ident) => {
        request!($client, crate::serialize::ser_tweet_request($client.authentication.clone(), $id, $twopts), $f)
    }
}

impl Client {
    pub async fn get_tweet(&self, id: u64, twopts: request::TweetOptions) -> RequestResult<data::Tweet> {
        tweet_request!(self, id, twopts, get_tweet)
    }

    pub async fn get_tweets(&self, ids: Vec<u64>, twopts: request::TweetOptions) -> RequestResult<Vec<data::Tweet>> {
        let req = serialize::ser_tweets_request(self.authentication.clone(), ids, twopts);
        request!(self, req, get_tweets)
    }

    pub async fn search_tweets(&self, search_opts: request::SearchOptions, tweet_opts: request::TweetOptions, timeline_opts: request::TimelineOptions) -> RequestResult<Vec<data::Tweet>> {
        let req = serialize::ser_search_request(self.authentication.clone(), search_opts, tweet_opts, timeline_opts);
        request!(self, req, search_tweets)
    }

    pub async fn like(&self, id: u64, twopts: request::TweetOptions) -> RequestResult<data::Tweet> {
        tweet_request!(self, id, twopts, like_tweet)
    }

    pub async fn unlike(&self, id: u64, twopts: request::TweetOptions) -> RequestResult<data::Tweet> {
        tweet_request!(self, id, twopts, unlike_tweet)
    }

    pub async fn retweet(&self, id: u64, twopts: request::TweetOptions) -> RequestResult<data::Tweet> {
        tweet_request!(self, id, twopts, retweet_tweet)
    }

    pub async fn unretweet(&self, id: u64, twopts: request::TweetOptions) -> RequestResult<data::Tweet> {
        tweet_request!(self, id, twopts, unretweet_tweet)
    }

    pub async fn delete_tweet(&self, id: u64, twopts: request::TweetOptions) -> RequestResult<data::Tweet> {
        tweet_request!(self, id, twopts, delete_tweet)
    }

    pub async fn home_timeline(&self, tlopts: request::TimelineOptions, twopts: request::TweetOptions, replies: bool) -> RequestResult<Vec<data::Tweet>> {
        let req = serialize::ser_home_timeline_request(self.authentication.clone(), tlopts, twopts, replies);
        request!(self, req, get_home_timeline)
    }

    pub async fn mention_timeline(&self, tlopts: request::TimelineOptions, twopts: request::TweetOptions) -> RequestResult<Vec<data::Tweet>> {
        let req = serialize::ser_mention_timeline_request(self.authentication.clone(), tlopts, twopts);
        request!(self, req, get_mention_timeline)
    }

    pub async fn user_timeline(&self, user: request::UserIdentifier, tlopts: request::TimelineOptions, twopts: request::TweetOptions, replies: bool, retweets: bool) -> RequestResult<Vec<data::Tweet>> {
        let req = serialize::ser_user_timeline_request(self.authentication.clone(), user, tlopts, twopts, replies, retweets);
        request!(self, req, get_user_timeline)
    }

    pub async fn publish(&self, tweet: request::TweetBuilder, twopts: request::TweetOptions) -> RequestResult<data::Tweet> {
        let req = serialize::ser_publish_tweet_request(self.authentication.clone(), tweet, twopts);
        request!(self, req, publish_tweet)
    }

    pub async fn update_profile(&self, profile: request::ProfileBuilder, entities: bool, statuses: bool) -> RequestResult<data::User> {
        let req = serialize::ser_update_profile_request(self.authentication.clone(), profile, entities, statuses);
        request!(self, req, update_profile)
    }

    async fn request<R, S, T, F, Fut>(&self, req: R, rf: F) -> RequestResult<T>
    where
        Fut: Future<Output=Result<tonic::Response<S>, tonic::Status>>,
        R: Clone,
        S: response::Response<T>,
        F: Fn(TwitterClient<Channel>, tonic::Request<R>) -> Fut,
    {
        let deadline = self.wait_timeout.map(|t| Utc::now() + t);

        loop {
            let resp = rf(self.au_client.clone(), tonic::Request::new(req.clone())).await?;
            let meta = resp.metadata().clone();

            match resp.into_inner().response_result() {
                None => return Err(ClientError::InvalidResponse.into()),

                Some(Ok(msg)) => return Ok(msg),

                Some(Err(err)) => {
                    if err.code == twitter1::error::Code::RateLimit as i32 {
                        let retry = meta.get("retry")
                            .ok_or(ClientError::RetryUnknown)?;
                        let retry = std::str::from_utf8(retry.as_bytes())
                            .map_err(|_| ClientError::RetryUnknown)?
                            .parse::<i64>()
                            .map_err(|_| ClientError::RetryUnknown)?;
                        let retry = DateTime::<Utc>::from_utc(NaiveDateTime::from_timestamp(retry, 0), Utc);
                        
                        if deadline.is_some() && deadline.unwrap() < retry {
                            return Err(ClientError::RetryTimeout.into());
                        }

                        match (retry - Utc::now()).to_std() {
                            Ok(wait_time) => time::sleep(wait_time).await,
                            Err(_)        => (), //no need to wait if duration was negative
                        }

                        continue
                    }

                    return Err(if err.code == twitter1::error::Code::TwitterError as i32 {
                        ClientError::TwitterError(err.message).into()
                    } else if err.code == twitter1::error::Code::BadRequest as i32 {
                        ClientError::BadRequest(err.message).into()
                    } else if err.code == twitter1::error::Code::BadResponse as i32 {
                        ClientError::BadResponse(err.message).into()
                    } else {
                        ClientError::UnknownError(err.code).into()
                    })
                },
            }
        }
    }
}
