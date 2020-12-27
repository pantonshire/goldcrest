use crate::twitter1::{self, twitter_client::TwitterClient};
use crate::response;
use crate::data::{self, Deserialize};

use std::sync::{Arc, Mutex};
use std::future::Future;
use tonic::transport::{Endpoint, Channel};
use chrono::prelude::*;

pub struct Client {
    au_client: Arc<Mutex<TwitterClient<Channel>>>,
    wait_timeout: chrono::Duration,
}

impl Client {
    pub async fn new(uri: &str, request_timeout: chrono::Duration, wait_timeout: chrono::Duration) -> Result<Client, Box<dyn std::error::Error>> {
        let channel = Endpoint::from_shared(uri.to_string())?
            .timeout(request_timeout.to_std()?)
            .connect()
            .await?;
        Ok(Client{
            au_client: Arc::new(Mutex::new(TwitterClient::new(channel))),
            wait_timeout: wait_timeout,
        })
    }

    async fn request<R, S, T, F, Fut>(&mut self, req: R, rf: F) -> Result<T, Box<dyn std::error::Error>>
    where
        Fut: Future<Output=Result<tonic::Response<S>, tonic::Status>>,
        R: Clone,
        S: response::Response<T>,
        F: Fn(Arc<Mutex<TwitterClient<Channel>>>, tonic::Request<R>) -> Fut,
    {
        //TODO: retries, deadline

        let deadline = Utc::now() + self.wait_timeout;

        loop {
            let resp = rf(Arc::clone(&self.au_client), tonic::Request::new(req.clone())).await?;

            let meta = resp.metadata().clone();

            match resp.into_inner().response_result() {
                None => return Err(Box::new(ClientError::InvalidResponse)),
                Some(Ok(msg)) => return Ok(msg),
                Some(Err(err)) => {
                    if err.code == twitter1::error::Code::RateLimit as i32 {

                    } else {
                        
                    }
                },
            }
        }

        // let meta = resp.metadata();
        // let msg = resp.into_inner();
        // let msg = resp.get_ref();
        // meta.get("");

        panic!("Not implemented")
    }

    pub async fn get_tweet(&mut self, req: twitter1::TweetRequest) -> Result<data::Tweet, Box<dyn std::error::Error>> {
        self.request(req, |t, r| async move { t.lock().unwrap().get_tweet(r).await })
            .await
            .and_then(|t| Ok(t.des()?))
    }
}

#[derive(Debug)]
pub enum ClientError {
    UnknownError(u32),
    InvalidResponse,
}

impl std::fmt::Display for ClientError {
    fn fmt(&self, fmt: &mut std::fmt::Formatter) -> std::fmt::Result {
        write!(fmt, "Client error: {}", match self {
            ClientError::InvalidResponse => "invalid response",
        })
    }
}

impl std::error::Error for ClientError {}
