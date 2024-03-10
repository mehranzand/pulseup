# Build client
FROM node:21-alpine as frontend

WORKDIR /pulseup/app

COPY package.json ./
COPY package-lock.json ./

# Copy package.json and install dependencies
RUN npm install react-scripts@5.0.1 -g --silent

# Copy all files
COPY . ./
RUN npm run build

FROM golang:1.22.1-alpine AS backend

WORKDIR /pulseup/server

# Copy go mod files
COPY go.* ./
RUN go mod download

COPY main.go ./

EXPOSE 8080
ENTRYPOINT ["/pulseup"]