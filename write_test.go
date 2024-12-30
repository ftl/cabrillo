package cabrillo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTestdataRoundtrip(t *testing.T) {
	entries, err := os.ReadDir("testdata")
	require.NoError(t, err)

	for _, entry := range entries {
		t.Run(entry.Name(), func(t *testing.T) {
			originalFile, err := os.Open("testdata/" + entry.Name())
			require.NoError(t, err)
			defer originalFile.Close()

			originalLog, err := Read(originalFile)
			require.NoError(t, err)
			for tag, value := range originalLog.Custom {
				if value == "" {
					delete(originalLog.Custom, tag)
				}
			}

			outputFile, err := os.CreateTemp("", "cabrillo-roundtrip-"+entry.Name()+"-*")
			assert.NoError(t, err)
			processedFilename := outputFile.Name()
			err = Write(outputFile, originalLog, true)
			assert.NoError(t, err)
			outputFile.Close()

			processedFile, err := os.Open(processedFilename)
			require.NoError(t, err)
			defer processedFile.Close()

			processedLog, err := Read(processedFile)
			assert.NoError(t, err)

			assert.Equal(t, originalLog, processedLog)
		})
	}
}

func TestWrapRows(t *testing.T) {
	tests := []struct {
		name     string
		tag      Tag
		value    string
		expected []string
	}{
		{
			name:     "empty",
			tag:      SoapboxTag,
			value:    "",
			expected: []string{"SOAPBOX:"},
		},
		{
			name:     "short",
			tag:      SoapboxTag,
			value:    "1234567890",
			expected: []string{"SOAPBOX: 1234567890"},
		},
		{
			name:  "long without whitespace",
			tag:   SoapboxTag,
			value: "12345678901234567890123456789012345678901234567890123456789012345678901234567890",
			expected: []string{
				"SOAPBOX: 123456789012345678901234567890123456789012345678901234567890123456",
				"SOAPBOX: 78901234567890",
			},
		},
		{
			name:  "long with whitespace",
			tag:   SoapboxTag,
			value: "123456789 123456789 123456789 123456789 123456789 123456789 123456789 1234567890",
			expected: []string{
				"SOAPBOX: 123456789 123456789 123456789 123456789 123456789 123456789",
				"SOAPBOX: 123456789 1234567890",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := wrapRows(tt.tag, tt.value, true)
			require.Equal(t, len(tt.expected), len(actual))
			for i, row := range actual {
				line := row.String()
				assert.True(t, len(line) <= 75, "len %d: %s = %d", i, line, len(line))
				assert.Equal(t, tt.expected[i], line)
			}
		})
	}
}
