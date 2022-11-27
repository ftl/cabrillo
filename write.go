package cabrillo

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

func Write(w io.Writer, l *Log, appendTX bool) error {
	tags := make([]Tag, 0, len(rowGenerators)+len(l.Custom))
	tags = append(tags,
		CreatedByTag, ContestTag, CallsignTag, OperatorsTag, GridLocatorTag, LocationTag,
		ClaimedScoreTag, OfftimeTag, CategoryAssistedTag, CategoryBandTag, CategoryModeTag,
		CategoryOperatorTag, CategoryPowerTag, CategoryStationTag, CategoryTimeTag,
		CategoryTransmitterTag, CategoryOverlayTag, CertificateTag, ClubTag, NameTag, EmailTag,
		AddressTag, AddressCityTag, AddressStateProvinceTag, AddressPostalcodeTag, AddressCountryTag,
		SoapboxTag,
	)
	for tag := range l.Custom {
		tags = append(tags, tag)
	}

	return WriteWithTags(w, l, appendTX, true, tags...)
}

func WriteWithTags(w io.Writer, l *Log, appendTX bool, ommitIfEmpty bool, tags ...Tag) error {
	err := writeRows(w, row{StartOfLogTag, l.CabrilloVersion, false})
	if err != nil {
		return err
	}

	for _, tag := range tags {
		generator, ok := rowGenerators[tag]
		var rows []row
		if ok {
			rows = generator.ToRow(l, ommitIfEmpty)
		} else {
			rows = []row{{tag, l.Custom[tag], ommitIfEmpty}}
		}

		if rows == nil {
			continue
		}

		err = writeRows(w, rows...)
		if err != nil {
			return err
		}
	}

	err = writeQSOs(w, QSOTag, l.QSOData, appendTX)
	if err != nil {
		return err
	}

	err = writeQSOs(w, XQSOTag, l.IgnoredQSOs, appendTX)
	if err != nil {
		return err
	}

	err = writeRows(w, row{EndOfLogTag, "", false})
	return err
}

type rowGenerator interface {
	ToRow(*Log, bool) []row
}

type rowGeneratorFunc func(*Log, bool) []row

func (f rowGeneratorFunc) ToRow(l *Log, ommitIfEmpty bool) []row {
	return f(l, ommitIfEmpty)
}

var rowGenerators = map[Tag]rowGenerator{
	CallsignTag: rowGeneratorFunc(func(l *Log, ommitIfEmpty bool) []row {
		return []row{{CallsignTag, l.Callsign.String(), ommitIfEmpty}}
	}),
	ContestTag: rowGeneratorFunc(func(l *Log, ommitIfEmpty bool) []row {
		return []row{{ContestTag, string(l.Contest), ommitIfEmpty}}
	}),
	CategoryAssistedTag: rowGeneratorFunc(func(l *Log, ommitIfEmpty bool) []row {
		return []row{{CategoryAssistedTag, string(l.Category.Assisted), ommitIfEmpty}}
	}),
	CategoryBandTag: rowGeneratorFunc(func(l *Log, ommitIfEmpty bool) []row {
		return []row{{CategoryBandTag, string(l.Category.Band), ommitIfEmpty}}
	}),
	CategoryModeTag: rowGeneratorFunc(func(l *Log, ommitIfEmpty bool) []row {
		return []row{{CategoryModeTag, string(l.Category.Mode), ommitIfEmpty}}
	}),
	CategoryOperatorTag: rowGeneratorFunc(func(l *Log, ommitIfEmpty bool) []row {
		return []row{{CategoryOperatorTag, string(l.Category.Operator), ommitIfEmpty}}
	}),
	CategoryPowerTag: rowGeneratorFunc(func(l *Log, ommitIfEmpty bool) []row {
		return []row{{CategoryPowerTag, string(l.Category.Power), ommitIfEmpty}}
	}),
	CategoryStationTag: rowGeneratorFunc(func(l *Log, ommitIfEmpty bool) []row {
		return []row{{CategoryStationTag, string(l.Category.Station), ommitIfEmpty}}
	}),
	CategoryTimeTag: rowGeneratorFunc(func(l *Log, ommitIfEmpty bool) []row {
		return []row{{CategoryTimeTag, string(l.Category.Time), ommitIfEmpty}}
	}),
	CategoryTransmitterTag: rowGeneratorFunc(func(l *Log, ommitIfEmpty bool) []row {
		return []row{{CategoryTransmitterTag, string(l.Category.Transmitter), ommitIfEmpty}}
	}),
	CategoryOverlayTag: rowGeneratorFunc(func(l *Log, ommitIfEmpty bool) []row {
		return []row{{CategoryOverlayTag, string(l.Category.Overlay), ommitIfEmpty}}
	}),
	CertificateTag: rowGeneratorFunc(certificateRow),
	ClaimedScoreTag: rowGeneratorFunc(func(l *Log, ommitIfEmpty bool) []row {
		return []row{{ClaimedScoreTag, strconv.Itoa(l.ClaimedScore), ommitIfEmpty}}
	}),
	ClubTag: rowGeneratorFunc(func(l *Log, ommitIfEmpty bool) []row {
		return []row{{ClubTag, l.Club, ommitIfEmpty}}
	}),
	CreatedByTag: rowGeneratorFunc(func(l *Log, ommitIfEmpty bool) []row {
		return []row{{CreatedByTag, l.CreatedBy, ommitIfEmpty}}
	}),
	EmailTag: rowGeneratorFunc(func(l *Log, ommitIfEmpty bool) []row {
		return []row{{EmailTag, l.Email, ommitIfEmpty}}
	}),
	GridLocatorTag: rowGeneratorFunc(func(l *Log, ommitIfEmpty bool) []row {
		return []row{{GridLocatorTag, l.GridLocator.String(), ommitIfEmpty}}
	}),
	LocationTag: rowGeneratorFunc(func(l *Log, ommitIfEmpty bool) []row {
		return []row{{LocationTag, l.Location, ommitIfEmpty}}
	}),
	NameTag: rowGeneratorFunc(func(l *Log, ommitIfEmpty bool) []row {
		return []row{{NameTag, l.Name, ommitIfEmpty}}
	}),
	AddressTag: rowGeneratorFunc(func(l *Log, ommitIfEmpty bool) []row {
		return []row{{AddressTag, l.Address.Text, ommitIfEmpty}}
	}),
	AddressCityTag: rowGeneratorFunc(func(l *Log, ommitIfEmpty bool) []row {
		return []row{{AddressCityTag, l.Address.City, ommitIfEmpty}}
	}),
	AddressStateProvinceTag: rowGeneratorFunc(func(l *Log, ommitIfEmpty bool) []row {
		return []row{{AddressStateProvinceTag, l.Address.StateProvince, ommitIfEmpty}}
	}),
	AddressPostalcodeTag: rowGeneratorFunc(func(l *Log, ommitIfEmpty bool) []row {
		return []row{{AddressPostalcodeTag, l.Address.Postalcode, ommitIfEmpty}}
	}),
	AddressCountryTag: rowGeneratorFunc(func(l *Log, ommitIfEmpty bool) []row {
		return []row{{AddressCountryTag, l.Address.Country, ommitIfEmpty}}
	}),
	OperatorsTag: rowGeneratorFunc(operatorsRow),
	OfftimeTag:   rowGeneratorFunc(offtimeRow),
	SoapboxTag:   rowGeneratorFunc(soapboxRows),
}

