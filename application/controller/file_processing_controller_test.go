package controller

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestProcessCSV_Success(t *testing.T) {
	// Example CSV data
	csvData := []byte(`Id,Date,Transaction
0,7/15,+60.5
1,7/28,-10.3
2,8/2,-20.46
3,8/13,+10`)

	transactions, err := processCSV(csvData)

	assert.Nil(t, err)
	assert.Equal(t, 4, len(transactions))

	assert.Equal(t, 60.5, transactions[0].Amount)
	assert.Equal(t, time.Date(0, time.July, 15, 0, 0, 0, 0, time.UTC), transactions[0].Date)
	assert.Equal(t, true, transactions[0].IsCredit)

	assert.Equal(t, -10.3, transactions[1].Amount)
	assert.Equal(t, time.Date(0, time.July, 28, 0, 0, 0, 0, time.UTC), transactions[1].Date)
	assert.Equal(t, false, transactions[1].IsCredit)

}

func TestProcessCSV_InvalidCSVData(t *testing.T) {
	invalidCSVData := []byte(`Invalid,Format
0,7/15,+60.5`)

	_, err := processCSV(invalidCSVData)

	assert.Error(t, err)
	assert.Equal(t, "record on line 2: wrong number of fields", err.Error())
}

func TestProcessCSV_InvalidDateFormat(t *testing.T) {
	invalidDateFormatCSV := []byte(`Id,Date,Transaction
0,15/7,+60.5`)

	_, err := processCSV(invalidDateFormatCSV)

	assert.Error(t, err)
	assert.Equal(t, "error line csv: Error reading date field", err.Error())
}

func TestProcessCSV_InvalidAmountFormat(t *testing.T) {
	invalidAmountFormatCSV := []byte(`Id,Date,Transaction
0,7/15,+60.5
1,7/28,-10.3
2,8/2,-20.46
3,8/13,invalid`)

	_, err := processCSV(invalidAmountFormatCSV)

	assert.Error(t, err)
	assert.Equal(t, "error line csv: Error reading amount field", err.Error())
}
