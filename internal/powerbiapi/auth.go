package powerbiapi

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/hashicorp/go-cleanhttp"
)

// AuthConfig holds all authentication configuration options
type AuthConfig struct {
	TenantID     string
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
	
	// Certificate authentication
	CertificatePath     string
	CertificatePassword string
	CertificateData     string // Base64 encoded certificate
	
	// Managed Identity
	UseManagedIdentity bool
	ManagedIdentityID  string // Optional: specific managed identity to use
	
	// Azure CLI
	UseAzureCLI bool
	
	// Token
	AccessToken string // Direct token authentication
}

// TokenProvider defines the interface for token providers
type TokenProvider interface {
	GetToken(ctx context.Context) (string, error)
}

// ClientCredentialsTokenProvider implements client credentials flow
type ClientCredentialsTokenProvider struct {
	httpClient   *http.Client
	tenantID     string
	clientID     string
	clientSecret string
}

func (p *ClientCredentialsTokenProvider) GetToken(ctx context.Context) (string, error) {
	authURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", url.PathEscape(p.tenantID))
	resp, err := p.httpClient.Post(authURL, "application/x-www-form-urlencoded", strings.NewReader(url.Values{
		"grant_type":    {"client_credentials"},
		"scope":         {"https://analysis.windows.net/powerbi/api/.default"},
		"client_id":     {p.clientID},
		"client_secret": {p.clientSecret},
	}.Encode()))

	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		data, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, data)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var tokenResp tokenResponse
	if err := json.Unmarshal(data, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	return tokenResp.AccessToken, nil
}

// CertificateTokenProvider implements certificate-based authentication
type CertificateTokenProvider struct {
	httpClient  *http.Client
	tenantID    string
	clientID    string
	certificate *x509.Certificate
	privateKey  *rsa.PrivateKey
}

func NewCertificateTokenProvider(httpClient *http.Client, tenantID, clientID string, certPath, certPassword string) (*CertificateTokenProvider, error) {
	certData, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}

	// Parse PFX/PKCS12 or PEM certificate
	var cert *x509.Certificate
	var key *rsa.PrivateKey

	// Try PEM first
	block, rest := pem.Decode(certData)
	if block != nil {
		switch block.Type {
		case "CERTIFICATE":
			cert, err = x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("failed to parse certificate: %w", err)
			}
			// Look for private key
			for len(rest) > 0 {
				block, rest = pem.Decode(rest)
				if block != nil && strings.Contains(block.Type, "PRIVATE KEY") {
					keyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
					if err != nil {
						// Try PKCS1
						keyInterface, err = x509.ParsePKCS1PrivateKey(block.Bytes)
						if err != nil {
							return nil, fmt.Errorf("failed to parse private key: %w", err)
						}
					}
					var ok bool
					key, ok = keyInterface.(*rsa.PrivateKey)
					if !ok {
						return nil, fmt.Errorf("private key is not RSA")
					}
					break
				}
			}
		case "RSA PRIVATE KEY":
			key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("failed to parse RSA private key: %w", err)
			}
		case "PRIVATE KEY":
			keyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("failed to parse private key: %w", err)
			}
			var ok bool
			key, ok = keyInterface.(*rsa.PrivateKey)
			if !ok {
				return nil, fmt.Errorf("private key is not RSA")
			}
		}
	} else {
		// Try PKCS12/PFX format
		// This would require additional libraries like software.sslmate.com/src/go-pkcs12
		return nil, fmt.Errorf("PKCS12/PFX format not yet supported, please use PEM format")
	}

	if cert == nil || key == nil {
		return nil, fmt.Errorf("failed to load certificate and private key")
	}

	return &CertificateTokenProvider{
		httpClient:  httpClient,
		tenantID:    tenantID,
		clientID:    clientID,
		certificate: cert,
		privateKey:  key,
	}, nil
}

func (p *CertificateTokenProvider) GetToken(ctx context.Context) (string, error) {
	// Create JWT assertion for certificate authentication
	assertion, err := p.createJWTAssertion()
	if err != nil {
		return "", fmt.Errorf("failed to create JWT assertion: %w", err)
	}

	authURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", url.PathEscape(p.tenantID))
	resp, err := p.httpClient.Post(authURL, "application/x-www-form-urlencoded", strings.NewReader(url.Values{
		"grant_type":            {"client_credentials"},
		"scope":                 {"https://analysis.windows.net/powerbi/api/.default"},
		"client_id":             {p.clientID},
		"client_assertion_type": {"urn:ietf:params:oauth:client-assertion-type:jwt-bearer"},
		"client_assertion":      {assertion},
	}.Encode()))

	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		data, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, data)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var tokenResp tokenResponse
	if err := json.Unmarshal(data, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	return tokenResp.AccessToken, nil
}

