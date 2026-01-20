package ksef

import (
	"fmt"

	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
)

// Address defines the XML structure for KSeF addresses
type Address struct {
	CountryCode string `xml:"KodKraju"`
	AddressL1   string `xml:"AdresL1"`
	AddressL2   string `xml:"AdresL2,omitempty"`
}

// Seller defines the XML structure for KSeF seller
type Seller struct {
	NIP     string          `xml:"DaneIdentyfikacyjne>NIP"`
	Name    string          `xml:"DaneIdentyfikacyjne>Nazwa"`
	Address *Address        `xml:"Adres"`
	Contact *ContactDetails `xml:"DaneKontaktowe,omitempty"`
}

// ContactDetails defines the XML structure for KSeF contact
type ContactDetails struct {
	Email string `xml:"Email,omitempty"`
	Phone string `xml:"Telefon,omitempty"`
}

// Buyer defines the XML structure for KSeF buyer
type Buyer struct {
	NIP string `xml:"DaneIdentyfikacyjne>NIP,omitempty"`
	// or
	UECode      string `xml:"DaneIdentyfikacyjne>KodUE,omitempty"`   // Country code when in European Union
	UEVatNumber string `xml:"DaneIdentyfikacyjne>NrVatUE,omitempty"` // EU VAT number
	// or
	CountryCode string `xml:"DaneIdentyfikacyjne>KodKraju,omitempty"` // Country code outside European Union
	IDNumber    string `xml:"DaneIdentyfikacyjne>NrID,omitempty"`     // Tax ID number outside European Union
	// or
	NoID int `xml:"DaneIdentyfikacyjne>BrakID,omitempty"`

	Name    string          `xml:"DaneIdentyfikacyjne>Nazwa,omitempty"`
	Address *Address        `xml:"Adres,omitempty"`
	Contact *ContactDetails `xml:"DaneKontaktowe,omitempty"`

	JST int `xml:"JST"` // JST (Jednostka Samorządu Terytorialnego = local government unit) 1 = Yes, 2 = No
	GV  int `xml:"GV"`  // GV (Group VAT) 1 = Yes, 2 = No
}

// newAddress gets the address data from GOBL address
func newAddress(address *org.Address) *Address {
	adres := &Address{
		CountryCode: string(address.Country),
		AddressL1:   addressLine1(address),
		AddressL2:   addressLine2(address),
	}

	return adres
}

// nameToString get the seller name out of the organization
func nameToString(name *org.Name) string {
	return name.Prefix + nameMaybe(name.Given) +
		nameMaybe(name.Middle) + nameMaybe(name.Surname) +
		nameMaybe(name.Surname2) + nameMaybe(name.Suffix)
}

// NewSeller converts a GOBL Party into a KSeF seller
func NewSeller(supplier *org.Party) *Seller {
	var name string
	if supplier.Name != "" {
		name = supplier.Name
	} else {
		name = nameToString(supplier.People[0].Name)
	}
	seller := &Seller{
		Address: newAddress(supplier.Addresses[0]),
		NIP:     string(supplier.TaxID.Code),
		Name:    name,
	}
	if len(supplier.Telephones) > 0 {
		seller.Contact = &ContactDetails{
			Phone: supplier.Telephones[0].Number,
		}
	}
	if len(supplier.Emails) > 0 {
		if seller.Contact == nil {
			seller.Contact = &ContactDetails{}
		}
		seller.Contact.Email = supplier.Emails[0].Address
	}

	return seller
}

// EU countries
var euCountriesCodes = []string{
	"AT", // Austria
	"BE", // Belgium
	"BG", // Bulgaria
	"CY", // Cyprus
	"CZ", // Czech Republic
	"DE", // Germany
	"DK", // Denmark
	"EE", // Estonia
	"EL", // Greece
	"ES", // Spain
	"FI", // Finland
	"FR", // France
	"HR", // Croatia
	"HU", // Hungary
	"IE", // Ireland
	"IT", // Italy
	"LT", // Lithuania
	"LU", // Luxembourg
	"LV", // Latvia
	"MT", // Malta
	"NL", // Netherlands
	"PL", // Poland
	"PT", // Portugal
	"RO", // Romania
	"SE", // Sweden
	"SI", // Slovenia
	"SK", // Slovakia
	"XI", // Northern Ireland (listed in XML schema but not member of EU)
}

// NewBuyer converts a GOBL Party into a KSeF buyer
func NewBuyer(customer *org.Party) *Buyer {

	buyer := &Buyer{
		Name: customer.Name,
		JST:  2, // hardcoded as "No"
		GV:   2, // hardcoded as "No"
	}

	if customer.TaxID == nil {
		// Buyer is a private individual
		buyer.NoID = 1
	} else if customer.TaxID.Country == l10n.PL.Tax() {
		// Buyer is a Polish business entity
		buyer.NIP = string(customer.TaxID.Code)
	} else if isEUCountry(string(customer.TaxID.Country)) {
		// Buyer is an EU business entity (non-Polish)
		buyer.UECode = string(customer.TaxID.Country)
		buyer.UEVatNumber = string(customer.TaxID.Code)
	} else {
		// Buyer is a business entity from outside the EU
		buyer.CountryCode = string(customer.TaxID.Country)
		if len(customer.TaxID.Code) > 0 {
			buyer.IDNumber = string(customer.TaxID.Code)
		}
	}

	fmt.Println(buyer.Name)

	if len(customer.Addresses) > 0 {
		buyer.Address = newAddress(customer.Addresses[0])
	}

	if len(customer.Telephones) > 0 {
		buyer.Contact = &ContactDetails{
			Phone: customer.Telephones[0].Number,
		}
	}
	if len(customer.Emails) > 0 {
		if buyer.Contact == nil {
			buyer.Contact = &ContactDetails{}
		}
		buyer.Contact.Email = customer.Emails[0].Address
	}

	return buyer
}

func addressLine1(address *org.Address) string {
	if address.PostOfficeBox != "" {
		return address.PostOfficeBox
	}

	return address.Street +
		", " + address.Number +
		addressMaybe(address.Block) +
		addressMaybe(address.Floor) +
		addressMaybe(address.Door)
}

func addressLine2(address *org.Address) string {
	return address.Code.String() + ", " + address.Locality
}

func addressMaybe(element string) string {
	if element != "" {
		return ", " + element
	}
	return ""
}

func nameMaybe(element string) string {
	if element != "" {
		return " " + element
	}
	return ""
}

func isEUCountry(code string) bool {
	for _, c := range euCountriesCodes {
		if c == code {
			return true
		}
	}
	return false
}
