FROM ubuntu:18.04

EXPOSE 443

LABEL \
  maintainer="https://github.com/olblak"

RUN \
  addgroup --gid 101 openvpn && \
  useradd -d /var/lib/ldap/ -g openvpn -m -u 101 openvpn

COPY config/server.conf /etc/openvpn/server/server.conf
COPY config/auth-ldap.conf /etc/openvpn/server/auth-ldap.conf
COPY entrypoint.sh /entrypoint.sh

RUN \
  apt-get update &&\
  apt-get install -y openvpn openvpn-auth-ldap dnsmasq &&\
  apt-get clean &&\
  rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN chmod 0750 /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh" ] 
