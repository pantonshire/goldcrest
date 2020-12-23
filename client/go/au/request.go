package au

import pb "github.com/pantonshire/goldcrest/protocol"

type authentication struct {
  consumerKey, accessToken, secretKey, secretToken string
}

func (auth authentication) ser() *pb.Authentication {
  return &pb.Authentication{
    ConsumerKey: auth.consumerKey,
    AccessToken: auth.accessToken,
    SecretKey:   auth.secretKey,
    SecretToken: auth.secretToken,
  }
}

type TweetMode string

const (
  CompatibilityMode = "Compatibility"
  ExtendedMode      = "Extended"
)

func (m TweetMode) ser() pb.TweetOptions_Mode {
  if m == ExtendedMode {
    return pb.TweetOptions_EXTENDED
  }
  return pb.TweetOptions_COMPAT
}

type TweetOptions struct {
  trimUser          bool
  includeMyRetweet  bool
  includeEntities   bool
  includeExtAltText bool
  includeCardURI    bool
  mode              TweetMode
}

func NewTweetOptions() TweetOptions {
  return TweetOptions{
    trimUser:          false,
    includeMyRetweet:  true,
    includeEntities:   true,
    includeExtAltText: true,
    includeCardURI:    true,
    mode:              ExtendedMode,
  }
}

func (opts TweetOptions) WithTrimUser(b bool) TweetOptions {
  opts.trimUser = b
  return opts
}

func (opts TweetOptions) WithMyRetweet(b bool) TweetOptions {
  opts.includeMyRetweet = b
  return opts
}

func (opts TweetOptions) WithEntities(b bool) TweetOptions {
  opts.includeEntities = b
  return opts
}

func (opts TweetOptions) WithAltText(b bool) TweetOptions {
  opts.includeExtAltText = b
  return opts
}

func (opts TweetOptions) WithCardURI(b bool) TweetOptions {
  opts.includeCardURI = b
  return opts
}

func (opts TweetOptions) WithMode(m TweetMode) TweetOptions {
  opts.mode = m
  return opts
}

func (opts TweetOptions) ser() *pb.TweetOptions {
  return &pb.TweetOptions{
    TrimUser:          opts.trimUser,
    IncludeMyRetweet:  opts.includeMyRetweet,
    IncludeEntities:   opts.includeEntities,
    IncludeExtAltText: opts.includeExtAltText,
    IncludeCardUri:    opts.includeCardURI,
    Mode:              opts.mode.ser(),
  }
}

type UserIdentifier interface {
  serIntoUserTimelineRequest(req *pb.UserTimelineRequest)
}

type userIdentifierID uint64

func UserID(id uint64) UserIdentifier {
  return userIdentifierID(id)
}

func (uid userIdentifierID) serIntoUserTimelineRequest(req *pb.UserTimelineRequest) {
  req.User = &pb.UserTimelineRequest_UserId{UserId: uint64(uid)}
}

type userIdentifierHandle string

func UserHandle(handle string) UserIdentifier {
  return userIdentifierHandle(handle)
}

func (uid userIdentifierHandle) serIntoUserTimelineRequest(req *pb.UserTimelineRequest) {
  req.User = &pb.UserTimelineRequest_UserHandle{UserHandle: string(uid)}
}

type TweetComposer struct {
  text              string
  replyID           *uint64
  excludeUserIDs    []uint64
  attachmentURL     *string
  mediaIDs          []uint64
  possiblySensitive bool
  enableDMCommands  bool
  failDMCommands    bool
}

func NewTweetComposer(text string) TweetComposer {
  return TweetComposer{
    text: text,
  }
}

func (com TweetComposer) ReplyTo(tweetID uint64, excludeUserIDs ...uint64) TweetComposer {
  com.replyID = new(uint64)
  *com.replyID = tweetID
  com.excludeUserIDs = make([]uint64, len(excludeUserIDs))
  copy(com.excludeUserIDs, excludeUserIDs)
  return com
}

