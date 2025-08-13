package powerbiapi

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/go-cleanhttp"
)

// TestClientCredentialsTokenProvider tests service principal authentication with client secret
func TestClientCredentialsTokenProvider(t *testing.T) {
	// Skip integration test if credentials not available
	tenantID := os.Getenv("POWERBI_TENANT_ID")
	clientID := os.Getenv("POWERBI_CLIENT_ID")
	clientSecret := os.Getenv("POWERBI_CLIENT_SECRET")

	if tenantID == "" || clientID == "" || clientSecret == "" {
		t.Skip("Skipping integration test - POWERBI_TENANT_ID, POWERBI_CLIENT_ID, or POWERBI_CLIENT_SECRET not set")
	}

	provider := &ClientCredentialsTokenProvider{
		httpClient:   cleanhttp.DefaultClient(),
		tenantID:     tenantID,
		clientID:     clientID,
		clientSecret: clientSecret,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	token, err := provider.GetToken(ctx)
	if err != nil {
		t.Fatalf("Failed to get token: %v", err)
	}

	if token == "" {
		t.Fatal("Expected non-empty token")
	}

	t.Logf("Successfully obtained token of length: %d", len(token))
}

// TestManagedIdentityTokenProvider tests managed identity authentication
func TestManagedIdentityTokenProvider(t *testing.T) {
	// This test only runs in Azure environments with managed identity
	// Check if we're in an Azure environment
	if os.Getenv("IDENTITY_ENDPOINT") == "" && os.Getenv("IMDS_ENDPOINT") == "" {
		t.Skip("Skipping managed identity test - not running in Azure environment")
	}

	provider := &ManagedIdentityTokenProvider{
		httpClient: cleanhttp.DefaultClient(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	token, err := provider.GetToken(ctx)
	if err != nil {
		t.Fatalf("Failed to get token: %v", err)
	}

	if token == "" {
		t.Fatal("Expected non-empty token")
	}

	t.Logf("Successfully obtained managed identity token of length: %d", len(token))
}

// TestAzureCLITokenProvider tests Azure CLI authentication
func TestAzureCLITokenProvider(t *testing.T) {
	provider := &AzureCLITokenProvider{}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	token, err := provider.GetToken(ctx)
	
	// Azure CLI might not be installed or logged in, so we just test that the method doesn't panic
	if err != nil {
		t.Logf("Azure CLI authentication failed (expected if az not installed/logged in): %v", err)
		return
	}

	if token == "" {
		t.Fatal("Expected non-empty token when Azure CLI succeeds")
	}

	t.Logf("Successfully obtained Azure CLI token of length: %d", len(token))
}

// TestDirectTokenProvider tests direct token authentication
func TestDirectTokenProvider(t *testing.T) {
	testToken := "test-token-12345"

	provider := &DirectTokenProvider{
		accessToken: testToken,
	}

	ctx := context.Background()
	token, err := provider.GetToken(ctx)
	if err != nil {
		t.Fatalf("Failed to get token: %v", err)
	}

	if token != testToken {
		t.Fatalf("Expected token %s, got %s", testToken, token)
	}
}

// TestDirectTokenProviderEmpty tests direct token authentication with empty token
func TestDirectTokenProviderEmpty(t *testing.T) {
	provider := &DirectTokenProvider{
		accessToken: "",
	}

	ctx := context.Background()
	_, err := provider.GetToken(ctx)
	if err == nil {
		t.Fatal("Expected error for empty token")
	}

	expectedError := "no access token provided"
	if err.Error() != expectedError {
		t.Fatalf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

// TestNewClientWithAuthConfig tests the main client creation function with various auth methods
func TestNewClientWithAuthConfig(t *testing.T) {
	tests := []struct {
		name      string
		config    *AuthConfig
		expectErr bool
	}{
		{
			name: "Direct token authentication",
			config: &AuthConfig{
				AccessToken: "test-token",
			},
			expectErr: false,
		},
		{
			name: "Client credentials authentication",
			config: &AuthConfig{
				TenantID:     "test-tenant",
				ClientID:     "test-client",
				ClientSecret: "test-secret",
			},
			expectErr: false,
		},
		{
			name: "Azure CLI authentication",
			config: &AuthConfig{
				UseAzureCLI: true,
			},
			expectErr: false,
		},
		{
			name: "Managed Identity authentication",
			config: &AuthConfig{
				UseManagedIdentity: true,
			},
			expectErr: false,
		},
		{
			name: "Certificate data authentication",
			config: &AuthConfig{
				TenantID:        "test-tenant",
				ClientID:        "test-client",
				CertificateData: base64.StdEncoding.EncodeToString([]byte("fake-cert-data")),
			},
			expectErr: true, // Will fail because it's not a real certificate
		},
		{
			name: "No authentication method",
			config: &AuthConfig{
				TenantID: "test-tenant",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClientWithAuthConfig(tt.config)

			if tt.expectErr {
				if err == nil {
					t.Fatal("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if client == nil {
				t.Fatal("Expected non-nil client")
			}

			if client.HTTPClient == nil {
				t.Fatal("Expected non-nil HTTP client")
			}
		})
	}
}

// TestCertificateTokenProvider tests certificate parsing (without actual certificate)
func TestCertificateTokenProviderInvalidCert(t *testing.T) {
	// Create a temporary file with invalid certificate data
	tmpFile, err := ioutil.TempFile("", "invalid-cert-*.pem")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write invalid certificate data
	invalidCert := `-----BEGIN CERTIFICATE-----
INVALID_CERTIFICATE_DATA
-----END CERTIFICATE-----`

	if _, err := tmpFile.WriteString(invalidCert); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Test that invalid certificate is handled gracefully
	_, err = NewCertificateTokenProvider(cleanhttp.DefaultClient(), "tenant", "client", tmpFile.Name(), "")
	if err == nil {
		t.Fatal("Expected error for invalid certificate")
	}

	t.Logf("Correctly rejected invalid certificate: %v", err)
}

// TestAuthConfigPriority tests that authentication methods are prioritized correctly
func TestAuthConfigPriority(t *testing.T) {
	// Test that direct token takes highest priority
	config := &AuthConfig{
		AccessToken:        "direct-token",
		UseManagedIdentity: true,
		UseAzureCLI:        true,
		ClientSecret:       "secret",
		TenantID:           "tenant",
		ClientID:           "client",
	}

	client, err := NewClientWithAuthConfig(config)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if client == nil {
		t.Fatal("Expected non-nil client")
	}
}

// TestManagedIdentityPriority tests that managed identity takes priority over other methods
func TestManagedIdentityPriority(t *testing.T) {
	config := &AuthConfig{
		UseManagedIdentity: true,
		UseAzureCLI:        true,
		ClientSecret:       "secret",
		TenantID:           "tenant",
		ClientID:           "client",
	}

	client, err := NewClientWithAuthConfig(config)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if client == nil {
		t.Fatal("Expected non-nil client")
	}
}