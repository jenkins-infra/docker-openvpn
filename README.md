# README

This project contents everything related to Jenkins infrastructure vpn.
It includes the following elements:

* Build an Openvpn docker image integrated with [openldap](https://github.com/jenkins-infra/ldap).
* Manage client configuration and certificate

## CONNECTION
To connect to this vpn, your vpn client must be configured with your Jenkins account username/password and certificate authentication.
Certificate authentication requires the following files:

* The **[ca.crt](https://github.com/jenkins-infra/openvpn/blob/master/cert/pki/ca.crt)**
* **username.key** cfr [HowTo Get client access](#howto-get-client-access) !! Your private key **must** remain **secret**,
* **username.crt** is located in [cert/pki/issued](https://github.com/jenkins-infra/openvpn/tree/master/cert/pki/issued), once an administrator sign  your request and publish it.

```
client
remote vpn.jenkins.io 443
ca "~/.cert/jenkins/ca.crt"
cert "~/.cert/jenkins/username.crt"
key "~/.cert/jenkins/username.key"
auth-user-pass
dev tun
proto tcp
nobind
auth-nocache
script-security 2
persist-key
persist-tun
user nobody
group nobody
```

**With network manager client, you must enable the option :**

`Use this connection only for resources on its network`

### Windows only
If you want to use multiple VPN connections at the same time with OpenVPN, you have to install a new TAP adapter. This can be very easily by running *as Admin* the `C:\Program Files\TAP-Windows\bin\addtap.bat`. The TAP-Windows tool is installed in parallel of OpenVPN.

## CERTIFICATES
This project holds vpn keys for connecting on Jenkins infrastructure vpn.

If you think that you should have access to this network, feel free to read [HowTo Get client access](#howto-get-client-access).

### Client
#### HowTo get client access
To access the Jenkins infrastructure private network, you need a certificate containing your Jenkins username as CN.
Then this certificate must be signed by an administrator who also assign you a static IP configuration.

Feel free to follow the next action points:

* Fork this repository on your own Github account: [fork a repo](https://help.github.com/articles/fork-a-repo/)
* Build easyvpn binary by running one of the following commands depending on your
  * `make init_osx`
  * `make init_linux`
  * `make init_windows` then copy utils/easyvpn/easyvpn.exe at the root of this repository
* Generate your private key and certificate request: `./easyvpn request <your username>`
  Your private key will be generated in `cert/pki/private`, this key **must** remain **secret**.
* Create a new Pull Request on jenkinsinfra/openvpn, staging branch: [How to Create a pull request](https://help.github.com/articles/creating-a-pull-request/)
* Open an INFRA ticket on [JIRA](https://issues.jenkins-ci.org) referencing your PR
* Grab a cup of coffee and wait patiently for an administrator to sign your certificate request.
* Once an admin notify you that everything is right, you can then retrieve your certificate from `./cert/pki/issued/<your_username>.crt`

#### HowTo show request information

* Enter in the vpn network directory: `cd cert`
* Run `make show-req name=<username>`

#### HowTo show certificate information

* Enter in the vpn network directory: `cd cert`
* Run `make show-certs name=<username>`

### Administrator
#### HowTo become an administrator
To add/revoke certificates, you must be allowed to decrypt `cert/pki/private/ca.key.enc`.
This file is encrypted with [sops](https://github.com/mozilla/sops) and you are public gpg key must be added to .sops.yaml by an existing administrator in order to be allow to run `make decrypt`.

This repository relies on [easy-rsa](https://github.com/OpenVPN/easy-rsa/blob/master/README.quickstart.md).

#### HowTo approve client access?
To validate and sign a client certificate, you are going to execute the following actions

* Build easyvpn binary by running one of the following commands depending on your
 * `make init_osx`
 * `make init_linux`
 * `make init_windows` then copy utils/easyvpn/easyvpn.exe at the root of this repository
* Git checkout on the right branch "staging"
* Sign certificate request: `./easyvpn sign <CN_to_sign>`
* Update Docker image tag in the [puppet](https://github.com/jenkins-infra/jenkins-infra/blob/staging/dist/profile/manifests/openvpn.pp) configuration.

#### HowTo revoke client access?

* Build easyvpn binary by running one of the following commands depending on your
 * `make init_osx`
 * `make init_linux`
 * `make init_windows` and copy utils/easyvpn/easyvpn.exe at the root of this repository
* Revoke certificate: `./easyvpn revoke <CN_to_sign>`
* Update Docker image tag in the [puppet](https://github.com/jenkins-infra/jenkins-infra/blob/staging/dist/profile/manifests/openvpn.pp) configuration.

## DOCKER
### CONFIGURATION
This image can be configured at runtime with different environment variables.

* `AUTH_LDAP_BINDDN` **Define user dn used to query the ldap database**
* `AUTH_LDAP_URL` **Define ldap endpoint url**
* `AUTH_LDAP_PASSWORD` **Define user dn password**
* `AUTH_LDAP_GROUPS_MEMBER` **Define required group member to authenticate**

Some examples can be found inside [docker-compose.yaml](docker/docker-compose.yaml)

### TESTING
To test this image, you need a "mock" ldap and SSL certificates.
Then go in directory `docker` and run one of the following commands

! Certificates must be readable by UID 101
`make start` - **Start the ldap and vpn service**

## INFRASTRUCTURE

This project is designed to work with the following requirements:

* Machine provisioned by [Terraform](https://github.com/jenkins-infra/azure)
* Service configured and orchestrated by [Puppet](https://github.com/jenkins-infra/jenkins-infra/blob/staging/dist/profile/manifests/openvpn.pp)

## CONTRIBUTING
Feel free to contribute to this image by:

1. Fork this project into your account
2. Make your changes in your local fork
3. Submit a pull request with a description and a link to a Jira ticket
4. Ask for a review

## ISSUE
Please report any issue on the Jenkins infrastructure [project](https://issues.jenkins-ci.org/secure/Dashboard.jspa)

## LINKS
* [How to contribute to OSS?](https://opensource.guide/how-to-contribute/)
* [jenkins-infra/azure](https://github.com/jenkins-infra/azure)
* [jenkins-infra/jenkins-infra](https://github.com/jenkins-infra/jenkins-infra/blob/staging/dist/profile/manifests/openvpn.pp)
* [jenkins-infra/openvpn](https://github.com/jenkins-infra/openvpn)
* [mozilla/sops](https://github.com/mozilla/sops)
* [openvpn/easy-rsa](https://github.com/OpenVPN/easy-rsa)
