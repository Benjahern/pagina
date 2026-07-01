// Package dto holds request/response shapes for the HTTP layer.
//
// auth.go covers the unified login flow: a single POST /auth/login
// accepts a User or Admin email; the service decides the role by
// looking up Admins first and falling back to Users. The security
// reasoning behind that order (timing attacks, user enumeration)
// lives in services/auth_service.go once we get there.
package dto

import "time"

// LoginRequest is the body for POST /auth/login. Unified for both User
// and Admin — the same endpoint, same shape, role resolved server-side.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

// RegisterRequest is the body for POST /auth/register. Creates a User
// only — Admins are seeded, never self-registered.
//
// PrivacyVersion is required because the consent is recorded by-version:
// the service stores (user_id, version, accepted_at, ip) so audits can
// answer which document the user actually saw at registration time.
// The frontend UX is implicit consent ("by clicking Continue you
// accept the privacy policy vX.Y") — there is no checkbox, but the
// backend still needs the version string the user was shown.
type RegisterRequest struct {
	Email          string `json:"email" binding:"required,email,max=255"`
	Password       string `json:"password" binding:"required,min=8,max=72"`
	Name           string `json:"name" binding:"required,min=1,max=200"`
	Phone          string `json:"phone" binding:"required,min=8,max=20"`
	Rut            string `json:"rut" binding:"required,min=7,max=12"`
	PrivacyVersion string `json:"privacy_version" binding:"required,min=1,max=50"`
}

// AuthResponse is what POST /auth/login and POST /auth/register return.
//
// Token is a signed JWT; the client stores it and sends it back in the
// Authorization: Bearer header. The role claim ("user" | "admin") is
// embedded inside the token, so middleware can authorize without
// hitting the DB.
//
// ExpiresAt mirrors the JWT's exp claim so the UI can render a
// countdown without decoding the token client-side.
type AuthResponse struct {
	Token     string          `json:"token"`
	ExpiresAt time.Time       `json:"expires_at"`
	Role      string          `json:"role"`
	Profile   ProfileResponse `json:"profile"`
}

// ProfileResponse is what GET /auth/me returns. Self-contained: the
// caller can render the header/nav without further round-trips.
//
// Phone is a pointer because the response is the union of User and
// Admin: Admins don't have a phone column, so the field is nil for
// admin sessions and serializes as JSON null. For user sessions the
// field is always populated (User.Phone is NOT NULL).
type ProfileResponse struct {
	ID        uint64    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Phone     *string   `json:"phone,omitempty"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// PasswordChangeRequest is the body for PUT /auth/password. Requires
// the current password to prevent account takeover via stolen session
// tokens — an attacker holding only the JWT cannot change the password
// without also knowing the plaintext.
type PasswordChangeRequest struct {
	CurrentPassword string `json:"current_password" binding:"required,min=8,max=72"`
	NewPassword     string `json:"new_password" binding:"required,min=8,max=72,nefield=CurrentPassword"`
}