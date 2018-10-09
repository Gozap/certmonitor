FROM alpine:3.8

LABEL maintainer="mritd <mritd1234@gmail.com>"

ARG TZ="Asia/Shanghai"
ARG VERSION="v1.0.0"

ENV TZ ${TZ}
ENV VERSION ${VERSION}
ENV DOWNLOAD_URL https://github.com/Gozap/certmonitor/releases/download/${VERSION}/certmonitor_linux_amd64

RUN apk upgrade \
    && apk add bash tzdata wget \
    && wget ${DOWNLOAD_URL} -O /usr/bin/certmonitor \
    && chmod +x /usr/bin/certmonitor \
    && ln -sf /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone \
    && apk del wget \
    && rm -rf /var/cache/apk/*

CMD ["certmonitor"]
