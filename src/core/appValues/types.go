package appValues

// UserRole godoc
// @description UserRole is the role of the user.
type UserRole string // @name UserRole

// JWTClaimsInfo is a general structure that can be filled by jwt
// claims fields.
type JWTClaimsInfo struct {
	UserId   string `json:"user_id"`
	Refresh  bool   `json:"refresh"`
	AuthHash string `json:"auth_hash"`
	Exp      int64  `json:"exp"`
}
