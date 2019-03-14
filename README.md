# README

This project contents everything related to Jenkins infrastructure vpn. 
It includes following elements:
* Build an openvpn docker image integrated with [openldap](https://github.com/jenkins-infra/ldap).
* Manage client configuration and certificate

## CERTIFICATES
This projects holds vpn keys for connecting on Jenkins Infrastructure.

If you think that you should have access to it or a specific network, feel free to read [HowTo Get client access](#howto-get-client-access).

## Client
### HowTo get client access
In order to Jenkins infrastructure private networks, you need a certificate containing your jenkins username as CN.
Then this certificate must be signed by an administrator who also assign you a static IP.

Feel free to follow next action points:

* Fork this repository on your own Github account: [fork a repo](https://help.github.com/articles/fork-a-repo/)
* Build easyvpn cli: `make init`
* Create your private key and certificate request: `./easyvpn request <your username>`
* Create a new Pull Request on jenkinsinfra/openvpn, staging branch: [create a pull request](https://help.github.com/articles/creating-a-pull-request/)
* Grab a cup of coffee and wait patiently until an administrator issues your certificate.
* Once an admin notify you that everything is right, your can then retrieve your certificate from `./cert/pki/issued/<your_username>.crt`

### HowTo show request information

* Enter in the vpn network directory: `cd cert`
* Run `make show-req name=<username>`

### HowTo show certificate information

* Enter in the vpn network directory: `cd cert`
* Run `make show-certs name=<username>`

## Administrator
### HowTo become an administrator
In order to add/revoke certificates, you must be allowed to decrypt `cert/pki/private/ca.key.enc`.
This file is encrypted with [sops](https://github.com/mozilla/sops) and you are public gpg key must be added to .sops.yaml by an existing administrator in order to be allow to run `make decrypt`.

This repository relies on [easy-rsa](https://github.com/OpenVPN/easy-rsa/blob/master/README.quickstart.md).

### HowTo approve client access?
In order to validate and sign a client certificate, your are going to do following actions

* Build easyvpn cli: `make init`
* Git checkout on the right branch "staging"
* Sign certificate request: `./easyvpn sign <CN_to_sign>`
* Update docker image in the [puppet](https://github.com/jenkins-infra/jenkins-infra/blob/staging/dist/profile/manifests/openvpn.pp) configuration.

### HowTo revoke client access?

* Build easyvpn cli: `make init`
* Revoke certificate: `./easyvpn revoke <CN_to_sign>`
* Update docker image in the [puppet](https://github.com/jenkins-infra/jenkins-infra/blob/staging/dist/profile/manifests/openvpn.pp) configuration.

## DOCKER
### CONFIGURATION
This image can be configured at runtime with different environment variables.

* `AUTH_LDAP_BINDDN` **Define user dn used to query the ldap database**
* `AUTH_LDAP_URL` **Define ldap endpoint url**
* `AUTH_LDAP_PASSWORD` **Define user dn password**
* `AUTH_LDAP_GROUPS_MEMBER` **Define required group member to authenticate**

Some examples can be found inside [docker-compose.yaml](docker/docker-compose.yaml)

### TESTING
In order to test this image, you need a "mock" ldap and SSL certificates.
Then go in directory `docker` and run one of the following commands

! Certificates must be readable by UID 101
`make start` - **Start the ldap and vpn service**

## INFRASTRUCTURE

This project is designed to work with following pieces:

* Machine provisioned by [Terraform](https://github.com/jenkins-infra/azure)
* Service configured and orchestrated by [Puppet](https://github.com/jenkins-infra/jenkins-infra/blob/staging/dist/profile/manifests/openvpn.pp)

## CONTRIBUTING
Feel free to contribute to this image by:

1. Fork this project into your account
2. Make your changes in your local fork
3. Submit a pull request with a description and a link to a Jira ticket 
4. Ask for a review

## ISSUE
Please report any issue on the jenkins infrastructure [project](https://issues.jenkins-ci.org/secure/Dashboard.jspa)

## LINKS
* [How to contribute to OSS?](https://opensource.guide/how-to-contribute/)
* [jenkins-infra/azure](https://github.com/jenkins-infra/azure)
* [jenkins-infra/jenkins-infra](https://github.com/jenkins-infra/jenkins-infra/blob/staging/dist/profile/manifests/openvpn.pp)
* [jenkins-infra/openvpn](https://github.com/jenkins-infra/openvpn)
* [mozilla/sops](https://github.com/mozilla/sops)
* [openvpn/easy-rsa](https://github.com/OpenVPN/easy-rsa)
