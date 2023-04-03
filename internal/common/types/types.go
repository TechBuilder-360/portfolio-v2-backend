package types

type (
	ApiError struct {
		Field string `json:"field"`
		Msg   string `json:"msg"`
	}

	Registration struct {
		Email string `json:"email" binding:"required"`
	}

	Authentication struct {
		Firstname string `json:"firstname" binding:"required"`
		Lastname  string `json:"lastname" binding:"required"`
		Email     string `json:"email" binding:"required"`
		Password  string `json:"password" binding:"required"`
	}

	RegisterResponse struct {
		Email string `json:"email"`
	}

	Token struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		LifeSpan     uint64 `json:"life_span"`
	}

	UserLogin struct {
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		MiddleName string `json:"middle_name"`
		Bio        string `json:"bio"`
		UserName   string `json:"user_name"`
		ProfilePix string `json:"profile_pix"`
		Profession string `json:"profession"`
	}

	LoginResponse struct {
		Auth       *Token    `json:"auth"`
		Profile    UserLogin `json:"profile"`
		HasProfile bool      `json:"has_profile"`
	}
)

type ENVIRONMENT string
type AuthType string
