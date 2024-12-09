FROM --platform=$BUILDPLATFORM golang:latest AS builder

WORKDIR /app

COPY . .

RUN GOOS=$TARGETOS GOARCH=$GOARCH make

FROM scratch

WORKDIR /app

COPY --from=builder /app/bin/test-psql-conn /app/exec

ENV IS_DOCKER=true

EXPOSE 8080

ENTRYPOINT [ "/app/exec" ]
