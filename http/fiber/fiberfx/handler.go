package fiberfx

import "fmt"

func GetFiberApp(appName string) string {
	return fmt.Sprintf(`name:"fiber-%s"`, appName)
}

func fiberHandlers(appName, method, prefix, path string) string {
	return fmt.Sprintf(`name:"fiber-handler-%s-%s-%s-%s"`, appName, method, prefix, path)
}

func fiberHandlerRoutes(appName string) string {
	return fmt.Sprintf(`group:"fiber-handlers-%s"`, appName)
}

func routerCallbacksName(appName string) string {
	return fmt.Sprintf(`name:"fiber-%s-router-callbacks"`, appName)
}