type row struct {
	tag          Tag
	value        string
	ommitIfEmpty bool
}

func (r row) Write(w io.Writer) error {
	if r.ommitIfEmpty && r.value == "" {
		return nil
	}
	_, err := fmt.Fprintf(w, "%s: %s\n", r.tag, r.value)
	return err
}

func writeRows(w io.Writer, rows ...row) error {
	for _, row := range rows {
		err := row.Write(w)
		if err != nil {
			return err
		}
	}
	return nil
}

func certificateRow(l *Log, ommitIfEmpty bool) []row {
	value := "YES"
	if !l.Certificate {
		value = "NO"
	}
	return []row{{CertificateTag, value, ommitIfEmpty}}
}

func operatorsRow(l *Log, ommitIfEmpty bool) []row {
	operators := make([]string, 0, len(l.Operators)+1)
	if l.Host.String() != "" {
		operators = append(operators, "@"+l.Host.String())
	}
	for _, op := range l.Operators {
		if op == l.Host {
			continue
		}
		operators = append(operators, op.String())
	}

	value := strings.Join(operators, ", ")
	return wrapRows(OperatorsTag, value, ommitIfEmpty)
}

func offtimeRow(l *Log, ommitIfEmpty bool) []row {
	var value string
	if l.Offtime.Begin.IsZero() || l.Offtime.End.IsZero() {
		value = ""
	} else {
		value = fmt.Sprintf("%s %s",
			formatTimestamp(l.Offtime.Begin),
			formatTimestamp(l.Offtime.End),
		)
	}
	return []row{{OfftimeTag, value, ommitIfEmpty}}
}

func formatTimestamp(timestamp time.Time) string {
	return timestamp.UTC().Format(TimestampLayout)
}

func soapboxRows(l *Log, ommitIfEmpty bool) []row {
	lines := strings.Split(l.Soapbox, "\n")
	result := make([]row, 0, len(lines))
	for _, line := range lines {
		result = append(result, wrapRows(SoapboxTag, line, ommitIfEmpty)...)
	}
	return result
}

func wrapRows(tag Tag, value string, ommitIfEmpty bool) []row {
	const maxLength = 75
	result := make([]row, 0, (len(value)/maxLength)+1)

	for len(value) > maxLength {
		wrapIndex := strings.LastIndexAny(value[:maxLength], " \n\t")
		if wrapIndex == -1 {
			wrapIndex = maxLength - 1
		}
		result = append(result, row{tag, value[:wrapIndex+1], ommitIfEmpty})
		value = value[wrapIndex+1:]
	}
	result = append(result, row{tag, value, ommitIfEmpty})
	return result
}

func writeQSOs(w io.Writer, tag Tag, data []QSO, appendTX bool) error {
	for _, qso := range data {
		err := writeQSO(w, tag, qso, appendTX)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeQSO(w io.Writer, tag Tag, data QSO, appendTX bool) error {
	_, err := fmt.Fprintf(w, "%s: %s %s %s %s %s %s %s",
		tag,
		data.Frequency,
		data.Mode,
		formatTimestamp(data.Timestamp),
		data.Sent.Call.String(),
		strings.Join(data.Sent.Exchange, " "),
		data.Received.Call.String(),
		strings.Join(data.Received.Exchange, " "),
	)
	if err != nil {
		return err
	}
	if appendTX {
		_, err = fmt.Fprintf(w, " %d", data.Transmitter)
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprintln(w)
	return err
}
