package service

type IAccountService interface {
}

type accountService struct {
}

// NewAccountService instantiates Account Service
func NewAccountService() IAccountService {
	return &accountService{}
}

//func (s accountService) name()  {
//
//}