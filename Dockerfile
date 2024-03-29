FROM golang:1.16-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
COPY prod.env test.env ./
COPY keys ./
RUN go mod download

# Copy local code to the container image.
COPY . .

# Definition of this variable is used by 'skaffold debug' to identify a golang binary.
# Default behavior - a failure prints a stack trace for the current goroutine.
# See https://golang.org/pkg/runtime/
ENV GOTRACEBACK=single

# Skaffold passes in debug-oriented compiler flags
ARG SKAFFOLD_GO_GCFLAGS
RUN echo "Go gcflags: ${SKAFFOLD_GO_GCFLAGS}"
# Build go for linux
RUN CGO_ENABLED=0 GOOS=linux go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -ldflags="-w -s" -o coiner -v

# New image to deploy
FROM alpine:latest

WORKDIR /app

ENV GOTRACEBACK=single
# install certificates
# RUN apk update && apk add ca-certificates && update-ca-certificates && rm -rf /var/cache/apk/*
# RUN echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] http://packages.cloud.google.com/apt cloud-sdk main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list && curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key --keyring /usr/share/keyrings/cloud.google.gpg  add - && apt-get update -y && apt-get install google-cloud-sdk -y

# Copy the binary to the production image from the builder stage.d
COPY --from=builder /app/ .

ENTRYPOINT ["./coiner"]
