package repositories

const secret = ""

type UserAuthRepository interface {
	AuthenticateUser(username string, password string) (bool, error)
}

type UserBearerAuthRepository interface {
	AuthorizeBearerUser()
}
