package khaleejiapi

import "context"

// ValidationResource provides access to validation APIs.
type ValidationResource struct {
	client *Client
}

// EmailResult represents an email validation result.
type EmailResult struct {
	Valid      bool        `json:"valid"`
	Email      string      `json:"email"`
	Checks     EmailChecks `json:"checks,omitempty"`
	Suggestion string      `json:"suggestion,omitempty"`
	Domain     string      `json:"domain,omitempty"`
}

type EmailChecks struct {
	Format     *bool `json:"format,omitempty"`
	Syntax     *bool `json:"syntax,omitempty"`
	MX         bool  `json:"mx"`
	Disposable bool  `json:"disposable"`
	Role       bool  `json:"role"`
}

// ValidateEmail validates an email address.
func (r *ValidationResource) ValidateEmail(ctx context.Context, email string) (*EmailResult, error) {
	return doGet[*EmailResult](r.client, ctx, "/email/validate", map[string]string{"email": email})
}

// PhoneResult represents a phone validation result.
type PhoneResult struct {
	Valid       bool   `json:"valid"`
	Phone       string `json:"phone"`
	Formatted   string `json:"formatted,omitempty"`
	CountryCode string `json:"countryCode,omitempty"`
	Type        string `json:"type,omitempty"`
	Carrier     string `json:"carrier,omitempty"`
}

// ValidatePhone validates a phone number.
func (r *ValidationResource) ValidatePhone(ctx context.Context, phone string, country string) (*PhoneResult, error) {
	params := map[string]string{"phone": phone}
	if country != "" {
		params["country"] = country
	}
	return doGet[*PhoneResult](r.client, ctx, "/phone/validate", params)
}

// IBANResult represents an IBAN validation result.
type IBANResult struct {
	Valid    bool   `json:"valid"`
	IBAN     string `json:"iban"`
	Country  string `json:"country,omitempty"`
	BankName string `json:"bankName,omitempty"`
	BankCode string `json:"bankCode,omitempty"`
}

// ValidateIBAN validates an IBAN.
func (r *ValidationResource) ValidateIBAN(ctx context.Context, iban string) (*IBANResult, error) {
	return doGet[*IBANResult](r.client, ctx, "/iban/validate", map[string]string{"iban": iban})
}

// VATResult represents a VAT/TRN validation result.
type VATResult struct {
	Valid   bool   `json:"valid"`
	TRN     string `json:"trn"`
	Country string `json:"country,omitempty"`
}

// ValidateVAT validates a VAT/TRN number.
func (r *ValidationResource) ValidateVAT(ctx context.Context, trn string) (*VATResult, error) {
	return doGet[*VATResult](r.client, ctx, "/vat/validate", map[string]string{"trn": trn})
}

// EmiratesIDResult represents an Emirates ID validation result.
type EmiratesIDResult struct {
	Valid           bool   `json:"valid"`
	ID              string `json:"id"`
	NationalityCode string `json:"nationalityCode,omitempty"`
	BirthYear       int    `json:"birthYear,omitempty"`
}

// ValidateEmiratesID validates a UAE Emirates ID.
func (r *ValidationResource) ValidateEmiratesID(ctx context.Context, id string) (*EmiratesIDResult, error) {
	return doGet[*EmiratesIDResult](r.client, ctx, "/validation/emirates-id", map[string]string{"id": id})
}

// SaudiIDResult represents a Saudi ID validation result.
type SaudiIDResult struct {
	ID          string   `json:"id"`
	Valid       bool     `json:"valid"`
	Type        string   `json:"type,omitempty"`
	TypeAr      string   `json:"typeAr,omitempty"`
	Nationality string   `json:"nationality,omitempty"`
	Errors      []string `json:"errors,omitempty"`
}

// SaudiIDBatchResult represents a batch Saudi ID validation result.
type SaudiIDBatchResult struct {
	Results []SaudiIDResult `json:"results"`
	Summary BatchSummary    `json:"summary"`
}

// BatchSummary contains summary counts for batch operations.
type BatchSummary struct {
	Total   int `json:"total"`
	Valid   int `json:"valid"`
	Invalid int `json:"invalid"`
}

// ValidateSaudiID validates a Saudi National ID or Iqama.
func (r *ValidationResource) ValidateSaudiID(ctx context.Context, id string) (*SaudiIDResult, error) {
	return doGet[*SaudiIDResult](r.client, ctx, "/validation/saudi-id", map[string]string{"id": id})
}

