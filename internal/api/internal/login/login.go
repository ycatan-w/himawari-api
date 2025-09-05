package login

import (
	"io"
	"net/http"

	"github.com/ycatan-w/himawari-api/internal/api/internal/utils"
	"github.com/ycatan-w/himawari-api/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type UserFormData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetRequestUser(w http.ResponseWriter, body io.ReadCloser) *UserFormData {
	return utils.GetFormData[UserFormData](w, body)
}

func ValidateUser(w http.ResponseWriter, requestUser *UserFormData) bool {
	var errs []utils.APIError

	if len(requestUser.Username) <= 0 {
		errs = append(errs, utils.BuildApiError("username", "REQUIRED_USERNAME", "Username is required."))
	}
	if len(requestUser.Password) <= 0 {
		errs = append(errs, utils.BuildApiError("password", "REQUIRED_PASSWORD", "Password is required."))
	}
	if len(requestUser.Username) > 0 && len(requestUser.Password) > 0 {
		user := db.FindUserByUsername(requestUser.Username)

		if user == nil {
			errs = append(errs, utils.BuildApiError("", "INVALID_LOGIN", "Username or Password is invalid."))
		} else {
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestUser.Password))
			if err != nil {
				errs = append(errs, utils.BuildApiError("", "INVALID_LOGIN", "Username or Password is invalid."))
			}
		}
	}

	if len(errs) > 0 {
		utils.RespondJSON(w, http.StatusBadRequest, utils.CreateResponseError("validation_failed", errs))
		return false
	}

	return true
}

func Login(w http.ResponseWriter, requestUser *UserFormData) (string, bool) {
	user := db.FindUserByUsername(requestUser.Username)
	if user == nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.CreateResponseError("user_not_found", []utils.APIError{utils.BuildApiError("", "INVALID_LOGIN", "Username or Password is invalid.")}))

		return "", false
	}
	db.PurgeExpiredTokens()

	return utils.GenerateAndSaveToken(w, user.ID)
}
