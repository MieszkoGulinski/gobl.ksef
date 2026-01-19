# Unsupported fields (hardcoded or not mapped)

Here are features not supported by the converter.

## Hardcoded values

| XML field | Struct field | Value | Notes |
| --------- | ---------- | ----- | ----- |
| `KodFormularza` | `FormCode` | `FA` | |
| `KodFormularza@kodSystemowy` (attribute) | `FormCode.SystemCode` | `FA (3)` | |
| `KodFormularza@wersjaSchemy` (attribute) | `FormCode.SchemaVersion` | `1-0E` | |
| `WariantFormularza` | `FormVariant` | `3` | |
| `SystemInfo` | `SystemInfo` | `GOBL.KSEF` | |
| `Podmiot2>JST` | `JST` | `2` | Jednostka Samorządu Terytorialnego = local government unit, special case where the buyer is a local government unit |
| `Podmiot2>GV` | `GV` | `2` | Grupa VAT = VAT group, special case |
| `Fa>P_16` | `CashAccounting` | `2` | |
| `Fa>P_18` | `ReverseCharge` | `2` | |
| `Fa>P_18A` | `SplitPaymentMechanism` | `2` | |
| `Fa>Adnotacje>Zwolnienie>P_19N` | `NoTaxExemptGoods` | `1` | For tax exempt goods, set `P_19` to 1, otherwise set `P_19N` to 1 |
| `Fa>NoweSrodkiTransportu>P_22N` | `NoNewTransportIntraCommunitySupply` | `1` | For new transport intra-community supply, set `P_22` to 1 (rare special case), otherwise set `P_22N` to 1 |
| `Fa>P_23` | `SimplifiedProcedureBySecondTaxpayer` | `2` | For simplified procedure by second taxpayer (for three-party transactions inside the European Union), set `P_23` to 1, otherwise set `P_23` to 2 |
| `Fa>PMarzy>P_PMarzyN` | `NoMarginProcedures` | `1` | For margin procedure (applies to specific types of goods and services), set `P_PMarzy` to 1, otherwise set `P_PMarzyN` to 1 |
| `Fa>Platnosc>RachunekBankowyFaktora` | `FactorBankAccounts` | `[]` | Bank account of the factor (third party) |

## Not mapped

The following fields are not present in the structs to be converted to XML:

