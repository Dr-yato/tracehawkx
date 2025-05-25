package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tracehawk/tracehawkx/modules"
)

const (
	DefaultConfigName = ".tracehawkx"
	DefaultConfigType = "yaml"
)

// Load loads configuration from file and command line flags
func Load(cmd *cobra.Command) (*modules.ScanConfig, error) {
	// Initialize viper
	viper.SetConfigName(DefaultConfigName)
	viper.SetConfigType(DefaultConfigType)
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")

	// Set default values
	setDefaults()

	// Read config file if it exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Bind command line flags
	if err := bindFlags(cmd); err != nil {
		return nil, fmt.Errorf("failed to bind flags: %w", err)
	}

	// Create config struct
	config := &modules.ScanConfig{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Override with command line flags
	if err := overrideWithFlags(cmd, config); err != nil {
		return nil, fmt.Errorf("failed to override with flags: %w", err)
	}

	// Set up logging level
	setupLogging(cmd)

	return config, nil
}

// Initialize sets up default configuration and templates
func Initialize() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, DefaultConfigName+"."+DefaultConfigType)

	// Create default config if it doesn't exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := createDefaultConfig(configPath); err != nil {
			return fmt.Errorf("failed to create default config: %w", err)
		}
		logrus.Infof("Created default configuration at %s", configPath)
	}

	// Initialize templates and signatures
	if err := initializeTemplates(); err != nil {
		logrus.Warnf("Failed to initialize templates: %v", err)
	}

	logrus.Info("TraceHawk X initialization completed")
	return nil
}

func setDefaults() {
	viper.SetDefault("output", "results")
	viper.SetDefault("report_dir", "reports")
	viper.SetDefault("threads", 50)
	viper.SetDefault("rate_limit", 150)
	viper.SetDefault("deep", false)
	viper.SetDefault("bleeding_edge", false)
	viper.SetDefault("stealth", false)
	viper.SetDefault("aggressive", false)
	viper.SetDefault("no_throttle", false)
	viper.SetDefault("temperature", 0.7)
	viper.SetDefault("generate_patch", false)
	viper.SetDefault("dep_drift", false)
	viper.SetDefault("timing_map", false)
	viper.SetDefault("blue_team", false)
}

func bindFlags(cmd *cobra.Command) error {
	flags := []string{
		"target", "target-file", "exclude", "output", "report", "format",
		"threads", "rate-limit", "timeout", "deep", "bleeding-edge",
		"stealth", "aggressive", "no-throttle", "llm-model", "temp",
		"generate-patch", "shadow-clone", "dep-drift", "timing-map", "blue-team",
	}

	for _, flag := range flags {
		if cmd.Flags().Lookup(flag) != nil {
			key := strings.ReplaceAll(flag, "-", "_")
			if err := viper.BindPFlag(key, cmd.Flags().Lookup(flag)); err != nil {
				return err
			}
		}
	}

	return nil
}

func overrideWithFlags(cmd *cobra.Command, config *modules.ScanConfig) error {
	// Handle targets
	if targets, err := cmd.Flags().GetStringSlice("target"); err == nil && len(targets) > 0 {
		config.Targets = append(config.Targets, targets...)
	}

	// Handle target files
	if targetFiles, err := cmd.Flags().GetStringSlice("target-file"); err == nil && len(targetFiles) > 0 {
		for _, file := range targetFiles {
			targets, err := loadTargetsFromFile(file)
			if err != nil {
				logrus.Warnf("Failed to load targets from file %s: %v", file, err)
				continue
			}
			config.Targets = append(config.Targets, targets...)
		}
	}

	// Handle exclude list
	if exclude, err := cmd.Flags().GetStringSlice("exclude"); err == nil {
		config.Exclude = exclude
	}

	// Handle output directories
	if output, err := cmd.Flags().GetString("output"); err == nil && output != "" {
		config.Output = output
	}
	if reportDir, err := cmd.Flags().GetString("report"); err == nil && reportDir != "" {
		config.ReportDir = reportDir
	}

	// Handle scan mode flags
	if deep, err := cmd.Flags().GetBool("deep"); err == nil {
		config.Deep = deep
	}
	if bleedingEdge, err := cmd.Flags().GetBool("bleeding-edge"); err == nil {
		config.BleedingEdge = bleedingEdge
	}
	if stealth, err := cmd.Flags().GetBool("stealth"); err == nil {
		config.Stealth = stealth
	}
	if aggressive, err := cmd.Flags().GetBool("aggressive"); err == nil {
		config.Aggressive = aggressive
	}
	if noThrottle, err := cmd.Flags().GetBool("no-throttle"); err == nil {
		config.NoThrottle = noThrottle
	}

	// Handle module-specific flags
	if llmModel, err := cmd.Flags().GetString("llm-model"); err == nil && llmModel != "" {
		config.LLMModel = llmModel
	}
	if temp, err := cmd.Flags().GetFloat64("temp"); err == nil {
		config.Temperature = temp
	}
	if generatePatch, err := cmd.Flags().GetBool("generate-patch"); err == nil {
		config.GeneratePatch = generatePatch
	}
	if shadowClone, err := cmd.Flags().GetString("shadow-clone"); err == nil && shadowClone != "" {
		config.ShadowClone = shadowClone
	}
	if depDrift, err := cmd.Flags().GetBool("dep-drift"); err == nil {
		config.DepDrift = depDrift
	}
	if timingMap, err := cmd.Flags().GetBool("timing-map"); err == nil {
		config.TimingMap = timingMap
	}
	if blueTeam, err := cmd.Flags().GetBool("blue-team"); err == nil {
		config.BlueTeam = blueTeam
	}

	return nil
}

func setupLogging(cmd *cobra.Command) {
	if debug, _ := cmd.Flags().GetBool("debug"); debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(logrus.WarnLevel)
	}

	if noColor, _ := cmd.Flags().GetBool("no-color"); noColor {
		logrus.SetFormatter(&logrus.TextFormatter{DisableColors: true})
	}
}

func createDefaultConfig(path string) error {
	defaultConfig := `# TraceHawk X Configuration
output: results
report_dir: reports
threads: 50
rate_limit: 150
deep: false
bleeding_edge: false
stealth: false
aggressive: false
no_throttle: false

# Module-specific settings
temperature: 0.7
generate_patch: false
dep_drift: false
timing_map: false
blue_team: false

# Exclude patterns (optional)
exclude: []

# Global timeouts (optional)
# timeout: 30m
`

	return os.WriteFile(path, []byte(defaultConfig), 0644)
}

func loadTargetsFromFile(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var targets []string
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			targets = append(targets, line)
		}
	}

	return targets, nil
}

func initializeTemplates() error {
	// Create template directories
	templateDirs := []string{
		"assets/templates/nuclei",
		"assets/templates/nmap",
		"assets/templates/reports",
		"assets/signatures",
	}

	for _, dir := range templateDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// TODO: Download/update template files from official sources
	logrus.Info("Template directories created")
	return nil
}
