package helpers

import (
	"fmt"
	"net/http"
	"strings"
	"tutr-backend/constants"
	"tutr-backend/middlewares"
	"tutr-backend/models"
	u "tutr-backend/utils"
)

func VerifyTokenAndGetUserID(w http.ResponseWriter, r *http.Request) (data models.TokenDataModel) {
	resp := u.ResponseStr{
		Status:     "failed",
		Message:    "something went wrong",
		MyResponse: nil,
	}
	tokenData := models.TokenDataModel{
		UserID:    "",
		UserEmail: "",
		UserType:  "",
	}
	tokenString := r.Header.Get(constants.Authrization)
	if tokenString == "" {
		u.SendResponseWithUnauthorizedError(w)
		return
	}
	tokenString = tokenString[len("Bearer "):]
	claims, err := middlewares.VerifyToken(tokenString)
	if err != nil {
		resp.Message = "unable to verify your token"
		if strings.Contains(err.Error(), "expired") {
			resp.Message = "session expired. please login again"
			fmt.Println(err)
			u.SendResponseWithServerError(w, resp)
			return
		}
		fmt.Println(err)
		u.SendResponseWithServerError(w, resp)
		return
	}
	tokenData.UserID = claims["user_id"].(string)
	tokenData.UserEmail = claims["email"].(string)
	tokenData.UserType = claims["user_type"].(string)
	return tokenData
}
