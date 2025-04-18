package profiledata_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	profiledata "tutr-backend/controller/profile_data"
	u "tutr-backend/utils"

	"github.com/stretchr/testify/assert"
)

func TestGetUserProfileDataMissingToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/get_user_profile", nil)
	w := httptest.NewRecorder()

	profiledata.GetUserProfileData(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	assert.Contains(t, w.Body.String(), "unauthorized")
}

func TestGetUserProfileDataInvalidToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/get_user_profile", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	w := httptest.NewRecorder()
	profiledata.GetUserProfileData(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)

	var resp u.ResponseStr
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "failed", resp.Status)
	assert.Empty(t, resp.MyResponse)
}

// func TestHandleDBResultSuccess(t *testing.T) {
// 	w := httptest.NewRecorder()
// 	mockData := map[string]string{"name": "test"}

// 	handleDBResult(w, nil, "Fetched", mockData)

// 	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
// 	assert.Contains(t, w.Body.String(), "Fetched")
// 	assert.Contains(t, w.Body.String(), "test")
// }

// func handleDBResult(w http.ResponseWriter, err error, successMsg string, data interface{}) {
// 	resp := u.ResponseStr{}
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			resp.Status = "failed"
// 			resp.Message = "No user found associated with this user id."
// 			fmt.Println(err)
// 			u.SendResponseWithStatusNotFound(w, resp)
// 		} else {
// 			resp.Status = "failed"
// 			resp.Message = "unable to fetch user profile at this time."
// 			fmt.Println(err)
// 			u.SendResponseWithServerError(w, resp)
// 		}
// 		return
// 	}

// 	resp.Status = "success"
// 	resp.Message = successMsg
// 	resp.MyResponse = data
// 	u.SendResponseWithOK(w, resp)
// }
