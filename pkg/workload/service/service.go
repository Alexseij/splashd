package service

import "github.com/splashd/gen/workload/service"

var _ service.ServiceServer = &Service{}

type Service struct {
	service.UnimplementedServiceServer
}
