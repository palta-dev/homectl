package discovery

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/palta-dev/homectl/apps/server/internal/config"
)

// WebInfo contains scraped website information
type WebInfo struct {
	Title       string
	FaviconURL  string
	Description string
}

// HostConfig represents a host to scan
type HostConfig struct {
	Name    string   `yaml:"name,omitempty"`
	Address string   `yaml:"address"`
	Ports   []int    `yaml:"ports,omitempty"`
	Tags    []string `yaml:"tags,omitempty"`
}

// Discoverer handles service discovery
type Discoverer struct {
	hosts   []HostConfig
	timeout time.Duration
}

// NewDiscoverer creates a new service discoverer
func NewDiscoverer(hosts []HostConfig, timeout time.Duration) *Discoverer {
	if timeout == 0 {
		timeout = 2 * time.Second
	}
	return &Discoverer{
		hosts:   hosts,
		timeout: timeout,
	}
}

// DiscoverServices finds services from configured hosts
func (d *Discoverer) DiscoverServices(ctx context.Context) ([]config.Service, error) {
	if len(d.hosts) == 0 {
		d.hosts = []HostConfig{
			{Name: "Localhost", Address: "localhost", Tags: []string{"local"}},
		}
	}

	var allServices []config.Service
	var mu sync.Mutex

	for _, host := range d.hosts {
		services := d.scanHost(ctx, host)
		mu.Lock()
		allServices = append(allServices, services...)
		mu.Unlock()
	}

	return allServices, nil
}

// scanHost scans a single host for open ports
func (d *Discoverer) scanHost(ctx context.Context, host HostConfig) []config.Service {
	var services []config.Service
	var mu sync.Mutex
	var wg sync.WaitGroup

	var portsToScan []int
	if len(host.Ports) > 0 {
		portsToScan = host.Ports
	} else {
		portsToScan = make([]int, 65535)
		for i := 1; i <= 65535; i++ {
			portsToScan[i-1] = i
		}
	}

	var scanned int32
	totalPorts := len(portsToScan)
	numWorkers := 500
	portsPerWorker := (totalPorts + numWorkers - 1) / numWorkers

	for w := 0; w < numWorkers; w++ {
		startIdx := w * portsPerWorker
		endIdx := startIdx + portsPerWorker
		if endIdx > totalPorts {
			endIdx = totalPorts
		}
		if startIdx >= totalPorts {
			break
		}

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i < end; i++ {
				select {
				case <-ctx.Done():
					return
				default:
				}

				port := portsToScan[i]
				if d.checkPort(ctx, host.Address, port) {
					serviceName := d.identifyService(port)
					
					// Try to get web info for HTTP/HTTPS ports
					var webInfo *WebInfo
					if isWebPort(port) {
						webInfo = d.fetchWebInfo(ctx, host.Address, port)
						if webInfo != nil && webInfo.Title != "" {
							serviceName = webInfo.Title
						}
					}
					
					service := config.Service{
						Name:        serviceName,
						URL:         d.buildURL(host.Address, port),
						Description: fmt.Sprintf("Auto-discovered on %s:%d", host.Address, port),
						Icon:        d.getIconForService(serviceName, webInfo),
						Tags:        append([]string{host.Address, "discovered"}, host.Tags...),
						NewTab:      true,
					}
					
					mu.Lock()
					services = append(services, service)
					mu.Unlock()
				}
				atomic.AddInt32(&scanned, 1)
			}
		}(startIdx, endIdx)
	}

	wg.Wait()
	return services
}

// isWebPort checks if a port is commonly used for web services
func isWebPort(port int) bool {
	webPorts := map[int]bool{
		80: true, 443: true, 3000: true, 3001: true, 4000: true,
		5000: true, 5173: true, 5174: true, 5175: true, 5176: true,
		5177: true, 8000: true, 8001: true, 8008: true, 8080: true,
		8081: true, 8082: true, 8083: true, 8088: true, 8123: true,
		8181: true, 8200: true, 8443: true, 8888: true, 8999: true,
		9000: true, 9001: true, 9090: true, 10000: true,
	}
	return webPorts[port]
}

// fetchWebInfo scrapes website title and favicon
func (d *Discoverer) fetchWebInfo(ctx context.Context, host string, port int) *WebInfo {
	protocol := "http"
	if port == 443 || port == 8443 || port == 8920 {
		protocol = "https"
	}
	
	webURL := fmt.Sprintf("%s://%s:%d", protocol, host, port)
	
	req, err := http.NewRequestWithContext(ctx, "GET", webURL, nil)
	if err != nil {
		return nil
	}
	
	client := &http.Client{
		Timeout: 5 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 3 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}
	
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024)) // 1MB max
	if err != nil {
		return nil
	}
	
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil
	}
	
	webInfo := &WebInfo{}
	
	// Get title
	webInfo.Title = strings.TrimSpace(doc.Find("title").First().Text())
	if len(webInfo.Title) > 50 {
		webInfo.Title = webInfo.Title[:50] + "..."
	}
	
	// Get description
	webInfo.Description, _ = doc.Find("meta[name='description']").Attr("content")
	
	// Get favicon URL
	faviconURL := ""
	doc.Find("link[rel='icon'], link[rel='shortcut icon'], link[rel='apple-touch-icon']").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			faviconURL = href
		}
	})
	
	// Fallback to /favicon.ico
	if faviconURL == "" {
		faviconURL = "/favicon.ico"
	}
	
	// Make absolute URL if relative
	if !strings.HasPrefix(faviconURL, "http") {
		if strings.HasPrefix(faviconURL, "//") {
			faviconURL = protocol + ":" + faviconURL
		} else if strings.HasPrefix(faviconURL, "/") {
			faviconURL = fmt.Sprintf("%s://%s:%d%s", protocol, host, port, faviconURL)
		} else {
			faviconURL = fmt.Sprintf("%s://%s:%d/%s", protocol, host, port, faviconURL)
		}
	}
	
	webInfo.FaviconURL = faviconURL
	
	return webInfo
}

