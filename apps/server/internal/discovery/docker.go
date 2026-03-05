package discovery

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/palta-dev/homectl/apps/server/internal/config"
)

type DockerDiscoverer struct {
	client      *client.Client
	labelPrefix string
}

func NewDockerDiscoverer(socket string, labelPrefix string) (*DockerDiscoverer, error) {
	var opts []client.Opt
	if socket != "" {
		if !strings.HasPrefix(socket, "unix://") {
			socket = "unix://" + socket
		}
		opts = append(opts, client.WithHost(socket))
	}
	opts = append(opts, client.WithAPIVersionNegotiation())

	cli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		return nil, err
	}

	if labelPrefix == "" {
		labelPrefix = "homectl"
	}

	return &DockerDiscoverer{
		client:      cli,
		labelPrefix: labelPrefix,
	}, nil
}

func (d *DockerDiscoverer) DiscoverServices(ctx context.Context) ([]config.Service, error) {
	containers, err := d.client.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		log.Printf("[DOCKER ERROR] Failed to list containers: %v", err)
		return nil, fmt.Errorf("listing containers: %w", err)
	}

	log.Printf("[DOCKER DEBUG] Found %d total containers on host", len(containers))

	var services []config.Service
	for _, container := range containers {
		name := "unknown"
		if len(container.Names) > 0 {
			name = strings.TrimPrefix(container.Names[0], "/")
		}
		
		log.Printf("[DOCKER DEBUG] Processing container: %s (State: %s, Ports: %d)", name, container.State, len(container.Ports))

		var port int
		if len(container.Ports) > 0 {
			for _, p := range container.Ports {
				if p.PublicPort != 0 {
					port = int(p.PublicPort)
					break
				}
			}
			if port == 0 {
				port = int(container.Ports[0].PrivatePort)
			}
		}

		var service config.Service
		service.Name = name
		
		// If it's running and has a port, use localhost. Otherwise, just a placeholder or no URL.
		if port != 0 {
			service.URL = fmt.Sprintf("http://localhost:%d", port)
		} else {
			service.URL = "#" // Placeholder for stopped/unmapped containers
		}

		service.Tags = []string{"docker", container.State}
		service.NewTab = true
		
		// Labels override with basic sanitization
		if val, ok := container.Labels[d.labelPrefix+".name"]; ok { 
			service.Name = sanitize(val) 
		}
		if val, ok := container.Labels[d.labelPrefix+".url"]; ok { 
			service.URL = sanitize(val) 
		}
		
		services = append(services, service)
		log.Printf("[DOCKER DEBUG] Successfully discovered: %s (State: %s)", service.Name, container.State)
	}

	return services, nil
}

func sanitize(s string) string {
	// Remove basic HTML tags to prevent dashboard pollution/XSS if rendered unsafely elsewhere
	s = strings.ReplaceAll(s, "<", "")
	s = strings.ReplaceAll(s, ">", "")
	return strings.TrimSpace(s)
}

func (d *DockerDiscoverer) Close() {
	if d.client != nil {
		d.client.Close()
	}
}