// ValidateSaudiIDBatch validates multiple Saudi IDs (max 100).
func (r *ValidationResource) ValidateSaudiIDBatch(ctx context.Context, ids []string) (*SaudiIDBatchResult, error) {
	body := struct {
		IDs []string `json:"ids"`
	}{IDs: ids}
	return doPost[*SaudiIDBatchResult](r.client, ctx, "/validation/saudi-id", body)
}

// GeoResource provides access to geolocation APIs.
type GeoResource struct {
	client *Client
}

// IPResult represents an IP lookup result.
type IPResult struct {
	IP          string  `json:"ip"`
	Country     string  `json:"country,omitempty"`
	CountryName string  `json:"countryName,omitempty"`
	City        string  `json:"city,omitempty"`
	Latitude    float64 `json:"latitude,omitempty"`
	Longitude   float64 `json:"longitude,omitempty"`
	ISP         string  `json:"isp,omitempty"`
}

// IPLookup looks up IP geolocation data.
func (r *GeoResource) IPLookup(ctx context.Context, ip string) (*IPResult, error) {
	params := map[string]string{}
	if ip != "" {
		params["ip"] = ip
	}
	return doGet[*IPResult](r.client, ctx, "/ip/lookup", params)
}

// TimezoneResult represents a timezone lookup result.
type TimezoneResult struct {
	Location    string `json:"location,omitempty"`
	Timezone    string `json:"timezone"`
	UTCOffset   string `json:"utcOffset,omitempty"`
	DSTActive   bool   `json:"dstActive,omitempty"`
	CurrentTime string `json:"currentTime,omitempty"`
}

// GetTimezone gets timezone data for a location.
func (r *GeoResource) GetTimezone(ctx context.Context, location string) (*TimezoneResult, error) {
	return doGet[*TimezoneResult](r.client, ctx, "/timezone", map[string]string{"location": location})
}

