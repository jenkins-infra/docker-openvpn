# README

This project contents everything related to Jenkins infrastructure VPN.
It includes the following elements:

* Build an Openvpn docker image integrated with [openldap](https://github.com/jenkins-infra/ldap).
* Manage client configuration and certificate

## CONNECTION
To connect to this VPN, your VPN client must be configured with your Jenkins account username/password and certificate authentication.
Certificate authentication requires the following files:

* The **[ca.crt](https://github.com/jenkins-infra/openvpn/blob/master/cert/pki/ca.crt)**
* **username.key** cfr [HowTo Get client access](#howto-get-client-access) !! Your private key **must** remain **secret**,
* **username.crt** is located in [cert/pki/issued](https://github.com/jenkins-infra/openvpn/tree/master/cert/pki/issued), once an administrator signs  your request and publish it.

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

#### DNS Problems

If you are having issues connecting to resources behind the VPN, but the VPN appears to be working correctly, check your DNS settings.  Some providers seem to filter out requests to the zone.  To test, try `dig release.ci.jenkins.io`.

To enable a different DNS provider only when connected to the VPN you can add the following to you OpenVPN config file

```dhcp-option DNS 8.8.8.8```

### Windows only
If you want to use multiple VPN connections at the same time with OpenVPN, you have to install a new TAP adapter. This can be very easily by running *as Admin* the `C:\Program Files\TAP-Windows\bin\addtap.bat`. The TAP-Windows tool is installed in parallel with OpenVPN.

## CERTIFICATES
This project holds VPN keys for connecting on Jenkins infrastructure VPN.

If you think that you should have access to this network, feel free to read [HowTo Get client access](#howto-get-client-access).

### Client
#### HowTo get client access
To access the Jenkins infrastructure private network, you need a certificate containing your Jenkins username as CN.
Then this certificate must be signed by an administrator who also assigns you a static IP configuration.

Feel free to follow the next action points:

* Fork this repository on your own Github account: [fork a repo](https://help.github.com/articles/fork-a-repo/)
* Build EASYVPN binary by running one of the following commands depending on your
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

* Enter in the VPN network directory: `cd cert`
* Run `make show-req name=<username>`

#### HowTo show certificate information

* Enter in the VPN network directory: `cd cert`
* Run `make show-cert name=<username>`

#### Howto validate your vpn access

You can test if your private key matches your certificate and certificate request by running following commands:

```
openssl pkey -in <your_private_key> -pubout -outform pem | sha256sum
==
openssl x509 -in <your_certificate> -pubkey -noout -outform pem | sha256sum
==
openssl req -in <your_certificate_request> -pubkey -noout -outform pem | sha256sum
```

### Administrator
#### HowTo become an administrator
To add/revoke certificates, you must be allowed to decrypt `cert/pki/private/ca.key.enc`.
This file is encrypted with [sops](https://github.com/mozilla/sops) and you are public gpg key must be added to .sops.yaml by an existing administrator to be allowed to run `make decrypt`.

This repository relies on [easy-rsa](https://github.com/OpenVPN/easy-rsa/blob/master/README.quickstart.md).

#### HowTo approve client access?
To validate and sign a client certificate, you are going to execute the following actions

* Build EASYVPN binary by running one of the following commands depending on your
 * `make init_osx`
 * `make init_linux`
 * `make init_windows` then copy utils/easyvpn/easyvpn.exe at the root of this repository
* Git checkout on the right branch "staging"
* Sign certificate request: `./easyvpn sign <CN_to_sign>`
* Merge staging into master
* Update Docker image tag in the [puppet](https://github.com/jenkins-infra/jenkins-infra/blob/staging/dist/profile/manifests/openvpn.pp) configuration.

#### HowTo revoke client access?

* Build EASYVPN binary by running one of the following commands depending on your
 * `make init_osx`
 * `make init_linux`
 * `make init_windows` and copy utils/easyvpn/easyvpn.exe at the root of this repository
* Revoke certificate: `./easyvpn revoke <CN_to_sign>`
* Update Docker image tag in the [puppet](https://github.com/jenkins-infra/jenkins-infra/blob/staging/dist/profile/manifests/openvpn.pp) configuration.

#### HowTo review certificate revocation list

If the certificate revocation list expired, then the openvpn logs will contains errors like 'VERIFY ERROR: depth=0, error=CRL has expired:...'
We can run `openssl crl -in cert/pki/crl.pem -noout -text` to validate that the crl expired and that we need to generate a new one.

To generate a new one:
* Decrypt ca.key `sops -d cert/pki/private/ca.key.enc > cert/pki/private/ca.key`
* Generate a new crl.pem - `cd cert ; ./easyrsa gen-crl`
* Publish the new crl.pem - `git add cert/pki/crl.pem && git commit cert/pki/crl.pem -s -m 'Renew revocation list certificate'`
* Delete local ca.key - `rm cert/pki/private/ca.key`

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
