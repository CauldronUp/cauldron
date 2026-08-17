# Cauldron in a container, for dropping into an environment that already exists.
#
# The binary is static and the Recipes are embedded, so the runtime stage needs
# nothing but the binary and a set of CA certificates. It is not run as root:
# nothing here needs it, and a container in somebody else's compose file should
# not be the one that does.

FROM golang:1.26-alpine AS build

WORKDIR /src

# Dependencies first, so a change to the Recipes does not re-download them.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO off keeps the binary static, so the runtime stage can be scratch.
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/cauldron ./cmd/cauldron

FROM alpine:3.21

RUN adduser -D -u 10001 cauldron

COPY --from=build /out/cauldron /usr/local/bin/cauldron

USER cauldron

# 0.0.0.0 rather than the loopback default: a container that binds loopback is
# reachable from nothing but itself, which makes the container pointless. The
# host's own default stays loopback, where being hard to reach is the safer
# mistake.
EXPOSE 4600

ENTRYPOINT ["cauldron"]
CMD ["serve", "--headless", "--host", "0.0.0.0"]
