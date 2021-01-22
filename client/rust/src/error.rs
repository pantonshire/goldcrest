use std::error::Error;
use std::fmt::{self, Display, Formatter};

pub type ConnectionResult<T> = Result<T, ConnectionError>;

#[derive(Debug)]
pub enum ConnectionError {
    Unauthenticated,
    InvalidUri,
    Channel(Box<tonic::transport::Error>),
}

impl Display for ConnectionError {
    fn fmt(&self, fmt: &mut Formatter) -> fmt::Result {
        match self {
            ConnectionError::Unauthenticated => write!(fmt, "not authenticated"),
            ConnectionError::InvalidUri      => write!(fmt, "invalid URI"),
            ConnectionError::Channel(err)    => err.fmt(fmt),
        }
    }
}

impl Error for ConnectionError {}

impl From<tonic::transport::Error> for ConnectionError {
    fn from(err: tonic::transport::Error) -> Self {
        ConnectionError::Channel(Box::new(err))
    }
}

pub type RequestResult<T> = Result<T, RequestError>;

#[derive(Debug)]
pub enum RequestError {
    Tonic(Box<tonic::Status>),
    Deserialization(Box<DeserializationError>),
    Client(Box<TwitterError>),
    RetryTimeout,
    RetryUnknown,
}

impl Display for RequestError {
    fn fmt(&self, fmt: &mut Formatter) -> fmt::Result {
        match self {
            RequestError::Tonic(err)           => err.fmt(fmt),
            RequestError::Deserialization(err) => err.fmt(fmt),
            RequestError::Client(err)          => err.fmt(fmt),
            RequestError::RetryTimeout         => write!(fmt, "deadline exceeded waiting for rate limit to reset"),
            RequestError::RetryUnknown         => write!(fmt, "rate limit reached, but retry time unknown"),
        }
    }
}

impl Error for RequestError {}

impl From<tonic::Status> for RequestError {
    fn from(err: tonic::Status) -> Self {
        RequestError::Tonic(Box::new(err))
    }
}

impl From<DeserializationError> for RequestError {
    fn from(err: DeserializationError) -> Self {
        RequestError::Deserialization(Box::new(err))
    }
}

impl From<TwitterError> for RequestError {
    fn from(err: TwitterError) -> Self {
        RequestError::Client(Box::new(err))
    }
}

pub(crate) type DeserializationResult<T> = Result<T, DeserializationError>;

#[derive(Debug)]
pub enum DeserializationError {
    FieldMissing,
    FieldOutOfRange,
}

impl Display for DeserializationError {
    fn fmt(&self, fmt: &mut Formatter) -> fmt::Result {
        write!(fmt, "[Goldcrest deserialization] {}", match self {
            DeserializationError::FieldMissing    => "missing field",
            DeserializationError::FieldOutOfRange => "field out of range",
        })
    }
}

impl Error for DeserializationError {}

pub(crate) trait Exists<T> {
    fn exists(self) -> DeserializationResult<T>;
}

impl<T> Exists<T> for Option<T> {
    fn exists(self) -> DeserializationResult<T> {
        self.ok_or(DeserializationError::FieldMissing)
    }
}

#[derive(Debug)]
pub enum TwitterError {
    UnknownError(i32),
    InvalidResponse,
    TwitterError(String),
    BadRequest(String),
    BadResponse(String),
}

impl Display for TwitterError {
    fn fmt(&self, fmt: &mut Formatter) -> fmt::Result {
        write!(fmt, "Goldcrest Twitter error: {}", match self {
            TwitterError::UnknownError(code) => format!("unknown error code {}", code),
            TwitterError::InvalidResponse    => "invalid response".to_owned(),
            TwitterError::TwitterError(s)    => format!("Twitter server-side error: {}", s),
            TwitterError::BadRequest(s)      => format!("bad request to Twitter: {}", s),
            TwitterError::BadResponse(s)     => format!("bad response from Twitter: {}", s),
        })
    }
}

impl Error for TwitterError {}
