// Package cabrillo implements the Cabrillo V3 file format as defined by the [WWROF]
//
// [WWROF] https://wwrof.org/cabrillo/
package cabrillo

import (
	"strconv"
	"strings"
	"time"

	"github.com/ftl/hamradio/callsign"
	"github.com/ftl/hamradio/locator"
)

const TimestampLayout = "2006-01-02 1504"

const (
	Version2 = "2.0"
	Version3 = "3.0"
)

func NewLog() *Log {
	return &Log{
		CabrilloVersion: Version3,
		Custom:          make(map[Tag]string),
		QSOData:         []QSO{},
		IgnoredQSOs:     []QSO{},
	}
}

type Log struct {
	CabrilloVersion string
	Callsign        callsign.Callsign
	Contest         ContestIdentifier
	Category        Category
	Certificate     bool
	ClaimedScore    int
	Club            string
	CreatedBy       string
	Email           string
	GridLocator     locator.Locator
	Location        string
	Name            string
	Address         Address
	Operators       []callsign.Callsign
	Host            callsign.Callsign
	Offtime         Offtime
	Soapbox         string
	Debug           int
	Custom          map[Tag]string
	QSOData         []QSO
	IgnoredQSOs     []QSO
}

type Tag string

func (t Tag) IsCustom() bool {
	if t == XQSOTag {
		return false
	}
	return strings.HasPrefix(string(t), XPrefix)
}

const (
	StartOfLogTag           Tag = "START-OF-LOG"
	EndOfLogTag             Tag = "END-OF-LOG"
	CallsignTag             Tag = "CALLSIGN"
	ContestTag              Tag = "CONTEST"
	CategoryAssistedTag     Tag = "CATEGORY-ASSISTED"
	CategoryBandTag         Tag = "CATEGORY-BAND"
	CategoryModeTag         Tag = "CATEGORY-MODE"
	CategoryOperatorTag     Tag = "CATEGORY-OPERATOR"
	CategoryPowerTag        Tag = "CATEGORY-POWER"
	CategoryStationTag      Tag = "CATEGORY-STATION"
	CategoryTimeTag         Tag = "CATEGORY-TIME"
	CategoryTransmitterTag  Tag = "CATEGORY-TRANSMITTER"
	CategoryOverlayTag      Tag = "CATEGORY-OVERLAY"
	CertificateTag          Tag = "CERTIFICATE"
	ClaimedScoreTag         Tag = "CLAIMED-SCORE"
	ClubTag                 Tag = "CLUB"
	CreatedByTag            Tag = "CREATED-BY"
	EmailTag                Tag = "EMAIL"
	GridLocatorTag          Tag = "GRID-LOCATOR"
	LocationTag             Tag = "LOCATION"
	NameTag                 Tag = "NAME"
	AddressTag              Tag = "ADDRESS"
	AddressCityTag          Tag = "ADDRESS-CITY"
	AddressStateProvinceTag Tag = "ADDRESS-STATE-PROVINCE"
	AddressPostalcodeTag    Tag = "ADDRESS-POSTALCODE"
	AddressCountryTag       Tag = "ADDRESS-COUNTRY"
	OperatorsTag            Tag = "OPERATORS"
	OfftimeTag              Tag = "OFFTIME"
	SoapboxTag              Tag = "SOAPBOX"
	QSOTag                  Tag = "QSO"
	XQSOTag                 Tag = "X-QSO"
	XPrefix                     = "X-"
)

type ContestIdentifier string

type Category struct {
	Assisted    CategoryAssisted
	Band        CategoryBand
	Mode        CategoryMode
	Operator    CategoryOperator
	Power       CategoryPower
	Station     CategoryStation
	Time        CategoryTime
	Transmitter CategoryTransmitter
	Overlay     CategoryOverlay
}

type CategoryAssisted string

func (c CategoryAssisted) Bool() bool {
	return c == Assisted
}

