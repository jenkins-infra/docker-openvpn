FROM golang:1.26.1 AS easyvpn

WORKDIR /easyvpn

COPY utils/easyvpn .
RUN make build_linux
COPY config.yaml .
RUN mkdir -p ./cert/ccd/private
COPY ./cert/ ./cert
RUN ./easyvpn --commit=false --push=false clientconfig --all

FROM ubuntu:22.04 AS vpn-server

# We want to use the latest available packages
# hadolint ignore=DL3008
RUN addgroup --gid 101 openvpn \
  && useradd -d /var/lib/ldap/ -g openvpn -m -u 101 openvpn \
  && apt-get update -q \
  && apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    dnsmasq \
    openvpn \
    openvpn-auth-ldap \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY --chown=openvpn cert/pki/ca.crt /etc/openvpn/server/ca.crt
COPY --chown=openvpn cert/pki/crl.pem /etc/openvpn/server/crl.pem
COPY --chown=openvpn --from=easyvpn /easyvpn/cert/ccd /home/openvpn/available-ccds

COPY docker/config/server.conf /etc/openvpn/server/server.conf
COPY docker/config/auth-ldap.conf /etc/openvpn/server/auth-ldap.conf

COPY docker/entrypoint.sh /entrypoint.sh

LABEL io.jenkins-infra.tools="curl,dnsmasq,openvpn"

EXPOSE 443

ENTRYPOINT ["bash", "/entrypoint.sh" ]
