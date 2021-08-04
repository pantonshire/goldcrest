![latest](https://img.shields.io/github/v/release/pantonshire/goldcrest?include_prereleases&label=latest)

# Goldcrest Twitter API Proxy
Goldcrest is a proxy server using gRPC for interacting with the Twitter API v1.1. Its main focus is on providing centralised
rate-limit tracking so that several processes can concurrently use the Twitter API without having to worry about rate-limits.

Currently, there are clients in [Go](client/go/au) and [Rust](client/rust).

## Supported endpoints
Goldcrest currently supports the following Twitter API endpoints:  

| Twitter API endpoint         | gRPC method          |
|------------------------------|----------------------|
| `statuses/update`            | `PublishTweet`       |
| `statuses/retweet`           | `RetweetTweet`       |
| `search/tweets`              | `SearchTweets`       |
| `statuses/unretweet`         | `UnretweetTweet`     |
| `statuses/show`              | `GetTweet`           |
| `statuses/lookup`            | `GetTweets`          |
| `statuses/destroy`           | `DeleteTweet`        |
| `statuses/home_timeline`     | `GetHomeTimeline`    |
| `statuses/mentions_timeline` | `GetMentionTimeline` |
| `statuses/user_timeline`     | `GetUserTimeline`    |
| `favorites/create`           | `LikeTweet`          |
| `favorites/destroy`          | `UnlikeTweet`        |
| `account/update_profile`     | `UpdateProfile`      |

## Setup
### Docker
Pre-built images are available on [Docker Hub](https://hub.docker.com/r/pantonshire/goldcrest).

```sh
docker run -d -p 127.0.0.1:8000:80 pantonshire/goldcrest
```

Currently, images are available for amd64 and arm64. If you're using a different architecture, you'll probably need to build the image yourself:

```sh
docker build -t goldcrest https://github.com/Pantonshire/goldcrest.git#main
docker run -d -p 127.0.0.1:8000:80 goldcrest
```

The configuration file is located at `/etc/goldcrest/config.yaml` in the container.

### Building from source
1. To compile Goldcrest from source, you will first need the following:
    - [Go](https://golang.org/dl/) ≥ 1.15
    - [Protocol Buffers Compiler v3](https://developers.google.com/protocol-buffers/docs/downloads)
    - [protoc-gen-go](https://pkg.go.dev/google.golang.org/protobuf/cmd/protoc-gen-go)  
    (`go install google.golang.org/protobuf/cmd/protoc-gen-go`)
2. Run `make proto` from the repository root.
3. Run `make` from the repository root.
4. `cp default.goldcrest.yaml goldcrest.yaml` to get a correctly-named config file.
