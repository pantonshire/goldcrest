use crate::twitter1::{self, twitter_client::TwitterClient};
use crate::response;
use crate::data::{self, Deserialize};

// use std::sync::{Arc, Mutex};
use std::future::Future;
use tonic::transport::{Endpoint, Channel};
use chrono::prelude::*;
use tokio::time;

#[derive(Clone)]
pub struct Client {
    // au_client: Arc<Mutex<TwitterClient<Channel>>>,
    au_client: TwitterClient<Channel>,
    wait_timeout: chrono::Duration,
}

macro_rules! request {
    // ($f:ident) => { |t,r| async move { t.lock().unwrap().$f(r).await } }
    ($f:ident) => { |mut t, r| async move { t.$f(r).await } }
}

impl Client {
    pub async fn new(uri: &str, request_timeout: chrono::Duration, wait_timeout: chrono::Duration) -> Result<Client, Box<dyn std::error::Error>> {
        let channel = Endpoint::from_shared(uri.to_string())?
            .timeout(request_timeout.to_std()?)
            .connect()
            .await?;
        Ok(Client{
            // au_client: Arc::new(Mutex::new(TwitterClient::new(channel))),
            au_client: TwitterClient::new(channel),
            wait_timeout: wait_timeout,
        })
    }

    async fn request<R, S, T, F, Fut>(&self, req: R, rf: F) -> Result<T, Box<dyn std::error::Error>>
    where
        Fut: Future<Output=Result<tonic::Response<S>, tonic::Status>>,
        R: Clone,
        S: response::Response<T>,
        // F: Fn(Arc<Mutex<TwitterClient<Channel>>>, tonic::Request<R>) -> Fut,
        F: Fn(TwitterClient<Channel>, tonic::Request<R>) -> Fut,
    {
        let deadline = Utc::now() + self.wait_timeout;

        loop {
            // let resp = rf(Arc::clone(&self.au_client), tonic::Request::new(req.clone())).await?;
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
                        
                        if deadline < retry {
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

    pub async fn get_tweet(&self, req: twitter1::TweetRequest) -> Result<data::Tweet, Box<dyn std::error::Error>> {
        // self.request(req, |t, r| async move { t.lock().unwrap().get_tweet(r).await })
        //     .await
        //     .and_then(|t| Ok(t.des()?))
        self.request(req, request!(get_tweet))
            .await
            .and_then(|t| Ok(t.des()?))
        // self.request(req, |mut t, r| async move { t.get_tweet(r).await })
        //     .await
        //     .and_then(|t| Ok(t.des()?))
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
}

impl std::fmt::Display for ClientError {
    fn fmt(&self, fmt: &mut std::fmt::Formatter) -> std::fmt::Result {
        write!(fmt, "[Goldcrest client] {}", match self {
            ClientError::UnknownError(code) => format!("unknown error code {}", code),
            ClientError::InvalidResponse    => "invalid response".to_string(),
            ClientError::RetryTimeout       => "timed out waiting for rate limit to reset".to_string(),
            ClientError::RetryUnknown       => "rate limit reached, but reset time unknown".to_string(),
            ClientError::TwitterError(s)    => format!("Twitter server-side error: {}", s),
            ClientError::BadRequest(s)      => format!("bad request to Twitter: {}", s),
            ClientError::BadResponse(s)     => format!("bad response from Twitter: {}", s),
        })
    }
}

impl std::error::Error for ClientError {}
