package brain

import (
	"github.com/taask/taask-server/auth"
	"github.com/taask/taask-server/partner"
)

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

// AttemptPartnerAuth allows a partner to auth
func (m *Manager) AttemptPartnerAuth(attempt *auth.Attempt) (*auth.EncMemberSession, error) {
	// TODO: figure out a less gross way to tell the partnerManager about the partner's UUID
	// we're doing it this way right now because we want AuthManager to be generalized,
	// but this is a special case of needing to know the internals of auth

	m.PartnerManager.SetPartner(&partner.Partner{UUID: attempt.MemberUUID})
	encMemberSession, err := m.PartnerManager.Auth.AttemptAuth(attempt)
	if err != nil {
		m.PartnerManager.SetPartner(nil)
		return nil, err
	}

	return encMemberSession, nil
}

// CheckPartnerAuth checks partner auth
func (m *Manager) CheckPartnerAuth(session *auth.Session) error {
	return m.PartnerManager.Auth.CheckAuth(session)
}
