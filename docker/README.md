# README

This project contents everything needed to build an openvpn docker image integrated with [openldap](https://github.com/jenkins-infra/ldap).
This image is designed to be used inside the Jenkins Infrastructure Project.

## CONFIGURATION
This image can be configured at runtime with different environment variables.

* `AUTH_LDAP_BINDDN` **Define user dn used to query the ldap database**
* `AUTH_LDAP_URL` **Define ldap endpoint url**
* `AUTH_LDAP_PASSWORD` **Define user dn password**
* `AUTH_LDAP_GROUPS_MEMBER` **Define required group member to authenticate**

Some examples can be found inside [docker-compose.yaml](docker-compose.yaml)

## TESTING
In order to test this image, you need an mock-ldap and ssl certificates.

`make build` - **Build the vpn docker image** 

`make gen_cert` - **Generate SSL Certificates used by the vpn and ldap**

! Certificates must be readable by uid 101

`make up` - **Start the ldap and vpn service**

## CONTRIBUTING
Feel free to contribute to this image by:

1. Fork this project into your account
2. Make your changes in your local fork
3. Submit a pull request with a description and a link to a jira ticket 
4. Ask for a review

## KEY MANAGEMENT
In order to authentication on any of the jenkins infra vpn, you are going to need a client certificate and a private key.
While the certificate can be public, the private key must remain private.

In order to simplify key management, all clients certificates are stored in this repository and respect following files structure

./ssl
./ssl/private: never pushed on the repository and contain your private key
./ssl/ca.pem: Represent ca public certificate
./ssl/keys/: Contains every generated client certificates, remarks those certificate must be encrypted with ca cert
  Keys stored in this directory, must respect following rules
  emailAddress= must represent your jenkins account email address
  cn= must represent your the jenkins account username
./ssl/revoke_keys/ Contains revoked keys
./ssl/scripts contains scripts used to populate ssl directory


## ISSUE
Please report any issue on the jenkins infrastructure [project](https://issues.jenkins-ci.org/secure/Dashboard.jspa)

## LINKS
* [openvpn-auth-ldap](https://github.com/threerings/openvpn-auth-ldap)
* [How to contribute to OSS?](https://opensource.guide/how-to-contribute/)
