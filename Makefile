BUILD_VERSION   := $(version)
BUILD_TIME      := $(shell date "+%F %T")
COMMIT_SHA1     := $(shell git rev-parse HEAD)

all:
	gox -osarch="darwin/amd64 linux/386 linux/amd64" \
		-output="dist/{{.Dir}}_{{.OS}}_{{.Arch}}" \
    	-ldflags   "-X 'github.com/Gozap/certmonitor/cmd.Version=${BUILD_VERSION}' \
					-X 'github.com/Gozap/certmonitor/cmd.BuildTime=${BUILD_TIME}' \
					-X 'github.com/Gozap/certmonitor/cmd.CommitID=${COMMIT_SHA1}'"

release: all
	ghr -u gozap -t $(GITHUB_RELEASE_TOKEN) -replace -recreate --debug $(version) dist

docker: all
	docker build -t gozap/certmonitor:$(version) .

clean:
	rm -rf dist

install:
	go install

.PHONY : all release docker clean install
