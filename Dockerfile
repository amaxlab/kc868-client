FROM golang
WORKDIR /go/src/github.com/amaxlab/kc868-client
ADD . /go/src/github.com/amaxlab/kc868-client

RUN go get -u github.com/amaxlab/go-lib/log github.com/amaxlab/go-lib/config github.com/go-chi/chi
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -installsuffix cgo -o /go/bin/kc868-client

FROM scratch
COPY --from=0 /go/bin/kc868-client /go/bin/kc868-client

ENV APP_DEBUG="false"
ENV APP_PORT="8080"
ENV APP_KC868_HOST="192.168.0.1"
ENV APP_KC868_PORT="4196"

EXPOSE 8080

ENTRYPOINT ["/go/bin/kc868-client"]
