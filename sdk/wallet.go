package sdk

import (
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
	"strings"

	gw "github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

type JobWallet struct {
	UserID string
	OrgMSP string
	Domain string
}

func (jw *JobWallet) Create() bool {
	wallet := gw.NewInMemoryWallet()
	cert := jw.getCert()
	key := jw.getKey()
	identity := gw.NewX509Identity(jw.OrgMSP, cert, key)

	err := wallet.Put(jw.UserID, identity)
	if err != nil {
		log.Fatalf("Failed to put identity into wallet: %v", err)
		return false
	}

	return true
}

func (jw *JobWallet) getCert() string {
	certPath := "./crypto/users/" + jw.UserID + "@" + jw.OrgMSP + "." + jw.Domain + "/msp/signcerts/cert.pem"
	// Read the certificate
	certBytes, err := os.ReadFile(certPath)
	if err != nil {
		log.Fatalf("Failed to read certificate file: %v", err)
	}
	block, _ := pem.Decode(certBytes)
	if block == nil || block.Type != "CERTIFICATE" {
		log.Fatalf("Failed to parse PEM block containing the certificate")
	}
	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatalf("Failed to parse certificate: %v", err)
		return ""
	}

	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certificate.Raw,
	})

	return string(certPEM)
}

func (jw *JobWallet) getKey() string {
	keyPath := "./crypto/users/" + jw.UserID + "@" + jw.OrgMSP + "." + jw.Domain + "/msp/keystore/priv_sk"
	// Read the private key
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("Failed to read private key file: %v", err)
	}
	keyBlock, _ := pem.Decode(keyBytes)
	if keyBlock == nil || !strings.HasSuffix(keyBlock.Type, "PRIVATE KEY") {
		log.Fatalf("Failed to parse PEM block containing the private key")
	}

	// Convert the private key to PEM format string
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyBlock.Bytes,
	})

	return string(keyPEM)
}

func Remove() {

}
