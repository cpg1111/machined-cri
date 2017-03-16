package manager

import (
	"context"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	kubelet "k8s.io/kubernetes/pkg/kubelet/api/v1alpha1/runtime"
	"k8s.io/kubernetes/pkg/kubelet/server/streaming"
	uitlexec "k8s.io/kubernetes/pkg/util/exec"

	"github.com/cpg1111/machined-cri/runtime"
)

type MachinedManager struct {
	server          *grpc.Server
	streamingServer streaming.Server
	runtimeService  runtime.RuntimeService
	imageService    runtime.ImageService
}

func NewMachinedManager(rtSrv runtime.RuntimeService, iSrv runtime.ImageService, streamingServer streaming.Server) (*MachinedManager, error) {
	m := &MachinedManager{
		server:          grpc.NewServer(),
		runtimeService:  rtSrv,
		imageService:    iSrv,
		streamingServer: streamingServer,
	}
	m.registerServer()
	return m, nil
}

func (m *MachinedManager) registerServer() {
	kubelet.RegisterRuntimeServiceServer(m.server, m)
	kubelet.RegisterImageServiceServer(m.server, m)
}

func (m *MachinedManager) Version() (*kubelet.VersionResponse, error) {
	resp, err := s.runtimeService.Version()
	if err != nil {
		glog.Errorf("Get version from runtime service failed: %v", err)
		return nil, err
	}
	return resp, nil
}

func (m *MachinedManager) RunPodSandbox(ctx context.Context, req *kubelet.RunPodSandboxRequest) (*kubelet.RunPodSandboxResponse, error) {
	glog.V(3).Infof("RunPodSandbox with request %s", req.String())
	podID, err := m.runtimeService.RunPodSandbox(req.Config)
	if err != nil {
		glog.Errorf("RunPodSandbox from runtime service failed: %v", err)
		return nil, err
	}
	return &kubelet.RunPodSandboxResponse{PodSandboxId: podID}, nil
}

func (m *MachinedManager) StopPodSandbox(ctx context.Context, req *kubelet.StopPodSandboxRequest) (*kubelet.StopPodSandboxResponse, error) {
	glog.V(3).Infof("StopPodSandbox with request %s", req.String())
	err := m.runtimeService.StopPodSandbox(req.PodSandboxId)
	if err != nil {
		glog.Errorf("StopPodSandbox from runtime service failed: %v", err)
		return nil, err
	}
	return &kubelet.StopPodSandboxResponse{}, nil
}

func (m *MachinedManager) RemovePodSandbox(ctx context.Context, req *kubelet.RemovePodSandboxRequest) (*kubelet.RemovePodSandboxResponse, error) {
	glog.V(3).Infof("RemovePodSandbox with request %s", req.String())
	err := m.runtimeService.RemovePodSandbox(req.PodSandboxId)
	if err != nil {
		glog.Errorf("RemovePodSandbox from runtime service failed: %v", err)
		return nil, err
	}
	return &kubelet.RemovePodSandboxResponse{}, nil
}

func (m *MachinedManager) PodSandboxStatus(ctx context.Context, req *kubelet.PodSandboxStatusRequest) (*kubelet.PodSandboxStatusResponse, error) {
	glog.V(3).Infof("PodSandboxStatus with request %s", req.String())
	podStatus, err := m.runtimeService.PodSandboxStatus(req.PodSandboxId)
	if err != nil {
		glog.Errorf("PodSandboxStatus from runtime service failed: %s", err)
		return nil, err
	}
	return &kubelet.PodSandboxStatusResponse{Status: podStatus}, nil
}

func (m *MachinedManager) ListPodSandbox(ctx context.Context, req *kubelet.ListPodSandboxRequest) (*kubelet.ListPodSandboxRequest, error) {
	glog.V(3).Infof("ListPodSandbox with request %s", req.String())
	pods, err := m.runtimeService.ListPodSandbox(req.GetFilter())
	if err != nil {
		glog.Errorf("ListPodSandbox from runtime service failed: %s", err)
		return nil, err
	}
	return &kubelet.ListPodSandboxResponse{Items: pods}, nil
}