| XML field |  Notes |
| --------- |  ----- |
| `Podmiot2>NrKlienta`  | customer number, if the supplier uses such a number in the contract or order document - not required in schema |
| `Podmiot2>IDNabywcy`  | unique key of the customer, if customer's data changed between base invoice and correction invoice - not required in schema |
| `Fa>FaWiersz>UU_ID` | unique identifier for the line item - not required in schema |
| `Fa>WarunkiTransakcji` | transaction conditions, containing contract date and number and/or order date and number - not required in schema |
| `Stopka` | footer information, including information identifying the parties in various national databases - not required in schema |
| `Podmiot3>Rola` | role of the third party - required if a third party is present, enum with values from 1 to 11 |
| `Podmiot3>RolaInna` | role of the third party - required if a third party is present, set to "1" for "other" role, not included in options above  |
| `Podmiot3>OpisRoli` | role of the third party - required if a third party is present, fill with description of the role if "RolaInna" is set to "1" |
| `Podmiot3>Udzial` | percentage share of the third party (when there are two buyers) |
| `Podmiot3>NrKlienta` | similar to `Podmiot2>NrKlienta` but for the third party |
| `Fa>WZ` | warehouse document(s) number - not required in schema |
| `Fa>FaWiersz>CN` | Combined Nomenclature product type code - not required in schema |
| `Fa>Rozliczenie` | additional discounts and charges - not required in schema |
| `Fa>Podmiot2K` | in case of correction invoice where customer's data changed, contains the customer data from the old invoice - not required in schema |
| `Fa>OkresFaKorygowanej` | when the correction invoice indicates a discount, period of the discount - not required in schema |
| `Fa>Adnotacje>PMarzy>P_PMarzy` | set to "1" when using margin scheme, otherwise set `P_PMarzyN` to "1" - scheme requires either one to be set |
| `Fa>Adnotacje>PMarzy>P_PMarzy_2` | margin scheme for travel agencies |
| `Fa>Adnotacje>PMarzy>P_PMarzy_3_1` | margin scheme for used goods |
| `Fa>Adnotacje>PMarzy>P_PMarzy_3_2` | margin scheme for works of art |
| `Fa>Adnotacje>PMarzy>P_PMarzy_3_3` | margin scheme for antiques and collectibles |
| `Fa>FaWiersz>P_11A` | price including tax (as opposite to `P_11` which does not include tax), to be used in cases described in appropriate law (see below) |
| `Fa>Platnosc>ZaplataCzesciowa>FormaPlatnosci` | Payment means for partial payment |
| `Fa>Platnosc>TerminPlatnosci>TerminOpis` | Alternative format for payment deadline. In our code there's `TerminPlatnosci>Termin` which contains a date, but this field can contain a textual description of the payment deadline (e.g. "within 30 days from the date of receiving invoice") |
| `Fa>OkresFa` | Period of the invoice - alternative to `P_6`. Mapped incorrectly in code (elements `P_6_Od` and `P_6_Do` should be wrapped in `OkresFa` tag) |
| `Fa>Adnotacje>Zwolnienie>P_19` | For tax exempt goods, set `P_19` to 1 |
| `Fa>Adnotacje>Zwolnienie>P_19A` | For tax exempt goods (Polish VAT law), text of legal basis |
| `Fa>Adnotacje>Zwolnienie>P_19B` | For tax exempt goods (directive 2006/112/EC), text of legal basis |
| `Fa>Adnotacje>Zwolnienie>P_19C` | For tax exempt goods (other legal basis), text of legal basis |
| `Fa>Zamowienie` | Information about the order (for advance invoices), contains total order value (`WartoscZamowienia`) and order line items (`ZamowienieWiersz`), required for `ZAL` (advance payment) type invoices |
| `Fa>P_15ZK` | For KOR_ZAL, amount to pay before correction, for other cases, remaining amount to pay before correction |

`WarunkiTransakcji` (transaction conditions) may contain (taken from example 4):
- `Umowy` - contract(s) date and number
- `Zamowienia` - order(s) date and number
- `NrPartiiTowaru` - product batch number
- `WarunkiDostawy` - conditions of delivery
- `Transport` - how the goods will be transported (contains many nested fields specifying e.g. transport company, destination address, etc.)

## Unset fields

The following fields are present in the struct to be converted to XML, but are not set anywhere in our code:

| XML field | Struct field | Notes |
| --------- | ------------ | ----- |
| `Fa>P_2` | `IssuePlace` | issue place - not required in schema |
| `Fa>P_14_1W` | `StandardRateTaxConvertedToPln` |
| `Fa>P_14_2W` | `ReducedRateTaxConvertedToPln` |
| `Fa>P_14_3W` | `SuperReducedRateTaxConvertedToPln` |
| `Fa>FP` | `FP` | indicates a case where an invoice is issued in addition to a regular receipt - not required in schema |
| `Fa>FaWiersz>StanPrzed` | `BeforeCorrectionMarker` | in a correction invoice, indicates that the line describes the state before the correction |
| `Fa>P_13_11` | `MarginNetSale` |
| `Fa>FaWiersz>GTU` | `SpecialGoodsCode` | Code identifying certain classes of goods and services (01 = alcoholic beverages, 02 = vehicle fuels...), values GTU_01 to GTU_13
| `Fa>P_6_Od` | Start of the invoice period. Mapped incorrectly in code (elements `P_6_Od` and `P_6_Do` should be wrapped in `OkresFa` tag) |
| `Fa>P_6_Do` | End of the invoice period. Mapped incorrectly in code (elements `P_6_Od` and `P_6_Do` should be wrapped in `OkresFa` tag) |

