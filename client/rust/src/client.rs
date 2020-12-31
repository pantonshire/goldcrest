use crate::{response, request, deserialize::Deserialize};
use crate::data::{self};
use crate::twitter1::{self, twitter_client::TwitterClient};

use std::future::Future;
use chrono::prelude::*;
use tonic::transport::{Endpoint, Channel};
use tokio::time;

pub type ReqResult<T> = Result<T, Box<dyn std::error::Error>>;

pub struct ClientBuilder {
    _scheme: String,
    _host: String,
    _port: u32,
    _request_timeout: chrono::Duration,
    _wait_timeout: chrono::Duration,
    _authentication: Option<request::Authentication>,
}

impl ClientBuilder {
    pub fn new() -> ClientBuilder {
        ClientBuilder{
            _scheme: "http".to_owned(),
            _host: "localhost".to_owned(),
            _port: 8000,
            _request_timeout: chrono::Duration::zero(),
            _wait_timeout: chrono::Duration::zero(),
            _authentication: None,
        }
    }

    pub async fn connect(self) -> Result<Client, Box<dyn std::error::Error>> {
        if self._authentication.is_none() {
            return Err(Box::new(ClientError::Unauthenticated));
        }
        let uri = format!("{}://{}:{}", self._scheme, self._host, self._port);
        let mut ep = Endpoint::from_shared(uri)?;
        if !self._request_timeout.is_zero() {
            ep = ep.timeout(self._request_timeout.to_std()?);
        }
        let channel = ep.connect().await?;
        Ok(Client{
            au_client: TwitterClient::new(channel),
            wait_timeout: if self._wait_timeout.is_zero() {
                None
            } else {
                Some(self._wait_timeout)
            },
            authentication: self._authentication.unwrap(),
        })
    }

    pub fn scheme(&mut self, scheme: &str) -> &mut Self {
        self._scheme = scheme.to_owned();
        self
    }

    pub fn host(&mut self, host: &str) -> &mut Self {
        self._host = host.to_owned();
        self
    }

    pub fn port(&mut self, port: u32) -> &mut Self {
        self._port = port;
        self
    }

    pub fn socket(&mut self, host: &str, port: u32) -> &mut Self {
        self.host(host)
            .port(port)
    }

    pub fn authenticate(&mut self, auth: request::Authentication) -> &mut Self {
        self._authentication = Some(auth);
        self
    }

    pub fn request_timeout(&mut self, timeout: chrono::Duration) -> &mut Self {
        self._request_timeout = timeout;
        self
    }

    pub fn wait_timeout(&mut self, timeout: chrono::Duration) -> &mut Self {
        self._wait_timeout = timeout;
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
        request!($client, request::new_tweet_request($client.authentication.clone(), $id, $twopts), $f)
    }
}

impl Client {
    pub async fn get_tweet(&self, id: u64, twopts: request::TweetOptions) -> ReqResult<data::Tweet> {
        tweet_request!(self, id, twopts, get_tweet)
    }

    pub async fn get_tweets(&self, ids: Vec<u64>, twopts: request::TweetOptions) -> ReqResult<Vec<data::Tweet>> {
        let req = request::new_tweets_request(self.authentication.clone(), ids, twopts);
        request!(self, req, get_tweets)
    }

    pub async fn like(&self, id: u64, twopts: request::TweetOptions) -> ReqResult<data::Tweet> {
        tweet_request!(self, id, twopts, like_tweet)
    }

    pub async fn unlike(&self, id: u64, twopts: request::TweetOptions) -> ReqResult<data::Tweet> {
        tweet_request!(self, id, twopts, unlike_tweet)
    }

    pub async fn retweet(&self, id: u64, twopts: request::TweetOptions) -> ReqResult<data::Tweet> {
        tweet_request!(self, id, twopts, retweet_tweet)
    }

    pub async fn unretweet(&self, id: u64, twopts: request::TweetOptions) -> ReqResult<data::Tweet> {
        tweet_request!(self, id, twopts, unretweet_tweet)
    }

