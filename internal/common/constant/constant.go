package constant

import "github.com/TechBuilder-360/portfolio-v2-backend/internal/common/types"

const (
	RequestIdentifier = "Request_Id"

	Expiration    string = "exp"
	Authorized    string = "authorized"
	AccountID     string = "account_id"
	VerifiedEmail string = "verified_email"
	AccountStatus string = "account_status"

	//Account string = "account"

	UnexpectedError string = "an error occurred"

	Activation types.CacheKey = "ACTIVATION-"
	JWT        types.CacheKey = "JWT-"

	EMAILPASSWORD types.AuthType = "EMAIL_PASSWORD"

	Disabled types.AccountStatus = "DISABLED"
	Active   types.AccountStatus = "ACTIVE"
)
