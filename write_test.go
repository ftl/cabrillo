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
