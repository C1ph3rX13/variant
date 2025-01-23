package persistence

import (
	"github.com/kardianos/service"
)

func (p *Program) Create(name, displayName, description string) (service.Service, error) {
	svcConfig := &service.Config{
		Name:        name,
		DisplayName: displayName,
		Description: description,
	}

	prg := &Program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// Start should not block. Do the actual work async.
func (p *Program) Start(s service.Service) error {

	go p.run()
	return nil
}

// Do work here
func (p *Program) run() {

}
func (p *Program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}
