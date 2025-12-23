package GoMock

//go:generate mockgen -source=api.go -destination=./api_mock.go -package=GoMock

type UserProvider interface {
	GetUser(id int) (*User, error)
}
