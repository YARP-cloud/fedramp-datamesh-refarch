package security

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/frocore/fedramp-data-mesh/cli/internal/config"
)

type SecurityContext struct {
	cfg            *config.Config
	awsCredentials *credentials.Credentials
	lastRefresh    time.Time
	sessionToken   string
	role           string
}

func NewSecurityContext(cfg *config.Config) (*SecurityContext, error) {
	ctx := &SecurityContext{
		cfg:  cfg,
		role: cfg.DefaultRole,
	}
	
	// Initialize AWS credentials
	if err := ctx.initAWSCredentials(); err != nil {
		return nil, err
	}
	
	return ctx, nil
}

func (s *SecurityContext) initAWSCredentials() error {
	// Check for profile in AWS config
	if s.cfg.AWSProfile != "" {
		s.awsCredentials = credentials.NewSharedCredentials("", s.cfg.AWSProfile)
		return nil
	}
	
	// Check environment variables
	envCreds := credentials.NewEnvCredentials()
	if _, err := envCreds.Get(); err == nil {
		s.awsCredentials = envCreds
		return nil
	}
	
	// Try to get credentials from AWS IAM Identity Center
	if err := s.refreshFromSSO(); err != nil {
		return fmt.Errorf("failed to get AWS credentials: %w", err)
	}
	
	return nil
}

func (s *SecurityContext) refreshFromSSO() error {
	// Check if SSO token exists and is valid
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not get user home directory: %w", err)
	}
	
	ssoDir := filepath.Join(homeDir, ".aws", "sso", "cache")
	if _, err := os.Stat(ssoDir); os.IsNotExist(err) {
		return fmt.Errorf("AWS SSO cache directory not found, please run 'aws sso login'")
	}
	
	// Find the most recent SSO token file
	var newestFile string
	var newestTime time.Time
	
	err = filepath.Walk(ssoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		
		if newestFile == "" || info.ModTime().After(newestTime) {
			newestFile = path
			newestTime = info.ModTime()
		}
		
		return nil
	})
	
	if err != nil {
		return fmt.Errorf("error walking SSO cache directory: %w", err)
	}
	
	if newestFile == "" {
		return fmt.Errorf("no SSO token found, please run 'aws sso login'")
	}
	
	// Read token file to check expiration
	data, err := os.ReadFile(newestFile)
	if err != nil {
		return fmt.Errorf("error reading SSO token file: %w", err)
	}
	
	var tokenData struct {
		ExpiresAt string `json:"expiresAt"`
	}
	
	if err := json.Unmarshal(data, &tokenData); err != nil {
		return fmt.Errorf("error parsing SSO token file: %w", err)
	}
	
	expiresAt, err := time.Parse(time.RFC3339, tokenData.ExpiresAt)
	if err != nil {
		return fmt.Errorf("error parsing expiration time: %w", err)
	}
	
	if time.Now().After(expiresAt) {
		return fmt.Errorf("SSO token expired, please run 'aws sso login'")
	}
	
	// Get temporary credentials for the role
	cmd := exec.Command("aws", "sso", "get-role-credentials",
		"--profile", s.cfg.AWSProfile,
		"--role-name", s.role,
		"--account-id", s.cfg.AWSAccountID)
	
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error getting role credentials: %w", err)
	}
	
	var result struct {
		RoleCredentials struct {
			AccessKeyId     string `json:"accessKeyId"`
			SecretAccessKey string `json:"secretAccessKey"`
			SessionToken    string `json:"sessionToken"`
			Expiration      int64  `json:"expiration"`
		} `json:"roleCredentials"`
	}
	
	if err := json.Unmarshal(output, &result); err != nil {
		return fmt.Errorf("error parsing role credentials: %w", err)
	}
	
	// Create static credentials from SSO
	s.awsCredentials = credentials.NewStaticCredentials(
		result.RoleCredentials.AccessKeyId,
		result.RoleCredentials.SecretAccessKey,
		result.RoleCredentials.SessionToken,
	)
	
	s.lastRefresh = time.Now()
	s.sessionToken = result.RoleCredentials.SessionToken
	
	return nil
}

func (s *SecurityContext) GetAWSCredentials() (string, string, string, error) {
	// Check if credentials need refresh
	if time.Since(s.lastRefresh) > 55*time.Minute {
		if err := s.refreshFromSSO(); err != nil {
			return "", "", "", err
		}
	}
	
	// Get credential values
	creds, err := s.awsCredentials.Get()
	if err != nil {
		return "", "", "", fmt.Errorf("error getting AWS credentials: %w", err)
	}
	
	return creds.AccessKeyID, creds.SecretAccessKey, creds.SessionToken, nil
}

func (s *SecurityContext) GetAWSSession() (*session.Session, error) {
	accessKey, secretKey, sessionToken, err := s.GetAWSCredentials()
	if err != nil {
		return nil, err
	}
	
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(s.cfg.AWSRegion),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, sessionToken),
	})
	
	if err != nil {
		return nil, fmt.Errorf("error creating AWS session: %w", err)
	}
	
	return sess, nil
}

func (s *SecurityContext) AssumeRole(role string) error {
	// Save current role
	prevRole := s.role
	s.role = role
	
	// Try to get credentials with new role
	err := s.refreshFromSSO()
	if err != nil {
		// Restore previous role on failure
		s.role = prevRole
		return fmt.Errorf("failed to assume role %s: %w", role, err)
	}
	
	return nil
}

func (s *SecurityContext) GetCurrentRole() string {
	return s.role
}

// Function to validate if a user has access to a specific data product
func (s *SecurityContext) CanAccessDataProduct(dataProduct string) (bool, error) {
	// Get AWS session
	sess, err := s.GetAWSSession()
	if err != nil {
		return false, err
	}
	
	// Call STS GetCallerIdentity to get current identity
	stsClient := sts.New(sess)
	identity, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return false, fmt.Errorf("error getting caller identity: %w", err)
	}
	
	// In a real implementation, we would check against the data catalog
	// to determine if the current user/role has access to the data product
	// This is a simplified implementation that always returns true
	
	// Log the access attempt for audit purposes
	fmt.Printf("User %s is requesting access to data product %s\n", *identity.Arn, dataProduct)
	
	return true, nil
}
