FROM --platform=$BUILDPLATFORM  golang:buster

ARG TARGETARCH

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY *.go ./

RUN mkdir -p migrations
COPY migrations/* ./migrations/

RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg

RUN GOOS=linux GOARCH=$TARGETARCH go build -o /go-binary  /app/cmd/main.go


CMD [ "/go-binary" ]

