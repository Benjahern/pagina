package models

import "time"

// PrivacyConsent records a user's acceptance of a privacy policy
// version. Append-only: rows are inserted, never updated or deleted
// by the app. A correction is a new row (or a supersession entry
// in a separate audit log), never a mutation.
//
// Why a separate table instead of a column on User:
//   - History: a user accepts the policy multiple times across
//     versions; storing every row preserves "what was agreed to when".
//   - Defense: in a privacy complaint, the trail is the evidence.
//     A single "consent_accepted" flag on User loses the version
//     history and the exact timestamps.
//   - Audit independence: consents live even if the User is hard-
//     deleted (admin operation can anonymize user_id to NULL).
//
// (UserID, DocumentVersion) is indexed but NOT unique: the same
// pair can legitimately repeat if a policy rolls back or if a
// duplicate acceptance is recorded (each row keeps its own
// AcceptedAt — the timestamp is what counts, not uniqueness).
//
// No UpdatedAt: append-only by design.
type PrivacyConsent struct {
	ConsentID       uint64    `gorm:"primaryKey;autoIncrement"`
	UserID          uint64    `gorm:"not null;index:idx_privacy_consent_user,priority:1;index:idx_privacy_consent_user_version,priority:1"`
	DocumentVersion string    `gorm:"not null;size:50;index:idx_privacy_consent_user_version,priority:2"`
	AcceptedAt      time.Time `gorm:"not null;autoCreateTime"`
	IPAddress       string    `gorm:"not null;size:45"`
}

// TableName returns the exact DBML table name. Plural form to match
// the convention used across the schema (Users, Orders, etc.).
func (PrivacyConsent) TableName() string {
	return "Privacy_Consent"
}