// GeocodeResult represents a geocoding result.
type GeocodeResult struct {
	Address   string  `json:"address,omitempty"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Country   string  `json:"country,omitempty"`
}

// Geocode converts an address to coordinates.
func (r *GeoResource) Geocode(ctx context.Context, address string) (*GeocodeResult, error) {
	return doGet[*GeocodeResult](r.client, ctx, "/geocode", map[string]string{"address": address})
}

// FinanceResource provides access to finance APIs.
type FinanceResource struct {
	client *Client
}

// ExchangeRatesResult represents exchange rate data.
type ExchangeRatesResult struct {
	Base   string             `json:"base"`
	Rates  map[string]float64 `json:"rates"`
	Source string             `json:"source,omitempty"`
}

// GetExchangeRates gets exchange rates.
func (r *FinanceResource) GetExchangeRates(ctx context.Context, base string, symbols string) (*ExchangeRatesResult, error) {
	params := map[string]string{"base": base}
	if symbols != "" {
		params["symbols"] = symbols
	}
	return doGet[*ExchangeRatesResult](r.client, ctx, "/exchange/rates", params)
}

// VATCalcResult represents a VAT calculation result.
type VATCalcResult struct {
	Country     string  `json:"country"`
	VATRate     float64 `json:"vatRate"`
	InputAmount float64 `json:"inputAmount"`
	BaseAmount  float64 `json:"baseAmount"`
	VATAmount   float64 `json:"vatAmount"`
	TotalAmount float64 `json:"totalAmount"`
	Currency    string  `json:"currency"`
}

// CalculateVAT calculates VAT.
func (r *FinanceResource) CalculateVAT(ctx context.Context, amount float64, country string, inclusive bool) (*VATCalcResult, error) {
	params := map[string]string{
		"amount":    fmt.Sprintf("%g", amount),
		"country":   country,
		"inclusive": fmt.Sprintf("%t", inclusive),
	}
	return doGet[*VATCalcResult](r.client, ctx, "/vat/calculate", params)
}

// HolidaysResult represents public holidays data.
type HolidaysResult struct {
	Country  string    `json:"country"`
	Year     int       `json:"year"`
	Holidays []Holiday `json:"holidays"`
}

// Holiday represents a single public holiday.
type Holiday struct {
	Name   string `json:"name"`
	NameAr string `json:"nameAr,omitempty"`
	Date   string `json:"date"`
	Type   string `json:"type"`
}

// GetHolidays gets public holidays for a GCC country.
func (r *FinanceResource) GetHolidays(ctx context.Context, country string, year int) (*HolidaysResult, error) {
	params := map[string]string{"country": country}
	if year > 0 {
		params["year"] = fmt.Sprintf("%d", year)
	}
	return doGet[*HolidaysResult](r.client, ctx, "/holidays", params)
}

// BusinessDaysResult represents business days calculation.
type BusinessDaysResult struct {
	BusinessDays  *int   `json:"businessDays,omitempty"`
	TotalDays     *int   `json:"totalDays,omitempty"`
	IsBusinessDay *bool  `json:"isBusinessDay,omitempty"`
	From          string `json:"from,omitempty"`
	To            string `json:"to,omitempty"`
	ResultDate    string `json:"resultDate,omitempty"`
}

// BusinessDaysParams contains parameters for the business days API.
type BusinessDaysParams struct {
	Country string
	Date    string
	From    string
	To      string
	Add     int
}

// GetBusinessDays calculates business days.
func (r *FinanceResource) GetBusinessDays(ctx context.Context, p BusinessDaysParams) (*BusinessDaysResult, error) {
	params := map[string]string{}
	if p.Country != "" {
		params["country"] = p.Country
	}
	if p.Date != "" {
		params["date"] = p.Date
	}
	if p.From != "" {
		params["from"] = p.From
	}
	if p.To != "" {
		params["to"] = p.To
	}
	if p.Add != 0 {
		params["add"] = fmt.Sprintf("%d", p.Add)
	}
	return doGet[*BusinessDaysResult](r.client, ctx, "/business-days", params)
}

// CommunicationResource provides access to communication APIs.
type CommunicationResource struct {
	client *Client
}

// TranslationResult represents a translation result.
type TranslationResult struct {
	Text             string `json:"text,omitempty"`
	Translated       string `json:"translated,omitempty"`
	From             string `json:"from,omitempty"`
	To               string `json:"to,omitempty"`
	DetectedLanguage string `json:"detectedLanguage,omitempty"`
}

// TranslateParams contains parameters for the translate API.
type TranslateParams struct {
	Text      string `json:"text"`
	Target    string `json:"target"`
	Source    string `json:"source,omitempty"`
	Formality string `json:"formality,omitempty"`
	Dialect   string `json:"dialect,omitempty"`
}

// Translate translates text using AI (Google Gemini).
func (r *CommunicationResource) Translate(ctx context.Context, p TranslateParams) (*TranslationResult, error) {
	return doPost[*TranslationResult](r.client, ctx, "/translate", p)
}

// IslamicResource provides access to Islamic APIs.
type IslamicResource struct {
	client *Client
}

// HijriResult represents a Hijri conversion result.
type HijriResult struct {
	Gregorian HijriDate `json:"gregorian"`
	Hijri     HijriDate `json:"hijri"`
	Direction string    `json:"direction"`
}

// HijriDate represents a date in either calendar.
type HijriDate struct {
	Date        string `json:"date"`
	Year        int    `json:"year"`
	Month       int    `json:"month"`
	Day         int    `json:"day"`
	MonthName   string `json:"monthName,omitempty"`
	MonthNameAr string `json:"monthNameAr,omitempty"`
	DayOfWeek   string `json:"dayOfWeek,omitempty"`
	DayOfWeekAr string `json:"dayOfWeekAr,omitempty"`
}

// ConvertHijriParams contains params for Hijri conversion.
type ConvertHijriParams struct {
	Date  string
	Hijri string
	Today bool
}

// ConvertHijri converts between Gregorian and Hijri calendars.
func (r *IslamicResource) ConvertHijri(ctx context.Context, p ConvertHijriParams) (*HijriResult, error) {
	params := map[string]string{}
	if p.Date != "" {
		params["date"] = p.Date
	}
	if p.Hijri != "" {
		params["hijri"] = p.Hijri
	}
	if p.Today {
		params["today"] = "true"
	}
	return doGet[*HijriResult](r.client, ctx, "/hijri/convert", params)
}

// PrayerTimesResult represents prayer times data.
type PrayerTimesResult struct {
	Location *PrayerLocation `json:"location,omitempty"`
	Date     string          `json:"date"`
	Prayers  Prayers         `json:"prayers"`
	Qibla    *Qibla          `json:"qibla,omitempty"`
	Method   *PrayerMethod   `json:"method,omitempty"`
	School   string          `json:"school,omitempty"`
}

type PrayerLocation struct {
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	City    string  `json:"city,omitempty"`
	Country string  `json:"country,omitempty"`
}

type Prayers struct {
	Fajr    string `json:"fajr"`
	Sunrise string `json:"sunrise"`
	Dhuhr   string `json:"dhuhr"`
	Asr     string `json:"asr"`
	Maghrib string `json:"maghrib"`
	Isha    string `json:"isha"`
}

type Qibla struct {
	Direction float64 `json:"direction"`
	Compass   string  `json:"compass,omitempty"`
}

type PrayerMethod struct {
	Name string `json:"name"`
}

// PrayerTimesParams contains parameters for the prayer times API.
type PrayerTimesParams struct {
	City   string
	Lat    float64
	Lng    float64
	Date   string
	Method string // default: "mwl"
	School string // default: "shafi"
}

// GetPrayerTimes gets prayer times for a location.
func (r *IslamicResource) GetPrayerTimes(ctx context.Context, p PrayerTimesParams) (*PrayerTimesResult, error) {
	params := map[string]string{}
	if p.City != "" {
		params["city"] = p.City
	}
	if p.Lat != 0 {
		params["lat"] = fmt.Sprintf("%f", p.Lat)
	}
	if p.Lng != 0 {
		params["lng"] = fmt.Sprintf("%f", p.Lng)
	}
	if p.Date != "" {
		params["date"] = p.Date
	}
	if p.Method != "" {
		params["method"] = p.Method
	} else {
		params["method"] = "mwl"
	}
	if p.School != "" {
		params["school"] = p.School
	} else {
		params["school"] = "shafi"
	}
	return doGet[*PrayerTimesResult](r.client, ctx, "/prayer-times", params)
}

// ArabicResult represents an Arabic text processing result.
type ArabicResult struct {
	Operation    string   `json:"operation"`
	Original     string   `json:"original,omitempty"`
	Text         string   `json:"text,omitempty"`
	Result       string   `json:"result,omitempty"`
	RemovedCount int      `json:"removedCount,omitempty"`
	Count        int      `json:"count,omitempty"`
	Words        []string `json:"words,omitempty"`
	Score        float64  `json:"score,omitempty"`
	Label        string   `json:"label,omitempty"`
	Script       string   `json:"script,omitempty"`
}

// ProcessArabicParams contains parameters for Arabic text processing.
type ProcessArabicParams struct {
	Text      string `json:"text"`
	Operation string `json:"operation"`
	Direction string `json:"direction,omitempty"`
}

// ProcessArabic processes Arabic text.
func (r *IslamicResource) ProcessArabic(ctx context.Context, p ProcessArabicParams) (*ArabicResult, error) {
	body := map[string]any{
		"text":      p.Text,
		"operation": p.Operation,
	}
	if p.Direction != "" {
		body["options"] = map[string]string{"direction": p.Direction}
	}
	return doPost[*ArabicResult](r.client, ctx, "/arabic/process", body)
}

// UtilityResource provides access to utility APIs.
type UtilityResource struct {
	client *Client
}

// WeatherResult represents weather data.
type WeatherResult struct {
	City        string  `json:"city,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
	Humidity    int     `json:"humidity,omitempty"`
	Condition   string  `json:"condition,omitempty"`
	WindSpeed   float64 `json:"windSpeed,omitempty"`
}

