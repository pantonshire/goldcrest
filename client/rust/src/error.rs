use std::error::Error;

pub type RequestResult<T> = Result<T, RequestError>;

#[derive(Debug)]
pub enum RequestError {
    Tonic(Box<tonic::Status>),
    Deserialization(Box<DeserializationError>),
    Client(Box<ClientError>),
}

impl std::fmt::Display for RequestError {
    fn fmt(&self, fmt: &mut std::fmt::Formatter) -> std::fmt::Result {
        match self {
            RequestError::Tonic(err)           => err.fmt(fmt),
            RequestError::Deserialization(err) => err.fmt(fmt),
            RequestError::Client(err)          => err.fmt(fmt),
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

impl From<ClientError> for RequestError {
    fn from(err: ClientError) -> Self {
        RequestError::Client(Box::new(err))
    }
}

pub(crate) type DeserializationResult<T> = Result<T, DeserializationError>;

#[derive(Debug)]
pub enum DeserializationError {
    FieldMissing,
    FieldOutOfRange,
}

impl std::fmt::Display for DeserializationError {
    fn fmt(&self, fmt: &mut std::fmt::Formatter) -> std::fmt::Result {
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
