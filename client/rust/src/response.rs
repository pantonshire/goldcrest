use crate::twitter1;

pub trait Response<T> {
    fn response_result(self) -> Option<Result<T, twitter1::Error>>;
}

impl Response<twitter1::Tweet> for twitter1::TweetResponse {
    fn response_result(self) -> Option<Result<twitter1::Tweet, twitter1::Error>> {
        Some(match self.response? {
            twitter1::tweet_response::Response::Tweet(msg) => Ok(msg),
            twitter1::tweet_response::Response::Error(err) => Err(err),
        })
    }
}

impl Response<twitter1::Tweets> for twitter1::TweetsResponse {
    fn response_result(self) -> Option<Result<twitter1::Tweets, twitter1::Error>> {
        Some(match self.response? {
            twitter1::tweets_response::Response::Tweets(msg) => Ok(msg),
            twitter1::tweets_response::Response::Error(err) => Err(err),
        })
    }
}

impl Response<twitter1::User> for twitter1::UserResponse {
    fn response_result(self) -> Option<Result<twitter1::User, twitter1::Error>> {
        Some(match self.response? {
            twitter1::user_response::Response::User(msg) => Ok(msg),
            twitter1::user_response::Response::Error(err) => Err(err),
        })
    }
}
