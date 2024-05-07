FROM  node:21.7.3-alpine as node

WORKDIR /build

# Copy package.json and install dependencies
COPY web ./
RUN yarn install

# Copy assets to build
COPY web ./web
COPY web/public ./web/public

# Build web
RUN yarn build 

FROM golang:1.22.2-alpine AS builder

RUN mkdir /pulseup

WORKDIR /pulseup

# Copy go mod files
COPY go.* ./
RUN go mod download

# Copy assets built with node
COPY --from=node /build/dist ./web/dist

# Copy all other files
COPY internal ./internal
COPY main.go ./
COPY .env ./

# Args
ARG TAG=head
ARG TARGETOS TARGETARCH

# Build binary
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w -X main.version=$TAG" -o pulseup

FROM scratch

COPY --from=builder /pulseup ./

EXPOSE 7070

ENTRYPOINT ["/pulseup"]