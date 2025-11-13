# TPM Attestation Security Overview

This document explains how the Trusted Platform Module (TPM) is used for attestation in this project, and describes the roles of the TPM Quote, Signature, and TPM Public Key.

## How the TPM is Used

The TPM is a hardware security module that provides cryptographic operations and secure storage. In this project, it is used to:

1. **Open the TPM Device**
   - The code attempts to open the TPM device at `/dev/tpmrm0` or `/dev/tpm0`.
   - A nonce (random challenge) is provided to bind the attestation to a specific request.

2. **Read PCRs (Platform Configuration Registers)**
   - PCRs are special TPM registers that store hashes representing the system's boot and runtime state.
   - The code reads selected PCRs using SHA-256 as the hash algorithm. These PCRs reflect the measured state of the system.

3. **Generate a Quote**
   - The TPM generates a cryptographic "quote" over the selected PCR values and the provided nonce.
   - This quote proves the system's state at a specific time and prevents replay attacks.

4. **Sign the Quote**
   - The TPM signs the quote using its Attestation Key with ECC (Elliptic Curve Cryptography), a private key securely stored in the TPM.
   - The signature ensures the authenticity and integrity of the quote.

5. **Provide the TPM Public Key**
   - The public part of the Attestation Key is exported in PEM format.
   - This key is used by verifiers to check the signature on the quote.

## Key Concepts

| Concept           | Description |
|-------------------|-------------|
| **Quote**         | A TPM quote is a signed structure containing the selected PCR values and the provided nonce. It proves the system's state at the time of attestation and binds the result to a specific challenge. |
| **Signature**     | The cryptographic signature over the quote, created by the TPM's Attestation Key using ECC as signing algorithm. It allows external parties to verify the authenticity and integrity of the quote. |
| **TPM Public Key**| The public part of the Attestation Key, distributed as part of the attestation report. It is used by verifiers to check the signature on the quote. |

## Security Considerations

- **Nonce Usage:** The nonce ensures that each attestation is unique and cannot be replayed.
- **PCR Selection:** Only trusted PCRs should be selected for attestation, as they reflect the system's measured state.
- **Key Management:** The Attestation Key is generated and stored securely within the TPM. Only the public part is exported.
- **Signature Verification:** Verifiers must use the provided TPM public key to check the quote's signature and validate the attestation.

For implementation details, see `pkg/tpm/tpm.go`.
