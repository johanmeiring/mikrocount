TARGET_OS = linux
TARGET_ARCH = amd64
OUT_FILE = mikrocount
SOURCE_FILE = mikrocount.go
DOCKER_TAG = "johanmeiring/mikrocount"

.PHONY: deps
deps:
	dep ensure

.PHONY: build
build: deps
	env GOOS=$(TARGET_OS) GOARCH=$(TARGET_ARCH) go build -o $(OUT_FILE) $(SOURCE_FILE)

.PHONY: docker
docker: build
	docker build -t $(DOCKER_TAG) .
