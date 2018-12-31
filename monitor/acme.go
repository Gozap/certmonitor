/*
 * Copyright 2018 Gozap, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package monitor

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/xenolf/lego/certificate"

	"github.com/xenolf/lego/providers/dns/alidns"

	"github.com/xenolf/lego/challenge"
	"github.com/xenolf/lego/providers/dns/godaddy"

	"github.com/xenolf/lego/certcrypto"

	"github.com/xenolf/lego/lego"

	"github.com/gozap/certmonitor/conf"

	"github.com/xenolf/lego/registration"
)

type ACMEUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *ACMEUser) GetEmail() string {
	return u.Email
}
func (u ACMEUser) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *ACMEUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

var acmeOnce sync.Once

func ReNew(website conf.WebsiteConfig) error {

	// Create a user. New accounts need an email and private key to start.
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}

	acmeUser := ACMEUser{
		Email: conf.ACME.Email,
		key:   privateKey,
	}

	config := lego.NewConfig(&acmeUser)
	config.KeyType = certcrypto.RSA2048

	// A client facilitates communication with the CA server.
	client, err := lego.NewClient(config)
	if err != nil {
		return err
	}

	// get dns provider
	var provider challenge.Provider
	var providerConfig conf.ACMEProviderConfig

	for _, p := range conf.ACME.Providers {
		if p.Name == website.DNSProvider {
			providerConfig = p
		}
	}

	switch website.DNSProvider {
	case "alidns":
		alidnsConfig := alidns.NewDefaultConfig()
		alidnsConfig.APIKey = providerConfig.APIKey
		alidnsConfig.SecretKey = providerConfig.APISecret
		provider, err = alidns.NewDNSProviderConfig(alidnsConfig)
		if err != nil {
			return err
		}

	case "godaddy":
		godaddyConfig := godaddy.NewDefaultConfig()
		godaddyConfig.APIKey = providerConfig.APIKey
		godaddyConfig.APISecret = providerConfig.APISecret
		provider, err = godaddy.NewDNSProviderConfig(godaddyConfig)
		if err != nil {
			return err
		}

	default:
		return errors.New("unsupported dns provider")
	}

	// set dns provider
	err = client.Challenge.SetDNS01Provider(provider)
	if err != nil {
		return err
	}

	// exclude HTTP01、TLSALPN01 challenge
	client.Challenge.Exclude([]challenge.Type{challenge.HTTP01, challenge.TLSALPN01})

	// New users will need to register
	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return err
	}
	acmeUser.Registration = reg

	// parse domain
	var domain string
	var domains []string
	domainArr := strings.Split(website.Domain, ".")
	if len(domainArr) >= 2 {
		domain = domainArr[len(domainArr)-2] + "." + domainArr[len(domainArr)-1]
		domains = append(domains, domain, "*."+domain)
	} else {
		return errors.New(fmt.Sprintf("parse domain [%s] error", website.Domain))
	}

	// create request
	request := certificate.ObtainRequest{
		Domains: domains,
		Bundle:  true,
	}

	// obtain cert
	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		return err
	}

	// save cert
	baseName := strings.Replace(domain, ".", "-", -1)
	certFileName := baseName + ".cer"
	csrFileName := baseName + ".csr"
	privateKeyFileName := baseName + ".key"
	CAFileName := baseName + "-ca.cer"

	// create files
	certFile, err := os.OpenFile(filepath.Join(conf.ACME.CertDir, certFileName), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	csrFile, err := os.OpenFile(filepath.Join(conf.ACME.CertDir, csrFileName), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	privateKeyFile, err := os.OpenFile(filepath.Join(conf.ACME.CertDir, privateKeyFileName), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	CAFile, err := os.OpenFile(filepath.Join(conf.ACME.CertDir, CAFileName), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	// save to file
	err = pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certificates.Certificate})
	if err != nil {
		return err
	}

	err = pem.Encode(csrFile, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: certificates.CSR})
	if err != nil {
		return err
	}

	err = pem.Encode(privateKeyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: certificates.PrivateKey})
	if err != nil {
		return err
	}

	err = pem.Encode(CAFile, &pem.Block{Type: "CERTIFICATE", Bytes: certificates.IssuerCertificate})
	if err != nil {
		return err
	}

	return nil

}

func ACMEInit() {
	acmeOnce.Do(func() {
		if _, err := os.Stat(conf.ACME.CertDir); os.IsNotExist(err) {
			err = os.MkdirAll(conf.ACME.CertDir, 0755)
			if err != nil {
				logrus.Panic(err)
			}
		} else if err != nil {
			logrus.Panic(err)
		}
	})
}
