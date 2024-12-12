package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func RandomCreateTagAPI(t *testing.T, CusData CreateTagRequest) UserRsp {
	var tagData CreateTagRequest
	if CusData.PostID != uuid.Nil && CusData.TagComments != "" {
		tagData = CusData
	} else {
		tagData = CreateTagRequest{
			PostID: uuid.New(),
			// TagComments: gofakeit.,
		}
	}

	var createdUser UserRsp

	t.Run("RandomTag", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		server := newTestServer(t)
		require.NotEmpty(t, server)

		data, err := json.Marshal(tagData)
		require.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(data))
		require.NoError(t, err)
		require.NotEmpty(t, request)

		server.router.ServeHTTP(recorder, request)
		require.NotEmpty(t, recorder)

		user, err := io.ReadAll(recorder.Body)
		require.NoError(t, err)
		err = json.Unmarshal(user, &createdUser)
		require.NoError(t, err)
		fmt.Println(" ")
	})

	return createdUser
}
