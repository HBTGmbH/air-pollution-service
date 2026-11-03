package csv_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/HBTGmbH/air-pollution-service/internal/csv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const header = "Entity,Code,Year,Nitrogen oxide (NOx),Sulphur dioxide (SO₂) emissions," +
	"Carbon monoxide (CO) emissions,Organic carbon (OC) emissions," +
	"Non-methane volatile organic compounds (NMVOC) emissions," +
	"Black carbon (BC) emissions,Ammonia (NH₃) emissions\n"

func writeFile(t *testing.T, content string) *csv.File {
	t.Helper()

	fileName := filepath.Join(t.TempDir(), "air-pollution.csv")
	require.NoError(t, os.WriteFile(fileName, []byte(content), 0o600))

	return csv.New(fileName)
}

func TestReadRows(t *testing.T) {
	t.Parallel()

	file := writeFile(t, header+
		"Germany,DEU,2000,1.5,2.5,3.5,4.5,5.5,6.5,7.5\n"+
		"Africa,,2001,1,2,3,4,5,6,7\n")

	rows, err := file.ReadRows()

	require.NoError(t, err)
	assert.Equal(t, []*csv.Row{
		{
			Entity: "Germany", Code: "DEU", Year: 2000,
			NOxEmissions: 1.5, SO2Emissions: 2.5, COEmissions: 3.5, OCEmissions: 4.5,
			NMVOCEmissions: 5.5, BCEmissions: 6.5, NH3Emissions: 7.5,
		},
		{
			Entity: "Africa", Year: 2001,
			NOxEmissions: 1, SO2Emissions: 2, COEmissions: 3, OCEmissions: 4,
			NMVOCEmissions: 5, BCEmissions: 6, NH3Emissions: 7,
		},
	}, rows)
}

func TestReadRowsInvalidContent(t *testing.T) {
	t.Parallel()

	file := writeFile(t, header+"Germany,DEU,last year,1,2,3,4,5,6,7\n")

	_, err := file.ReadRows()

	require.Error(t, err)
}

func TestReadRowsMissingFile(t *testing.T) {
	t.Parallel()

	_, err := csv.New(filepath.Join(t.TempDir(), "missing.csv")).ReadRows()

	require.ErrorIs(t, err, fs.ErrNotExist)
}
