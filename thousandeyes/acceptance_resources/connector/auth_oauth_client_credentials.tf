resource "thousandeyes_connector" "test_auth" {
  name   = "UAT - Connector Auth OAuth Client Creds"
  target = "https://example.com/webhooks/thousandeyes/auth-oauth-cc"

  authentication {
    type                = "oauth-client-credentials"
    oauth_client_id     = "test-client-id"
    oauth_client_secret = "test-client-secret"
    oauth_token_url     = "https://auth.example.com/oauth/token"
  }
}
