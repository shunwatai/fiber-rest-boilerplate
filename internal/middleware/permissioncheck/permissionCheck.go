package permissioncheck

type PermissionChecker struct{}

func (pc *PermissionChecker) CheckAccess(resourceName string) fiber.Handler {
		return c.Next()
}
