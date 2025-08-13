package powerbi

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}


// TestProviderSchemaValidation tests that the provider schema is valid
func TestProviderSchemaValidation(t *testing.T) {
	provider := Provider()
	
	// Test that the provider schema is valid
	if err := provider.InternalValidate(); err != nil {
		t.Fatalf("Provider schema validation failed: %v", err)
	}
	
	// Test that required fields are properly defined
	schema := provider.Schema
	
	// Check that authentication-related fields exist
	requiredFields := []string{
		"tenant_id", "client_id", "client_secret", "username", "password",
		"certificate_path", "certificate_data", "certificate_password",
		"use_managed_identity", "managed_identity_id", "use_azure_cli", "access_token",
	}
	
	for _, field := range requiredFields {
		if _, exists := schema[field]; !exists {
			t.Fatalf("Required field '%s' not found in provider schema", field)
		}
	}
}

// TestProviderConfigureFunction tests that the provider configure function works
func TestProviderConfigureFunction(t *testing.T) {
	provider := Provider()
	
	// Test with empty configuration - should fail with validation error
	resourceData := schema.TestResourceDataRaw(t, provider.Schema, map[string]interface{}{})
	
	_, err := provider.ConfigureFunc(resourceData)
	if err == nil {
		t.Fatal("Expected error for empty configuration, but got none")
	}
	
	// Test with valid managed identity configuration
	resourceData = schema.TestResourceDataRaw(t, provider.Schema, map[string]interface{}{
		"use_managed_identity": true,
	})
	
	// This will fail because we're not in an Azure environment, but it should pass validation
	_, err = provider.ConfigureFunc(resourceData)
	// We expect an error here because we're not actually in Azure, but it should be an API error, not a validation error
	if err != nil {
		// Make sure it's not a validation error
		if err.Error() == "no authentication method configured. Please configure one of: access_token, managed_identity, azure_cli, certificate, client_secret, or username/password" {
			t.Fatalf("Unexpected validation error: %v", err)
		}
	}
}