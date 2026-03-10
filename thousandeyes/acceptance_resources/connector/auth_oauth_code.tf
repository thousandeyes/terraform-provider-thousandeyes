resource "thousandeyes_connector" "test_auth" {
  name   = "UAT - Connector Auth OAuth Code"
  target = "https://example.com/webhooks/thousandeyes/auth-oauth-code"

  authentication {
    type                = "oauth-auth-code"
    oauth_client_id     = "test-client-id"
    oauth_client_secret = "test-client-secret"
    oauth_token_url     = "https://auth.example.com/oauth/token"
    oauth_auth_url      = "https://auth.example.com/oauth/authorize"
    code                = "test-auth-code"
    redirect_uri        = "https://app.example.com/callback"
  }
}
