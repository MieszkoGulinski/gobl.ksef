package api

import (
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/beevik/etree"
	"github.com/invopop/xmldsig"
)

var ErrCertificatePrivateKeyNotRSA = errors.New("certificate private key is not RSA, goxades only supports RSA")

func (c *Client) buildSignedAuthorizationRequest(challenge *authorizationChallengeResponse, contextIdentifier *ContextIdentifier) ([]byte, error) {
	// I tried to use the github.com/invopop/xmldsig library, but it doesn't work, as it has many options hardcoded that aren't compatible with the KSEF API

	// 1. Assembly the XML request - the signing library requires XML as an etree object

	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="utf-8"`)

	root := doc.CreateElement("AuthTokenRequest")
	root.CreateAttr("xmlns:xsi", "http://www.w3.org/2001/XMLSchema-instance")
	root.CreateAttr("xmlns:xsd", "http://www.w3.org/2001/XMLSchema")
	root.CreateAttr("xmlns", "http://ksef.mf.gov.pl/auth/token/2.0")

	root.CreateElement("Challenge").SetText(challenge.Challenge)

	ctx := root.CreateElement("ContextIdentifier")
	if contextIdentifier.Nip != "" {
		ctx.CreateElement("Nip").SetText(contextIdentifier.Nip)
	}
	if contextIdentifier.NipVatUe != "" {
		ctx.CreateElement("NipVatUe").SetText(contextIdentifier.NipVatUe)
	}
	if contextIdentifier.InternalId != "" {
		ctx.CreateElement("InternalId").SetText(contextIdentifier.InternalId)
	}
	if contextIdentifier.PeppolId != "" {
		ctx.CreateElement("PeppolId").SetText(contextIdentifier.PeppolId)
	}

	subjectIdentifierType := "certificateSubject"
	if contextIdentifier != nil && contextIdentifier.NipVatUe != "" {
		subjectIdentifierType = "certificateFingerprint"
	}
	root.CreateElement("SubjectIdentifierType").SetText(subjectIdentifierType)

	unsignedXML, err := doc.WriteToString()
	if err != nil {
		return nil, err
	}

	// Sign
	cert, _ := xmldsig.LoadCertificate(c.certificatePath, "")
	signature, _ := xmldsig.Sign([]byte(unsignedXML),
		xmldsig.WithCertificate(cert),
		xmldsig.WithKSeF(),
	)

	// attach signature to XML
	signatureXML, err := xml.Marshal(signature)
	if err != nil {
		return nil, err
	}
	sigDoc := etree.NewDocument()
	if err := sigDoc.ReadFromBytes(signatureXML); err != nil {
		return nil, err
	}
	root.AddChild(sigDoc.Root())

	signedXML, err := doc.WriteToString()
	if err != nil {
		return nil, err
	}

	fmt.Println(signedXML)

	return []byte(signedXML), nil
}
