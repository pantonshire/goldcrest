use crate::twitter1;

trait Serialise<T: Sized> {
    fn ser(self) -> T;
}

#[derive(Clone, Debug)]
pub struct Authentication {
    consumer_key: String,
    consumer_secret: String,
    token: String,
    token_secret: String,
}

impl Serialise<twitter1::Authentication> for Authentication {
    fn ser(self) -> twitter1::Authentication {
        twitter1::Authentication{
            consumer_key: self.consumer_key,
            secret_key: self.consumer_secret,
            access_token: self.token,
            secret_token: self.token_secret,
        }
    }
}
