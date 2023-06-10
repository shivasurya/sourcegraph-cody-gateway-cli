package main

import (
	"cody-gateway-cli/app"
	"cody-gateway-cli/config"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/howeyc/gopass"
	"github.com/jessevdk/go-flags"
)

// options are command-line options that are provided by the user.
type options struct {
	Verbose          bool   `short:"V" long:"verbose" description:"Enable verbose output"`
	GatewayHost      string `long:"host" description:"Define an alternate SSH Port" default:"22"`
	GatewayToken     bool   `long:"accesstoken" description:"Ask for a password to use for cody gateway authentication"`
	DebugSecretToken string `short:"s" long:"debugtoken" description:"Define a bearer secret token to use" env:"SECRET_TOKEN"`
	GatewayMode      string `short:"m" long:"mode" description:"Define chat or code completion mode"`
}

func main() {
	var opts options
	args, err := flags.ParseArgs(&opts, os.Args[1:])
	// print args
	fmt.Println(args)
	if err != nil {
		os.Exit(1)
	}

	// Convert to internal config
	cfg := config.New()
	cfg.Verbose = opts.Verbose
	cfg.GatewayHost = opts.GatewayHost
	if opts.GatewayMode != "" {
		cfg.GatewayMode = opts.GatewayMode
	}
	if opts.DebugSecretToken != "" {
		cfg.DebugSecretToken = opts.DebugSecretToken
	}

	if opts.GatewayToken {
		color.White("Enter Cody Gateway Access Token for %s: ", cfg.GatewayToken)
		p, err := gopass.GetPasswd()
		if err != nil {
			color.Red("Unable to obtain Cody Gateway Access Token: %s", err)
		}
		cfg.GatewayToken = string(p)
	}

	// Run the App
	err = app.Run(cfg)
	if err != nil {
		if cfg.Verbose {
			color.Red("Error executing: %s", err)
		}
		os.Exit(1)
	}
	if cfg.Verbose {
		color.Green("Execution completed successfully")
	}
}