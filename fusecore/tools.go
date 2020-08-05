package fusecore

import "fmt"

func getRemotePath(name string) string {
	return fmt.Sprintf("%s/%s", copyDir, name)
}
