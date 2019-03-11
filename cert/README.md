# Jenkins Vpn Keys

This projects holds vpn keys for connecting on Jenkins Infrastructure vpn and designed to work with following pieces:

* Machine provisioned by [Terraform](https://github.com/jenkins-infra/azure)
* Docker image defined [here](https://github.com/jenkins-infra/openvpn)
* Service configured and orchestrated by [Puppet](https://github.com/jenkins-infra/jenkins-infra/blob/staging/dist/profile/manifests/openvpn.pp) 

If you think that you should have access to a specific vpn network, feel free to read [HowTo Get client access](#howto-get-client-access).

If you want to help with the administration task, everything is explained in section [HowTo become an administrator](#howto-become-an-administrator).

## Client
### HowTo get client access
In order to connect on one of the Jenkins infrastructure vpn, you need a certificate containing your jenkins username account, 
Then this certificate must be signed by an administrator.

Feel free to follow next action points:

* Fork this repository on your own Github account: [fork a repo](https://help.github.com/articles/fork-a-repo/)
* Enter in the needed vpn network directory: `cd cert`
* Create your private key and certificate request: `make request name=<your username>`
* Create a new Pull Request on master branch: [create a pull request](https://help.github.com/articles/creating-a-pull-request/)
* Grab a cup of coffee and wait patiently until an administrator issues your certificate.
* Once ready your certificate can be retrieve from `./cert/pki/issued/<your_username>.crt`

### HowTo show request information

* Enter in the vpn network directory: `cd cert`
* Run `make show-req name=<username>`

### HowTo show certificate information

* Enter in the vpn network directory: `cd cert`
* Run `make show-certs name=<username>`

## Administrator
### HowTo become an administrator
In order to do any administrative tasks, you must be allow to decrypt `cert/pki/private/ca.key.enc`.
This file is encrypted with [sops](https://github.com/mozilla/sops) and you are public gpg key must be added to .sops.yaml by an existing administrator in order to be allow to run `make decrypt`.

This repository relies on [easy-rsa](https://github.com/OpenVPN/easy-rsa/blob/master/README.quickstart.md).

### HowTo approve client access?
In order to validate and sign a client certificate, your are going to do following actions

* Enter in the vpn network directory: `cd cert`
* Decrypt ca.key: `make decrypt`
* Sign certificate request: `make sign`

### HowTo revoke client access?

* Enter in the vpn network directory: `cd cert`
* Revoke certificate: `make revoke`

## Links
* [jenkins-infra/azure](https://github.com/jenkins-infra/azure)
* [jenkins-infra/jenkins-infra](https://github.com/jenkins-infra/jenkins-infra/blob/staging/dist/profile/manifests/openvpn.pp)
* [jenkins-infra/openvpn](https://github.com/jenkins-infra/openvpn)
* [mozilla/sops](https://github.com/mozilla/sops)
* [openvpn/easy-rsa](https://github.com/OpenVPN/easy-rsa)
