package core

import (
	"net"

	"github.com/sirupsen/logrus"
	"github.com/splashd/gen/workload/service"
	"github.com/splashd/internal/util/unix"
	svc "github.com/splashd/pkg/workload/service"
	"google.golang.org/grpc"
)

const DefaultSocketValue string = "8080"

type Daemon struct {
	socket string
	log    *logrus.Entry
}

func NewDaemon(
	socket string,
	log *logrus.Entry,
) *Daemon {
	return &Daemon{
		log:    log,
		socket: socket,
	}
}

func (d *Daemon) Run() error {
	var (
		doneCh = make(chan struct{})
		errCh  = make(chan error)
	)

	d.log.Info("Starting handle system signals")
	unix.HandleSignals(
		d.log,
		doneCh,
	)

	conn, err := net.Listen("unix", d.socket)
	if err != nil {
		return err
	}

	d.log.Info("Create new grpc server")
	s := grpc.NewServer()
	svcServer := &svc.Service{}

	d.log.Info("Register service server")
	service.RegisterServiceServer(s, svcServer)

	d.log.Info("Start listening")
	go func() {
		defer d.log.Info("End listening")
		if err := s.Serve(conn); err != nil {
			errCh <- err
		}
		doneCh <- struct{}{}
	}()

	select {
	case <-doneCh:
		d.log.Info("Perform graceful stop")
		s.GracefulStop()
		d.log.Info("Safe shutdown")
		return nil
	case err := <-errCh:
		d.log.Info("Error while listening")
		return err
	}

}
