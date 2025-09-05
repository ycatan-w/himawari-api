package register

import (
	"io"
	"net/http"
	"regexp"

	"github.com/ycatan-w/himawari-api/internal/api/internal/utils"
	"github.com/ycatan-w/himawari-api/internal/db"
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
	if len(requestUser.Username) < 4 || len(requestUser.Username) > 50 {
		errs = append(errs, utils.BuildApiError("username", "INVALID_USERNAME_LENGTH", "Username length must be between 4 and 50 characters."))
	}

	if matchUsr, _ := regexp.MatchString("^[a-zA-Z][a-zA-Z0-9\\._\\-@]*$", requestUser.Username); !matchUsr {
		errs = append(errs, utils.BuildApiError("username", "INVALID_USERNAME_CHARS", "Username must start with a letter and contain only letters, numbers, dots (.), underscores (_), dashes (-), or @"))
	}
	if len(errs) == 0 {
		if exists, err := db.UserExists(requestUser.Username); err != nil {
			utils.RespondJSON(w, http.StatusInternalServerError, utils.CreateResponseError(
				"db_connection_failed",
				[]utils.APIError{utils.BuildApiError("", "DB_CONNECTION_FAILED", err.Error())},
			))
			return false
		} else if exists {
			errs = append(errs, utils.BuildApiError("username", "ALREADY_EXISTS_USERNAME", "Username already exists."))
		}
	}
	if len(requestUser.Password) < 8 {
		errs = append(errs, utils.BuildApiError("password", "INVALID_PASSWORD_LENGTH", "Password must be at least 8 characters."))
	}

	if len(errs) > 0 {
		utils.RespondJSON(w, http.StatusBadRequest, utils.CreateResponseError("validation_failed", errs))
		return false
	}
	return true
}

func Register(w http.ResponseWriter, requestUser *UserFormData) (string, bool) {
	hash, ok := utils.EncryptPassword(w, requestUser.Password)
	if !ok {
		return "", false
	}
	user := db.User{
		Username: requestUser.Username,
		Password: hash,
	}
	createdId, err := db.CreateUser(user)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.CreateResponseError("db_connection_failed", []utils.APIError{utils.BuildApiError("", "DB_CONNECTION_FAILED", err.Error())}))

		return "", false
	}

	return utils.GenerateAndSaveToken(w, createdId)
}
