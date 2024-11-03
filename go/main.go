package main

import (
	"crypto/x509"
	"fmt"
	"time"

	"github.com/tsaarni/certyaml"
)

func main() {
	// Define the root CA (self-signed certificate).
	ca := certyaml.Certificate{Subject: "cn=ca"}

	// Define the intermediate CAs, marked as CA certificates and issued by the root CA.
	isCa := true
	serverSubCa := certyaml.Certificate{Subject: "cn=server-sub-ca", IsCA: &isCa, Issuer: &ca}
	clientSubCa := certyaml.Certificate{Subject: "cn=client-sub-ca", IsCA: &isCa, Issuer: &ca}

	// Define the server certificate, signed by the server sub-CA and with localhost as SAN.
	server := certyaml.Certificate{
		Subject:         "cn=test-server",
		SubjectAltNames: []string{"DNS:localhost"},
		Issuer:          &serverSubCa,
	}

	// Define the client certificate, signed by the client sub-CA.
	client := certyaml.Certificate{
		Subject: "cn=test-client",
		Issuer:  &clientSubCa,
	}

	// Print information about the server and client certificate.
	cert, _ := server.X509Certificate()
	printCert("Server certificate:", &cert)
	cert, _ = client.X509Certificate()
	printCert("Client certificate:", &cert)

	// Create a CertPool and add the trusted root CA certificate for verification purposes.
	trustedCertPool := x509.NewCertPool()
	trustedCertPool.AppendCertsFromPEM(ca.CertPEM())

	chainCertPool := x509.NewCertPool()
	chainCertPool.AppendCertsFromPEM(server.CertPEM())

	// Verify the server certificate using the cert pool.
	opts := x509.VerifyOptions{
		Roots:         trustedCertPool, // Trusted root CAs.
		Intermediates: chainCertPool,   // Additional intermediate CAs required to build full certificate chain.
	}
	serverCert, _ := server.X509Certificate()
	if _, err := serverCert.Verify(opts); err != nil {
		fmt.Println("Failed to verify server certificate:", err)
		return
	}
	fmt.Println("Server certificate verified successfully")

	// Define an expired client certificate.
	expiredAt := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	expired := certyaml.Certificate{
		Subject:  "cn=expired-client",
		Issuer:   &clientSubCa,
		NotAfter: &expiredAt,
	}
	chainCertPool.AppendCertsFromPEM(expired.CertPEM())

	// Verify the expired client certificate using the cert pool.
	expiredCert, _ := expired.X509Certificate()
	_, err := expiredCert.Verify(opts)
	fmt.Println("Expired client certificate verification failed as expected:", err)
}

func printCert(header string, cert *x509.Certificate) {
	fmt.Println(header)
	fmt.Println("  Subject:", cert.Subject)
	fmt.Println("  Issuer:", cert.Issuer)
	fmt.Println("  NotBefore:", cert.NotBefore)
	fmt.Println("  NotAfter:", cert.NotAfter)
	fmt.Println("  SubjectAltNames:", cert.DNSNames)
}
