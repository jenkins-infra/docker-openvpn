
ARG VERSION=1.12
FROM golang:$VERSION

ARG UID=1000
ARG GID=1000

ENV GOPATH=/home/gopher

RUN \
  groupadd -g $GID gopher  && \
  useradd -g $GID -u $UID -m gopher && \
  apt-get update && \
    apt-get -y install --no-install-recommends make && \
    apt-get clean && \
    find /var/lib/apt/lists -type f -delete

USER gopher
