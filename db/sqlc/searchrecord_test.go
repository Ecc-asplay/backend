package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateRandomSearchedRecord(t *testing.T) Searchrecord {
	srData := CreateSearchedRecordParams{
		SearchContent: "test_search_content",
		IsUser:        true,
	}

	sr, err := testQueries.CreateSearchedRecord(context.Background(), srData)
	require.NoError(t, err)
	require.NotEmpty(t, sr)
	require.Equal(t, sr.SearchContent, srData.SearchContent)

	return sr
}

func TestCreateSearchedRecord(t *testing.T) {
	CreateRandomSearchedRecord(t)
}

func TestGetSearchedRecords(t *testing.T) {
	for i := 0; i < 100; i++ {
		CreateRandomSearchedRecord(t)
	}

	records, err := testQueries.GetSearchedRecordList(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, records)
	require.GreaterOrEqual(t, len(records), 201)
}
