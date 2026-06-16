package aws

import (
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

type awsConfig struct {
	Regions               []string `hcl:"regions,optional"`
	DefaultRegion         *string  `hcl:"default_region"`
	Profile               *string  `hcl:"profile"`
	AccessKey             *string  `hcl:"access_key"`
	SecretKey             *string  `hcl:"secret_key"`
	SessionToken          *string  `hcl:"session_token"`
	MaxErrorRetryAttempts *int     `hcl:"max_error_retry_attempts"`
	MinErrorRetryDelay    *int     `hcl:"min_error_retry_delay"`
	IgnoreErrorMessages   []string `hcl:"ignore_error_messages,optional"`
	IgnoreErrorCodes      []string `hcl:"ignore_error_codes,optional"`
	EndpointUrl           *string  `hcl:"endpoint_url"`
	S3ForcePathStyle      *bool    `hcl:"s3_force_path_style"`

	// Inline IAM AssumeRole. Upstream the plugin can only assume a role via a
	// named `profile` resolved from ~/.aws/config; this fork lets the role be
	// passed in the connection config so the source process needs no shared
	// config file (see service.go getBaseClientForAccountUncached). The base
	// credentials that perform the assume come from whatever else resolves into
	// the config (IMDS, env, or a `profile`).
	RoleArn         *string `hcl:"role_arn"`
	ExternalId      *string `hcl:"external_id"`
	RoleSessionName *string `hcl:"role_session_name"`
}

func ConfigInstance() interface{} {
	return &awsConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) awsConfig {
	if connection == nil {
		return awsConfig{}
	}
	raw := connection.GetConfig()
	if raw == nil {
		return awsConfig{}
	}
	config, _ := raw.(awsConfig)

	if config.Regions != nil {
		if len(config.Regions) == 0 {
			// Setting "regions = []" in the connection config is not valid
			errorMessage := fmt.Sprintf("connection %s has invalid value for \"regions\", it must contain at least 1 region.", connection.Name)
			panic(errorMessage)
		}

		for i, r := range config.Regions {
			config.Regions[i] = NormalizeRegion(r)
		}
	}

	return config
}

func NormalizeRegion(region string) string {
	// ensure regions are lower case, to work consistently in matching
	// and comparisons
	return strings.ToLower(region)
}
