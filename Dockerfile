FROM ubuntu:18.04

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
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* \
  && mkdir /etc/ldap/ssl

COPY --chown=openvpn cert/pki/ca.crt /etc/openvpn/server/ca.crt
COPY --chown=openvpn cert/pki/crl.pem /etc/openvpn/server/crl.pem
COPY --chown=openvpn cert/ccd /etc/openvpn/server/ccd

COPY docker/config/server.conf /etc/openvpn/server/server.conf
COPY docker/config/auth-ldap.conf /etc/openvpn/server/auth-ldap.conf

COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod 0750 /entrypoint.sh

LABEL io.jenkins-infra.tools="curl,dnsmasq,openvpn"

EXPOSE 443

ENTRYPOINT ["bash", "/entrypoint.sh" ]