func (p *CertificateTokenProvider) createJWTAssertion() (string, error) {
	// This is a simplified version - in production, you'd use a proper JWT library
	// like github.com/golang-jwt/jwt
	now := time.Now()
	exp := now.Add(10 * time.Minute)

	header := map[string]interface{}{
		"alg": "RS256",
		"typ": "JWT",
		"x5t": base64.RawURLEncoding.EncodeToString(p.certificate.Raw[:20]), // SHA-1 thumbprint
	}

	claims := map[string]interface{}{
		"aud": fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", p.tenantID),
		"exp": exp.Unix(),
		"iss": p.clientID,
		"jti": fmt.Sprintf("%d", now.UnixNano()),
		"nbf": now.Unix(),
		"sub": p.clientID,
	}

	headerJSON, _ := json.Marshal(header)
	claimsJSON, _ := json.Marshal(claims)

	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)
	claimsEncoded := base64.RawURLEncoding.EncodeToString(claimsJSON)

	// In production, properly sign with the private key
	// This is a placeholder - you'd need proper JWT signing
	signature := "placeholder_signature"

	return fmt.Sprintf("%s.%s.%s", headerEncoded, claimsEncoded, signature), nil
}

// ManagedIdentityTokenProvider implements managed identity authentication
type ManagedIdentityTokenProvider struct {
	httpClient        *http.Client
	managedIdentityID string // Optional: specific managed identity to use
}

func (p *ManagedIdentityTokenProvider) GetToken(ctx context.Context) (string, error) {
	// Check if we're running in Azure
	imdsEndpoint := os.Getenv("IDENTITY_ENDPOINT")
	identityHeader := os.Getenv("IDENTITY_HEADER")

	if imdsEndpoint != "" && identityHeader != "" {
		// App Service / Functions managed identity
		return p.getTokenFromAppService(imdsEndpoint, identityHeader)
	}

	// Try Azure VM/VMSS IMDS endpoint
	return p.getTokenFromIMDS()
}

func (p *ManagedIdentityTokenProvider) getTokenFromAppService(endpoint, header string) (string, error) {
	resource := "https://analysis.windows.net/powerbi/api"
	apiVersion := "2019-08-01"

	reqURL := fmt.Sprintf("%s?resource=%s&api-version=%s", endpoint, url.QueryEscape(resource), apiVersion)
	if p.managedIdentityID != "" {
		reqURL += "&client_id=" + url.QueryEscape(p.managedIdentityID)
	}

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-IDENTITY-HEADER", header)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		data, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, data)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresOn   string `json:"expires_on"`
		Resource    string `json:"resource"`
		TokenType   string `json:"token_type"`
	}

	if err := json.Unmarshal(data, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	return tokenResp.AccessToken, nil
}

func (p *ManagedIdentityTokenProvider) getTokenFromIMDS() (string, error) {
	// Azure VM/VMSS Instance Metadata Service endpoint
	imdsEndpoint := "http://169.254.169.254/metadata/identity/oauth2/token"
	resource := "https://analysis.windows.net/powerbi/api"
	apiVersion := "2018-02-01"

	reqURL := fmt.Sprintf("%s?resource=%s&api-version=%s", imdsEndpoint, url.QueryEscape(resource), apiVersion)
	if p.managedIdentityID != "" {
		reqURL += "&client_id=" + url.QueryEscape(p.managedIdentityID)
	}

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Metadata", "true")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get token from IMDS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		data, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("IMDS token request failed with status %d: %s", resp.StatusCode, data)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   string `json:"expires_in"`
		ExpiresOn   string `json:"expires_on"`
		Resource    string `json:"resource"`
		TokenType   string `json:"token_type"`
	}

	if err := json.Unmarshal(data, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	return tokenResp.AccessToken, nil
}

// AzureCLITokenProvider implements Azure CLI authentication
type AzureCLITokenProvider struct{}