const (
	Assisted    CategoryAssisted = "ASSISTED"
	NonAssisted CategoryAssisted = "NON-ASSISTED"
)

type CategoryBand string

const (
	BandAll        CategoryBand = "ALL"
	Band160m       CategoryBand = "160M"
	Band80m        CategoryBand = "80M"
	Band40m        CategoryBand = "40M"
	Band20m        CategoryBand = "20M"
	Band15m        CategoryBand = "15M"
	Band10m        CategoryBand = "10M"
	Band6m         CategoryBand = "6M"
	Band4m         CategoryBand = "4M"
	Band2m         CategoryBand = "2M"
	Band222        CategoryBand = "222"
	Band432        CategoryBand = "432"
	Band902        CategoryBand = "902"
	Band1_2G       CategoryBand = "1.2G"
	Band2_3G       CategoryBand = "2.3G"
	Band3_4G       CategoryBand = "3.4G"
	Band5_6G       CategoryBand = "5.7G"
	Band10G        CategoryBand = "10G"
	Band24G        CategoryBand = "24G"
	Band47G        CategoryBand = "47G"
	Band75G        CategoryBand = "75G"
	Band122G       CategoryBand = "122G"
	Band134G       CategoryBand = "134G"
	Band241G       CategoryBand = "241G"
	BandLight      CategoryBand = "LIGHT"
	BandVHF_3Band  CategoryBand = "VHF-3-BAND"
	BandVHF_FMOnly CategoryBand = "VHF-FM-ONLY"
)

type CategoryMode string

const (
	ModeCW    CategoryMode = "CW"
	ModeDIGI  CategoryMode = "DIGI"
	ModeFM    CategoryMode = "FM"
	ModeRTTY  CategoryMode = "RTTY"
	ModeSSB   CategoryMode = "SSB"
	ModeMIXED CategoryMode = "MIXED"
)

type CategoryOperator string

const (
	SingleOperator CategoryOperator = "SINGLE-OP"
	MultiOperator  CategoryOperator = "MULTI-OP"
	Checklog       CategoryOperator = "CHECKLOG"
)

type CategoryPower string

const (
	HighPower CategoryPower = "HIGH"
	LowPower  CategoryPower = "LOW"
	QRP       CategoryPower = "QRP"
)

type CategoryStation string

const (
	DistributedStation    CategoryStation = "DISTRIBUTED"
	FixedStation          CategoryStation = "FIXED"
	MobileStation         CategoryStation = "MOBILE"
	PortableStation       CategoryStation = "PORTABLE"
	RoverStation          CategoryStation = "ROVER"
	RoverLimitedStation   CategoryStation = "ROVER-LIMITED"
	RoverUnlimitedStation CategoryStation = "ROVER-UNLIMITED"
	ExpeditionStation     CategoryStation = "EXPEDITION"
	HQStation             CategoryStation = "HQ"
	SchoolStation         CategoryStation = "SCHOOL"
	ExplorerStation       CategoryStation = "EXPLORER"
)

type CategoryTime string

var durations = map[CategoryTime]time.Duration{
	Hours6:  6 * time.Hour,
	Hours8:  8 * time.Hour,
	Hours12: 12 * time.Hour,
	Hours24: 24 * time.Hour,
}

func (c CategoryTime) Duration() time.Duration {
	return durations[c]
}

const (
	Hours6  CategoryTime = "6-HOURS"
	Hours8  CategoryTime = "8-HOURS"
	Hours12 CategoryTime = "12-HOURS"
	Hours24 CategoryTime = "24-HOURS"
)

type CategoryTransmitter string

const (
	OneTransmitter       CategoryTransmitter = "ONE"
	TwoTransmitter       CategoryTransmitter = "TWO"
	LimitedTransmitter   CategoryTransmitter = "LIMITED"
	UnlimitedTransmitter CategoryTransmitter = "UNLIMITED"
	SWL                  CategoryTransmitter = "SWL"
)

