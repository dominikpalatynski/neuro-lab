package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
	"types"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Devices       []DeviceWrapper     `yaml:"devices" mapstructure:"devices"`
	CurrentDevice string              `yaml:"current_device,omitempty" mapstructure:"current_device"`
	APIEndpoint   string              `yaml:"api_endpoint,omitempty" mapstructure:"api_endpoint"`
	APIResources  []types.APIResource `yaml:"api_resources,omitempty" mapstructure:"api_resources"`
}

// DeviceWrapper wraps a device for YAML structure
type DeviceWrapper struct {
	Device DeviceInfo `yaml:"device" mapstructure:"device"`
}

// DeviceInfo contains device details
type DeviceInfo struct {
	Name     string `yaml:"name" mapstructure:"name"`
	DeviceID uint   `yaml:"device_id" mapstructure:"device_id"`
}

// DiscoveryCache represents cached API resource discovery data
type DiscoveryCache struct {
	CachedAt   time.Time           `json:"cached_at"`
	ExpiresAt  time.Time           `json:"expires_at"`
	APIVersion string              `json:"api_version"`
	Resources  []types.APIResource `json:"resources"`
}

const (
	// DefaultAPIEndpoint is the default API base URL
	DefaultAPIEndpoint = "http://localhost:3002/api/v1"
	// ConfigDirName is the config directory name in user's home
	ConfigDirName = ".neurolab"
	// ConfigFileName is the name of the config file
	ConfigFileName = "config"
	// CacheDirName is the cache directory name
	CacheDirName = "cache"
	// DiscoveryCacheFileName is the name of the discovery cache file
	DiscoveryCacheFileName = "api_resources.json"
	// DiscoveryCacheTTL is how long the cache is valid
	DiscoveryCacheTTL = 5 * time.Minute
)

// Initialize sets up the configuration system
// This should be called from the root command's PersistentPreRunE
func Initialize(cfgFile string) error {
	// 1. Set up Viper to use environment variables
	viper.SetEnvPrefix("NEUROLAB")
	// Allow for nested keys in environment variables (e.g. NEUROLAB_API_ENDPOINT)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	// 2. Set defaults
	viper.SetDefault("api_endpoint", DefaultAPIEndpoint)

	// 3. Handle the configuration file
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Search for a config file in default location
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}

		configDir := home + "/" + ConfigDirName
		// Ensure config directory exists
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}

		viper.AddConfigPath(configDir)
		viper.SetConfigName(ConfigFileName)
		viper.SetConfigType("yaml")
	}

	// 4. Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if the config file doesn't exist
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	return nil
}

// GetConfig returns the current configuration
func GetConfig() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return &cfg, nil
}

// SaveConfig writes the current configuration to disk
func SaveConfig(cfg *Config) error {
	// Set all values in Viper
	viper.Set("devices", cfg.Devices)
	viper.Set("current_device", cfg.CurrentDevice)
	viper.Set("api_endpoint", cfg.APIEndpoint)

	// Get the config file path
	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		// Config file doesn't exist yet, create it in default location
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		configFile = home + "/" + ConfigDirName + "/" + ConfigFileName + ".yaml"
	}

	// Write the config file
	if err := viper.WriteConfigAs(configFile); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetCurrentDevice returns the currently selected device
func GetCurrentDevice() (string, error) {
	currentDevice := viper.GetString("current_device")
	if currentDevice == "" {
		return "", fmt.Errorf("no device selected. Use 'init' to fetch devices and 'use' to select one")
	}
	return currentDevice, nil
}

// GetCurrentDeviceInfo returns the full device info for the current device
func GetCurrentDeviceInfo() (*DeviceInfo, error) {
	currentDeviceName, err := GetCurrentDevice()
	if err != nil {
		return nil, err
	}

	cfg, err := GetConfig()
	if err != nil {
		return nil, err
	}

	for _, dw := range cfg.Devices {
		if dw.Device.Name == currentDeviceName {
			return &dw.Device, nil
		}
	}

	return nil, fmt.Errorf("current device '%s' not found in config", currentDeviceName)
}

