package fiberfx

import "fmt"

func fiberHandlers(appName, method, prefix, path string) string {
	return fmt.Sprintf(`name:"fiber-handler-%s-%s-%s-%s"`, appName, method, prefix, path)
}

func fiberHandlerRoutes(appName string) string {
	return fmt.Sprintf(`group:"fiber-handlers-%s"`, appName)
}
