package repositories

type (
	Repository interface {
		// GetAll() (map[string]string, error)
		GetByCode(string) (string, error)
		Insert(string) (string, error)
	}
)
