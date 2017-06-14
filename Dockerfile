FROM       alpine:latest
MAINTAINER Tim Vaillancourt <tim.vaillancourt@percona.com>
EXPOSE     9104

ENV GOPATH /go
ENV APPPATH $GOPATH/src/github.com/percona/mongodb_exporter
COPY . $APPPATH
RUN apk add --update -t build-deps go git mercurial libc-dev gcc libgcc curl \
    && cd $APPPATH && go get github.com/Masterminds/glide && /go/bin/glide up -v && go get -d && go build -o /bin/mongodb_exporter \
    && apk del --purge build-deps && rm -rf $GOPATH

ENTRYPOINT [ "/bin/mongodb_exporter" ]
