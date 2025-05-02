package main

import (
	"fmt"
	"os"

	"github.com/frocore/fedramp-data-mesh/cli/cmd"
	"github.com/frocore/fedramp-data-mesh/cli/internal/config"
	"github.com/frocore/fedramp-data-mesh/cli/internal/logging"
	"github.com/frocore/fedramp-data-mesh/cli/internal/security"
)

func main() {
	// Initialize logger
	log := logging.NewLogger()
	
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Errorf("Failed to load configuration: %v", err)
		os.Exit(1)
	}
	
	// Initialize security context
	secCtx, err := security.NewSecurityContext(cfg)
	if err != nil {
		log.Errorf("Failed to initialize security context: %v", err)
		os.Exit(1)
	}
	
	// Execute root command
	if err := cmd.Execute(cfg, secCtx, log); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
