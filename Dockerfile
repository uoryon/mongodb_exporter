FROM       alpine:latest
MAINTAINER Tim Vaillancourt <tim.vaillancourt@percona.com>
EXPOSE     9104

ENV GOPATH /go
ENV GOBIN /go/bin
ENV LGOBIN /go/bin
ENV APPPATH $GOPATH/src/github.com/percona/mongodb_exporter
RUN apk add --update -t build-deps go git mercurial libc-dev gcc libgcc curl
COPY . $APPPATH
RUN cd $APPPATH && mkdir -p /go/bin && curl https://glide.sh/get | sh && /go/bin/glide up -v && go get -d && go build -o /bin/mongodb_exporter \
    && apk del --purge build-deps && rm -rf $GOPATH

ENTRYPOINT [ "/bin/mongodb_exporter" ]
