package openapi

//go:generate go run github.com/go-swagger/go-swagger/cmd/swagger@v0.33.1 generate client -f https://raw.githubusercontent.com/kommodity-io/kommodity/refs/heads/main/openapi/attestation/swagger.yaml -A attestation -t ./attestation -c attestationclient -m attestationmodels