// checkPort tests if a TCP port is open
func (d *Discoverer) checkPort(ctx context.Context, host string, port int) bool {
	ctx, cancel := context.WithTimeout(ctx, d.timeout)
	defer cancel()

	dialer := &net.Dialer{
		Timeout:   d.timeout,
		KeepAlive: -1,
	}

	address := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return false
	}
	defer conn.Close()

	return true
}

// identifyService returns a service name based on port
func (d *Discoverer) identifyService(port int) string {
	wellKnownPorts := map[int]string{
		20: "FTP-Data", 21: "FTP", 22: "SSH", 23: "Telnet", 25: "SMTP",
		53: "DNS", 67: "DHCP", 68: "DHCP", 69: "TFTP", 110: "POP3",
		119: "NNTP", 123: "NTP", 135: "RPC", 137: "NetBIOS", 138: "NetBIOS",
		139: "NetBIOS", 143: "IMAP", 161: "SNMP", 162: "SNMP", 389: "LDAP",
		445: "SMB", 465: "SMTPS", 514: "Syslog", 587: "SMTP",
		636: "LDAPS", 993: "IMAPS", 995: "POP3S",
		
		80: "HTTP", 443: "HTTPS", 3000: "Node-App", 3001: "Dev-Server",
		4000: "HTTP-Proxy", 5000: "HTTP-Alt", 5173: "Vite", 5174: "Vite",
		5175: "Vite", 5176: "Vite", 5177: "Vite", 8000: "HTTP-Dev",
		8001: "HTTP-Dev", 8008: "HTTP-Alt", 8080: "HTTP-Proxy",
		8081: "HTTP-Admin", 8082: "HTTP-Service", 8088: "HTTP-Alt",
		8443: "HTTPS-Alt", 9000: "Portainer", 9090: "Prometheus",
		
		1433: "MSSQL", 1521: "Oracle", 3306: "MySQL", 5432: "PostgreSQL",
		6379: "Redis", 9200: "Elasticsearch", 27017: "MongoDB",
		
		5672: "RabbitMQ", 9092: "Kafka", 15672: "RabbitMQ-UI",
		
		8096: "Jellyfin", 8920: "Jellyfin-Secure", 32400: "Plex",
		
		8123: "Home-Assistant", 1883: "MQTT", 8883: "MQTT-Secure",
		
		8006: "Proxmox",
		
		5001: "Synology-Secure", 9001: "MinIO",
		
		3389: "RDP", 5900: "VNC",
		
		8333: "Bitcoin", 11211: "Memcached",
	}

	if name, ok := wellKnownPorts[port]; ok {
		return name
	}
	return fmt.Sprintf("Port-%d", port)
}

// buildURL constructs a URL for the service
func (d *Discoverer) buildURL(host string, port int) string {
	securePorts := map[int]bool{
		443: true, 8443: true, 8920: true, 8883: true, 5001: true,
	}

	protocol := "http"
	if securePorts[port] {
		protocol = "https"
	}

	if port == 80 || port == 443 {
		return fmt.Sprintf("%s://%s", protocol, host)
	}

	return fmt.Sprintf("%s://%s:%d", protocol, host, port)
}

// getIconForService returns an icon URL or name
func (d *Discoverer) getIconForService(serviceName string, webInfo *WebInfo) string {
	// If we have a favicon from the website, use it
	if webInfo != nil && webInfo.FaviconURL != "" {
		return webInfo.FaviconURL
	}
	
	serviceName = strings.ToLower(serviceName)
	
	iconMap := map[string]string{
		"ftp": "file", "ssh": "terminal", "telnet": "terminal",
		"smtp": "mail", "dns": "network", "dhcp": "network",
		"pop3": "mail", "imap": "mail", "snmp": "network",
		"ldap": "directory", "https": "globe", "http": "globe",
		"smb": "network", "rpc": "network", "netbios": "network",
		"node": "server", "vite": "code", "dev": "code",
		"mysql": "database", "postgresql": "database", "mongodb": "database",
		"redis": "database", "elasticsearch": "database", "mssql": "database",
		"oracle": "database", "memcached": "database",
		"rabbitmq": "network", "kafka": "network",
		"plex": "film", "jellyfin": "film",
		"home": "home", "mqtt": "network",
		"proxmox": "server", "synology": "server", "minio": "cloud",
		"portainer": "docker", "prometheus": "graph",
		"rdp": "desktop", "vnc": "desktop",
		"bitcoin": "bitcoin",
	}

	for key, icon := range iconMap {
		if strings.Contains(serviceName, key) {
			return icon
		}
	}

	return "server"
}

// GenerateDefaultHosts creates a list of hosts to scan for a subnet
func GenerateDefaultHosts(subnet string) []HostConfig {
	var hosts []HostConfig
	
	scanRange := []int{1, 2, 3, 4, 5, 10, 20, 50, 100, 150, 200, 254}
	for _, i := range scanRange {
		addr := fmt.Sprintf("%s.%d", subnet, i)
		hosts = append(hosts, HostConfig{
			Address: addr,
			Tags:    []string{"lan"},
		})
	}

	return hosts
}
