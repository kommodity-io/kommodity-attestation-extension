// Package main implements the Kommodity Attestation Extension command-line tool.
package main

import (
	"github.com/kommodity-io/kommodity-attestation-extension/pkg/cmdline"
	"github.com/kommodity-io/kommodity-attestation-extension/pkg/exec"
	"github.com/kommodity-io/kommodity/pkg/logging"
	"go.uber.org/zap"
)

func main() {
	logger := logging.NewLogger()
	zap.ReplaceGlobals(logger)

	logger.Info("Kommodity Attestation Extension")

	args, err := cmdline.ParseProcCmdline()
	if err != nil {
		logger.Error("Error parsing command line", zap.Error(err))

		return
	}

	logger.Info("Parsed command line arguments", zap.Any("args", args))

	err = exec.Execute(args)
	if err != nil {
		logger.Error("Error executing attestation", zap.Error(err))

		return
	}

	logger.Info("Attestation execution completed successfully")
}
