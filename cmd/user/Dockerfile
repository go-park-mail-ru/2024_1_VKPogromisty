# Start from golang base image
FROM golang:1.22.1-alpine3.19 AS build-stage

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git && apk add --no-cache bash && apk add build-base

# Copy the source from the current directory to the Working Directory inside the container
COPY . .
COPY .env /

# RUN go install github.com/golang/mock/mockgen@v1.6.0

# RUN make mocks

# RUN make test

# Build the Go app
RUN CGO_ENABLED=0 go build -o /build/user ./cmd/user/main.go

FROM gcr.io/distroless/base-debian11 AS build-release-stage

COPY --from=build-stage /build/user /build/user
COPY --from=build-stage .env .

EXPOSE 8082

# Run the executable
CMD ["/build/user"]