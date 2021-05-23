FROM golang:1.15
ARG GOARCH=amd64

RUN echo $GOARCH > /goarch

WORKDIR /go/src/github.com/doctori/opencycledatabase
COPY . .
RUN make GOARCH=$(cat /goarch.txt)

FROM scratch
copy dist/ocd /
expose 8080
VOLUME ["/tmp"]
ENTRYPOINT ["/ocd"]