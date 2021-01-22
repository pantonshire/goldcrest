pub mod client;
pub mod data;
pub mod request;
pub mod error;

mod optional;
mod deserialize;
mod serialize;
mod response;

pub use client::{Client, ClientBuilder};
pub use request::*;

pub mod twitter1 {
    tonic::include_proto!("twitter1");
}
