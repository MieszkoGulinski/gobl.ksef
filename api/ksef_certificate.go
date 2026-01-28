package api

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"net/url"
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

var (
	ErrCertificateEnrollmentPollingCountExceeded = errors.New("certificate enrollment polling count exceeded")
	ErrCertificateEnrollmentFailed               = errors.New("certificate enrollment failed")
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

// CertificateCreationResult bundles the responses returned during certificate creation.
// CertificateEnrollmentStatusResponse describes the status returned for an enrollment reference number.
type CertificateEnrollmentStatusResponse struct {
	RequestDate             time.Time                    `json:"requestDate"`
	Status                  *CertificateEnrollmentStatus `json:"status"`
	CertificateSerialNumber string                       `json:"certificateSerialNumber,omitempty"`
}

// CertificateEnrollmentStatus mirrors the StatusInfo structure returned by the API.
type CertificateEnrollmentStatus struct {
	Code        int      `json:"code"`
	Description string   `json:"description"`
	Details     []string `json:"details,omitempty"`
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
// validFrom is optional, if not provided, the certificate will be valid from the current date.
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

// CreateKsefCertificate orchestrates the full flow of requesting a KSeF certificate using the provided key material.
// Returns the certificate serial number once the request succeeds.
func (c *Client) CreateKsefCertificate(ctx context.Context, certificateName string, certificateType CertificateType, privateKey *ecdsa.PrivateKey, validFrom *time.Time) (string, error) {
	if privateKey == nil {
		return "", fmt.Errorf("private key is required")
	}

	enrollmentData, err := c.GetCertificateEnrollmentData(ctx)
	if err != nil {
		return "", err
	}

	csr, err := enrollmentData.GenerateCSR(privateKey)
	if err != nil {
		return "", err
	}

	enrollmentResp, err := c.EnrollCertificate(ctx, certificateName, certificateType, csr, validFrom)
	if err != nil {
		return "", err
	}

	statusResp, err := c.PollCertificateEnrollmentStatus(ctx, enrollmentResp.ReferenceNumber)
	if err != nil {
		return "", err
	}

	if statusResp.CertificateSerialNumber == "" {
		return "", fmt.Errorf("certificate serial number missing in enrollment status response")
	}
	return statusResp.CertificateSerialNumber, nil
}

// RevokeCertificate submits a revocation request for the provided certificate serial number.
func (c *Client) RevokeCertificate(ctx context.Context, certificateSerialNumber string) error {
	if certificateSerialNumber == "" {
		return fmt.Errorf("certificateSerialNumber is required")
	}

	token, err := c.getAccessToken(ctx)
	if err != nil {
		return err
	}

	resp, err := c.client.R().
		SetContext(ctx).
		SetAuthToken(token).
		Post(c.url + "/certificates/" + url.PathEscape(certificateSerialNumber) + "/revoke")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return newErrorResponse(resp)
	}

	return nil
}

// GetCertificateEnrollmentStatus returns processing status for the specified certificate request reference number.
func (c *Client) GetCertificateEnrollmentStatus(ctx context.Context, referenceNumber string) (*CertificateEnrollmentStatusResponse, error) {
	if referenceNumber == "" {
		return nil, fmt.Errorf("referenceNumber is required")
	}

	token, err := c.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	response := &CertificateEnrollmentStatusResponse{}
	resp, err := c.client.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetResult(response).
		Get(c.url + "/certificates/enrollments/" + url.PathEscape(referenceNumber))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, newErrorResponse(resp)
	}

	return response, nil
}

// PollCertificateEnrollmentStatus keeps querying enrollment status until it succeeds or fails.
func (c *Client) PollCertificateEnrollmentStatus(ctx context.Context, referenceNumber string) (*CertificateEnrollmentStatusResponse, error) {
	attempt := 0
	for {
		attempt++
		if attempt > 30 {
			return nil, ErrCertificateEnrollmentPollingCountExceeded
		}

		statusResp, err := c.GetCertificateEnrollmentStatus(ctx, referenceNumber)
		if err != nil {
			return nil, err
		}
		if statusResp == nil || statusResp.Status == nil {
			return nil, fmt.Errorf("certificate enrollment status response missing status")
		}

		switch statusResp.Status.Code {
		case 200:
			return statusResp, nil
		case 100:
			time.Sleep(2 * time.Second)
			continue
		default:
			return nil, fmt.Errorf("%w: %s", ErrCertificateEnrollmentFailed, statusResp.Status.Description)
		}
	}
}