type CategoryOverlay string

const (
	ClassicOverlay    CategoryOverlay = "CLASSIC"
	RookieOverlay     CategoryOverlay = "ROOKIE"
	TBWiresOverlay    CategoryOverlay = "TB-WIRES"
	YouthOverlay      CategoryOverlay = "YOUTH"
	NoviceTechOverlay CategoryOverlay = "NOVICE-TECH"
	Over50Overlay     CategoryOverlay = "OVER-50"
)

type Address struct {
	Text          string
	City          string
	StateProvince string
	Postalcode    string
	Country       string
}

type Offtime struct {
	Begin time.Time
	End   time.Time
}

func (o Offtime) Duration() time.Duration {
	return o.End.Sub(o.Begin)
}

type QSO struct {
	Frequency   QSOFrequency
	Mode        QSOMode
	Timestamp   time.Time
	Sent        QSOInfo
	Received    QSOInfo
	Transmitter int
}

type QSOFrequency string

func (f QSOFrequency) IsFrequency() bool {
	return f.ToKilohertz() != 0
}

func (f QSOFrequency) ToKilohertz() int {
	kHz, err := strconv.Atoi(string(f))
	if err != nil {
		return 0
	}
	if kHz < 1800 {
		return 0
	}
	return kHz
}

func (f QSOFrequency) ToBand() CategoryBand {
	if f.IsFrequency() {
		kHz := f.ToKilohertz()
		switch {
		case kHz < 3_500:
			return Band160m
		case kHz < 7_000:
			return Band80m
		case kHz < 14_000:
			return Band40m
		case kHz < 21_000:
			return Band20m
		case kHz < 28_000:
			return Band15m
		case kHz < 50_000_000:
			return Band10m
		case kHz < 70_000_000:
			return Band6m
		case kHz < 144_000_000:
			return Band4m
		case kHz < 222_000_000:
			return Band2m
		case kHz < 430_000_000:
			return Band222
		case kHz < 900_000_000:
			return Band432
		case kHz < 1_200_000_000:
			return Band902
		case kHz < 2_300_000_000:
			return Band1_2G
		default:
			return Band2_3G
		}
	}
	switch f {
	case Frequency50MHz:
		return Band6m
	case Frequency70MHz:
		return Band4m
	case Frequency144MHz:
		return Band2m
	default:
		return CategoryBand(strings.ToUpper(string(f)))
	}
}

const (
	Frequency50MHz  QSOFrequency = "50"
	Frequency70MHz  QSOFrequency = "70"
	Frequency144MHz QSOFrequency = "144"
	Frequency222MHz QSOFrequency = "222"
	Frequency432MHz QSOFrequency = "432"
	Frequency902MHz QSOFrequency = "902"
	Frequency1_2GHz QSOFrequency = "1.2G"
	Frequency2_3GHz QSOFrequency = "2.3G"
	Frequency3_4GHz QSOFrequency = "3.4G"
	Frequency5_7GHz QSOFrequency = "5.7G"
	Frequency10GHz  QSOFrequency = "10G"
	Frequency24GHz  QSOFrequency = "24G"
	Frequency47GHz  QSOFrequency = "47G"
	Frequency75GHz  QSOFrequency = "75G"
	Frequency122GHz QSOFrequency = "122G"
	Frequency134GHz QSOFrequency = "134G"
	Frequency241GHz QSOFrequency = "241G"
	FrequencyLight  QSOFrequency = "LIGHT"
)

type QSOMode string

const (
	QSOModeCW    QSOMode = "CW"
	QSOModePhone QSOMode = "PH"
	QSOModeFM    QSOMode = "FM"
	QSOModeRTTY  QSOMode = "RY"
	QSOModeDigi  QSOMode = "DG"
)

type QSOInfo struct {
	Call     callsign.Callsign
	Exchange []string
}
