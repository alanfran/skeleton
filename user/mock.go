package user

// MockStore stores Users in memory for testing purposes.
type MockStore struct {
	mem map[int]User
}
