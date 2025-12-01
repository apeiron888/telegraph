package acl

// Role represents RBAC roles in the system
type Role string

const (
	RoleMember    Role = "member"
	RoleModerator Role = "moderator"
	RoleAdmin     Role = "admin"
)

// roleHierarchy defines role precedence (higher index = more powerful)
var roleHierarchy = []Role{RoleMember, RoleModerator, RoleAdmin}

// HasRole checks if the user has the specified role
func HasRole(userRole string, requiredRole Role) bool {
	return Role(userRole) == requiredRole
}

// HasRoleOrHigher checks if user has the required role or higher
func HasRoleOrHigher(userRole string, requiredRole Role) bool {
	userRoleLevel := getRoleLevel(Role(userRole))
	requiredRoleLevel := getRoleLevel(requiredRole)
	return userRoleLevel >= requiredRoleLevel
}

// getRoleLevel returns the hierarchy level of a role
func getRoleLevel(role Role) int {
	for i, r := range roleHierarchy {
		if r == role {
			return i
		}
	}
	return -1
}

// Permission represents specific permissions in the system
type Permission string

const (
	PermissionCreateChannel   Permission = "channel:create"
	PermissionDeleteChannel   Permission = "channel:delete"
	PermissionManageMembers   Permission = "channel:manage_members"
	PermissionBroadcast       Permission = "channel:broadcast"
	PermissionDeleteMessage   Permission = "message:delete"
	PermissionDeleteAnyMessage Permission = "message:delete_any"
	PermissionViewAuditLogs   Permission = "audit:view"
)

// rolePermissions defines which roles have which permissions
var rolePermissions = map[Role][]Permission{
	RoleMember: {
		PermissionCreateChannel,
		PermissionDeleteMessage,
	},
	RoleModerator: {
		PermissionCreateChannel,
		PermissionDeleteMessage,
		PermissionDeleteAnyMessage,
		PermissionManageMembers,
		PermissionDeleteChannel,
	},
	RoleAdmin: {
		PermissionCreateChannel,
		PermissionDeleteMessage,
		PermissionDeleteAnyMessage,
		PermissionManageMembers,
		PermissionDeleteChannel,
		PermissionBroadcast,
		PermissionViewAuditLogs,
	},
}

// HasPermission checks if a role has a specific permission
func HasPermission(userRole string, permission Permission) bool {
	role := Role(userRole)
	perms, exists := rolePermissions[role]
	if !exists {
		return false
	}

	for _, p := range perms {
		if p == permission {
			return true
		}
	}
	return false
}
