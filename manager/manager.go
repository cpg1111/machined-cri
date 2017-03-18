package manager

import (
	"context"
	"net"
	"os"
	"syscall"

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

func (m *MachinedManager) Serve(addr string) error {
	glog.V(3).Infof("Starting Machined on %s", addr)
	if err := syscall.Unlink(addr); err != nil && os.IsNotExist(err) {
		return err
	}
	if m.streamingServer != nil {
		go func() {
			err = m.streamingServer.Start(true)
			if err != nil {
				glog.Fatalf("Failed to start streaming server: %v", err)
			}
		}()
		listener, err := net.Listen("unix", addr)
		if err != nil {
			glog.Fatalf("Failed to listen %s: %v", addr, err)
			return err
		}
		defer listener.Close()
		return m.server.Serve(listener)
	}
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

func (m *MachinedManager) CreateContainer(ctx context.Context, req *kubelet.CreateContainerRequest) (*kubelet.CreateContainerResponse, error) {
	glog.V(3).Infof("CreateContainer with request %s", req.String())
	containerID, err := m.runtimeService.CreateContainer(req.PodSandboxId, req.Config, req.SandboxConfig)
	if err != nil {
		glog.Errorf("CreateContainer from runtime service failed: %v", err)
		return nil, err
	}
	return &kubelet.CreateContainerResponse{ContainerId: containerID}, nil
}

func (m *MachinedManager) StartContainer(ctx context.Context, req *kubelet.StartContainerRequest) (*kubelet.StartContainerResponse, error) {
	glog.V(3).Infof("StartContainer with request %s", req.String())
	err := m.runtimeService.StartContainer(req.ContainerId)
	if err != nil {
		glog.Errorf("StartContainer from runtime service failed: %v", err)
		return err
	}
	return &kubelet.StartContainerResponse{}, nil
}

func (m *MachinedManager) StopContainer(ctx context.Context, req *kubelet.StopContainerRequest) (*kubelet.StopContainerResponse, error) {
	glog.V(3).Infof("StopContainer with request %s", req.String())
	err := m.runtimeService.StopContainer(req.ContainerId, req.Timeout)
	if err != nil {
		glog.Errorf("StopContainer from runtime service failed: %v", err)
		return nil, err
	}
	return &kubelet.StopContainerResponse{}, nil
}

func (m *MachinedManager) RemoveContainer(ctx context.Context, req *kubelet.RemoveContainerRequest) (*kubelet.RemoveContainerResponse, error) {
	glog.V(3).Infof("RemoveContainer with request %s", req.String())
	err := m.runtimeService.RemoveContainer(req.ContainerId)
	if err != nil {
		glog.Errorf("RemoveContainer from runtime service failed: %v", err)
		return nil, err
	}
	return &kubelet.RemoveContainerResponse{}, nil
}

func (m *MachinedManager) ListContainers(ctx context.Context, req *kubelet.ListContainersRequest) (*kubelet.ListContainersResponse, error) {
	glog.V(3).Infof("ListContainers with request %s", req.String())
	containers, err := m.runtimeService.ListContainers(req.GetFilter())
	if err != nil {
		glog.Errorf("ListContainers from runtime service failed: %v", err)
		return nil, err
	}
	return &kubelet.ListContainersResponse{Containers: container}, nil
}

func (m *MachinedManager) ContainerStatus(ctx context.Context, req *kubelet.ContainerStatusRequest) (*kubelet.ContainerStatusResponse, error) {
	glog.V(3).Infof("ContainerStatus with request %s", req.String())
	containerStatus, err := m.runtimeService.ContainerStatus(req.ContainerId)
	if err != nil {
		glog.Errorf("ContainerStatus from runtime service failed: %v", err)
		return nil, err
	}
	return &kubelet.ContainerStatusResponse{Status: containerStatus}, nil
}
