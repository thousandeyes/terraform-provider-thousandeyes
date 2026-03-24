package schemas

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestClientCertificateDiffSuppress_sameCertKeyPlusCertVsCertOnly(t *testing.T) {
	const certOnly = "-----BEGIN CERTIFICATE-----MIICUTCCAfugAwIBAgIBADANBgkqhkiG9w0BAQQFADBXMQswCQYDVQQGEwJDTjEL\n-----END CERTIFICATE-----\n"
	const withKey = "-----BEGIN PRIVATE KEY-----MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC7VJTUt9Us8cKB\n-----END PRIVATE KEY-----\n" + certOnly

	if !ClientCertificateDiffSuppress("client_certificate", certOnly, withKey, nil) {
		t.Fatal("expected suppress when certificate payload matches")
	}
	if !ClientCertificateDiffSuppress("client_certificate", withKey, certOnly, nil) {
		t.Fatal("expected suppress when arguments reversed")
	}
}

func TestClientCertificateDiffSuppress_sameLineBeginHeader(t *testing.T) {
	// No newline between -----BEGIN CERTIFICATE----- and base64 (API / user style).
	const certTight = "-----BEGIN CERTIFICATE-----MIICUTCCAfugAwIBAgIBADANBgkqhkiG9w0BAQQFADBXMQswCQYDVQQGEwJDTjEL\n-----END CERTIFICATE-----"
	const withKey = "-----BEGIN PRIVATE KEY-----abc\n-----END PRIVATE KEY-----\n" + certTight

	if !ClientCertificateDiffSuppress("client_certificate", certTight, withKey, nil) {
		t.Fatal("expected suppress for tight PEM header layout")
	}
}

func TestClientCertificateDiffSuppress_wrappedVsSingleLineBase64(t *testing.T) {
	const singleLine = "-----BEGIN CERTIFICATE-----\nMIICUTCCAfugAwIBAgIBADANBgkqhkiG9w0BAQQFADBXMQswCQYDVQQGEwJDTjEL\n-----END CERTIFICATE-----"
	const wrapped = "-----BEGIN CERTIFICATE-----\nMIICUTCCAfugAwIBAgIBADA\nNBgkqhkiG9w0BAQQFADBXMQswCQYDVQQGEwJDTjEL\n-----END CERTIFICATE-----"

	if !ClientCertificateDiffSuppress("client_certificate", singleLine, wrapped, nil) {
		t.Fatal("expected suppress when only PEM wrapping differs")
	}
}

func TestClientCertificateDiffSuppress_differentCerts(t *testing.T) {
	const a = "-----BEGIN CERTIFICATE-----\nQQ==\n-----END CERTIFICATE-----"
	const b = "-----BEGIN CERTIFICATE-----\nQg==\n-----END CERTIFICATE-----"

	if ClientCertificateDiffSuppress("client_certificate", a, b, nil) {
		t.Fatal("expected diff when certificate differs")
	}
}

func TestClientCertificateDiffSuppress_identicalStrings(t *testing.T) {
	s := "any"
	if !ClientCertificateDiffSuppress("client_certificate", s, s, &schema.ResourceData{}) {
		t.Fatal("expected suppress for identical strings")
	}
}

func TestClientCertificateDiffSuppress_bothEmpty(t *testing.T) {
	if !ClientCertificateDiffSuppress("client_certificate", "", "", nil) {
		t.Fatal("expected suppress when both empty")
	}
}
