GO := go
GO_BUILD := $(GO) build
PROTOC := protoc

MAIN := cmd/main/main.go
BUILD := build
EXEC := goldcrest

BUILD_TARGETS := linux-amd64 linux-arm linux-arm64 darwin-amd64 windows-amd64 windows-arm
FULL_BUILD_TARGETS := $(addprefix build-,$(BUILD_TARGETS))

PROTO_SRC := proto
PROTO_BUILD := rpc

PROTO_SOURCES := $(wildcard $(PROTO_SRC)/*.proto)
PROTO_GO := $(addprefix $(PROTO_BUILD)/,$(notdir $(PROTO_SOURCES:.proto=.go)))

build: buildpath
	$(GO_BUILD) -v -o $(BUILD)/$(EXEC) $(MAIN)
	cp $(BUILD)/$(EXEC) $(EXEC)

dist: $(FULL_BUILD_TARGETS)

build-%: buildpath
	GOOS=$(word 1,$(subst -, ,$*)) GOARCH=$(word 2,$(subst -, ,$*)) $(GO_BUILD) -o $(BUILD)/$(EXEC)-$* $(MAIN)

buildpath:
	@ mkdir -p $(BUILD)

clean:
	rm -rf $(BUILD)

proto: $(PROTO_GO)

$(PROTO_BUILD)/%.go: $(PROTO_SRC)/%.proto
	$(PROTOC) -I $(PROTO_SRC) --go_out=plugins=grpc,paths=source_relative:$(PROTO_BUILD) $<

clean-proto:
	rm $(wildcard $(PROTO_BUILD)/*.pb.go)

.PHONY: buildpath build dist build-% clean proto clean-proto