BUILD_VERSION   := $(shell cat version)
BUILD_TIME      := $(shell date "+%F %T")
COMMIT_SHA1     := $(shell git rev-parse HEAD)

all:clean
	gox -osarch="darwin/amd64 linux/386 linux/amd64" \
        -output="dist/{{.Dir}}_{{.OS}}_{{.Arch}}" \
    	-ldflags   "-X 'github.com/gozap/certmonitor/cmd.Version=${BUILD_VERSION}' \
                    -X 'github.com/gozap/certmonitor/cmd.BuildTime=${BUILD_TIME}' \
                    -X 'github.com/gozap/certmonitor/cmd.CommitID=${COMMIT_SHA1}'"

release: all
	ghr -u gozap -t ${GITHUB_RELEASE_TOKEN} -replace -recreate --debug ${BUILD_VERSION} dist

docker: all
	docker build -t gozap/certmonitor:${BUILD_VERSION} .

clean:
	rm -rf dist

install:
	go install

.PHONY : all release docker clean install

.EXPORT_ALL_VARIABLES:

GO111MODULE = on
GOPROXY = https://goproxy.io
GOSUMDB = sum.golang.google.cn
