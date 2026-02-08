package diagnose

import "context"

type DiagnoseService interface {
	DiagnoseFilemanagerService() (string, error)
	DiagnoseAuthService() (string, error)
	DiagnoseLifecycleService() (string, error)
}

type diagnoseServiceRMQ struct {
	// Connections to rabbitmq here
}

func NewDiagnoseServiceRMQ(ctx context.Context) DiagnoseService {
	_ = ctx // use for connection init with timeout

	// Establish connections to rabbitmq here

	return &diagnoseServiceRMQ{
		// Connections to rabbitmq here
	}
}

func (s *diagnoseServiceRMQ) DiagnoseFilemanagerService() (string, error) {
	// TODO: complete RPC logic
	return "", nil
}

func (s *diagnoseServiceRMQ) DiagnoseAuthService() (string, error) {
	// TODO: complete RPC logic
	return "", nil
}

func (s *diagnoseServiceRMQ) DiagnoseLifecycleService() (string, error) {
	// TODO: complete RPC logic
	return "", nil
}
