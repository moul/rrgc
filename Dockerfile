# dynamic config
ARG             BUILD_DATE
ARG             VCS_REF
ARG             VERSION

# build
FROM            golang:1.23.2-alpine as builder
RUN             apk add --no-cache git gcc musl-dev make
ENV             GO111MODULE=on
WORKDIR         /go/src/moul.io/rrgc
COPY            go.* ./
RUN             go mod download
COPY            . ./
RUN             make install

# minimalist runtime
FROM alpine:3.20.3
LABEL           org.label-schema.build-date=$BUILD_DATE \
                org.label-schema.name="rrgc" \
                org.label-schema.description="" \
                org.label-schema.url="https://moul.io/rrgc/" \
                org.label-schema.vcs-ref=$VCS_REF \
                org.label-schema.vcs-url="https://github.com/moul/rrgc" \
                org.label-schema.vendor="Manfred Touron" \
                org.label-schema.version=$VERSION \
                org.label-schema.schema-version="1.0" \
                org.label-schema.cmd="docker run -i -t --rm moul/rrgc" \
                org.label-schema.help="docker exec -it $CONTAINER rrgc --help"
COPY            --from=builder /go/bin/rrgc /bin/
ENTRYPOINT      ["/bin/rrgc"]
#CMD             []
