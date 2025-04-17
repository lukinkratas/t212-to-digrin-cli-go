package dataframe

import (
	"slices"

	"github.com/gocarina/gocsv"
)


type Schema struct {
	Action                               string  `csv:"Action"`
	Time                                 string  `csv:"Time"`
	ISIN                                 string  `csv:"ISIN"`
	Ticker                               string  `csv:"Ticker"`
	Name                                 string  `csv:"Name"`
	Notes                                string  `csv:"Notes"`
	Id                                   string  `csv:"ID"`
	NoOfShares                           float64 `csv:"No. of shares"`
	PricePerShare                        float64 `csv:"Price / share"`
	CurrencyPricePerShare                string  `csv:"Currency (Price / share)"`
	ExchangeRate                         string  `csv:"Exchange rate"`
	CurrencyResult                       string  `csv:"Currency (Result)"`
	Total                                float64 `csv:"Total"`
	CurrencyTotal                        string  `csv:"Currency (Total)"`
	WithholdingTax                       float64 `csv:"Withholding tax"`
	CurrencyWithholdingTax               string  `csv:"Currency (Withholding tax)"`
	CurrencyConversionFromAmount         float64 `csv:"Currency conversion from amount"`
	CurrencyCurrencyConversionFromAmount string  `csv:"Currency (Currency conversion from amount)"`
	CurrencyConversionToAmount           float64 `csv:"Currency conversion to amount"`
	CurrencyCurrencyConversionToAmount   string  `csv:"Currency (Currency conversion to amount)"`
	CurrencyConversionFee                float64 `csv:"Currency conversion fee"`
	CurrencyCurrencyConversionFee        string  `csv:"Currency (Currency conversion fee)"`
	FrenchTransactionTax                 float64 `csv:"French transaction tax"`
	CurrencyFrenchTransactionTax         string  `csv:"Currency (French transaction tax)"`
}

func DecodeCSV(csvEncoded []byte) []Schema {

	var dataFrame []Schema

	err := gocsv.UnmarshalBytes(csvEncoded, &dataFrame)
	if err != nil {
		panic(err)
	}

	return dataFrame
}

func Transform(dataFrame []Schema) []Schema {

	// Filter out blacklisted tickers
	tickerBlacklist := []string{
		"VNTRF",  // due to stock split
		"BRK.A",  // not available in digrin
	}
	
    dataFrame = slices.DeleteFunc(dataFrame, func(dataFrameRow Schema) bool {
		return slices.Contains(tickerBlacklist, dataFrameRow.Ticker)
    })
	
	// Filter only buys and sells
	allowedActions := []string{"Market buy", "Market sell"}
	
    dataFrame = slices.DeleteFunc(dataFrame, func(dataFrameRow Schema) bool {
        return !slices.Contains(allowedActions, dataFrameRow.Action)
    })

	// Apply the mapping to the ticker column
	tickerMap := map[string]string{
		"VWCE": "VWCE.DE",
		"VUAA": "VUAA.DE",
        "SXRV": "SXRV.DE",
        "ZPRV": "ZPRV.DE",
        "ZPRX": "ZPRX.DE",
        "MC":   "MC.PA",
        "ASML": "ASML.AS",
        "CSPX": "CSPX.L",
        "EISU": "EISU.L",
        "IITU": "IITU.L",
        "IUHC": "IUHC.L",
        "NDIA": "NDIA.L",
	}

	for _, dataFrameRow := range dataFrame {

		tickerSubstitute, ok := tickerMap[dataFrameRow.Ticker]
		if ok {
			dataFrameRow.Ticker = tickerSubstitute
		}

	}

	return dataFrame
}

func Encode(dataFrame []Schema) []byte {

	csvEncoded, err := gocsv.MarshalBytes(dataFrame)
	if err != nil {
		panic(err)
	}

	return csvEncoded
}