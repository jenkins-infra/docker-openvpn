# README

This project contents everything related to Jenkins infrastructure VPN.
It includes the following elements:

* Build an OpenVPN Docker image integrated with [openldap](https://github.com/jenkins-infra/ldap).
* Manage client configuration and certificate
* Hold VPN keys for connecting on Jenkins infrastructure VPN

If you think that you should have access to this network, feel free to read [HowTo Get client access](#howto-get-client-access).

## Connection

To connect to this VPN, your VPN client must be configured with your [Jenkins account](https://accounts.jenkins.io/) and certificate authentication, requiring the following files:

* The CertificateAuthority **[`ca.crt`](https://github.com/jenkins-infra/docker-openvpn/blob/main/cert/pki/ca.crt)**
* Your private key **`<your-jenkins-username>.key`**

  > ⚠️ your private key **must** remain **secret**! ⚠️

* Your certificate **`<your-jenkins-username>.crt`**

See [HowTo Get client access](#howto-get-client-access) below.

## Client

### How To get client access

To access the Jenkins infrastructure private network, you need a certificate containing your [Jenkins username](https://accounts.jenkins.io/) as CN ([commonName](https://docs.oracle.com/javase/7/docs/technotes/tools/solaris/keytool.html#DName)).
Then this certificate must be signed by an administrator who also assigns you a static IP configuration.

Feel free to follow the next action points:

* Open an issue on [jenkins-infra/helpdesk](https://github.com/jenkins-infra/helpdesk) describing the reason why you need an access to the VPN
  * If you need to access infra.ci.jenkins.io or release.ci.jenkins.io, mention it in your request to get access to the private VPN needed for these instances.
* [Fork](https://help.github.com/articles/fork-a-repo/) this repository on your own Github account: [fork the repo](https://github.com/jenkins-infra/docker-openvpn/fork)
* Clone your fork locally: `git clone https://github.com/<your-github-username>/docker-openvpn && cd docker-openvpn`
* Build EASYVPN binary by running one of the following commands depending on your operating system:
  * `make init_osx`
  * `make init_linux`
  * `make init_windows` then copy `./utils/easyvpn/easyvpn.exe` at the root of this repository

* Generate your private key and certificate request: `./easyvpn request <your-jenkins-username>`
  Your private key will be generated in `./cert/pki/private`

  > ⚠️ This key **must** remain **secret**! ⚠️

* Create a new pull request on [jenkins-infra/docker-openvpn](https://github.com/jenkins-infra/docker-openvpn)
  * From your local branch (usually the `main` branch)
  * Targeted to the remote `main` branch
  * References the helpdesk issue in the PR message
  * [GitHub documentation on how to create a pull request](https://help.github.com/articles/creating-a-pull-request/)

* Grab a cup of coffee and wait patiently for an administrator to sign your certificate request
* Once an admin notifies you that everything is setup, you can [sync your fork](https://docs.github.com/en/github/collaborating-with-pull-requests/working-with-forks/syncing-a-fork) then pull it to retrieve your certificate from `./cert/pki/issued/<your-jenkins-username>.crt`
* We recommend you to move the generated files and the ca.cert to an hidden folder in your home (`~/.cert`):

  ```bash
  mkdir -p ~/.cert/jenkins-infra
  mv ./cert/pki/issued/<your-jenkins-username>.crt ~/.cert/jenkins-infra/<your-jenkins-username>.crt
  mv ./cert/pki/private/<your-jenkins-username>.key ~/.cert/jenkins-infra/<your-jenkins-username>.key
  cp ./cert/pki/ca.crt ~/.cert/jenkins-infra/ca.crt
  ```

* Then, create the following configuration file (wether your are on Linux, macOS or Windows) `private-jenkins-infra.ovpn` on your Desktop:

  ```text
  client
  remote private.vpn.jenkins.io 443
  ca "/absolute/path/to/.cert/ca.crt"
  cert "/absolute/path/to/.cert/<your-jenkins-username>.crt"
  key "/absolute/path/to/.cert/<your-jenkins-username>.key"
  auth-user-pass
  dev tun
  proto tcp
  nobind
  auth-nocache
  script-security 2
  persist-key
  persist-tun
  remote-cert-tls server
  user nobody
  group nobody
  ```

  * Some important rules:
    * The file name does not matter but it MUST have an extension `.ovpn` to let your system detect it
    * The content of the file does not support the `~` shortcut, neither variables (`$HOME`/`%HOME%`). Please use absolute paths.
  * Then import this file (e.g. double click or use the appropriate command line) into  your VPN tool:
    * on macOS, we recommend using [Tunnelblick](https://tunnelblick.net/), an OpenVPN client
    * on Linux, we recommend using [NetworkManager](https://wiki.archlinux.org/title/NetworkManager) client. Note that in that case, **you must enable** the option `Use this connection only for resources on its network`
    * on Windows, we recommend using [OpenVPN Connect](https://openvpn.net/client-connect-vpn-for-windows/) client.

* ⚠️ When connecting, your VPN client requires a username and password. Use your Jenkins project account (same username + password as accounts.jenkins.io, issues.jenkins.io, ci.jenkins.io).

#### Windows only

If you want to use multiple VPN connections at the same time with OpenVPN, you have to install a new TAP adapter. This can be very easily by running *as Admin* the `C:\Program Files\TAP-Windows\bin\addtap.bat`. The [TAP-Windows](https://community.openvpn.net/openvpn/wiki/ManagingWindowsTAPDrivers) tool is installed in parallel with OpenVPN.

### HowTo show request information

* Enter in the VPN network directory: `cd ~/.cert`
* Run `make show-req name=<your-jenkins-username>`

### Howto validate your certificate

You can test if your private key matches your certificate and certificate request by running following commands:

```bash
openssl pkey -in ~/.cert/pki/private/<your-jenkins-username>.key -pubout -outform pem | sha256sum
# Should be equal to
openssl x509 -in ~/.cert/pki/issued/<your-jenkins-username>.crt -pubkey -noout -outform pem | sha256sum
# And also equal to
openssl req -in ~/.cert/pki/reqs/<your-jenkins-username>.req -pubkey -noout -outform pem | sha256sum
```

### DNS Problems

If you are having issues connecting to resources behind the VPN, but the VPN appears to be working correctly, check your DNS settings.  Some providers seem to filter out requests to the zone.  To test, try `dig release.ci.jenkins.io`, you should get something like this:

<!-- markdownlint-disable MD033 -->
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

To enable a different DNS provider only when connected to the VPN you can add the following to you OpenVPN config file:

```dhcp-option DNS 8.8.8.8```

## Administrator

### HowTo become an administrator

To add/revoke certificates, you must be allowed to decrypt sensitive files such as `./cert/pki/private/ca.key.enc`.

These files are encrypted with [sops](https://github.com/mozilla/sops), your public gpg key must be added to `./.sops.yaml` by an existing administrator to decrypt them.

This repository relies on [easy-rsa](https://github.com/OpenVPN/easy-rsa/blob/master/README.quickstart.md), used under the hood by a custom Golang CLI wrapper named `easyvpn`.

### HowTo Decrypt the Certificate Authority Key

* Ensure that you are an administrator (Check the section [HowTo become an administrator](#howto-become-an-administrator))
* Execute the command `make -C cert decrypt` from the root of the repository to decrypt the ca.key to `./cert/pki/private/ca.key` (which is a **secret** that **must remain git-ignored**)

### HowTo show certificate information

* Install [sops](https://github.com/mozilla/sops)
* Enter in the VPN network directory: `cd ~/.cert`
* Decrypt the required files as described in [HowTo Decrypt the Certificate Authority Key](#howto-decrypt-the-certificate-authority-key)
* Run `make show-cert name=<your-jenkins-username>`

#### HowTo approve client access?

To validate and sign a client certificate, you are going to execute the following actions:

* Build EASYVPN binary by running one of the following commands depending on your
  * `make init_osx`
  * `make init_linux`
  * `make init_windows` then copy `./utils/easyvpn/easyvpn.exe` at the root of this repository
* Using [the official GitHub command line `gh`](https://github.com/cli/cli), checkout the Pull Request of by the requester
  to retrieve their [CRL](https://en.wikipedia.org/wiki/Certificate_revocation_list) your local machine:

```shell
gh pr checkout <Pull Request ID>
```

* Sign the certificate request: `./easyvpn sign <CN_to_sign>`
  * by default this will create a Client Configuration file for the "private" VPN (private.vpn.jenkins.io), and store this file in  `./cert/ccd/private/`
* A git commit is automatically created on the local branch
* Push the approval commit on the current pull request with `git push` (the remote and local branch name are configured by the `gh` command line)
* Approve and merge the Pull Request to the `main` branch with the signed CRL
* Once merged, a new tag should be created automatically with automatic publishing of the image
* The Docker image tag should be automatically updated in the next 24h in the [puppet](https://github.com/jenkins-infra/jenkins-infra/blob/production/dist/profile/manifests/openvpn.pp) configuration.

### HowTo revoke client access?

* Build EASYVPN binary by running one of the following commands depending on your
  * `make init_osx`
  * `make init_linux`
  * `make init_windows` and copy `./utils/easyvpn/easyvpn.exe` at the root of this repository
* Revoke the certificate: `./easyvpn revoke <CN_to_sign>`
* A git commit is automatically created on the local branch
* Push the revocation commit (PR or branch, whatever you choose)
* The Docker image tag should be automatically updated in the next 24h in the [puppet](https://github.com/jenkins-infra/jenkins-infra/blob/production/dist/profile/manifests/openvpn.pp) configuration.

#### HowTo renew certificate revocation list

If the [CRL (Certificate Revocation list)](https://en.wikipedia.org/wiki/Certificate_revocation_list) expired, then the OpenVPN logs will contain errors like 'VERIFY ERROR: depth=0, error=CRL has expired:...'
We can run `openssl crl -in ./cert/pki/crl.pem -noout -text` to validate that the CRL expired and that we need to generate a new one.

To generate a new CRL:

* Decrypt the required files as described in [HowTo Decrypt the Certificate Authority Key](#howto-decrypt-the-certificate-authority-key)
* Generate a new crl.pem - `cd cert ; ./easyrsa gen-crl ; cd ..`
* Publish the new crl.pem - `git add ./cert/pki/crl.pem && git commit ./cert/pki/crl.pem -s -m 'Renew revocation list certificate'`
* Before pushing, Delete local ca.key - `rm ./cert/pki/private/ca.key`
* Check new expiration date with : `openssl crl -nextupdate -in /Users/smerle/code/docker-openvpn/cert/pki/crl.pem -noout`

### HowTo Renew Server-side Certificate?

* Build EASYVPN binary by running one of the following commands depending on your operating system:
  * `make init_osx`
  * `make init_linux`
  * `make init_windows` and copy `./utils/easyvpn/easyvpn.exe` at the root of this repository
* Decrypt the required files as described in [HowTo Decrypt the Certificate Authority Key](#howto-decrypt-the-certificate-authority-key)
* Revoke actual certificate (even if it is already expired): `./easyvpn revoke private.vpn.jenkins.io`
* Generate a new certificate + key, with the server DNS as argument: `./easyvpn request private.vpn.jenkins.io`

  > The generated key is in `./cert/pki/private/private.vpn.jenkins.io.key` **must** remain **secret**!

* Sign the request as a "server" request:

  ```shell
  cd ./certs # Running the signing command from this folder is mandatory.
  ./easyrsa --batch sign-req server private.vpn.jenkins.io
  ```

* Ensure that you git-added, git-commited and pushed the changes, without ANY secrets (which should be git-ignored)

* Update the secrets in the encrypted hieradata for OpenVPN in <https://github.com/jenkins-infra/jenkins-infra>

## Docker

### Configuration

This image can be configured at runtime with different environment variables:

* `AUTH_LDAP_BINDDN` **Define user dn used to query the ldap database**
* `AUTH_LDAP_URL` **Define ldap endpoint url**
* `AUTH_LDAP_PASSWORD` **Define user dn password**
* `AUTH_LDAP_GROUPS_MEMBER` **Define required group member to authenticate**
* `OPENVPN_NETWORK_NAME` **Define the network name from config.yaml to use**
* `OPENVPN_SERVER_SUBNET` **Define the VPN subnet**
* `OPENVPN_SERVER_NETMASK` **Define the netmask associated to the VPN subnet**

Some examples can be found inside [docker-compose.yaml](docker/docker-compose.yaml)

### Testing

To test this image, you need a "mock" ldap and SSL certificates, then go in the root folder and run `make start` to start the ldap and vpn service.

> ⚠️ Certificates must be readable by UID 101! ⚠️

## Infrastructure

This project is designed to work with the following requirements:

* Machine provisioned by [Terraform](https://github.com/jenkins-infra/azure)
* Service configured and orchestrated by [Puppet](https://github.com/jenkins-infra/jenkins-infra/blob/production/dist/profile/manifests/openvpn.pp)

## Contributing

Feel free to contribute to this image by:

1. Fork this project into your account
2. Make your changes in your local fork
3. Submit a pull request with a description and a link to a [jenkins-infra/helpdesk issue](https://github.com/jenkins-infra/helpdesk)
4. Ask for a review

## Issue

Please report any issue on the Jenkins infrastructure [jenkins-infra/helpdesk tracker](https://github.com/jenkins-infra/helpdesk)

## Links

* [How to contribute to OSS?](https://opensource.guide/how-to-contribute/)
* [jenkins-infra/azure](https://github.com/jenkins-infra/azure)
* [jenkins-infra/jenkins-infra](https://github.com/jenkins-infra/jenkins-infra/blob/production/dist/profile/manifests/openvpn.pp)
* [jenkins-infra/docker-openvpn](https://github.com/jenkins-infra/docker-openvpn)
* [mozilla/sops](https://github.com/mozilla/sops)
* [openvpn/easy-rsa](https://github.com/OpenVPN/easy-rsa)