## How the listing is done

The following check is done by this way:
- Take [official sample XML files provided by the Polish authorities](https://www.gov.pl/attachment/937002fa-c6b5-477d-8b56-22105fa728c2)
- For each XML field in the reference file, find place in our code that generates it
- If a reference file contains a field and we don't have a corresponding field in our code, we mark it as "Not mapped".
- If a field is present in our struct, but there exist one possible value for this field, such a field goes to "Hardcoded values" list.
- If a field is present in our struct, but there exist no code that sets this field, such a field goes to "Unset values" list.

Note that this particular check:
- Does not cover all possible fields in the XML schema, only the ones present in the sample files
- Does not detect fields in our code that are not present in the sample files, or in the XML schema at all

### List of test cases

1. FP invoice (invoice added to receipt)
2. Correction invoice with `StanPrzed` - two lines, one with `StanPrzed` set to 1 (shows the state before the correction), another without (after correction)
3. Correction invoice - with a single line containing negative numbers (alternative method to case 2)
4. Regular invoice with a third party, with multiple delivery dates
5. Correction invoice where customer's data changed, contains `Podmiot2K` and `IDNabywcy` fields, and no line items (`FaWiersz`) - it's only for correcting customer's identification data and nothing else
6. Correction invoice to give a discount, for multiple earlier corrected invoices (`DaneFaKorygowanej`), contains `OkresFaKorygowanej` field, and no line items (`FaWiersz`)
7. Differs from example 6 by having a line item (`FaWiersz`) - meaning that the invoice corrects only part of the earlier deliveries, not all of them
8. Sale of used goods - uses margin scheme, contains partial payment, contains `P_PMarzy_3_1` field, contains `P_11A` field
9. Invoice to a local government unit, contains both tax exempt items (`P_13_7`) and standard tax items (`P_13_1`)
10. Advance invoice, already paid in full, with two customers, each having 50% share
11. Corrects invoice from example 10 because of incorrect used tax rate, type `KOR_ZAL`, contains negative numbers for the incorrect tax rate in fields `P_13_1` and `P_14_1`, contains `P_13_2` and `P_14_2` fields for the correct tax rate

## References

What is FP:
- http://www.vademecumpodatnika.pl/artykul_narzedziowa,1190,0,19945,ujmowanie-w-ewidencji-vat-faktur-do-paragonow.html

GTU codes:
- https://www.podatki.biz/artykuly/jpkvat-z-deklaracja-oznaczenia-dostawy-i-swiadczenia-uslug-gtu_4_45232.htm

When should `P_11A` be used - art. 106e ust. 7 i 8:
- https://sip.lex.pl/akty-prawne/dzu-dziennik-ustaw/podatek-od-towarow-i-uslug-17086198/art-106-e

When should `P_23` be used:
- https://isp-modzelewski.pl/serwis/wewnatrzwspolnotowe-transakcje-trojstronne/#:~:text=ramach%20procedury%20uproszczonej.-,Zgodnie%20z%20art.,dodanej%20ostatniego%20w%20kolejno%C5%9Bci%20podatnika.
- https://www.krgroup.pl/procedura-uproszczona-rozliczenia-vat-w-wewnatrzwspolnotowej-transakcji-trojstronnej/#:~:text=Warunki%20zastosowania%20procedury,realizowanej%20w%20ramach%20procedury%20uproszczonej.
- https://sip.lex.pl/akty-prawne/dzu-dziennik-ustaw/podatek-od-towarow-i-uslug-17086198/dz-12-roz-8

When should `P_15ZK` be used (art. 106f ust. 3) for remaining amount to pay before correction:
https://sip.lex.pl/akty-prawne/dzu-dziennik-ustaw/podatek-od-towarow-i-uslug-17086198/art-106-f