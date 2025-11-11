package tpm

import (
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/google/go-tpm-tools/client"
	"github.com/google/go-tpm/legacy/tpm2"
	"github.com/google/go-tpm/tpmutil"
)

const (
	tpmDevicePath                = "/dev/tpm0"
	tpmResourceManagerDevicePath = "/dev/tpmrm0"
)

// Device represents a TPM device.
type Device struct {
	rwc      io.ReadWriteCloser
	nonce    []byte
	lastSig  []byte // cached signature for the most recent quote
	akPubPEM []byte // cached AK public key (PEM), matches the most recent quote
}

// OpenTPMDevice creates and opens a TPM device with the given nonce.
func OpenTPMDevice(nonce []byte) (*Device, error) {
	tpmDevice := Device{
		nonce: nonce,
	}

	err := tpmDevice.openTPM()
	if err != nil {
		return nil, fmt.Errorf("failed to open TPM device: %w", err)
	}

	defer func() {
		_ = tpmDevice.Close()
	}()

	return &tpmDevice, nil
}

// ReadPCRs reads the specified PCRs from the TPM device.
func (d *Device) ReadPCRs(pcrSelection *tpm2.PCRSelection) (map[int]string, error) {
	if pcrSelection == nil || len(pcrSelection.PCRs) == 0 {
		return nil, ErrInvalidPCRs
	}

	pcrs, err := client.ReadPCRs(d.rwc, *pcrSelection)
	if err != nil {
		return nil, fmt.Errorf("failed to read PCRs: %w", err)
	}

	mapPCRs := make(map[int]string, len(pcrSelection.PCRs))
	for _, pcrIndex := range pcrSelection.PCRs {
		//nolint:gosec // PCR index is controlled 0-15
		fetchedPCR, ok := pcrs.GetPcrs()[uint32(pcrIndex)]
		if !ok {
			return nil, fmt.Errorf("%w: %d", ErrPCRNotFound, pcrIndex)
		}

		mapPCRs[pcrIndex] = hex.EncodeToString(fetchedPCR)
	}

	return mapPCRs, nil
}

// Quote generates a TPM quote for the given PCRs.
func (d *Device) Quote(pcrSelection *tpm2.PCRSelection) ([]byte, error) {
	if d.rwc == nil {
		return nil, ErrTPMNotOpened
	}

	if len(d.nonce) == 0 {
		return nil, ErrInvalidNonce
	}

	if pcrSelection == nil || len(pcrSelection.PCRs) == 0 {
		return nil, ErrInvalidPCRs
	}

	attestationKey, err := client.AttestationKeyECC(d.rwc)
	if err != nil {
		return nil, fmt.Errorf("creating AK: %w", err)
	}

	defer attestationKey.Close()

	publicKey := attestationKey.PublicKey()

	pkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("marshal public key: %w", err)
	}

	d.akPubPEM = pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pkix})

	attest, sig, err := tpm2.Quote(
		d.rwc,
		attestationKey.Handle(),
		"",
		"unused", // Unused in downstream code
		d.nonce,
		*pcrSelection,
		tpm2.AlgECC,
	)
	if err != nil {
		return nil, fmt.Errorf("tpm2.Quote failed: %w", err)
	}

	sigBytes, err := sig.Encode()
	if err != nil {
		return nil, fmt.Errorf("marshal signature: %w", err)
	}

	d.lastSig = sigBytes

	return attest, nil
}

// Signature returns the signature for the TPM quote.
func (d *Device) Signature() ([]byte, error) {
	if len(d.lastSig) == 0 {
		return nil, ErrNoSignatureCached
	}

	return d.lastSig, nil
}

// Close closes the TPM device connection.
func (d *Device) Close() error {
	if d.rwc != nil {
		_ = d.rwc.Close()
	}

	return nil
}

// GetTPMPublicKey retrieves the TPM's public key.
func (d *Device) GetTPMPublicKey() ([]byte, error) {
	if len(d.akPubPEM) == 0 {
		return nil, ErrNoPublicKeyCached
	}

	return d.akPubPEM, nil
}

func (d *Device) openTPM() error {
	paths := []string{tpmResourceManagerDevicePath, tpmDevicePath}
	for _, p := range paths {
		_, err := os.Stat(p)
		if err == nil {
			tpm, err := tpmutil.OpenTPM(p)
			if err != nil {
				return fmt.Errorf("failed to open TPM device at %s: %w", p, err)
			}

			d.rwc = tpm

			return nil
		}
	}

	return ErrTPMNotFound
}

// GetPCRSelection constructs a PCR selection structure for the given PCR indices.
func GetPCRSelection(pcrIndices []int) *tpm2.PCRSelection {
	sort.Ints(pcrIndices)

	return &tpm2.PCRSelection{
		Hash: tpm2.AlgSHA256,
		PCRs: pcrIndices,
	}
}
