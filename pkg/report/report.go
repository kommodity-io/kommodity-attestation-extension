// Package report provides utilities for generating reports from all the attestation data.
package report

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/kommodity-io/kommodity-attestation-extension/pkg/apparmor"
	"github.com/kommodity-io/kommodity-attestation-extension/pkg/extensions"
	"github.com/kommodity-io/kommodity-attestation-extension/pkg/image"
	"github.com/kommodity-io/kommodity-attestation-extension/pkg/lockdown"
	"github.com/kommodity-io/kommodity-attestation-extension/pkg/openapi/attestation/attestationmodels"
	"github.com/kommodity-io/kommodity-attestation-extension/pkg/secureboot"
	"github.com/kommodity-io/kommodity-attestation-extension/pkg/selinux"
	"github.com/kommodity-io/kommodity-attestation-extension/pkg/squashfs"
	"github.com/kommodity-io/kommodity-attestation-extension/pkg/tpm"
	"github.com/kommodity-io/kommodity-attestation-extension/pkg/utils"
	"github.com/kommodity-io/kommodity-attestation-extension/pkg/version"
)

const (
	//nolint:lll
	// Talos populates PCR index 7 with its measurements: https://github.com/siderolabs/talos/blob/main/pkg/machinery/constants/constants.go#L733.
	talosPCRIndex = 7
)

// Attestable defines a single component that can produce a signed attestation.
type Attestable interface {
	Name() string
	Measure() (string, error)             // hash or digest of the component
	Evidence() (map[string]string, error) // any supporting metadata (e.g. versions)
}

// AttestableReport represents a complete attestation report composed of multiple attestable components.
type AttestableReport struct {
	Attestables []Attestable
}

// NewAllAttestableReport creates a new AttestableReport with all available attestable components.
func NewAllAttestableReport() *AttestableReport {
	return &AttestableReport{
		Attestables: []Attestable{
			&apparmor.Attestable{},
			&extensions.Attestable{},
			&image.Attestable{},
			&lockdown.Attestable{},
			&secureboot.Attestable{},
			&selinux.Attestable{},
			&squashfs.Attestable{},
			&version.Attestable{},
		},
	}
}

// NewAttestableReport creates a new empty AttestableReport.
func NewAttestableReport() *AttestableReport {
	return &AttestableReport{
		Attestables: make([]Attestable, 0),
	}
}

// AddAttestable adds a new Attestable component to the report.
func (r *AttestableReport) AddAttestable(a Attestable) *AttestableReport {
	r.Attestables = append(r.Attestables, a)

	return r
}

// Generate generates the attestation report by collecting measurements, quotes, and evidence.
func (r *AttestableReport) Generate(nonce []byte) (*attestationmodels.RestReport, error) {
	components := make([]*attestationmodels.RestComponentReport, 0)

	tpmDevice, err := tpm.OpenTPMDevice(nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to get TPM device: %w", err)
	}

	pcrSelection := tpm.GetPCRSelection([]int{talosPCRIndex})

	pcrs, err := tpmDevice.ReadPCRs(pcrSelection)
	if err != nil {
		return nil, fmt.Errorf("failed to read PCRs from TPM device: %w", err)
	}

	quote, err := tpmDevice.Quote(pcrSelection)
	if err != nil {
		return nil, fmt.Errorf("failed to get quote from TPM device: %w", err)
	}

	signature, err := tpmDevice.Signature()
	if err != nil {
		return nil, fmt.Errorf("failed to get signature from TPM device: %w", err)
	}

	publicKey, err := tpmDevice.GetTPMPublicKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get TPM public key: %w", err)
	}

	for _, attestable := range r.Attestables {
		measure, err := attestable.Measure()
		if err != nil {
			return nil, fmt.Errorf("failed to measure %s: %w", attestable.Name(), err)
		}

		evidence, err := attestable.Evidence()
		if err != nil {
			return nil, fmt.Errorf("failed to get evidence for %s: %w", attestable.Name(), err)
		}

		components = append(components, &attestationmodels.RestComponentReport{
			Evidence:    evidence,
			Measurement: measure,
			Name:        attestable.Name(),
		})
	}

	convPCRs := make(map[string]string)
	for k, v := range pcrs {
		convPCRs[strconv.Itoa(k)] = v
	}

	return &attestationmodels.RestReport{
		Components:   components,
		Pcrs:         convPCRs,
		Quote:        hex.EncodeToString(quote),
		Signature:    hex.EncodeToString(signature),
		Timestamp:    utils.UnixNowString(),
		TpmPublicKey: hex.EncodeToString(publicKey),
	}, nil
}
