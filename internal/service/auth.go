package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/common/constant"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/common/log"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/common/types"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/common/util"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/config"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/database/redis"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/infrastructure/sendgrid"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/model"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type IAuthService interface {
	RegisterUser(ctx context.Context, body types.Authentication, log log.Entry) (*types.RegisterResponse, error)
	Login(ctx context.Context, body types.LoginRequest, log log.Entry) (*types.LoginResponse, error)
	ActivateEmail(ctx context.Context, body types.AccountActivation, log log.Entry) error
	GenerateJWT(account *model.Account) (*types.Token, error)
	ValidateToken(encodedToken string) (*jwt.Token, error)
	RequestToken(ctx context.Context, account *model.Account, logger log.Entry) error
	//RefreshUserToken(body types.RefreshTokenRequest, token string, logger *log.Entry) (*types.Authentication, error)
}

type authService struct {
	accountRepo repository.IAccountRepository
	userRepo    repository.IUserRepository
	cache       *redis.Client
}

// NewAuthService instantiates Auth Service
func NewAuthService() IAuthService {
	return &authService{
		accountRepo: repository.NewAccountRepository(),
		userRepo:    repository.NewUserRepository(),
		cache:       redis.RedisClient(),
	}
}

func (a *authService) GenerateJWT(account *model.Account) (*types.Token, error) {
	var SecretKey = []byte(config.Instance.SecretKey)

	token := jwt.New(jwt.SigningMethodHS256)

	exp := time.Now().Add(1 * time.Hour).Unix()

	claims := token.Claims.(jwt.MapClaims)
	claims[constant.Expiration] = exp
	claims[constant.Authorized] = true
	claims[constant.AccountID] = account.ID
	claims[constant.VerifiedEmail] = account.EmailVerified
	claims[constant.AccountStatus] = account.Status

	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	rt := util.GenerateUUID()
	err = a.cache.Set(ctx, fmt.Sprintf("%s%s", constant.JWT, account.ID), rt, (time.Hour*24)*7)
	if err != nil {
		log.Error(err.Error())
	}

	return &types.Token{
		AccessToken:  tokenString,
		RefreshToken: rt,
		LifeSpan:     uint64(exp),
	}, nil
}

func (a *authService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, errors.New("invalid token")
		}
		return []byte(config.Instance.SecretKey), nil
	})
}

func (a *authService) RegisterUser(ctx context.Context, body types.Authentication, log log.Entry) (*types.RegisterResponse, error) {
	uw := repository.NewGormUnitOfWork()
	tx, err := uw.Begin()

	defer func() {
		if r := recover(); r != nil {
			log.Error("Panic: %v", r)
			tx.Rollback()
			panic(r)
		}
		if err != nil {
			tx.Rollback()
		}
	}()

	if err != nil {
		return nil, err
	}

	body.Email = util.ToLower(body.Email)
	if util.ValidateEmailAddress(body.Email) == false {
		log.Error("invalid email address %s", body.Email)
		return nil, errors.New("invalid email address")
	}

	account, err := a.accountRepo.GetByEmail(body.Email, ctx)
	if err != nil {
		return nil, err
	}

	if account != nil {
		return nil, errors.New("email address is already registered")
	}

	// validate password
	err = util.ValidatePassword(body.Password)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := util.HashPassword(body.Password)
	if err != nil {
		return nil, errors.New("an error occurred please retry")
	}
	verified := false

	if config.Instance.GetEnv() == config.SandboxEnv {
		verified = true
	}

	account = &model.Account{
		Email:         body.Email,
		Password:      &hashedPassword,
		AuthType:      constant.EMAILPASSWORD,
		EmailVerified: verified,
	}

	err = a.accountRepo.WithTx(tx).Create(account, ctx)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		AccountID: account.ID,
		Email:     account.Email,
		FirstName: util.CapitalizeFirstCharacter(body.FirstName),
		LastName:  util.CapitalizeFirstCharacter(body.LastName),
	}

	err = a.userRepo.WithTx(tx).Create(user, ctx)
	if err != nil {
		return nil, err
	}

	err = uw.Commit(tx)
	if err != nil {
		log.Error(err.Error())
		return nil, errors.New("an error occurred")
	}

	if config.Instance.GetEnv() != config.SandboxEnv {
		var token string
		token = util.GenerateUUID()

		err = a.cache.Set(ctx, fmt.Sprintf("%s%s", constant.Activation, account.ID), token, time.Hour*24)
		if err != nil {
			log.Error("Error occurred when when token %s", err)
		}

		// Send Activate email
		mailTemplate := &sendgrid.ActivationMailRequest{
			Token:    token,
			ToMail:   account.Email,
			ToName:   fmt.Sprintf("%s %s", user.LastName, user.FirstName),
			FullName: fmt.Sprintf("%s %s", user.LastName, user.FirstName),
			UID:      account.ID,
		}

		err = sendgrid.SendActivateMail(mailTemplate)
		if err != nil {
			log.Error("Error occurred when sending activation email. %s", err.Error())
		}
	}

	return &types.RegisterResponse{Email: account.Email}, nil
}

