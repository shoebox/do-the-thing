package signature

import (
	"bytes"
	"context"
	"dothething/internal/utiltest"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mozilla.org/pkcs7"
)

const validProvisioning = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>AppIDName</key>
	<string>dummy name</string>
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

const (
	validPath = "/path/to/file.provisioning"
)

var (
	subject provisioningService
)

func TestMain(m *testing.M) {
	//mockExec = new(utiltest.MockExecutor)
	//mockFs = new(utiltest.MockFileService)
	subject = provisioningService{}

	os.Exit(m.Run())
}

func getSignedReaderData(source string) io.ReadCloser {
	// Sign the valid provisioning datas
	d, _ := pkcs7.NewSignedData([]byte(source))

	// And retrieve the byytes
	b, _ := d.Finish()

	return ioutil.NopCloser(bytes.NewReader(b))
}

func TestDecode(t *testing.T) {
	// when:
	pp, err := subject.Decode(context.Background(), getSignedReaderData(validProvisioning))

	// then:
	assert.NoError(t, err)

	// and:
	assert.Equal(t, "Selfsigners united", pp.TeamName)
	assert.Equal(t, "12345ABCDE.*", pp.Entitlements.AppID)
	assert.Equal(t, "B5C2906D-D6EE-476E-AF17-D99AE14644AA", pp.UUID)
	assert.NoError(t, err)
}

func TestDecodeShouldHandleErrors(t *testing.T) {
	// when:
	pp, err := subject.Decode(context.Background(), strings.NewReader(""))

	// then:
	assert.EqualError(t, err, ErrorParsingPublicKey.Error())
	assert.Empty(t, pp)
}

func TestDecodeShouldHandleDecodingErrors(t *testing.T) {
	// when:
	_, err := subject.Decode(context.Background(), getSignedReaderData("invalid"))

	// then:
	assert.EqualError(t, err, "Failed to decode the provisioning file")
}

func TestDecodeCertShouldHandleDecodingErrors(t *testing.T) {
	// when:
	_, err := subject.Decode(context.Background(), strings.NewReader(invalidProvisioning))

	// then:
	assert.EqualError(t, err, "Failed to parse the provisioning file certificate")
}

func TestParseRawX509Certificates(t *testing.T) {
	t.Run("Should throw an error if failed to parse the provisiong public key", func(t *testing.T) {
		// setup:
		data := [][]byte{[]byte("Hello world")}

		// when: Parsing certificate datas
		res, err := parseRawX509Certificates(data)

		// then: An error should be returned
		assert.EqualValues(t, ErrorParsingPublicKey, err)
		assert.Nil(t, res)
	})
}

func TestFileDecodingProvisioningExecution(t *testing.T) {
	// setup:
	ctx := context.Background()

	// when:
	b, err := subject.decodeProvisioning(ctx, getSignedReaderData(validProvisioning))

	// then:
	assert.NoError(t, err)
	assert.EqualValues(t, validProvisioning, string(b))
}

func TestReadCert(t *testing.T) {
	// Making call to the target method
	pp, err := subject.decodeRawProvisioning(
		context.Background(),
		validPath,
		getSignedReaderData(validProvisioning),
	)

	// then: Asserting no errors
	assert.NoError(t, err)

	// and: The value in the channel should be populated
	assert.Equal(t, "Selfsigners united", pp.TeamName)
	assert.Equal(t, "12345ABCDE.*", pp.Entitlements.AppID)
	assert.Equal(t, "B5C2906D-D6EE-476E-AF17-D99AE14644AA", pp.UUID)
	assert.Equal(t, validPath, pp.FilePath)
}

func TestReadCertErrorHanding(t *testing.T) {
	// when:
	res, err := subject.decodeRawProvisioning(context.Background(),
		validPath,
		ioutil.NopCloser(strings.NewReader("")))

	// It should be empty
	assert.Empty(t, res)

	// and: A parsing error should have been raised
	assert.EqualError(t, err, "Failed to parse the provisioning file certificate")
}

func TestIsProvisioning(t *testing.T) {
	// setup:
	cases := []struct {
		fi    utiltest.MockFileInfo
		name  string
		valid bool
	}{
		{
			fi:    utiltest.NewMockFileInfo(0, false, ""),
			name:  "Invalid file mode",
			valid: false,
		},
		{
			fi:    utiltest.NewMockFileInfo(os.ModeAppend, true, "toto.mobileprovision"),
			name:  "Should not be a directory",
			valid: false,
		},
		{
			fi:    utiltest.NewMockFileInfo(os.ModeAppend, false, "toto.mobileprovision"),
			name:  "Valid mode",
			valid: true,
		},
		{
			fi:    utiltest.NewMockFileInfo(os.ModeAppend, false, "toto.prov"),
			name:  "Should have the right extension",
			valid: false,
		},
		{
			fi:    utiltest.NewMockFileInfo(os.ModeIrregular, false, "toto.mobileprovision"),
			name:  "Should be regular mode",
			valid: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// when:
			res := isProvisioningFile(c.fi)

			// then:
			assert.EqualValues(t, c.valid, res)
		})
	}
}
