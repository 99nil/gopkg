// Copyright © 2022 zc2638 <zc2638@qq.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cert

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"time"
)

func GenerateSelfSignedCertAndKey(keySize int, cert *x509.Certificate) (crt, key []byte, err error) {
	// Generate a key pair
	pk, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, nil, err
	}

	// This number represents a unique serial number issued by the CA,
	// which is represented here by a number
	if cert.SerialNumber == nil {
		serialNumber, err := rand.Int(rand.Reader, pk.N)
		if err != nil {
			return nil, nil, err
		}
		cert.SerialNumber = serialNumber
	}
	if cert.NotAfter.Before(time.Now()) {
		cert.NotAfter = time.Now().Add(time.Hour * 24 * 365)
	}
	if cert.KeyUsage == 0 {
		cert.KeyUsage = x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign
	}
	if len(cert.ExtKeyUsage) == 0 {
		cert.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	}

	// Create certificate.
	// Here the second parameter and the third parameter are the same,
	// which means that the certificate is a self-signed certificate,
	// and the return value is a DER encoded certificate.
	certificate, err := x509.CreateCertificate(rand.Reader, cert, cert, &pk.PublicKey, pk)
	if err != nil {
		return nil, nil, err
	}

	// Define the CERTIFICATE PEM block
	block := &pem.Block{
		Type:    "CERTIFICATE",
		Headers: nil,
		Bytes:   certificate,
	}
	// Generate PEM-encoded certificate content
	var crtBuf bytes.Buffer
	if err := pem.Encode(&crtBuf, block); err != nil {
		return nil, nil, err
	}

	// Define the PRIVATE KEY PEM block
	block = &pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(pk),
	}
	// Generate PEM-encoded key content
	var keyBuf bytes.Buffer
	if err := pem.Encode(&keyBuf, block); err != nil {
		return nil, nil, err
	}
	return crtBuf.Bytes(), keyBuf.Bytes(), nil
}
