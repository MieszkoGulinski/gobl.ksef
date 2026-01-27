package ksef

import (
	"github.com/invopop/gobl/addons/pl/favat"
	"github.com/invopop/gobl/bill"
)

// Payment defines the XML structure for KSeF payment
type Payment struct {
	PaidMarker             string            `xml:"Zaplacono,omitempty"`
	PaymentDate            string            `xml:"DataZaplaty,omitempty"`
	PartiallyPaidMarker    string            `xml:"ZnacznikZaplatyCzesciowej,omitempty"`
	AdvancePayments        []*AdvancePayment `xml:"ZaplataCzesciowa,omitempty"`
	DueDates               []*DueDate        `xml:"TerminPlatnosci,omitempty"`
	PaymentMean            string            `xml:"FormaPlatnosci,omitempty"` // enum: 1 = cash, 2 = card etc. (see KSeF documentation)
	OtherPaymentMeanMarker string            `xml:"PlatnoscInna,omitempty"`
	OtherPaymentMean       string            `xml:"OpisPlatnosci,omitempty"`
	BankAccounts           []*BankAccount    `xml:"RachunekBankowy,omitempty"`
	FactorBankAccounts     []*BankAccount    `xml:"RachunekBankowyFaktora,omitempty"`
	Discount               *Discount         `xml:"Skonto,omitempty"`
	PaymentLink            string            `xml:"LinkDoPlatnosci,omitempty"`
	KSeFPaymentID          string            `xml:"IPKSeF,omitempty"`
}

// AdvancePayment defines the XML structure for KSeF advance payments
type AdvancePayment struct {
	PaymentAmount          string `xml:"KwotaZaplatyCzesciowej,omitempty"`
	PaymentDate            string `xml:"DataZaplatyCzesciowej,omitempty"`
	PaymentMean            string `xml:"FormaPlatnosci,omitempty"`
	OtherPaymentMeanMarker int    `xml:"PlatnoscInna,omitempty"`
	OtherPaymentMean       string `xml:"OpisPlatnosci,omitempty"`
}

// DueDate defines the XML structure for KSeF due date
type DueDate struct {
	Date            string           `xml:"Termin,omitempty"`
	TermDescription *TermDescription `xml:"TerminOpis,omitempty"`
}

// TermDescription defines alternative payment term description
type TermDescription struct {
	Quantity      int    `xml:"Ilosc"`
	Unit          string `xml:"Jednostka"`
	StartingEvent string `xml:"ZdarzeniePoczatkowe"`
}

// BankAccount defines the XML structure for KSeF bank accounts
type BankAccount struct {
	AccountNumber         string `xml:"NrRB"`
	SWIFT                 string `xml:"SWIFT,omitempty"`
	BankSelfAccountMarker int    `xml:"RachunekWlasnyBanku,omitempty"` // enum - 1,2,3, not sure what exactly they mean
	BankName              string `xml:"NazwaBanku,omitempty"`
	AccountDescription    string `xml:"OpisRachunku,omitempty"`
}

// Discount defines the XML structure for KSeF early payment discount
type Discount struct {
	Conditions string `xml:"WarunkiSkonta,omitempty"`
	Amount     string `xml:"WysokoscSkonta,omitempty"`
}

// NewPayment gets payment data from GOBL invoice
func NewPayment(pay *bill.PaymentDetails, totals *bill.Totals) *Payment {
	if pay == nil {
		return nil
	}

	var payment = &Payment{
		DueDates:        []*DueDate{},
		AdvancePayments: []*AdvancePayment{},
	}

	if instructions := pay.Instructions; instructions != nil {
		paymentMeansCode := instructions.Ext.Get(favat.ExtKeyPaymentMeans).String()

		if paymentMeansCode == "" && instructions.Key != "" {
			payment.OtherPaymentMeanMarker = "1"
			payment.OtherPaymentMean = instructions.Key.String()
		} else if paymentMeansCode != "" {
			payment.PaymentMean = paymentMeansCode
		}

		payment.BankAccounts = []*BankAccount{}
		payment.FactorBankAccounts = []*BankAccount{}

		for _, account := range instructions.CreditTransfer {
			accountNumber := account.IBAN
			if accountNumber == "" {
				accountNumber = account.Number
			}
			payment.BankAccounts = append(payment.BankAccounts, &BankAccount{
				AccountNumber: accountNumber,
				SWIFT:         account.BIC,
				BankName:      account.Name,
			})
		}
	}

	if terms := pay.Terms; terms != nil {
		for _, dueDate := range pay.Terms.DueDates {
			payment.DueDates = append(payment.DueDates, &DueDate{
				Date: dueDate.Date.String(),
			})
		}
	}

	// According to FA_VAT v3 schema:
	// If an invoice is paid in full in one payment, PaidMarker should be "1"
	// Otherwise, set PartiallyPaidMarker with the following values:
	// 1 = invoice paid partially
	// 2 = paid in full after partial payments, and the last payment is the final one
	// If the invoice is not paid at all, do not add PaidMarker or PartiallyPaidMarker.

	if advances := pay.Advances; advances != nil {
		if len(advances) == 1 && totals.Due.IsZero() {
			// Invoice already paid in full in one payment
			payment.PaidMarker = "1"
			if advances[0].Date != nil {
				payment.PaymentDate = advances[0].Date.String()
			}
		} else {
			if totals.Due.IsZero() {
				// Invoice already paid in full in multiple payments
				payment.PartiallyPaidMarker = "2"
			}
			if !totals.Due.IsZero() && len(advances) > 0 {
				// Invoice paid partially
				payment.PartiallyPaidMarker = "1"
			}
			// Otherwise, not paid at all - no markers needed

			for _, advance := range advances {
				advancePayment := &AdvancePayment{
					PaymentAmount: advance.Amount.String(),
					PaymentDate:   advance.Date.String(),
				}

				if paymentMeansCode := advance.Ext.Get(favat.ExtKeyPaymentMeans).String(); paymentMeansCode != "" {
					advancePayment.PaymentMean = paymentMeansCode
				}
				payment.AdvancePayments = append(payment.AdvancePayments, advancePayment)
			}
		}
	}

	return payment
}
