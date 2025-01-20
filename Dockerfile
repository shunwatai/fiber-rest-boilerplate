ARG GO_VERSION=1.23
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine
WORKDIR /app

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
    appuser && \
    chown ${UID}:${UID} -R /go && \
    chown ${UID}:${UID} /app
    # chown ${UID}:${UID} -R /app /go

COPY --chown=appuser:appuser . /app/
RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add alpine-sdk musl-dev sqlite-dev \
        ca-certificates \
        tzdata \
        sqlite-dev \
        ffmpeg \
        bash && \
        update-ca-certificates

USER appuser 

RUN go install github.com/air-verse/air@latest && \
    go install github.com/swaggo/swag/cmd/swag@latest && \
    go install -tags 'postgres mysql sqlite3 mongodb libsqlite3 linux musl' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN swag init

# This is the architecture you’re building for, which is passed in by the builder.
# Placing it here allows the previous steps to be cached across architectures.
ARG TARGETARCH

RUN go mod tidy

# Expose the port that the application listens on.
EXPOSE 7000

# What the container should run when it is started.
CMD [ "air", "-c", ".air-alpine.toml" ]
