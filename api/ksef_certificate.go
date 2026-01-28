package api

import (
	"context"
	"fmt"
	"time"
)

// CertificateEnrollmentData stores the subject data required to prepare a certificate request.
type CertificateEnrollmentData struct {
	CommonName             string `json:"commonName"`
	Surname                string `json:"surname,omitempty"`
	SerialNumber           string `json:"serialNumber,omitempty"`
	CountryName            string `json:"countryName"`
	OrganizationName       string `json:"organizationName,omitempty"`
	GivenName              string `json:"givenName,omitempty"`
	UniqueIdentifier       string `json:"uniqueIdentifier"`
	OrganizationIdentifier string `json:"organizationIdentifier,omitempty"`
}

// CertificateType identifies the target certificate flavor.
type CertificateType string

const (
	CertificateTypeAuthentication CertificateType = "Authentication"
	CertificateTypeOffline        CertificateType = "Offline"
)

func (t CertificateType) isValid() bool {
	switch t {
	case CertificateTypeAuthentication, CertificateTypeOffline:
		return true
	default:
		return false
	}
}

type certificateEnrollmentRequest struct {
	CertificateName string          `json:"certificateName"`
	CertificateType CertificateType `json:"certificateType"`
	CSR             string          `json:"csr"`
	ValidFrom       *time.Time      `json:"validFrom,omitempty"`
}

// CertificateEnrollmentResponse stores the response metadata returned after submitting a CSR.
type CertificateEnrollmentResponse struct {
	ReferenceNumber string `json:"referenceNumber"`
	Timestamp       string `json:"timestamp"`
}

// GetCertificateEnrollmentData returns the identification data used when building a new CSR.
func (c *Client) GetCertificateEnrollmentData(ctx context.Context) (*CertificateEnrollmentData, error) {
	token, err := c.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	response := &CertificateEnrollmentData{}
	resp, err := c.client.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetResult(response).
		Get(c.url + "/certificates/enrollments/data")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, newErrorResponse(resp)
	}

	return response, nil
}

// EnrollCertificate submits the generated CSR and returns reference metadata.
func (c *Client) EnrollCertificate(ctx context.Context, certificateName string, certificateType CertificateType, csr string, validFrom *time.Time) (*CertificateEnrollmentResponse, error) {
	if certificateName == "" {
		return nil, fmt.Errorf("certificateName is required")
	}
	if !certificateType.isValid() {
		return nil, fmt.Errorf("certificateType is invalid: %s", certificateType)
	}
	if csr == "" {
		return nil, fmt.Errorf("csr is required")
	}

	token, err := c.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	request := &certificateEnrollmentRequest{
		CertificateName: certificateName,
		CertificateType: certificateType,
		CSR:             csr,
		ValidFrom:       validFrom,
	}

	response := &CertificateEnrollmentResponse{}
	resp, err := c.client.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetBody(request).
		SetResult(response).
		Post(c.url + "/certificates/enrollments")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, newErrorResponse(resp)
	}

	return response, nil
}
