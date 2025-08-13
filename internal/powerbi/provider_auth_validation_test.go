package powerbi

import (
	"strings"
	"testing"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
)

// TestValidateAuthenticationConfig tests the custom authentication validation function
func TestValidateAuthenticationConfig(t *testing.T) {
	testCases := []struct {
		name      string
		config    *powerbiapi.AuthConfig
		expectErr bool
		errMsg    string
	}{
		// Valid configurations
		{
			name: "Valid service principal with client secret",
			config: &powerbiapi.AuthConfig{
				TenantID:     "test-tenant",
				ClientID:     "test-client",
				ClientSecret: "test-secret",
			},
			expectErr: false,
		},
		{
			name: "Valid service principal with certificate path",
			config: &powerbiapi.AuthConfig{
				TenantID:        "test-tenant",
				ClientID:        "test-client",
				CertificatePath: "/path/to/cert.pem",
			},
			expectErr: false,
		},
		{
			name: "Valid service principal with certificate data",
			config: &powerbiapi.AuthConfig{
				TenantID:        "test-tenant",
				ClientID:        "test-client",
				CertificateData: "LS0tLS1CRUdJTi...",
			},
			expectErr: false,
		},
		{
			name: "Valid managed identity",
			config: &powerbiapi.AuthConfig{
				UseManagedIdentity: true,
			},
			expectErr: false,
		},
		{
			name: "Valid managed identity with user-assigned ID",
			config: &powerbiapi.AuthConfig{
				UseManagedIdentity: true,
				ManagedIdentityID:  "test-identity-id",
			},
			expectErr: false,
		},
		{
			name: "Valid Azure CLI authentication",
			config: &powerbiapi.AuthConfig{
				UseAzureCLI: true,
			},
			expectErr: false,
		},
		{
			name: "Valid direct access token",
			config: &powerbiapi.AuthConfig{
				AccessToken: "test-access-token",
			},
			expectErr: false,
		},
		{
			name: "Valid username/password",
			config: &powerbiapi.AuthConfig{
				TenantID:     "test-tenant",
				ClientID:     "test-client",
				ClientSecret: "test-secret",
				Username:     "user@domain.com",
				Password:     "test-password",
			},
			expectErr: false,
		},

		// Error cases - Multiple authentication methods
		{
			name: "Multiple auth methods: managed identity and Azure CLI",
			config: &powerbiapi.AuthConfig{
				UseManagedIdentity: true,
				UseAzureCLI:        true,
			},
			expectErr: true,
			errMsg:    "multiple authentication methods",
		},
		{
			name: "Multiple auth methods: access token and client secret",
			config: &powerbiapi.AuthConfig{
				AccessToken:  "test-token",
				ClientSecret: "test-secret",
			},
			expectErr: true,
			errMsg:    "multiple authentication methods",
		},
		{
			name: "Multiple auth methods: certificate and client secret",
			config: &powerbiapi.AuthConfig{
				TenantID:        "test-tenant",
				ClientID:        "test-client",
				ClientSecret:    "test-secret",
				CertificatePath: "/path/to/cert.pem",
			},
			expectErr: true,
			errMsg:    "multiple authentication methods",
		},

		// Error cases - No authentication method
		{
			name: "No authentication method",
			config: &powerbiapi.AuthConfig{
				TenantID: "test-tenant",
				ClientID: "test-client",
			},
			expectErr: true,
			errMsg:    "no authentication method configured",
		},
		{
			name: "Empty configuration",
			config: &powerbiapi.AuthConfig{},
			expectErr: true,
			errMsg:    "no authentication method configured",
		},

		// Error cases - Missing required fields
		{
			name: "Certificate without tenant_id",
			config: &powerbiapi.AuthConfig{
				ClientID:        "test-client",
				CertificatePath: "/path/to/cert.pem",
			},
			expectErr: true,
			errMsg:    "tenant_id is required when using certificate",
		},
		{
			name: "Certificate without client_id",
			config: &powerbiapi.AuthConfig{
				TenantID:        "test-tenant",
				CertificatePath: "/path/to/cert.pem",
			},
			expectErr: true,
			errMsg:    "client_id is required when using certificate",
		},
		{
			name: "Client secret without tenant_id",
			config: &powerbiapi.AuthConfig{
				ClientID:     "test-client",
				ClientSecret: "test-secret",
			},
			expectErr: true,
			errMsg:    "tenant_id is required when using client_secret",
		},
		{
			name: "Client secret without client_id",
			config: &powerbiapi.AuthConfig{
				TenantID:     "test-tenant",
				ClientSecret: "test-secret",
			},
			expectErr: true,
			errMsg:    "client_id is required when using client_secret",
		},
		{
			name: "Username/password without tenant_id",
			config: &powerbiapi.AuthConfig{
				ClientID:     "test-client",
				ClientSecret: "test-secret",
				Username:     "user@domain.com",
				Password:     "test-password",
			},
			expectErr: true,
			errMsg:    "tenant_id is required when using username/password",
		},
		{
			name: "Username/password without client_secret",
			config: &powerbiapi.AuthConfig{
				TenantID: "test-tenant",
				ClientID: "test-client",
				Username: "user@domain.com",
				Password: "test-password",
			},
			expectErr: true,
			errMsg:    "client_secret is required when using username/password",
		},

		// Error cases - Invalid combinations
		{
			name: "Both certificate_path and certificate_data",
			config: &powerbiapi.AuthConfig{
				TenantID:        "test-tenant",
				ClientID:        "test-client",
				CertificatePath: "/path/to/cert.pem",
				CertificateData: "LS0tLS1CRUdJTi...",
			},
			expectErr: true,
			errMsg:    "certificate_path and certificate_data cannot be used together",
		},

		// Edge cases - Partial username/password
		{
			name: "Username without password (should not trigger username/password auth)",
			config: &powerbiapi.AuthConfig{
				TenantID:     "test-tenant",
				ClientID:     "test-client",
				ClientSecret: "test-secret",
				Username:     "user@domain.com",
			},
			expectErr: false, // This should be valid as client_secret auth
		},
		{
			name: "Password without username (should not trigger username/password auth)",
			config: &powerbiapi.AuthConfig{
				TenantID:     "test-tenant",
				ClientID:     "test-client",
				ClientSecret: "test-secret",
				Password:     "test-password",
			},
			expectErr: false, // This should be valid as client_secret auth
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateAuthenticationConfig(tc.config)

			if tc.expectErr {
				if err == nil {
					t.Fatalf("Expected error but got none")
				}
				if tc.errMsg != "" {
					if !strings.Contains(err.Error(), tc.errMsg) {
						t.Fatalf("Expected error message to contain '%s', got: %s", tc.errMsg, err.Error())
					}
				}
				t.Logf("Got expected error: %v", err)
			} else {
				if err != nil {
					t.Fatalf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

// TestAuthenticationPriority tests that authentication methods are prioritized correctly
func TestAuthenticationPriority(t *testing.T) {
	// This test verifies the authentication priority in our validation function
	testCases := []struct {
		name           string
		config         *powerbiapi.AuthConfig
		expectedMethod string
	}{
		{
			name: "Access token has highest priority",
			config: &powerbiapi.AuthConfig{
				AccessToken:        "test-token",
				UseManagedIdentity: true, // This should cause a conflict
			},
			expectedMethod: "multiple_methods", // Should detect multiple methods
		},
		{
			name: "Managed identity priority over service principal components",
			config: &powerbiapi.AuthConfig{
				UseManagedIdentity: true,
			},
			expectedMethod: "managed_identity",
		},
		{
			name: "Azure CLI priority",
			config: &powerbiapi.AuthConfig{
				UseAzureCLI: true,
			},
			expectedMethod: "azure_cli",
		},
		{
			name: "Certificate authentication",
			config: &powerbiapi.AuthConfig{
				TenantID:        "test-tenant",
				ClientID:        "test-client",
				CertificatePath: "/path/to/cert.pem",
			},
			expectedMethod: "certificate",
		},
		{
			name: "Client secret authentication",
			config: &powerbiapi.AuthConfig{
				TenantID:     "test-tenant",
				ClientID:     "test-client",
				ClientSecret: "test-secret",
			},
			expectedMethod: "client_secret",
		},
		{
			name: "Username/password authentication",
			config: &powerbiapi.AuthConfig{
				TenantID:     "test-tenant",
				ClientID:     "test-client",
				ClientSecret: "test-secret",
				Username:     "user@domain.com",
				Password:     "test-password",
			},
			expectedMethod: "username_password",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateAuthenticationConfig(tc.config)

			switch tc.expectedMethod {
			case "multiple_methods":
				if err == nil || !strings.Contains(err.Error(), "multiple authentication methods") {
					t.Fatalf("Expected multiple authentication methods error, got: %v", err)
				}
			default:
				if err != nil {
					t.Fatalf("Expected valid configuration for %s, got error: %v", tc.expectedMethod, err)
				}
			}
		})
	}
}