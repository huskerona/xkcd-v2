package image_service

import "encoding/base64"

// Converts byte slice to string
func EncodeToBase64(image []byte) string {
	result := base64.StdEncoding.EncodeToString(image)

	return result
}
