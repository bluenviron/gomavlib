// Package main contains an example.
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
)

// this example shows how to:
// 1) create a node which communicates with a custom TCP/TLS endpoint in server mode.
// 2) print incoming messages.

func main() {
	// ensure the certificate and key exists
	if err := EnsureCertsExist(); err != nil {
		fmt.Println("Error ensuring certificates:", err)
		return
	}

	// create a node which communicates with a custom TCP/TLS endpoint in server mode
	node := &gomavlib.Node{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointCustomServer{
				Listen: func() (net.Listener, error) {
					// Loads the certificate and key from the generated certs dir
					cert, err := tls.LoadX509KeyPair("certs/cert.pem", "certs/key.pem")
					if err != nil {
						return nil, err
					}

					return tls.Listen("tcp", ":5600", &tls.Config{
						Certificates: []tls.Certificate{cert},
					})
				},
				Label: "TCP/TLS",
			},
		},
		Dialect:     ardupilotmega.Dialect,
		OutVersion:  gomavlib.V2, // change to V1 if you're unable to communicate with the target
		OutSystemID: 10,
	}
	err := node.Initialize()
	if err != nil {
		panic(err)
	}
	defer node.Close()

	// print incoming messages
	for evt := range node.Events() {
		if frm, ok := evt.(*gomavlib.EventFrame); ok {
			log.Printf("received: id=%d, %+v\n", frm.Message().GetID(), frm.Message())
		}
	}
}

// Below are just functions to check and generate certificate and private key
// they are just here to make this example simpler to run

// EnsureCertsExist checks if the cert.pem and key.pem exist in the certs directory,
// and if not, generates them.
func EnsureCertsExist() error {
	// Check if cert.pem exists
	if _, err := os.Stat("certs/cert.pem"); os.IsNotExist(err) {
		fmt.Println("cert.pem not found. Generating certificates...")
		return GenerateCertAndKey()
	}

	// Check if key.pem exists
	if _, err := os.Stat("certs/key.pem"); os.IsNotExist(err) {
		fmt.Println("key.pem not found. Generating certificates...")
		return GenerateCertAndKey()
	}

	return nil
}

// GenerateCertAndKey generates a self-signed certificate and private key, saving them to the certs/ directory.
func GenerateCertAndKey() error {
	// Create the certs directory if it doesn't exist
	err := os.MkdirAll("certs", os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create certs directory: %w", err)
	}

	// Generate RSA private key
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "gomavlib",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour), // valid for 1 year
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		BasicConstraintsValid: true,
	}

	// Create the certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %w", err)
	}

	// Save the certificate
	certOut, err := os.Create("certs/cert.pem")
	if err != nil {
		return fmt.Errorf("failed to create cert.pem: %w", err)
	}
	defer certOut.Close()

	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	if err != nil {
		return fmt.Errorf("failed to encode certificate to PEM: %w", err)
	}

	// Save the private key
	keyOut, err := os.Create("certs/key.pem")
	if err != nil {
		return fmt.Errorf("failed to create key.pem: %w", err)
	}
	defer keyOut.Close()

	err = pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	if err != nil {
		return fmt.Errorf("failed to encode private key to PEM: %w", err)
	}

	fmt.Println("cert.pem and key.pem generated in the 'certs/' directory.")
	return nil
}
