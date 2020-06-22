ARG GOLANG_VERSION=1
ARG UBUNTU_VERSION=latest
FROM library/golang:$GOLANG_VERSION AS builder

RUN apt-get update && \
    apt-get install -y build-essential libsystemd-dev && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /src
COPY cmd ./cmd
COPY exporter ./exporter
COPY go.* ./

RUN go install ./cmd/...
RUN go test ./...


FROM library/ubuntu:$UBUNTU_VERSION

RUN apt-get update && \
    apt-get install -y curl && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /go/bin/opensshd_exporter /opt/opensshd_exporter

# See https://github.com/tommie/container-setuputils#addusergroup
ARG RUN_USER=root
ARG RUN_UID=
ARG RUN_GID=
ARG RUN_SUPP_GIDS=

COPY --from=githubtommie/container-setuputils /addusergroup /sbin/addusergroup
RUN [ "$RUN_USER" = root ] || addusergroup -u "$RUN_UID" -g "$RUN_GID" -G "$RUN_SUPP_GIDS" "$RUN_USER"

USER $RUN_USER
ENTRYPOINT ["/opt/opensshd_exporter"]
