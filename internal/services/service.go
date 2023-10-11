package services

type StorageInteface interface {
	AuthInterface
}

type Services struct {
	Auth
}

func NewService(storage StorageInteface) Services {
	return Services{
		Auth{DB: storage},
	}
}
