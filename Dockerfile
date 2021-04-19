FROM golang:1.15.7-buster
RUN go get -u github.com/ruairigibney/wolfpack-file-server
ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor
ENV APP_HOME /go/src/wolfpack-file-server
WORKDIR $APP_HOME
EXPOSE $FS_ARCHIVE_PORT
EXPOSE $FS_PASSCODE_PORT
CMD ["./wolfpack-file-server"]
