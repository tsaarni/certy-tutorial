package com.example.app;

import java.security.KeyStore;
import java.security.cert.X509Certificate;
import java.time.Instant;
import java.util.Date;

import javax.net.ssl.TrustManager;
import javax.net.ssl.TrustManagerFactory;
import javax.net.ssl.X509TrustManager;

import org.junit.jupiter.api.Test;

import fi.protonode.certy.Credential;

public class AppTest {
    @Test
    public void testWithCerty() throws Exception {
        // Define the root CA (self-signed certificate).
        Credential ca = new Credential().subject("cn=ca").ca(true);

        // Define the intermediate CAs, marked as CA certificates and issued by the root CA.
        Credential serverSubCa = new Credential().subject("cn=server-sub-ca").ca(true).issuer(ca);
        Credential clientSubCa = new Credential().subject("cn=client-sub-ca").ca(true).issuer(ca);

        // Define the server certificate, signed by the server sub-CA and with localhost as SAN.
        Credential server = new Credential().subject("cn=test-server").issuer(serverSubCa)
                .subjectAltName("DNS:localhost");

        // Define the client certificate, signed by the client sub-CA.
        Credential client = new Credential().subject("cn=test-client").issuer(clientSubCa);

        // Print information about the server and client certificate.
        printCertificates("Server certificate:", server.getX509Certificate());
        printCertificates("Client certificate:", client.getX509Certificate());

        // Create a TrustStore and add the root CA certificate to it.
        KeyStore trustStore = KeyStore.getInstance(KeyStore.getDefaultType());
        trustStore.load(null, null);
        trustStore.setCertificateEntry("ca", ca.getX509Certificate());

        // Verify the server certificate with the TrustStore.
        TrustManagerFactory trustManagerFactory = TrustManagerFactory
                .getInstance(TrustManagerFactory.getDefaultAlgorithm());
        trustManagerFactory.init(trustStore);

        TrustManager[] trustManagers = trustManagerFactory.getTrustManagers();
        X509TrustManager x509TrustManager = (X509TrustManager) trustManagers[0];

        x509TrustManager.checkServerTrusted(server.getX509Certificates(), "ECDHE_ECDSA");
        System.out.println("Server certificate verified successfully");

        // Define an expired client certificate for testing purposes.
        Credential expired = new Credential().subject("cn=expired-client").issuer(clientSubCa)
                .notAfter(Date.from(Instant.parse("2020-01-01T00:00:00Z")));

        // Verify the expired certificate.
        try {
            x509TrustManager.checkClientTrusted(expired.getX509Certificates(), "ECDHE_ECDSA");
        } catch (Exception e) {
            System.out.println("Expired client certificate verification failed as expected: " + e.getMessage());
        }
    }

    void printCertificates(String heading, X509Certificate cert) {
        System.out.println(heading);
        System.out.println(cert);
    }
}
