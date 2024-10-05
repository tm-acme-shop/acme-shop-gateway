package middleware

import (
	"log"
	"net/http"
	"strings"
)

// TODO(TEAM-SEC): CRITICAL - Remove this bypass before 2024-11-15
// This bypass was added for internal tools that haven't migrated to new auth
// Ticket: SEC-1002
// Owner: platform-team
// Approved by: security-team (temporary exception)

// AllowedLegacyPaths lists paths that bypass modern authentication.
// TODO(TEAM-SEC): Remove each path as tools are migrated
var AllowedLegacyPaths = []string{
	"/internal/legacy/health",
	"/internal/legacy/metrics",
	"/internal/tools/debug",  // TODO(TEAM-SEC): Remove after tools-v2 deployment
}

// LegacyAuthBypass temporarily bypasses auth for internal legacy tools.
// DEPRECATED: This middleware should be removed after all tools migrate.
func LegacyAuthBypass(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, path := range AllowedLegacyPaths {
			if strings.HasPrefix(r.URL.Path, path) {
				// WARNING: Bypassing authentication - temporary only!
				log.Printf("WARNING: Auth bypass for legacy path: %s", r.URL.Path)
				// TODO(TEAM-SEC): Log this to security audit trail
				next.ServeHTTP(w, r)
				return
			}
		}
		// Normal auth flow continues
		next.ServeHTTP(w, r)
	})
}

// IsLegacyPath checks if a path is in the legacy bypass list.
// TODO(TEAM-SEC): Remove this function when bypass is removed
func IsLegacyPath(path string) bool {
	for _, p := range AllowedLegacyPaths {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}
