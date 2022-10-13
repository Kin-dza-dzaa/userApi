package validation

import "net/http"

func ParseError(err error) int {
	if err.Error() == "internal error" {
		return http.StatusInternalServerError
	} 
	return http.StatusBadRequest
}