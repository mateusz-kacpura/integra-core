// Package contracts contains wire and read-projection DTOs, not business rules.
package contracts

// UserRead is the shared read-replica projection of a user, consumed by
// services that maintain a local users_read table for cross-service queries.
type UserRead struct {
	Version   int64  `json:"version"`
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	AvatarURL string `json:"avatar_url"`
	Status    string `json:"status"`
}

// OrgRead is the shared read-replica projection of an organization, consumed by
// services that maintain a local orgs_read table for cross-service queries.
type OrgRead struct {
	Version int64  `json:"version"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	Status  string `json:"status"`
}
