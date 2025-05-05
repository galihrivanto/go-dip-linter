package example

// Concrete types
type Service struct{}
type Manager struct{}
type Helper struct{}

// Interfaces
type IService interface{}
type IManager interface{}

// Constructor returning a concrete type (should be flagged if it matches NamePatterns)
func NewService() Service {
	return Service{}
}

// Constructor returning an interface (should not be flagged)
func NewServiceInterface() IService {
	return &Service{}
}

// Constructor returning a concrete type (should be flagged if it matches NamePatterns)
func NewManager() Manager {
	return Manager{}
}

// Constructor returning an interface (should not be flagged)
func NewManagerInterface() IManager {
	return &Manager{}
}

// Constructor returning a concrete type but does not match NamePatterns (should not be flagged)
func CreateHelper() Helper {
	return Helper{}
}

// Constructor returning an interface but does not match NamePatterns (should not be flagged)
func CreateHelperInterface() IService {
	return &Service{}
}

func Add(a, b int) int {
	return a + b
}
