# Start from golang base image
FROM golang:1.22.1-alpine3.19 AS build-stage

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git && apk add --no-cache bash && apk add build-base

# Setup folders
RUN mkdir /app
WORKDIR /app

# Copy the source from the current directory to the working Directory inside the container
COPY . .
COPY .env .

# Download all the dependencies
RUN go get -d -v ./...

RUN go install github.com/golang/mock/mockgen@v1.6.0

RUN make mocks

RUN make test

WORKDIR /app/app
# Build the Go app
RUN CGO_ENABLED=0 go build -o /build

FROM gcr.io/distroless/base-debian11 AS build-release-stage

COPY --from=build-stage /build /build
COPY --from=build-stage /app/.env .env
COPY --from=build-stage /app/docs docs

EXPOSE 8080 8001

# Run the executable
CMD ["/build"]
