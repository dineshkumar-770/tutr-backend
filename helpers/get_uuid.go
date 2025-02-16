package helpers

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func GenerateUserID(idType string) string {
	id := uuid.New()
	idStr := id.String()
	idStr = strings.ReplaceAll(idStr, "-", "")

	if len(idStr) > 10 {
		idStr = idStr[:10]
	}

	finalID := idStr + "_" + idType
	fmt.Println(finalID)
	return finalID
}
