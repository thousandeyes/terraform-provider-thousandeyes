package schemas

import (
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var clientCertificatePEMBlocks = regexp.MustCompile(`-----BEGIN CERTIFICATE-----(?:\r?\n)?([\s\S]*?)-----END CERTIFICATE-----`)

// normalizeCertBodyForCompare trims the PEM body, then removes all whitespace so
// wrapped base64 (64-column lines) matches a single-line encoding.
func normalizeCertBodyForCompare(body string) string {
	body = strings.TrimSpace(body)
	if body == "" {
		return ""
	}

	var b strings.Builder
	b.Grow(len(body))
	for _, r := range body {
		switch r {
		case '\n', '\r', ' ', '\t':
			continue
		default:
			b.WriteRune(r)
		}
	}

	return b.String()
}

func normalizeCertificatePayload(pemBundle string) string {
	matches := clientCertificatePEMBlocks.FindAllStringSubmatch(pemBundle, -1)
	if len(matches) == 0 {
		return ""
	}

	var b strings.Builder
	for _, m := range matches {
		norm := normalizeCertBodyForCompare(m[1])
		if norm == "" {
			continue
		}
		b.WriteString(norm)
	}

	return b.String()
}

// ClientCertificateDiffSuppress returns true when the certificate PEM payload
// (between BEGIN/END CERTIFICATE) matches after trimming each body and
// collapsing whitespace, so wrapped and single-line base64 compare equal.
// Private keys and other PEM types are ignored.
func ClientCertificateDiffSuppress(_, oldVal, newVal string, _ *schema.ResourceData) bool {
	if oldVal == newVal {
		return true
	}

	oldNorm := normalizeCertificatePayload(oldVal)
	newNorm := normalizeCertificatePayload(newVal)
	if oldNorm == "" && newNorm == "" {
		return strings.TrimSpace(oldVal) == strings.TrimSpace(newVal)
	}
	return oldNorm == newNorm
}
