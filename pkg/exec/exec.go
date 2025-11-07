// Package exec provides the execution environment for Kommodity Attestation Extension.
package exec

import (
	"fmt"
	"net"

	"github.com/kommodity-io/kommodity-attestation-extension/pkg/openapi/attestation/attestationclient"
	"github.com/kommodity-io/kommodity-attestation-extension/pkg/openapi/attestation/attestationclient/attestation"
	"github.com/kommodity-io/kommodity-attestation-extension/pkg/openapi/attestation/attestationmodels"
	"github.com/kommodity-io/kommodity-attestation-extension/pkg/report"
	"github.com/kommodity-io/kommodity-attestation-extension/pkg/uuid"
)

const (
	cmdArgAttestationServer = "kommmodity.attestation.server"
)

// Execute runs the attestation process based on the provided command-line arguments.
func Execute(args map[string]string) error {
	server, ok := args[cmdArgAttestationServer]
	if !ok || server == "" {
		return fmt.Errorf("%w: argument=%s", ErrArgMissing, cmdArgAttestationServer)
	}

	client := attestationclient.NewHTTPClientWithConfig(nil,
		attestationclient.DefaultTransportConfig().WithHost(server),
	)

	nonce, err := client.Attestation.GetNonce(attestation.NewGetNonceParams())
	if err != nil {
		return fmt.Errorf("failed to get nonce: %w", err)
	}

	report := report.NewAllAttestableReport()

	responseReport, err := report.Generate([]byte(nonce.Payload.Nonce))
	if err != nil {
		return fmt.Errorf("failed to generate report: %w", err)
	}

	uuid, err := uuid.GetMachineUUID()
	if err != nil {
		return fmt.Errorf("failed to get machine UUID: %w", err)
	}

	publicIP, err := getPublicIP()
	if err != nil {
		return fmt.Errorf("failed to get public IP: %w", err)
	}

	_, err = client.Attestation.PostReport(attestation.NewPostReportParams().
		WithPayload(&attestationmodels.ReportAttestationReportRequest{
			Nonce:  nonce.Payload.Nonce,
			Report: responseReport,
			Node: &attestationmodels.ReportNodeInfo{
				UUID: uuid,
				IP:   publicIP,
			},
		}))
	if err != nil {
		return fmt.Errorf("failed to submit report: %w", err)
	}

	return nil
}

func getPublicIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("failed to get network interfaces: %w", err)
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok || ipNet.IP.IsLoopback() || ipNet.IP.IsLinkLocalUnicast() {
			continue
		}

		//nolint:varnamelen
		ip := ipNet.IP.To4()
		if ip == nil {
			continue // skip IPv6
		}

		// Skip private ranges
		if isPrivateIP(ip) {
			continue
		}

		return ip.String(), nil
	}

	return "", nil
}

func isPrivateIP(ip net.IP) bool {
	// Skip private ranges
	if ip[0] == 10 ||
		(ip[0] == 172 && ip[1] >= 16 && ip[1] <= 31) ||
		(ip[0] == 192 && ip[1] == 168) {
		return true
	}

	return false
}
