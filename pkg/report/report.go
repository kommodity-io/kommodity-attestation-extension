// Package report provides utilities for generating reports from all the attestation data.
package report

import (
	"encoding/hex"
	"fmt"

	"github.com/kommodity-io/kommodity-attestation-extension/pkg/apparmor"
	"github.com/kommodity-io/kommodity-attestation-extension/pkg/extensions"
	"github.com/kommodity-io/kommodity-attestation-extension/pkg/image"
	"github.com/kommodity-io/kommodity-attestation-extension/pkg/lockdown"
	"github.com/kommodity-io/kommodity-attestation-extension/pkg/openapi/attestation/attestationmodels"
	"github.com/kommodity-io/kommodity-attestation-extension/pkg/secureboot"
	"github.com/kommodity-io/kommodity-attestation-extension/pkg/selinux"
	"github.com/kommodity-io/kommodity-attestation-extension/pkg/squashfs"
	"github.com/kommodity-io/kommodity-attestation-extension/pkg/utils"
	"github.com/kommodity-io/kommodity-attestation-extension/pkg/version"
)

// Attestable defines a single component that can produce a signed attestation.
type Attestable interface {
	Name() string
	Measure() (string, error)             // hash or digest of the component
	GetPCRs() (map[int]string, error)     // PCR indices relevant to this component
	Quote(nonce []byte) ([]byte, error)   // TPM quote of relevant PCRs + nonce
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
func (r *AttestableReport) Generate(nonce []byte) (*attestationmodels.ReportReport, error) {
	components := make([]*attestationmodels.ReportComponentReport, 0)

	for _, attestable := range r.Attestables {
		measure, err := attestable.Measure()
		if err != nil {
			return nil, fmt.Errorf("failed to measure %s: %w", attestable.Name(), err)
		}

		quote, err := attestable.Quote(nonce)
		if err != nil {
			return nil, fmt.Errorf("failed to quote %s: %w", attestable.Name(), err)
		}

		evidence, err := attestable.Evidence()
		if err != nil {
			return nil, fmt.Errorf("failed to get evidence for %s: %w", attestable.Name(), err)
		}

		pcrs, err := attestable.GetPCRs()
		if err != nil {
			return nil, fmt.Errorf("failed to get PCRs for %s: %w", attestable.Name(), err)
		}

		components = append(components, &attestationmodels.ReportComponentReport{
			Evidence:    evidence,
			Measurement: measure,
			Pcrs:        pcrs,
			Name:        attestable.Name(),
			Quote:       hex.EncodeToString(quote),
			Signature:   hex.EncodeToString(quote),
		})
	}

	return &attestationmodels.ReportReport{
		Components: components,
		Timestamp:  utils.UnixNowString(),
	}, nil
}
