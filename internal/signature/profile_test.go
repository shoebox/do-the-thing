package signature

import (
	"context"
	"dothething/internal/utiltest"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const validProvisioning = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>AppIDName</key>
	<string>abc.def.ghi</string>
	<key>ApplicationIdentifierPrefix</key>
	<array>
	<string>12345ABCDE</string>
	</array>
	<key>CreationDate</key>
	<date>2018-01-01T13:00:00Z</date>
	<key>Platform</key>
	<array>
		<string>iOS</string>
	</array>
	<key>DeveloperCertificates</key>
	<array>
<data>MIIDaTCCAlGgAwIBAgIBATANBgkqhkiG9w0BAQsFADBdMTkwNwYDVQQDDDBpUGhvbmUgRGlzdHJpYnV0aW9uOiBEdW1teSBOYW1lIEx0ZCAoMTIzNDVBQkNERSkxEzAR BgNVBAsMClNFTEZTSUdORUQxCzAJBgNVBAYTAkdCMB4XDTE4MTAwMTIwMzcxNVoX DTI4MDkyODIwMzcxNVowXTE5MDcGA1UEAwwwaVBob25lIERpc3RyaWJ1dGlvbjog RHVtbXkgTmFtZSBMdGQgKDEyMzQ1QUJDREUpMRMwEQYDVQQLDApTRUxGU0lHTkVE MQswCQYDVQQGEwJHQjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAN5K m0F/jI7ryjErnjUxni0wIKZmLAnyRea9TXotCwzMk/WjuqQlQF2M/li6V8cc+ftP c0bCegmNbqOX2T9x86oBsbRQ5e81lFf/LtU6XVXEnc7DRJZ6Jgvx62PVEHZbY2eN 0Gn/sK2al4WD6hpf2k/olySyjNSsmyj/G712+OVZO6vfxEFyACYdS+3g+mSp5OpE IkV96Ze1y4RH9oXzluoeC0nEl36otoi/geG4w9XWsVK48Uz27JRqJSXmn8TAS+BW twl2P+pEOtDp2l9pf85lmeRPvNXLPGHEwrvUf2Hr6mLBwJQU1bpFhVCQxGBMq3ib I4ltf/EwKj4Li9qaAckCAwEAAaM0MDIwDgYDVR0PAQH/BAQDAgeAMCAGA1UdJQEB /wQWMBQGCCsGAQUFBwMEBggrBgEFBQcDAzANBgkqhkiG9w0BAQsFAAOCAQEASDBU 4AayZG5RPhXH8pyHN5AVm6upCsXhC8I4NI+lpO0hBBJ1svoUsWJeqgDncNEj/a0T wWG5NBqPU/EtT3hZcatLegy7X2cSIuvqYe2ZPrqfQaQuyyA3T/OMJOCElNAhsGsC VlwsnABbdZSgr9ts7+b/kWSXqQV5rOzkdRAHvNeIJ/kIBS08GQhOVGkAA3u7f4wp xXtSGMXGtd16VcaWFAhtXbQdBFCCfcz9pNWx61K0H585Ei8YqrJ7yMdtwDZD9QDH UtYddAvFAxv5jZ9N25PYE4LFjgwqDbK7+5rriqtGN7YuCCVKqwYUKtKIJusotcXF XKpC6uZNmz1/bvowpw==</data>
	</array>
	<key>Entitlements</key>
	<dict>
		<key>keychain-access-groups</key>
		<array>
			<string>12345ABCDE.*</string>
		</array>
		<key>get-task-allow</key>
		<false/>
		<key>application-identifier</key>
		<string>12345ABCDE.*</string>
		<key>com.apple.developer.associated-domains</key>
		<string>*</string>
		<key>com.apple.developer.team-identifier</key>
		<string>12345ABCDE</string>
		<key>aps-environment</key>
		<string>production</string>
	</dict>
	<key>ExpirationDate</key>
	<date>2028-01-01T14:00:00Z</date>
	<key>Name</key>
	<string>DO NOT USE: only for dummy signing</string>
	<key>ProvisionedDevices</key>
	<array>
	</array>
	<key>TeamIdentifier</key>
	<array>
		<string>12345ABCDE</string>
	</array>
	<key>TeamName</key>
	<string>Selfsigners united</string>
	<key>TimeToLive</key>
	<integer>3652</integer>
	<key>UUID</key>
	<string>B5C2906D-D6EE-476E-AF17-D99AE14644AA</string>
	<key>Version</key>
	<integer>1</integer>
</dict>
</plist>`

const invalidProvisioning = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>AppIDName</key>
	<string>abc.def.ghi</string>
	<key>ApplicationIdentifierPrefix</key>
	<array>
	<string>12345ABCDE</string>
	</array>
	<key>CreationDate</key>
	<date>2018-01-01T13:00:00Z</date>
	<key>Platform</key>
	<array>
		<string>iOS</string>
	</array>
	<key>DeveloperCertificates</key>
	<array>
		<data></data>
	</array>
	<key>Entitlements</key>
	<dict>
		<key>keychain-access-groups</key>
		<array>
			<string>12345ABCDE.*</string>
		</array>
		<key>get-task-allow</key>
		<false/>
		<key>application-identifier</key>
		<string>12345ABCDE.*</string>
		<key>com.apple.developer.associated-domains</key>
		<string>*</string>
		<key>com.apple.developer.team-identifier</key>
		<string>12345ABCDE</string>
		<key>aps-environment</key>
		<string>production</string>
	</dict>
	<key>ExpirationDate</key>
	<date>2028-01-01T14:00:00Z</date>
	<key>Name</key>
	<string>DO NOT USE: only for dummy signing</string>
	<key>ProvisionedDevices</key>
	<array>
	</array>
	<key>TeamIdentifier</key>
	<array>
		<string>12345ABCDE</string>
	</array>
	<key>TeamName</key>
	<string>Selfsigners united</string>
	<key>TimeToLive</key>
	<integer>3652</integer>
	<key>UUID</key>
	<string>B5C2906D-D6EE-476E-AF17-D99AE14644AA</string>
	<key>Version</key>
	<integer>1</integer>
</dict>
</plist>`

const validPath = "/path/to/file.provisioning"
const invalidPath = "/path/to/file2.provisioning"

var mockExec *utiltest.MockExecutor
var subject provisioningService

func TestMain(m *testing.M) {
	mockExec = new(utiltest.MockExecutor)
	subject = provisioningService{mockExec}

	mockExec.MockCommandContext(Security,
		[]string{Cms, ArgDecodeCMS, ArgInlineFile, validPath},
		validProvisioning,
		nil)

	os.Exit(m.Run())
}

func TestDecode(t *testing.T) {
	// when:
	pp, err := subject.Decode(context.Background(), validPath)

	// then:
	assert.NoError(t, err)
	assert.Equal(t, "abc.def.ghi", pp.AppID)
	assert.Equal(t, []string([]string{"12345ABCDE"}), pp.TeamIdentifier)
	assert.Equal(t, "Selfsigners united", pp.TeamName)
	assert.Equal(t, "B5C2906D-D6EE-476E-AF17-D99AE14644AA", pp.UUID)
	assert.NoError(t, err)
}

func TestDecodeShouldHandleErrors(t *testing.T) {
	// setup:
	mockExec.MockCommandContext(Security,
		[]string{Cms, ArgDecodeCMS, ArgInlineFile, invalidPath},
		"",
		errors.New("error text"))

	// when:
	pp, err := subject.Decode(context.Background(), invalidPath)

	// then:
	assert.EqualError(t, err, "error text")
	assert.Empty(t, pp)
}

func TestDecodeShouldHandleDecodingErrors(t *testing.T) {
	// setup:
	mockExec.MockCommandContext(Security,
		[]string{Cms, ArgDecodeCMS, ArgInlineFile, "/fake/path/fake-file.mobileprovision"},
		`{"json":"text"}`,
		nil)

	// when:
	pp, err := subject.Decode(context.Background(), "/fake/path/fake-file.mobileprovision")

	// then:
	assert.EqualError(t, err, "Failed to decode the provisioning file")
	assert.Empty(t, pp)
}

func TestDecodeCertShouldHandleDecodingErrors(t *testing.T) {
	// setup:
	mockExec.MockCommandContext(Security,
		[]string{Cms, ArgDecodeCMS, ArgInlineFile, "/toto/tutu/fake.mobileprovision"},
		invalidProvisioning,
		nil)

	// when:
	pp, err := subject.Decode(context.Background(), "/toto/tutu/fake.mobileprovision")

	// then:
	assert.EqualError(t, err, "Failed to parse the provisioning file certificate")

	// and:
	assert.Equal(t, "abc.def.ghi", pp.AppID)
	assert.Equal(t, []string([]string{"12345ABCDE"}), pp.TeamIdentifier)
	assert.Equal(t, "Selfsigners united", pp.TeamName)
	assert.Equal(t, "B5C2906D-D6EE-476E-AF17-D99AE14644AA", pp.UUID)
}

func TestFileDecodingProvisioningExecutation(t *testing.T) {
	// setup:
	ctx := context.Background()

	// when:
	b, err := subject.decodeProvisioning(ctx, validPath)

	// then:
	assert.NoError(t, err)
	assert.EqualValues(t, validProvisioning, string(b))

	// and:
	mockExec.AssertExpectations(t)
}

func TestParseRawX509Certificates(t *testing.T) {
	t.Run("Should throw an error if failed to parse the provisiong public key", func(t *testing.T) {
		// setup:
		data := [][]byte{[]byte("Hello world")}

		// when:
		res, err := parseRawX509Certificates(data)

		// then:
		assert.EqualValues(t, ErrorParsingPublicKey, err)
		assert.Nil(t, res)
	})
}
