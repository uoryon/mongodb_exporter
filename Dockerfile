FROM       alpine:latest
MAINTAINER Tim Vaillancourt <tim.vaillancourt@percona.com>
EXPOSE     9104

ENV GOPATH /go
ENV GOBIN /go/bin
ENV LGOBIN /go/bin
RUN export PATH=$PATH:$GOBIN
ENV APPPATH $GOPATH/src/github.com/percona/mongodb_exporter
RUN apk add --update -t build-deps go git mercurial libc-dev gcc libgcc curl
COPY . $APPPATH
RUN cd $APPPATH && go get github.com/Masterminds/glide && /go/bin/glide up -v && go get -d && go build -o /bin/mongodb_exporter \
    && apk del --purge build-deps && rm -rf $GOPATH

ENTRYPOINT [ "/bin/mongodb_exporter" ]
