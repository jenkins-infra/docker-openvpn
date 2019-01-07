# README

This project contents everything needed to build an openvpn docker image integrated with [openldap](https://github.com/jenkins-infra/ldap).
This image is designed to be used inside the Jenkins Infrastructure Project.

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

## ISSUE
Please report any issue on the jenkins infrastructure [project](https://issues.jenkins-ci.org/secure/Dashboard.jspa)

## LINKS
[How to contribute to OSS?](https://opensource.guide/how-to-contribute/)
