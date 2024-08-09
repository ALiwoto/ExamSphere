package appValues

// JWTClaimsInfo is a general structure that can be filled by jwt
// claims fields.
type JWTClaimsInfo struct {
	Username string `json:"username"`
	UserId   int64  `json:"user_id"`
	Refresh  bool   `json:"refresh"`
	AuthHash string `json:"auth_hash"`
	Exp      int64  `json:"exp"`
}
