package main

import (
	"code.vikunja.io/api/pkg/db"
	"code.vikunja.io/api/pkg/log"
	"code.vikunja.io/api/pkg/models"
	"code.vikunja.io/api/pkg/plugins"
	"code.vikunja.io/api/pkg/user"
)

type DefaultUserPlugin struct{}

func (p *DefaultUserPlugin) Name() string    { return "Default User" }
func (p *DefaultUserPlugin) Version() string { return "1.0.0" }
func (p *DefaultUserPlugin) Init() error {
	s := db.NewSession()
	defer s.Close()
	u := &user.User{
		Username: "default",
		Email:    "default@example.com",
		Password: "default",
	}
	newUser, err := user.CreateUser(s, u)
	if err != nil {
		_ = s.Rollback()
		if user.IsErrUsernameExists(err) || user.IsErrUserEmailExists(err) {
			return nil
		}
		log.Errorf("Error creating default user: %v", err)
		return nil
	}
	err = models.CreateNewProjectForUser(s, newUser)
	if err != nil {
		log.Errorf("Error when creating new project for user: %v", err)
		_ = s.Rollback()
		return nil
	}

	if err := s.Commit(); err != nil {
		log.Errorf("Error when committing: %v", err)
		return nil
	}
	log.Infof("New default user created with username: %s and password: %s", u.Username, u.Password)
	return nil
}
func (p *DefaultUserPlugin) Shutdown() error { return nil }

func NewPlugin() plugins.Plugin { return &DefaultUserPlugin{} }
