pub mod client;
pub mod data;
pub mod request;

mod response;
mod optional;

pub use client::{Client, ClientBuilder};
pub use request::Authentication;

pub mod twitter1 {
    tonic::include_proto!("twitter1");
}

#[cfg(test)]
mod tests {
    #[test]
    fn it_works() {
        assert_eq!(2 + 2, 4);
    }
}
