package xcconfig

import (
	"errors"
	"fmt"
)

const (
	CodeSignAllowed         = "CODE_SIGNING_ALLOWED"
	CodeSignIdentity        = "CODE_SIGN_IDENTITY"
	CodeSignEntitlements    = "CODE_SIGN_ENTITLEMENTS"
	DevelopmentTeam         = "DEVELOPMENT_TEAM"
	ProvisioningProfileId   = "PROVISIONING_PROFILE"
	ProvisioningProfileSpec = "PROVISIONING_PROFILE_SPECIFIER"
)

type SDKConfig struct {
	Name    string
	Version string
}

type EntryConfig struct {
	SDK     SDKConfig
	ARCH    string
	Config  string
	Variant string
}

type entry struct {
	Config EntryConfig
	Key    string
	Value  string
}

type Helper interface {
	Add(key string, value string, config EntryConfig) (bool, error)
	Generate() string
}

type helper struct {
	entries map[string]entry
}

func NewHelper() Helper {
	return helper{entries: map[string]entry{}}
}

func (h helper) Add(key string, value string, cfg EntryConfig) (bool, error) {
	// Do we key is already used ?
	if _, ok := h.entries[key]; ok {
		return false, errors.New(fmt.Sprintf("Key '%v' is already being used", key))
	}

	// Add entry
	h.entries[key] = entry{
		Config: cfg,
		Key:    key,
		Value:  value,
	}

	return true, nil

}

func (h helper) Generate() string {
	var res string
	for k, v := range h.entries {
		line := fmt.Sprintf("%v%v%v%v=%v", k,
			h.formatKeyValue("arch", v.Config.ARCH),
			h.formatKeyValue("config", v.Config.Config),
			h.formatSDKEntry(v.Config.SDK),
			v.Value)

		if res == "" {
			res = line
		} else {
			res = res + "\n" + line
		}
	}

	return res
}

func (h helper) formatKeyValue(key string, value string) string {
	var res string
	if value != "" {
		res = fmt.Sprintf("[%v=%v]", key, value)
	}

	return res
}

func (h helper) formatSDKEntry(c SDKConfig) string {
	var res string
	if c.Name != "" {
		res = fmt.Sprintf("[sdk=%v%v]", c.Name, c.Version)

	}
	return res
}
