# KhaleejiAPI Go SDK

Official Go SDK for [KhaleejiAPI](https://khaleejiapi.dev) — the MENA region's developer API platform.

## Requirements

- Go 1.22+

## Installation

```bash
go get github.com/khaleejiapi/sdk-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    khaleejiapi "github.com/khaleejiapi/sdk-go"
)

func main() {
    client := khaleejiapi.New("kapi_live_your_key_here")
    ctx := context.Background()

    // Validate an email
    email, err := client.Validation.ValidateEmail(ctx, "user@example.com")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Valid:", email.Valid)

    // Get prayer times
    prayers, err := client.Islamic.GetPrayerTimes(ctx, khaleejiapi.PrayerTimesParams{
        City: "Dubai",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Fajr:", prayers.Prayers.Fajr)

    // Exchange rates
    rates, err := client.Finance.GetExchangeRates(ctx, "AED", "USD,EUR,SAR")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Rates:", rates.Rates)
}
```

## API Reference

### Validation

```go
ctx := context.Background()

// Email validation
email, _ := client.Validation.ValidateEmail(ctx, "user@example.com")

// Phone validation
phone, _ := client.Validation.ValidatePhone(ctx, "+971501234567", "AE")

// IBAN validation
iban, _ := client.Validation.ValidateIBAN(ctx, "AE070331234567890123456")

// VAT/TRN validation
vat, _ := client.Validation.ValidateVAT(ctx, "100123456700003")

// Emirates ID validation
eid, _ := client.Validation.ValidateEmiratesID(ctx, "784-1990-1234567-1")

// Saudi ID validation
sid, _ := client.Validation.ValidateSaudiID(ctx, "1012345678")

// Saudi ID batch validation (max 100)
batch, _ := client.Validation.ValidateSaudiIDBatch(ctx, []string{"1012345678", "2098765432"})
```

### Geolocation

```go
// IP geolocation
ip, _ := client.Geo.IPLookup(ctx, "8.8.8.8")

// Timezone lookup
tz, _ := client.Geo.GetTimezone(ctx, "Dubai")

// Geocoding
geo, _ := client.Geo.Geocode(ctx, "Burj Khalifa, Dubai")
```

### Finance

```go
// Exchange rates
rates, _ := client.Finance.GetExchangeRates(ctx, "AED", "USD,EUR")

// VAT calculation
vat, _ := client.Finance.CalculateVAT(ctx, 100.0, "AE", false)

// Public holidays
holidays, _ := client.Finance.GetHolidays(ctx, "AE", 2026)

// Business days
days, _ := client.Finance.GetBusinessDays(ctx, khaleejiapi.BusinessDaysParams{
    Country: "AE",
    From:    "2026-01-01",
    To:      "2026-01-31",
})
```

### Communication

```go
// AI Translation
translation, _ := client.Communication.Translate(ctx, khaleejiapi.TranslateParams{
    Text:    "Hello, world!",
    Target:  "ar",
    Dialect: "gulf",
})
```

### Islamic

```go
// Hijri calendar conversion
hijri, _ := client.Islamic.ConvertHijri(ctx, khaleejiapi.ConvertHijriParams{Today: true})

// Prayer times
prayers, _ := client.Islamic.GetPrayerTimes(ctx, khaleejiapi.PrayerTimesParams{
    City:   "Mecca",
    Method: "umm_al_qura",
})

// Arabic text processing
arabic, _ := client.Islamic.ProcessArabic(ctx, khaleejiapi.ProcessArabicParams{
    Text:      "بِسْمِ اللَّهِ الرَّحْمَنِ الرَّحِيمِ",
    Operation: "removeDiacritics",
})
```

### Utility

```go
// Weather
weather, _ := client.Utility.GetWeather(ctx, "Dubai")

// Fraud check
fraud, _ := client.Utility.FraudCheck(ctx, khaleejiapi.FraudCheckParams{
    Email: "test@example.com",
    IP:    "1.2.3.4",
})

// URL shortener
short, _ := client.Utility.ShortenURL(ctx, "https://example.com/long-url", "")
```

## Configuration

```go
// Simple initialization
client := khaleejiapi.New("kapi_live_your_key")

// Full configuration
client := khaleejiapi.NewWithConfig(khaleejiapi.Config{
    APIKey:     "kapi_live_your_key",
    BaseURL:    "https://khaleejiapi.dev/api/v1",
    Timeout:    30 * time.Second,
    MaxRetries: 2,
})

// With custom HTTP client
client := khaleejiapi.NewWithConfig(khaleejiapi.Config{
    APIKey:     "kapi_live_your_key",
    HTTPClient: &http.Client{Timeout: 60 * time.Second},
})
```

## Error Handling

```go
result, err := client.Validation.ValidateEmail(ctx, "test@example.com")
if err != nil {
    var apiErr *khaleejiapi.APIError
    if errors.As(err, &apiErr) {
        switch apiErr.StatusCode {
        case 401:
            fmt.Println("Check your API key")
        case 429:
            fmt.Printf("Rate limited. Retry after %d seconds\n", apiErr.RateLimitInfo.Reset)
        default:
            fmt.Printf("Error %d: %s\n", apiErr.StatusCode, apiErr.Message)
        }
    }
}
```

## Context Support

All methods accept a `context.Context` for cancellation and timeout:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

result, err := client.Validation.ValidateEmail(ctx, "test@example.com")
```

## License

MIT — See [LICENSE](LICENSE) for details.

## Links

- [Documentation](https://khaleejiapi.dev/docs)
- [API Reference](https://khaleejiapi.dev/docs/v1)
- [Dashboard](https://khaleejiapi.dev/dashboard)
- [Status](https://khaleejiapi.dev/status)
