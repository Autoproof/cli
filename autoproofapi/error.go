package autoproofapi

import (
	"fmt"
)

type APIError struct {
	StatusCode int `json:"status_code"`
	Details    struct {
		Result  bool   `json:"result"`
		Message string `json:"message"`
	} `json:"details"`
}

func (e APIError) Error() string {
	return fmt.Sprintf(
		"Autoproof API request finished with non-OK status %d: %s", e.StatusCode, e.Details.Message)
}
