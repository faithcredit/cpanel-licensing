package utils

import "time"

func dateFormat() string {
	const layout = "2006-01-02 15:04:05"
	return layout
}
func Date(expiry string) string {

	expiryTime := time.Unix(StrToInt64(expiry), 0)
	formattedExpiry := expiryTime.Format(dateFormat())
	return formattedExpiry
}
func SoftwareKey() string {
	return "cpanelv5"
}
func GetURL() string {
	return "https://apiv2.cpanelseller.com/api/"
}