func (com TweetComposer) WithAttachment(url string) TweetComposer {
  com.attachmentURL = new(string)
  *com.attachmentURL = url
  return com
}

func (com TweetComposer) WithMedia(ids ...uint64) TweetComposer {
  com.mediaIDs = make([]uint64, len(ids))
  copy(com.mediaIDs, ids)
  return com
}

func (com TweetComposer) WithSensitive(sensitive bool) TweetComposer {
  com.possiblySensitive = sensitive
  return com
}

func (com TweetComposer) WithEnableDMCommands(enabled bool) TweetComposer {
  com.enableDMCommands = enabled
  return com
}

func (com TweetComposer) WithFailDMCommands(fail bool) TweetComposer {
  com.failDMCommands = fail
  return com
}

func (com TweetComposer) ser(auth authentication, twopts TweetOptions) *pb.PublishTweetRequest {
  req := pb.PublishTweetRequest{
    Auth:              auth.ser(),
    Text:              com.text,
    MediaIds:          com.mediaIDs,
    PossiblySensitive: com.possiblySensitive,
    EnableDmCommands:  com.enableDMCommands,
    FailDmCommands:    com.failDMCommands,
    Twopts:            twopts.ser(),
  }
  if com.replyID != nil {
    req.Reply = &pb.PublishTweetRequest_ReplyId{ReplyId: *com.replyID}
    req.AutoPopulateReplyMetadata = true
    if com.excludeUserIDs != nil {
      req.ExcludeReplyUserIds = com.excludeUserIDs
    }
  }
  if com.attachmentURL != nil {
    req.Attachment = &pb.PublishTweetRequest_AttachmentUrl{AttachmentUrl: *com.attachmentURL}
  }
  return &req
}

type ProfileUpdater struct {
  name             *string
  url              *string
  location         *string
  bio              *string
  profileLinkColor *string
}

func NewProfileUpdater() ProfileUpdater {
  return ProfileUpdater{}
}

func (pu ProfileUpdater) WithName(name string) ProfileUpdater {
  pu.name = new(string)
  *pu.name = name
  return pu
}

func (pu ProfileUpdater) WithURL(url string) ProfileUpdater {
  pu.url = new(string)
  *pu.url = url
  return pu
}

func (pu ProfileUpdater) WithLocation(location string) ProfileUpdater {
  pu.location = new(string)
  *pu.location = location
  return pu
}

func (pu ProfileUpdater) WithBio(bio string) ProfileUpdater {
  pu.bio = new(string)
  *pu.bio = bio
  return pu
}

func (pu ProfileUpdater) WithProfileLinkColor(color string) ProfileUpdater {
  pu.profileLinkColor = new(string)
  *pu.profileLinkColor = color
  return pu
}

func (pu ProfileUpdater) ser(auth authentication, includeEntities, includeStatuses bool) *pb.UpdateProfileRequest {
  req := pb.UpdateProfileRequest{
    Auth:            auth.ser(),
    IncludeEntities: includeEntities,
    IncludeStatuses: includeStatuses,
  }
  if pu.name != nil {
    req.UpdateName = &pb.UpdateProfileRequest_Name{Name: *pu.name}
  }
  if pu.url != nil {
    req.UpdateUrl = &pb.UpdateProfileRequest_Url{Url: *pu.url}
  }
  if pu.location != nil {
    req.UpdateLocation = &pb.UpdateProfileRequest_Location{Location: *pu.location}
  }
  if pu.bio != nil {
    req.UpdateBio = &pb.UpdateProfileRequest_Bio{Bio: *pu.bio}
  }
  if pu.profileLinkColor != nil {
    req.UpdateProfileLinkColor = &pb.UpdateProfileRequest_ProfileLinkColor{ProfileLinkColor: *pu.profileLinkColor}
  }
  return &req
}
