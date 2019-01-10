package brain

import "github.com/taask/taask-server/auth"

// AttemptRunnerAuth allows a runner to auth
func (m *Manager) AttemptRunnerAuth(attempt *auth.Attempt) (*auth.EncMemberSession, error) {
	return m.runnerAuth.AttemptAuth(attempt)
}

// CheckRunnerAuth checks runner auth
func (m *Manager) CheckRunnerAuth(session *auth.Session) error {
	return m.runnerAuth.CheckAuth(session)
}

// AttemptClientAuth allows a runner to auth
func (m *Manager) AttemptClientAuth(attempt *auth.Attempt) (*auth.EncMemberSession, error) {
	return m.clientAuth.AttemptAuth(attempt)
}

// CheckClientAuth checks client auth
func (m *Manager) CheckClientAuth(session *auth.Session) error {
	return m.clientAuth.CheckAuth(session)
}

// CheckClientAdminAuth checks client auth
func (m *Manager) CheckClientAdminAuth(session *auth.Session) error {
	return m.clientAuth.CheckAuthEnsureAdmin(session)
}
