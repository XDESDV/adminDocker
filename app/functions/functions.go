package functions

import (
	"strings"

	"github.com/gofrs/uuid"
)

// NewUUID ...
func NewUUID() string {
	uuid, _ := uuid.NewV4()
	return uuid.String()
}

// RemoveDuplicate : Deleting duplicates in a table
func RemoveDuplicate(xs *[]string) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *xs {
		x = strings.TrimPrefix(x, "-")
		if !found[strings.ToLower(x)] {
			found[strings.ToLower(x)] = true
			(*xs)[j] = (*xs)[i]
			j++
		}
	}
	*xs = (*xs)[:j]
}
