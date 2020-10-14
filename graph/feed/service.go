package feed

// Feed service constants
const (
	requestTimeoutSeconds = 30
)

// NewService creates a new MPESA Service
func NewService() *Service {

	return &Service{}
}

// Service organizes MPESA functionality
type Service struct {
}

func (s Service) checkPreconditions() {

}