// SetCurrentDevice updates the currently selected device
func SetCurrentDevice(deviceName string) error {
	cfg, err := GetConfig()
	if err != nil {
		return err
	}

	// Verify the device exists
	found := false
	for _, dw := range cfg.Devices {
		if dw.Device.Name == deviceName {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("device '%s' not found in config", deviceName)
	}

	cfg.CurrentDevice = deviceName
	return SaveConfig(cfg)
}

// GetDevices returns all available devices
func GetDevices() ([]DeviceWrapper, error) {
	cfg, err := GetConfig()
	if err != nil {
		return nil, err
	}
	return cfg.Devices, nil
}

// SetDevices updates the list of devices
func SetDevices(devices []DeviceWrapper) error {
	cfg, err := GetConfig()
	if err != nil {
		// If config doesn't exist yet, create a new one
		cfg = &Config{
			APIEndpoint: DefaultAPIEndpoint,
		}
	}

	cfg.Devices = devices
	return SaveConfig(cfg)
}

// GetAPIEndpoint returns the API base URL
func GetAPIEndpoint() string {
	return viper.GetString("api_endpoint")
}

// SetAPIEndpoint updates the API base URL
func SetAPIEndpoint(endpoint string) error {
	cfg, err := GetConfig()
	if err != nil {
		cfg = &Config{}
	}

	cfg.APIEndpoint = endpoint
	return SaveConfig(cfg)
}

// ConfigFileUsed returns the config file being used
func ConfigFileUsed() string {
	return viper.ConfigFileUsed()
}

// GetCachePath returns the full path to the cache directory
func GetCachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return home + "/" + ConfigDirName + "/" + CacheDirName, nil
}

// GetDiscoveryCache loads the cached discovery data from disk
// Returns nil if cache doesn't exist or is invalid
func GetDiscoveryCache() (*DiscoveryCache, error) {
	cachePath, err := GetCachePath()
	if err != nil {
		return nil, err
	}

	cacheFile := cachePath + "/" + DiscoveryCacheFileName

	// Check if file exists
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		return nil, nil // No cache file, not an error
	}

	// Read cache file
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache file: %w", err)
	}

	// Parse JSON
	var cache DiscoveryCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, fmt.Errorf("failed to parse cache file: %w", err)
	}

	return &cache, nil
}

// SetDiscoveryCache saves discovery data to cache file
func SetDiscoveryCache(resources []types.APIResource) error {
	cachePath, err := GetCachePath()
	if err != nil {
		return err
	}

	// Ensure cache directory exists
	if err := os.MkdirAll(cachePath, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Create cache data with timestamps
	now := time.Now()
	cache := DiscoveryCache{
		CachedAt:   now,
		ExpiresAt:  now.Add(DiscoveryCacheTTL),
		APIVersion: "v1",
		Resources:  resources,
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache data: %w", err)
	}

	// Write to file
	cacheFile := cachePath + "/" + DiscoveryCacheFileName
	if err := os.WriteFile(cacheFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	return nil
}

// IsDiscoveryCacheValid checks if the cached discovery data is still valid
func IsDiscoveryCacheValid() bool {
	cache, err := GetDiscoveryCache()
	if err != nil || cache == nil {
		return false
	}

	// Check if expired
	return time.Now().Before(cache.ExpiresAt)
}

// InvalidateDiscoveryCache removes the discovery cache file
func InvalidateDiscoveryCache() error {
	cachePath, err := GetCachePath()
	if err != nil {
		return err
	}

	cacheFile := cachePath + "/" + DiscoveryCacheFileName

	// Remove file (ignore error if doesn't exist)
	if err := os.Remove(cacheFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove cache file: %w", err)
	}

	return nil
}

// FetchAndCacheDiscovery fetches API resources from the discovery endpoint and caches them
func FetchAndCacheDiscovery(apiEndpoint string) error {
	// Make request to discovery endpoint (apiEndpoint already includes /api/v1)
	discoveryURL := apiEndpoint

	resp, err := http.Get(discoveryURL)
	if err != nil {
		return fmt.Errorf("failed to fetch discovery data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("discovery endpoint returned status %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read discovery response: %w", err)
	}

	// Parse APIResourceList
	var resourceList types.APIResourceList
	if err := json.Unmarshal(body, &resourceList); err != nil {
		return fmt.Errorf("failed to parse discovery response: %w", err)
	}

	// Cache the resources
	if err := SetDiscoveryCache(resourceList.Resources); err != nil {
		return fmt.Errorf("failed to cache discovery data: %w", err)
	}

	return nil
}

// GetDiscoveryResources returns cached API resources or empty slice if not cached
func GetDiscoveryResources() ([]types.APIResource, error) {
	cache, err := GetDiscoveryCache()
	if err != nil {
		return nil, err
	}

	if cache == nil {
		return []types.APIResource{}, nil
	}

	return cache.Resources, nil
}
