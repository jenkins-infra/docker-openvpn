package easyrsa

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
)

var (
	easyrsaDir = "cert"
	debug      = true
)

func easyrsa(args ...string) error {
	cmd := exec.Command("./easyrsa", args...)
	cmd.Dir = easyrsaDir

	fmt.Println(cmd.Dir)

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err := cmd.Run()

	if debug {
		fmt.Printf("Exec: %v", cmd.Args)
		fmt.Printf("%v\n", outb.String())
	}

	if err != nil {
		fmt.Println(errb.String())
	}
	return err
}

// GenerateRevocationListCert generate a certificate revocation List
func GenerateRevocationListCert() error {
	err := easyrsa("revoke")
	return err
}

// RequestClientCert generate a client private key and an associated certificate request
func RequestClientCert(CNs []string) []error {
	var err []error
	for _, CN := range CNs {
		err = append(err, easyrsa("--batch", "--req-cn="+CN, "gen-req", CN, "nopass"))
	}
	return err
}

// RevokeClientCert revoke a client certificate
func RevokeClientCert(CNs []string) []error {
	var errors []error
	for _, CN := range CNs {
		errors = append(errors, easyrsa("--batch", "revoke", CN))
	}
	return errors
}

// SignClientRequest sign client certificate requests
func SignClientRequest(CNs []string) []error {
	var errors []error
	for _, CN := range CNs {
		errors = append(errors, easyrsa("--batch", "sign-req", "client", CN))
	}

	// Delete useless files generated from easyrsa
	files := []string{
		"index.txt.attr.old",
		"index.txt.old",
		"serial.old",
		"extensions.temp"}

	for _, file := range files {
		err := os.Remove(path.Join(easyrsaDir, "pki", file))
		if err != nil {
			fmt.Println(err)
		}
	}

	return errors
}

// ShowClientCertificate display client certificate information
func ShowClientCertificate(CN string) error {
	err := easyrsa("show-cert", CN)
	return err
}

// ShowClientRequestCertificate display client certificate request information
func ShowClientRequestCertificate(CN string) error {
	err := easyrsa("show-req", CN)
	return err
}