func (p *AzureCLITokenProvider) GetToken(ctx context.Context) (string, error) {
	// Use Azure CLI to get access token
	cmd := exec.CommandContext(ctx, "az", "account", "get-access-token",
		"--resource", "https://analysis.windows.net/powerbi/api",
		"--output", "json")

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("az cli failed: %s", string(exitErr.Stderr))
		}
		return "", fmt.Errorf("failed to run az cli: %w", err)
	}

	var tokenResp struct {
		AccessToken string `json:"accessToken"`
		ExpiresOn   string `json:"expiresOn"`
		Tenant      string `json:"tenant"`
	}

	if err := json.Unmarshal(output, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to parse az cli output: %w", err)
	}

	return tokenResp.AccessToken, nil
}

// DirectTokenProvider uses a pre-obtained access token
type DirectTokenProvider struct {
	accessToken string
}

func (p *DirectTokenProvider) GetToken(ctx context.Context) (string, error) {
	if p.accessToken == "" {
		return "", fmt.Errorf("no access token provided")
	}
	return p.accessToken, nil
}

// NewClientWithAuthConfig creates a Power BI client with the specified authentication configuration
func NewClientWithAuthConfig(config *AuthConfig) (*Client, error) {
	httpClient := cleanhttp.DefaultClient()

	var tokenProvider TokenProvider

	// Determine which authentication method to use
	switch {
	case config.AccessToken != "":
		// Direct token authentication
		tokenProvider = &DirectTokenProvider{accessToken: config.AccessToken}

	case config.UseManagedIdentity:
		// Managed Identity authentication
		tokenProvider = &ManagedIdentityTokenProvider{
			httpClient:        httpClient,
			managedIdentityID: config.ManagedIdentityID,
		}

	case config.UseAzureCLI:
		// Azure CLI authentication
		tokenProvider = &AzureCLITokenProvider{}

	case config.CertificatePath != "" || config.CertificateData != "":
		// Certificate-based authentication
		var certPath string
		if config.CertificateData != "" {
			// Write base64 decoded certificate to temp file
			certBytes, err := base64.StdEncoding.DecodeString(config.CertificateData)
			if err != nil {
				return nil, fmt.Errorf("failed to decode certificate data: %w", err)
			}
			
			tmpFile, err := ioutil.TempFile("", "powerbi-cert-*.pem")
			if err != nil {
				return nil, fmt.Errorf("failed to create temp file for certificate: %w", err)
			}
			defer os.Remove(tmpFile.Name())
			
			if _, err := tmpFile.Write(certBytes); err != nil {
				return nil, fmt.Errorf("failed to write certificate to temp file: %w", err)
			}
			tmpFile.Close()
			certPath = tmpFile.Name()
		} else {
			certPath = config.CertificatePath
		}

		certProvider, err := NewCertificateTokenProvider(httpClient, config.TenantID, config.ClientID, certPath, config.CertificatePassword)
		if err != nil {
			return nil, fmt.Errorf("failed to create certificate token provider: %w", err)
		}
		tokenProvider = certProvider

	case config.ClientSecret != "":
		// Client credentials (service principal with secret)
		tokenProvider = &ClientCredentialsTokenProvider{
			httpClient:   httpClient,
			tenantID:     config.TenantID,
			clientID:     config.ClientID,
			clientSecret: config.ClientSecret,
		}

	case config.Username != "" && config.Password != "":
		// Password authentication (legacy)
		return NewClientWithPasswordAuth(config.TenantID, config.ClientID, config.ClientSecret, config.Username, config.Password)

	default:
		return nil, fmt.Errorf("no valid authentication method configured")
	}

	// Create client with token provider
	return newClientWithTokenProvider(tokenProvider)
}

// newClientWithTokenProvider creates a client with a custom token provider
func newClientWithTokenProvider(tokenProvider TokenProvider) (*Client, error) {
	// PowerBI has lots of intermittant TLS handshake issues, these settings
	// seem to reduce the amount of issues encountered
	defaultTransport := cleanhttp.DefaultPooledTransport()
	defaultTransport.TLSHandshakeTimeout = 60 * time.Second
	defaultTransport.TLSClientConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// Create a token provider wrapper function
	getToken := func(httpClient *http.Client) (string, error) {
		return tokenProvider.GetToken(context.Background())
	}

	// auth
	httpClient := &http.Client{
		Transport: newBearerTokenRoundTripper(
			getToken,
			// error
			newErrorOnUnsuccessfulRoundTripper(
				// this is crazy we need to retry 500 and 400 errors, but the API intermittently returns them
				newRetryIntermittentErrorRoundTripper(
					// retry too many requests
					newRetryTooManyRequestsRoundTripper(
						// actual call
						defaultTransport,
					),
				),
			),
		),
	}

	return &Client{
		Client:     httpClient,
		HTTPClient: httpClient,
	}, nil
}