FROM golang:1.15.7-buster
ENV GO111MODULE=on
ENV APP_USER app
ENV APP_HOME /go/src/wolfpack-file-server
RUN groupadd $APP_USER && useradd -m -g $APP_USER -l $APP_USER
RUN mkdir -p $APP_HOME && chown -R $APP_USER:$APP_USER $APP_HOME
WORKDIR $APP_HOME
USER $APP_USER
COPY src/ .
RUN go mod download
RUN go mod verify
COPY src/cmd/wolfpack-file-server/ .
RUN go build -o wolfpack-file-server
EXPOSE $FS_ARCHIVE_PORT
EXPOSE $FS_PASSCODE_PORT
CMD ["./wolfpack-file-server"]
