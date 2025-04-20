package cmd

import (
	cmtcfg "github.com/cometbft/cometbft/config"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
)

// initCometBFTConfig helps to override default CometBFT Config values.
// return cmtcfg.DefaultConfig if no custom configuration is required for the application.
func initCometBFTConfig() *cmtcfg.Config {
	cfg := cmtcfg.DefaultConfig()

	// these values put a higher strain on node memory
	// cfg.P2P.MaxNumInboundPeers = 100
	// cfg.P2P.MaxNumOutboundPeers = 40

	return cfg
}

// initAppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func initAppConfig() (string, interface{}) {
	// The following code snippet is just for reference.

	// CustomConfig defines an arbitrary custom config to extend app.toml.
	// If you don't need it, you can remove it.
	// If you wish to add fields that correspond to flags that aren't in the SDK server config,
	// this custom config can as well help.
	type CustomConfig struct {
		CustomField string `mapstructure:"custom-field"`
	}

	type CustomAppConfig struct {
		serverconfig.Config `mapstructure:",squash"`

		Custom CustomConfig `mapstructure:"custom"`
	}

	// Optionally allow the chain developer to overwrite the SDK's default
	// server config.
	srvCfg := serverconfig.DefaultConfig()
	// The SDK's default minimum gas price is set to "" (empty value) inside
	// app.toml. If left empty by validators, the node will halt on startup.
	// However, the chain developer can set a default app.toml value for their
	// validators here.
	//
	// In summary:
	// - if you leave srvCfg.MinGasPrices = "", all validators MUST tweak their
	//   own app.toml config,
	// - if you set srvCfg.MinGasPrices non-empty, validators CAN tweak their
	//   own app.toml to override, or use this default value.
	//
	// In simapp, we set the min gas prices to 0.
	srvCfg.MinGasPrices = "0stake"
	// srvCfg.BaseConfig.IAVLDisableFastNode = true // disable fastnode by default

	// Now we set the custom config default values.
	customAppConfig := CustomAppConfig{
		Config: *srvCfg,
		Custom: CustomConfig{
			CustomField: "anything",
		},
	}

	// The default SDK app template is defined in serverconfig.DefaultConfigTemplate.
	// We append the custom config template to the default one.
	// And we set the default config to the custom app template.
	customAppTemplate := serverconfig.DefaultConfigTemplate + `
[custom]
# That field will be parsed by server.InterceptConfigsPreRunHandler and held by viper.
# Do not forget to add quotes around the value if it is a string.
custom-field = "{{ .Custom.CustomField }}"`

	return customAppTemplate, customAppConfig
}
