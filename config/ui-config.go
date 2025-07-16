package config

import (
	"strings"
)

// UIConfig holds the user interface configuration
type UIConfig struct {
	UseModernUI     bool     `json:"use_modern_ui"`
	Theme           string   `json:"theme"` // "light", "dark", "auto"
	EnabledFeatures []string `json:"enabled_features"`
	BrandName       string   `json:"brand_name"`
	CustomLogo      string   `json:"custom_logo"`
}

// DefaultUIConfig returns the default UI configuration
func DefaultUIConfig() *UIConfig {
	return &UIConfig{
		UseModernUI:     true, // Default to modern UI for new installations
		Theme:           "light",
		EnabledFeatures: []string{"dashboard", "campaigns", "groups", "templates", "landing_pages", "sending_profiles"},
		BrandName:       "phi.shin",
		CustomLogo:      "/images/phishin_logo.svg",
	}
}

// GetTemplateName returns the appropriate template name based on UI configuration
func (ui *UIConfig) GetTemplateName(baseName string) string {
	if ui.UseModernUI {
		// Check if modern template exists, otherwise fallback to base
		modernName := baseName + "-modern"
		return modernName
	}
	return baseName
}

// IsFeatureEnabled checks if a specific feature is enabled
func (ui *UIConfig) IsFeatureEnabled(feature string) bool {
	for _, f := range ui.EnabledFeatures {
		if f == feature {
			return true
		}
	}
	return false
}

// GetBrandName returns the configured brand name
func (ui *UIConfig) GetBrandName() string {
	if ui.BrandName == "" {
		return "phi.shin"
	}
	return ui.BrandName
}

// GetLogoPath returns the path to the logo based on UI mode
func (ui *UIConfig) GetLogoPath(size string) string {
	if ui.CustomLogo != "" {
		return ui.CustomLogo
	}
	
	if ui.UseModernUI {
		switch size {
		case "small", "navbar":
			return "/images/phishin_logo_small.svg"
		default:
			return "/images/phishin_logo.svg"
		}
	}
	
	// Legacy logo paths
	switch size {
	case "small", "navbar":
		return "/images/logo_inv_small.png"
	default:
		return "/images/logo_purple.png"
	}
}

// GetCSSPaths returns the CSS files to load based on UI mode
func (ui *UIConfig) GetCSSPaths() []string {
	if ui.UseModernUI {
		return []string{
			"https://cdn.tailwindcss.com",
			"https://cdn.jsdelivr.net/npm/tw-elements/dist/css/tw-elements.min.css",
			"/css/phishin-modern.css",
		}
	}
	
	// Legacy CSS files
	return []string{
		"/css/dist/gophish.css",
	}
}

// GetJSPaths returns the JavaScript files to load based on UI mode
func (ui *UIConfig) GetJSPaths() []string {
	if ui.UseModernUI {
		return []string{
			"https://cdn.jsdelivr.net/npm/tw-elements/dist/js/tw-elements.umd.min.js",
			"/js/dist/vendor.min.js",
			"/js/dist/app/gophish.min.js",
		}
	}
	
	// Legacy JS files
	return []string{
		"/js/dist/vendor.min.js",
		"/js/dist/app/gophish.min.js",
	}
}

// GetNavigationItems returns the navigation structure based on enabled features
func (ui *UIConfig) GetNavigationItems(userRole string) []NavigationItem {
	items := []NavigationItem{
		{
			Name: "Dashboard",
			Icon: "fas fa-chart-line",
			Path: "/",
			Show: ui.IsFeatureEnabled("dashboard"),
		},
		{
			Name: "Campaigns",
			Icon: "fas fa-paper-plane",
			Path: "/campaigns",
			Show: ui.IsFeatureEnabled("campaigns"),
		},
		{
			Name: "Users & Groups",
			Icon: "fas fa-users",
			Path: "/groups",
			Show: ui.IsFeatureEnabled("groups"),
		},
		{
			Name: "Email Templates",
			Icon: "fas fa-envelope",
			Path: "/templates",
			Show: ui.IsFeatureEnabled("templates"),
		},
		{
			Name: "Landing Pages",
			Icon: "fas fa-globe",
			Path: "/landing_pages",
			Show: ui.IsFeatureEnabled("landing_pages"),
		},
		{
			Name: "Sending Profiles",
			Icon: "fas fa-server",
			Path: "/sending_profiles",
			Show: ui.IsFeatureEnabled("sending_profiles"),
		},
	}

	// Add advanced features section
	advancedItems := []NavigationItem{
		{
			Name:    "E-Learning",
			Icon:    "fas fa-graduation-cap",
			Path:    "/courses",
			Show:    ui.IsFeatureEnabled("e_learning"),
			Badge:   "New",
			Section: "Advanced",
		},
		{
			Name:    "Multi-Channel",
			Icon:    "fas fa-broadcast-tower",
			Path:    "/channels",
			Show:    ui.IsFeatureEnabled("multi_channel"),
			Badge:   "Pro",
			Section: "Advanced",
		},
		{
			Name:    "AI Generator",
			Icon:    "fas fa-robot",
			Path:    "/ai-generator",
			Show:    ui.IsFeatureEnabled("ai_generator"),
			Badge:   "AI",
			Section: "Advanced",
		},
	}

	// Add settings section
	settingsItems := []NavigationItem{
		{
			Name:    "Account Settings",
			Icon:    "fas fa-cog",
			Path:    "/settings",
			Show:    true,
			Section: "Settings",
		},
	}

	// Add admin items if user has admin role
	if strings.Contains(userRole, "admin") {
		adminItems := []NavigationItem{
			{
				Name:    "User Management",
				Icon:    "fas fa-user-shield",
				Path:    "/users",
				Show:    ui.IsFeatureEnabled("user_management"),
				Badge:   "Admin",
				Section: "Settings",
			},
			{
				Name:    "Webhooks",
				Icon:    "fas fa-webhook",
				Path:    "/webhooks",
				Show:    ui.IsFeatureEnabled("webhooks"),
				Badge:   "Admin",
				Section: "Settings",
			},
			{
				Name:    "License",
				Icon:    "fas fa-key",
				Path:    "/license",
				Show:    ui.IsFeatureEnabled("license_management"),
				Badge:   "Admin",
				Section: "Settings",
			},
		}
		settingsItems = append(settingsItems, adminItems...)
	}

	// Combine all items
	allItems := append(items, advancedItems...)
	allItems = append(allItems, settingsItems...)

	return allItems
}

