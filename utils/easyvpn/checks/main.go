package checks

import (
	"errors"
	"fmt"
	"path"
	"slices"
	"strings"

	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/config"
	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/easyrsa"
	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/helpers"
)

// IsAllCertsSigned validate that every requested certificate are signed
func IsAllCertsSigned(certDir string) (bool, []error) {

	var errs []error
	result := true
	req := helpers.GetUsernameFile(path.Join(certDir, "pki/reqs"), ".req")
	crt := helpers.GetUsernameFile(path.Join(certDir, "pki/issued"), ".crt")

	slices.Sort(req)
	slices.Sort(crt)

	if len(req) != len(crt) {
		result = false
		msg := fmt.Sprintf("Numbers of requested and signed certificates mismatch:\n\t* %v Requested (%v)\n\t* %v Signed (%v)\n", len(req), req, len(crt), crt)
		err := errors.New(msg)
		fmt.Println(err)

		errs = append(errs, err)
	}

	for i := range req {
		var err error
		reqCN := strings.Split(req[i], ".req")
		crtCN := strings.Split(crt[i], ".crt")

		err = easyrsa.ShowClientCertificate(crtCN[0])
		if err != nil {
			result = false
			fmt.Println(err)
			errs = append(errs, err)
		}
		err = easyrsa.ShowClientRequestCertificate(reqCN[0])
		if err != nil {
			result = false
			fmt.Println(err)
			errs = append(errs, err)
		}

		if reqCN[0] != crtCN[0] {
			result = false
			msg := fmt.Sprintf("Requested and signed certificate mismatch:\n\t%v should match %v\n", req[i], crt[i])
			err := errors.New(msg)
			fmt.Println(err)
			errs = append(errs, err)
		}
	}
	return result, errs
}

// IsAllClientConfigured validate that all signed certificate have client configuration
func IsAllClientConfigured(certDir string, net string, globalConfig config.Config) (bool, []error) {

	var errs []error
	result := true
	clientConfig := helpers.GetUsernameFile(path.Join(certDir, "ccd", net), "")
	crt := config.GetUsersWithCertificate(certDir, globalConfig)

	slices.Sort(clientConfig)

	if len(clientConfig) != len(crt) {
		result = false
		msg := fmt.Sprintf("Numbers of client configuration and signed certificates mismatch:\n\t* %v Requested (%v)\n\t* %v Signed (%v)\n", len(clientConfig), clientConfig, len(crt), crt)
		err := errors.New(msg)
		fmt.Println(err)

		errs = append(errs, err)
	}
	return result, errs
}
