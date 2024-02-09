ARG GO_VERSION=1.21
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine
WORKDIR /app

COPY . /app/
RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add alpine-sdk musl-dev sqlite-dev \
        ca-certificates \
        tzdata \
        sqlite-dev \
        && \
        update-ca-certificates

RUN go mod tidy && \
    go install github.com/cosmtrek/air@latest && \
    go install github.com/swaggo/swag/cmd/swag@latest && \
    go install -tags 'postgres mysql sqlite3 mongodb libsqlite3 linux musl' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN swag init

# This is the architecture you’re building for, which is passed in by the builder.
# Placing it here allows the previous steps to be cached across architectures.
ARG TARGETARCH

# Build the application.
# Leverage a cache mount to /go/pkg/mod/ to speed up subsequent builds.
# Leverage a bind mount to the current directory to avoid having to copy the
# source code into the container.
# CGO_ENABLED=1 -tags are required for sqlite3 & go-fitz
# RUN --mount=type=cache,target=/go/pkg/mod/ \
#     --mount=type=bind,target=. \
#     CGO_ENABLED=1 GOARCH=$TARGETARCH go build -tags "libsqlite3 linux musl" -o /bin/fiber-api .

# Create a non-privileged user that the app will run under.
# See https://docs.docker.com/go/dockerfile-user-best-practices/
# ARG UID=1000
ARG UID
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/app" \
    # --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser

# Expose the port that the application listens on.
EXPOSE 7000

# What the container should run when it is started.
CMD [ "air", "-c", ".air-alpine.toml" ]