func (a *authService) Login(ctx context.Context, body types.LoginRequest, log log.Entry) (*types.LoginResponse, error) {
	body.Email = util.ToLower(body.Email)
	if util.ValidateEmailAddress(body.Email) == false {
		log.Error("invalid email address %s", body.Email)
		return nil, errors.New("invalid email address")
	}

	account, err := a.accountRepo.GetByEmail(body.Email, ctx)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, errors.New("email address not registered")
	}

	if util.ComparePassword(body.Password, util.AddrToString(account.Password)) == false {
		return nil, errors.New("incorrect password")
	}

	//if account.EmailVerified == false {
	//	return nil, errors.New("account not verified")
	//}

	user, err := a.userRepo.GetByAccountID(account.ID, ctx)
	if err != nil {
		log.Error("error while fetching user %s", err.Error())
		return nil, errors.New("an error occurred")
	}

	go func() {
		account.LastLogin = time.Now()
		err = a.accountRepo.Update(account, ctx)
		if err != nil {
			log.Error("error updating last login %s", err.Error())
		}
	}()

	// Login
	response, err := a.GenerateJWT(account)
	if err != nil {
		return nil, err
	}

	return &types.LoginResponse{
		Auth: response,
		Profile: types.UserLogin{
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			MiddleName: user.MiddleName,
			Bio:        user.Bio,
			ProfilePix: user.ProfilePix,
			Profession: user.Profession,
		},
		HasProfile: user.HasProfile,
	}, nil
}

func (a *authService) ActivateEmail(ctx context.Context, body types.AccountActivation, logger log.Entry) error {
	account, err := a.accountRepo.GetById(ctx, body.UID)
	if err != nil {
		logger.Error(err.Error())
		return errors.New(constant.UnexpectedError)
	}
	if account == nil {
		logger.Error("account not found")
		return errors.New("invalid activation link")
	}

	if account.EmailVerified {
		return errors.New("account has been verified")
	}

	tk, err := a.cache.Get(ctx, fmt.Sprintf("%s%s", constant.Activation, account.ID))
	if err != nil {
		logger.Error(err.Error())
		return errors.New(constant.UnexpectedError)
	}

	if tk != body.Token {
		logger.Error("invalid activation token")
		return errors.New("invalid activation link")
	}

	account.EmailVerified = true
	err = a.accountRepo.Update(account, ctx)
	if err != nil {
		logger.Error(err.Error())
		return errors.New(constant.UnexpectedError)
	}

	go func() {
		err = a.cache.Del(ctx, fmt.Sprintf("%s%s", constant.Activation, account.ID))
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	return nil

}

func (a *authService) RequestToken(ctx context.Context, account *model.Account, logger log.Entry) error {
	if account.EmailVerified == true {
		return errors.New("account is verified")
	}

	user, err := a.userRepo.GetByAccountID(account.ID, ctx)
	if err != nil {
		logger.Error(err.Error())
		return errors.New("user not found")
	}

	if config.Instance.GetEnv() != config.SandboxEnv {
		var token string
		token = util.GenerateUUID()

		err := a.cache.Set(ctx, fmt.Sprintf("%s%s", constant.Activation, account.ID), token, time.Hour*24)
		if err != nil {
			logger.Error("Error occurred when when token %s", err)
		}

		// Send Activate email
		mailTemplate := &sendgrid.ActivationMailRequest{
			Token:    token,
			ToMail:   account.Email,
			ToName:   fmt.Sprintf("%s %s", user.LastName, user.FirstName),
			FullName: fmt.Sprintf("%s %s", user.LastName, user.FirstName),
			UID:      account.ID,
		}

		err = sendgrid.SendActivateMail(mailTemplate)
		if err != nil {
			logger.Error("Error occurred when sending activation email. %s", err.Error())
		}
	}

	return nil
}