    pub async fn delete_tweet(&self, id: u64, twopts: request::TweetOptions) -> ReqResult<data::Tweet> {
        tweet_request!(self, id, twopts, delete_tweet)
    }

    pub async fn home_timeline(&self, tlopts: request::TimelineOptions, twopts: request::TweetOptions, replies: bool) -> ReqResult<Vec<data::Tweet>> {
        let req = request::new_home_timeline_request(self.authentication.clone(), tlopts, twopts, replies);
        request!(self, req, get_home_timeline)
    }

    pub async fn mention_timeline(&self, tlopts: request::TimelineOptions, twopts: request::TweetOptions) -> ReqResult<Vec<data::Tweet>> {
        let req = request::new_mention_timeline_request(self.authentication.clone(), tlopts, twopts);
        request!(self, req, get_mention_timeline)
    }

    pub async fn user_timeline(&self, user: request::UserIdentifier, tlopts: request::TimelineOptions, twopts: request::TweetOptions, replies: bool, retweets: bool) -> ReqResult<Vec<data::Tweet>> {
        let req = request::new_user_timeline_request(self.authentication.clone(), user, tlopts, twopts, replies, retweets);
        request!(self, req, get_user_timeline)
    }

    pub async fn publish(&self, tweet: request::TweetBuilder, twopts: request::TweetOptions) -> ReqResult<data::Tweet> {
        let req = request::new_publish_tweet_request(self.authentication.clone(), tweet, twopts);
        request!(self, req, publish_tweet)
    }

    pub async fn update_profile(&self, profile: request::ProfileBuilder, entities: bool, statuses: bool) -> ReqResult<data::User> {
        let req = request::new_update_profile_request(self.authentication.clone(), profile, entities, statuses);
        request!(self, req, update_profile)
    }

    async fn request<R, S, T, F, Fut>(&self, req: R, rf: F) -> ReqResult<T>
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
                None => return Err(Box::new(ClientError::InvalidResponse)),

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
                            return Err(Box::new(ClientError::RetryTimeout));
                        }

                        match (retry - Utc::now()).to_std() {
                            Ok(wait_time) => time::delay_for(wait_time).await,
                            Err(_)        => (), //no need to wait if duration was negative
                        }

                        continue
                    }

                    return Err(Box::new(if err.code == twitter1::error::Code::TwitterError as i32 {
                        ClientError::TwitterError(err.message)
                    } else if err.code == twitter1::error::Code::BadRequest as i32 {
                        ClientError::BadRequest(err.message)
                    } else if err.code == twitter1::error::Code::BadResponse as i32 {
                        ClientError::BadResponse(err.message)
                    } else {
                        ClientError::UnknownError(err.code)
                    }))
                },
            }
        }
    }
}

#[derive(Debug)]
pub enum ClientError {
    UnknownError(i32),
    InvalidResponse,
    RetryTimeout,
    RetryUnknown,
    TwitterError(String),
    BadRequest(String),
    BadResponse(String),
    Unauthenticated,
}

impl std::fmt::Display for ClientError {
    fn fmt(&self, fmt: &mut std::fmt::Formatter) -> std::fmt::Result {
        write!(fmt, "[Goldcrest client] {}", match self {
            ClientError::UnknownError(code) => format!("unknown error code {}", code),
            ClientError::InvalidResponse    => "invalid response".to_owned(),
            ClientError::RetryTimeout       => "timed out waiting for rate limit to reset".to_owned(),
            ClientError::RetryUnknown       => "rate limit reached, but reset time unknown".to_owned(),
            ClientError::TwitterError(s)    => format!("Twitter server-side error: {}", s),
            ClientError::BadRequest(s)      => format!("bad request to Twitter: {}", s),
            ClientError::BadResponse(s)     => format!("bad response from Twitter: {}", s),
            ClientError::Unauthenticated    => "not authenticated".to_owned(),
        })
    }
}

impl std::error::Error for ClientError {}
