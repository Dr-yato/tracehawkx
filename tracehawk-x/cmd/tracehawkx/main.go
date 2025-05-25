package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/tracehawk/tracehawkx/internal/core"
	"github.com/tracehawk/tracehawkx/internal/core/config"
	"github.com/tracehawk/tracehawkx/modules"
	_ "github.com/tracehawk/tracehawkx/modules/bleeding-edge"
	_ "github.com/tracehawk/tracehawkx/modules/stable"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

const asciiArt = `                                                                                
                                       ::                                       
                                      -::-                                      
                              :::---- --:: ::---:-                              
                          ::---::-   ::::::   -:::::::                          
                       :::::        .:---::.        :::::                       
                     --:-           :-:  ::-           :::-                     
                   :::-            :-:    :::            ::::                   
                 ::--              :-:    ::-              ::::                 
                -::               -::      -::               :::                
               -::               -:-        -:-               ::-               
                                 :-:                                            
         -::-::::-:-:--::::-:::.:-: :::-:--:::::::::-::::::-:::::-:-:::         
           ---::=               -::          ::::              :::::-           
        -:-:  -:-  :-:  -:: -:-   -:: ::-:::-: :::::::  ::-::-.  ::::::-        
         :-   :::  ::   -:.  --:  .:- :  ::-    :-      ::-  :-- .  :-          
         ::   ---  -:.  :-.  :-::  :-    :::    :-      :::  :--   -:           
         :::-:-:-  -:.  ::.  :: -: :-    ::-    ::---   :::::--      -::        
         ::   ::-  ::.  -:.  -- .-:-:    ::-    ::      :-- ::       :::        
         --   -::  :--  -:   :-  .:::    ::-    -:   :: ::- -:-  ::  :--        
        -:::  :-:   --:::   :::-  -::   -:::   :::--:::---: --::  ::-:          
            .::            -::   -::-      ::-:   --:            ::             
             :::          -::.     ---: -::::     .:-:          :::             
              ::-         ::-        -::--         :::         ::-              
               :-:       :::       --::  :::-       :-:       -::               
                :::      :-:    -:-:-      ::::-    -:-     .:-:                
                 -:::   :::  -----            :::-.  :::   :---                 
                   :::-::: :--:                 .:::: :::-::-                   
                     - :-::::                      :::::: :                     
                      -::-:-:                     :: :::-                      
                     ::-   :::-::::-        -::-::::-   :::                     
                                ::----:--::-::::                                
                                                                                `

func main() {
	// Set up logging
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	logrus.SetOutput(os.Stderr)

	// Create root command
	rootCmd := &cobra.Command{
		Use:   "tracehawkx",
		Short: "TraceHawk X - Autonomous Terminal-Native Recon & Exploitation Agent",
		Long: color.HiCyanString(asciiArt) + "\n\n" +
			color.HiMagentaString("🦅 TraceHawk X - Autonomous Terminal-Native Recon & Exploitation Agent") + "\n\n" +
			color.HiWhiteString("A comprehensive penetration testing toolkit that performs end-to-end security\nassessments from enumeration through exploitation to auto-patch suggestions.\n\n") +
			color.HiGreenString("Features:") + "\n" +
			color.WhiteString("• Fully offline, no external APIs required\n") +
			color.WhiteString("• Modular plugin architecture\n") +
			color.WhiteString("• Beautiful HTML/Markdown reports\n") +
			color.WhiteString("• Advanced bleeding-edge modules\n") +
			color.WhiteString("• Network namespace sandboxing\n\n") +
			color.HiYellowString("🌐 Visit: ") + color.HiCyanString("https://hunter3.ninja") + "\n",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Global flags
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default: $HOME/.tracehawkx.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "debug output")
	rootCmd.PersistentFlags().Bool("no-color", false, "disable colored output")

	// Add subcommands
	rootCmd.AddCommand(
		newScanCommand(),
		newListModulesCommand(),
		newVersionCommand(),
		newInitCommand(),
		newTUICommand(),
	)

	// Handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logrus.Info("Received interrupt signal, shutting down gracefully...")
		cancel()
	}()

	// Execute command
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		logrus.Errorf("Command failed: %v", err)
		os.Exit(1)
	}
}

func newScanCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scan",
		Short: "Perform a security scan against target(s)",
		Long:  "Execute a comprehensive security scan including reconnaissance, vulnerability detection, and exploitation",
		RunE:  runScan,
	}

	// Target specification
	cmd.Flags().StringSliceP("target", "t", nil, "target hosts, domains, or CIDR ranges")
	cmd.Flags().StringSliceP("target-file", "T", nil, "file containing targets (one per line)")
	cmd.Flags().StringSlice("exclude", nil, "hosts/ranges to exclude from scan")

	// Scan modes
	cmd.Flags().Bool("deep", false, "enable deep scanning (slower but thorough)")
	cmd.Flags().Bool("bleeding-edge", false, "enable bleeding-edge modules")
	cmd.Flags().Bool("stealth", false, "use stealth scanning techniques")
	cmd.Flags().Bool("aggressive", false, "aggressive scanning (may be detected)")

	// Output options
	cmd.Flags().StringP("output", "o", "results", "output directory")
	cmd.Flags().StringP("report", "r", "reports", "report output directory")
	cmd.Flags().String("format", "all", "output format: json,markdown,html,all")
	cmd.Flags().Bool("live-report", false, "enable live report updates")

	// Performance tuning
	cmd.Flags().IntP("threads", "j", 50, "number of concurrent threads")
	cmd.Flags().IntP("rate-limit", "l", 150, "requests per second limit")
	cmd.Flags().Duration("timeout", 0, "overall scan timeout")
	cmd.Flags().Bool("no-throttle", false, "disable request throttling")

	// Module-specific flags
	cmd.Flags().String("llm-model", "", "path to LLM model file (.gguf)")
	cmd.Flags().Float64("temp", 0.7, "LLM temperature")
	cmd.Flags().Bool("generate-patch", false, "generate auto-patch recommendations")
	cmd.Flags().String("shadow-clone", "", "shadow clone proxy (format: port:9009)")
	cmd.Flags().Bool("dep-drift", false, "scan for dependency drift")
	cmd.Flags().Bool("timing-map", false, "perform side-channel timing analysis")
	cmd.Flags().Bool("blue-team", false, "generate blue team replay kit")

	cmd.MarkFlagRequired("target")
	return cmd
}

func newListModulesCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list-modules",
		Short: "List all available modules",
		Long:  "Display information about all available scanning modules",
		RunE: func(cmd *cobra.Command, args []string) error {
			modules.ListModules()
			return nil
		},
	}
}

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("%s\n", color.HiCyanString(asciiArt))
			fmt.Printf("%s\n", color.HiMagentaString("🦅 TraceHawk X - Autonomous Terminal-Native Recon & Exploitation Agent"))
			fmt.Printf("%s %s\n", color.HiGreenString("Version:"), color.WhiteString(version))
			fmt.Printf("%s %s\n", color.HiGreenString("Built:"), color.WhiteString(date))
			fmt.Printf("%s %s\n", color.HiGreenString("Commit:"), color.WhiteString(commit))
			fmt.Printf("%s %s\n", color.HiYellowString("Website:"), color.HiCyanString("https://hunter3.ninja"))
			fmt.Printf("\n%s\n", color.HiWhiteString("Made with ❤️ by the Hunter3.Ninja Team"))
			return nil
		},
	}
}

func newInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize TraceHawk X configuration and templates",
		Long:  "Set up default configuration files and download/update scanning templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			return config.Initialize()
		},
	}
}

func newTUICommand() *cobra.Command {
	return &cobra.Command{
		Use:   "tui",
		Short: "Launch interactive Terminal UI",
		Long:  "Start the interactive terminal interface for live scanning sessions",
		RunE: func(cmd *cobra.Command, args []string) error {
			return core.LaunchTUI()
		},
	}
}

func runScan(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load(cmd)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create scan orchestrator
	orchestrator, err := core.NewOrchestrator(cfg)
	if err != nil {
		return fmt.Errorf("failed to create orchestrator: %w", err)
	}

	// Execute scan
	return orchestrator.Execute(cmd.Context())
}
