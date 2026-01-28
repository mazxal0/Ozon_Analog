package utils

import "strings"

func objectNameFromURL(u, bucket string) (string, bool) {
	marker := "/" + bucket + "/"
	i := strings.Index(u, marker)
	if i == -1 {
		return "", false
	}
	obj := u[i+len(marker):]
	if obj == "" {
		return "", false
	}
	return obj, true
}
