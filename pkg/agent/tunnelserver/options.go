package tunnelserver

import (
	"dev.khulnasoft.com/api/v4/pkg/devspace"
	"dev.khulnasoft.com/pkg/devcontainer/config"
	"dev.khulnasoft.com/pkg/netstat"
	provider2 "dev.khulnasoft.com/pkg/provider"
)

type Option func(*tunnelServer) *tunnelServer

func WithWorkspace(workspace *provider2.Workspace) Option {
	return func(s *tunnelServer) *tunnelServer {
		s.workspace = workspace
		return s
	}
}

func WithForwarder(forwarder netstat.Forwarder) Option {
	return func(s *tunnelServer) *tunnelServer {
		s.forwarder = forwarder
		return s
	}
}

func WithAllowGitCredentials(allowGitCredentials bool) Option {
	return func(s *tunnelServer) *tunnelServer {
		s.allowGitCredentials = allowGitCredentials
		return s
	}
}

func WithAllowDockerCredentials(allowDockerCredentials bool) Option {
	return func(s *tunnelServer) *tunnelServer {
		s.allowDockerCredentials = allowDockerCredentials
		return s
	}
}

func WithAllowKubeConfig(allow bool) Option {
	return func(s *tunnelServer) *tunnelServer {
		s.allowKubeConfig = allow
		return s
	}
}

func WithMounts(mounts []*config.Mount) Option {
	return func(s *tunnelServer) *tunnelServer {
		s.mounts = mounts
		return s
	}
}

func WithPlatformOptions(options *devspace.PlatformOptions) Option {
	return func(s *tunnelServer) *tunnelServer {
		s.platformOptions = options
		return s
	}
}
