package dcgm

import (
	"fmt"

	"github.com/google/uuid"
)

func deviceName(index int) string {
	return fmt.Sprintf("nvidia%d", index)
}

// generateUUID returns a stable GPU UUID string
func generateUUID(hostname, deviceName string) string {
	str := fmt.Sprintf("%s-%s", hostname, deviceName)
	return "GPU" + uuid.NewMD5(uuid.NameSpaceURL, []byte(str)).String()
}
