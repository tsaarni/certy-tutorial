import datetime

from certy import Credential
from cryptography.x509.verification import PolicyBuilder, Store

def main():

    # Define the root CA (self-signed certificate).
    ca = Credential().subject("CN=ca")

    # Define the intermediate CAs, marked as CA certificates and issued by the root CA.
    server_sub_ca = Credential().subject("CN=server-sub-ca").ca().issuer(ca)
    client_sub_ca = Credential().subject("CN=client-sub-ca").ca().issuer(ca)

    # Define the server certificate, signed by the server sub-CA and with localhost as SAN.
    server = (
        Credential()
        .subject("CN=test-server")
        .issuer(server_sub_ca)
        .subject_alt_names("DNS:localhost")
    )

    # Define the client certificate, signed by the client sub-CA.
    client = (
        Credential()
        .subject("CN=test-client")
        .issuer(client_sub_ca)
    )

    # Print information about the server and client certificate.
    print_certificate("Server certificate:", server.get_certificate())
    print_certificate("Client certificate:", client.get_certificate())

    # Verify the server certificate with the trusted root CA.
    # Use ClientVerifier for convenience for this example.
    verifier = PolicyBuilder().store(Store(ca.get_certificates())).build_client_verifier()

    verifier.verify(server.get_certificate(), server.get_certificates())
    print("Server certificate verified successfully.")

    # Define an expired client certificate.
    expired = (
        Credential()
        .subject("CN=expired-client")
        .issuer(client_sub_ca)
        .not_before(datetime.datetime(2019, 1, 1, 0, 0, 0))
        .not_after(datetime.datetime(2020, 1, 1, 0, 0, 0))
    )

    # Verify the expired client certificate.
    try:
        verifier.verify(expired.get_certificate(), expired.get_certificates())
    except Exception as e:
        print("Expired client certificate verification failed as expected:", e)

def print_certificate(header, cert):
    output = f"""{header}
    Subject: {cert.subject.rfc4514_string()}
    Issuer: {cert.issuer.rfc4514_string()}
    Not valid before: {cert.not_valid_before_utc}
    Not valid after: {cert.not_valid_after_utc}
    """
    print(output)

if __name__ == "__main__":
    main()
