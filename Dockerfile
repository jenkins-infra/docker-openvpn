FROM ubuntu:18.04

EXPOSE 443

LABEL \
  maintainer="https://github.com/olblak"\
  project="https://github.com/jenkins-infra/openvpn"

RUN \
  addgroup --gid 101 openvpn && \
  useradd -d /var/lib/ldap/ -g openvpn -m -u 101 openvpn

COPY cert/pki/ca.crt /etc/openvpn/server/ca.crt
COPY cert/pki/crl.pem /etc/openvpn/server/crl.pem

COPY cert/ccd /etc/openvpn/server/ccd
# Missing COPY Revoking list
COPY docker/config/server.conf /etc/openvpn/server/server.conf
COPY docker/config/auth-ldap.conf /etc/openvpn/server/auth-ldap.conf
COPY docker/entrypoint.sh /entrypoint.sh

RUN \
  apt-get update &&\
  apt-get install -y openvpn openvpn-auth-ldap dnsmasq &&\
  apt-get clean &&\
  rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* &&\
  mkdir /etc/ldap/ssl

RUN chmod 0750 /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh" ] 
