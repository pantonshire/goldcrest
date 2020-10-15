GO := go
GO_BUILD := $(GO) build
GO_TEST := $(GO) test
PROTOC := protoc

PACKAGE := github.com/pantonshire/goldcrest/twitter1
MAIN := cmd/main/main.go
BUILD := build
EXEC := goldcrest

BUILD_TARGETS := linux-amd64 linux-arm linux-arm64 darwin-amd64 windows-amd64 windows-arm
FULL_BUILD_TARGETS := $(addprefix build-,$(BUILD_TARGETS))

PROTO_PATH := proto
PROTO_SOURCE := $(wildcard $(PROTO_PATH)/*.proto)
PROTO_BUILD := $(PROTO_SOURCE:.proto=.pb.go)

build: buildpath
	$(GO_BUILD) -v -o $(BUILD)/$(EXEC) $(MAIN)
	cp $(BUILD)/$(EXEC) $(EXEC)

test:
	$(GO_TEST) -v $(PACKAGE)

dist: $(FULL_BUILD_TARGETS)

build-%: buildpath
	GOOS=$(word 1,$(subst -, ,$*)) GOARCH=$(word 2,$(subst -, ,$*)) $(GO_BUILD) -o $(BUILD)/$(EXEC)-$* $(MAIN)

buildpath:
	@ mkdir -p $(BUILD)

clean:
	rm -rf $(BUILD)

proto: $(PROTO_BUILD)

$(PROTO_PATH)/%.pb.go: $(PROTO_PATH)/%.proto
	$(PROTOC) -I $(PROTO_PATH) --go_out=plugins=grpc,paths=source_relative:$(PROTO_PATH) $<

clean-proto:
	rm $(wildcard $(PROTO_PATH)/*.pb.go)

.PHONY: buildpath build test dist build-% clean proto clean-proto