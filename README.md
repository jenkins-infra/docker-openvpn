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

**With NetworkManager client, you must enable the option :**

`Use this connection only for resources on its network`

### DNS Problems

If you are having issues connecting to resources behind the VPN, but the VPN appears to be working correctly, check your DNS settings.  Some providers seem to filter out requests to the zone.  To test, try `dig release.ci.jenkins.io`, you should get something like this:

<details><summary>dig output (click to expand)</summary>

```text
; <<>> DiG 9.10.6 <<>> release.ci.jenkins.io
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 13457
;; flags: qr rd ra; QUERY: 1, ANSWER: 2, AUTHORITY: 0, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 1220
;; QUESTION SECTION:
;release.ci.jenkins.io.         IN      A

;; ANSWER SECTION:
release.ci.jenkins.io.  3600    IN      CNAME   private.aks.jenkins.io.
private.aks.jenkins.io. 3600    IN      A       10.0.2.5

;; Query time: 80 msec
;; SERVER: 192.168.1.254#53(192.168.1.254)
;; WHEN: Tue Oct 12 20:49:59 CEST 2021
;; MSG SIZE  rcvd: 92
```

</details>

To enable a different DNS provider only when connected to the VPN you can add the following to you OpenVPN config file

```dhcp-option DNS 8.8.8.8```

### Windows only

If you want to use multiple VPN connections at the same time with OpenVPN, you have to install a new TAP adapter. This can be very easily by running *as Admin* the `C:\Program Files\TAP-Windows\bin\addtap.bat`. The TAP-Windows tool is installed in parallel with OpenVPN.

## CERTIFICATES

This project holds VPN keys for connecting on Jenkins infrastructure VPN.

If you think that you should have access to this network, feel free to read [HowTo Get client access](#howto-get-client-access).

### Client

#### HowTo get client access

To access the Jenkins infrastructure private network, you need a certificate containing your [Jenkins username](https://accounts.jenkins.io/) as CN ([commonName](https://docs.oracle.com/javase/7/docs/technotes/tools/solaris/keytool.html#DName)).
Then this certificate must be signed by an administrator who also assigns you a static IP configuration.

Feel free to follow the next action points:

* [Fork](https://help.github.com/articles/fork-a-repo/) this repository on your own Github account: [fork the repo](https://github.com/jenkins-infra/openvpn/fork)
* Clone your fork locally: `git clone https://github.com/<your-github-username>/openvpn && cd openvpn`
* Build EASYVPN binary by running one of the following commands depending on your operating system:
  * `make init_osx`
  * `make init_linux`
  * `make init_windows` then copy utils/easyvpn/easyvpn.exe at the root of this repository
* Generate your private key and certificate request: `./easyvpn request <your-jenkins-username>`
  Your private key will be generated in `./cert/pki/private`, this key **must** remain **secret**.
* Create a new Pull Request on [jenkinsinfra/openvpn](https://github.com/jenkins-infra/openvpn), `staging` branch: [How to Create a pull request](https://help.github.com/articles/creating-a-pull-request/)
* Open an INFRA ticket on [JIRA](https://issues.jenkins-ci.org) referencing your PR
* Grab a cup of coffee and wait patiently for an administrator to sign your certificate request
* Once an admin notify you that everything is right, you can [sync your fork](https://docs.github.com/en/github/collaborating-with-pull-requests/working-with-forks/syncing-a-fork) then pull it to retrieve your certificate from `./cert/pki/issued/<your-jenkins-username>.crt`
* We recommend you to move the `./cert` folder to an hidden folder in your home (`~/.cert`)
* You can finally create the config file used by your VPN client (example here for Tunnelblick, an OSX VPN client, opening this file from the Finder should launch it)  

_jenkins-infra.ovpn_
```text
client
remote vpn.jenkins.io 443
ca "~/.cert/pki/ca.crt"
cert "~/.cert/pki/issued/<your-jenkins-username>.crt"
key "~/.cert/pki/private/<your-jenkins-username>.key"
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

#### HowTo show request information

* Enter in the VPN network directory: `cd cert`
* Run `make show-req name=<your-jenkins-username>`

#### HowTo show certificate information

* Install [sops](https://github.com/mozilla/sops)
* Enter in the VPN network directory: `cd ~/.cert`
* Run `make decrypt`
* Run `make show-cert name=<your-jenkins-username>`

#### Howto validate your certificate

You can test if your private key matches your certificate and certificate request by running following commands:

```bash
openssl pkey -in ~/.cert/pki/private/<your-jenkins-username>.key -pubout -outform pem | sha256sum
# Should be equal to
openssl x509 -in ~/.cert/pki/issued/<your-jenkins-username>.crt -pubkey -noout -outform pem | sha256sum
# And also equal to
openssl req -in ~/.cert/pki/reqs/<your-jenkins-username>.req -pubkey -noout -outform pem | sha256sum
```

### Administrator

#### HowTo become an administrator

To add/revoke certificates, you must be allowed to decrypt `cert/pki/private/ca.key.enc`.
This file is encrypted with [sops](https://github.com/mozilla/sops) and your public gpg key must be added to .sops.yaml by an existing administrator.

This repository relies on [easy-rsa](https://github.com/OpenVPN/easy-rsa/blob/master/README.quickstart.md).

#### HowTo approve client access?

To validate and sign a client certificate, you are going to execute the following actions

* Build EASYVPN binary by running one of the following commands depending on your
  * `make init_osx`
  * `make init_linux`
  * `make init_windows` then copy utils/easyvpn/easyvpn.exe at the root of this repository
* Merge the Pull Request of the requester to staging.
* Git checkout on the right branch "staging" to retrieve the CRL from the requester
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

To generate a new CRL:

* Decrypt ca.key `sops -d cert/pki/private/ca.key.enc > cert/pki/private/ca.key`
* Generate a new crl.pem - `cd cert ; ./easyrsa gen-crl; cd ..`
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