// GetWeather gets weather for a city.
func (r *UtilityResource) GetWeather(ctx context.Context, city string) (*WeatherResult, error) {
	return doGet[*WeatherResult](r.client, ctx, "/weather", map[string]string{"city": city})
}

// FraudResult represents a fraud check result.
type FraudResult struct {
	RiskScore      int    `json:"riskScore"`
	RiskLevel      string `json:"riskLevel"`
	Recommendation string `json:"recommendation,omitempty"`
}

// FraudCheckParams contains parameters for fraud checking.
type FraudCheckParams struct {
	IP    string `json:"ip,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

// FraudCheck checks for fraud.
func (r *UtilityResource) FraudCheck(ctx context.Context, p FraudCheckParams) (*FraudResult, error) {
	return doPost[*FraudResult](r.client, ctx, "/fraud/check", p)
}

// ShortenResult represents a URL shortening result.
type ShortenResult struct {
	Code        string `json:"code"`
	ShortURL    string `json:"shortUrl"`
	OriginalURL string `json:"originalUrl"`
}

// ShortenURL shortens a URL.
func (r *UtilityResource) ShortenURL(ctx context.Context, url string, customCode string) (*ShortenResult, error) {
	body := map[string]string{"url": url}
	if customCode != "" {
		body["customCode"] = customCode
	}
	return doPost[*ShortenResult](r.client, ctx, "/url/shorten", body)
}