// NavigationItem represents a single navigation menu item
type NavigationItem struct {
	Name    string `json:"name"`
	Icon    string `json:"icon"`
	Path    string `json:"path"`
	Show    bool   `json:"show"`
	Badge   string `json:"badge,omitempty"`
	Section string `json:"section,omitempty"`
}

// FeatureConfig holds feature-specific configuration
type FeatureConfig struct {
	ELearning    ELearningConfig    `json:"e_learning"`
	MultiChannel MultiChannelConfig `json:"multi_channel"`
	AIGenerator  AIGeneratorConfig  `json:"ai_generator"`
	License      LicenseConfig      `json:"license"`
}

// ELearningConfig holds e-learning platform configuration
type ELearningConfig struct {
	Enabled              bool `json:"enabled"`
	AutoEnrollment       bool `json:"auto_enrollment"`
	DefaultDeadlineDays  int  `json:"default_deadline_days"`
	CertificateGeneration bool `json:"certificate_generation"`
}

// MultiChannelConfig holds multi-channel configuration
type MultiChannelConfig struct {
	Enabled         bool     `json:"enabled"`
	EnabledChannels []string `json:"enabled_channels"` // "sms", "voice", "whatsapp", "telegram", "apk"
	DefaultChannel  string   `json:"default_channel"`
}

// AIGeneratorConfig holds AI generator configuration
type AIGeneratorConfig struct {
	Enabled    bool   `json:"enabled"`
	Provider   string `json:"provider"` // "openai", "anthropic", "local"
	Model      string `json:"model"`
	MaxTokens  int    `json:"max_tokens"`
	APIKey     string `json:"api_key"`
}

// LicenseConfig holds license management configuration
type LicenseConfig struct {
	Type                  string   `json:"type"` // "community", "professional", "enterprise"
	MaxUsers              int      `json:"max_users"`
	MaxCampaignsPerMonth  int      `json:"max_campaigns_per_month"`
	MaxTargetsPerCampaign int      `json:"max_targets_per_campaign"`
	EnabledFeatures       []string `json:"enabled_features"`
	ExpiresAt             string   `json:"expires_at"`
}

// DefaultFeatureConfig returns default feature configuration
func DefaultFeatureConfig() *FeatureConfig {
	return &FeatureConfig{
		ELearning: ELearningConfig{
			Enabled:               false, // Requires license
			AutoEnrollment:        true,
			DefaultDeadlineDays:   7,
			CertificateGeneration: true,
		},
		MultiChannel: MultiChannelConfig{
			Enabled:         false, // Requires license
			EnabledChannels: []string{"sms", "voice"},
			DefaultChannel:  "email",
		},
		AIGenerator: AIGeneratorConfig{
			Enabled:   false, // Requires license and API key
			Provider:  "openai",
			Model:     "gpt-3.5-turbo",
			MaxTokens: 2000,
		},
		License: LicenseConfig{
			Type:                  "community",
			MaxUsers:              1,
			MaxCampaignsPerMonth:  5,
			MaxTargetsPerCampaign: 50,
			EnabledFeatures:       []string{"dashboard", "campaigns", "groups", "templates", "landing_pages", "sending_profiles"},
		},
	}
}

// GetPageTitle returns the appropriate page title based on branding
func (ui *UIConfig) GetPageTitle(pageTitle string) string {
	brandName := ui.GetBrandName()
	if pageTitle == "" {
		return brandName
	}
	return pageTitle + " - " + brandName
}

// GetThemeClasses returns CSS classes for the current theme
func (ui *UIConfig) GetThemeClasses() string {
	switch ui.Theme {
	case "dark":
		return "dark"
	case "auto":
		return "auto-theme"
	default:
		return "light"
	}
}

// IsLegacyMode returns true if using legacy UI
func (ui *UIConfig) IsLegacyMode() bool {
	return !ui.UseModernUI
}

// GetFaviconPath returns the path to the favicon
func (ui *UIConfig) GetFaviconPath() string {
	if ui.UseModernUI {
		return "/images/favicon.ico" // TODO: Create phi.shin favicon
	}
	return "/images/favicon.ico"
}