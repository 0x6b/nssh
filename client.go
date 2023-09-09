package nssh

import (
	"encoding/json"
	"fmt"
	"github.com/0x6b/nssh/models"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

// A SoracomClient represents an API client for SORACOM API. See
// https://developers.soracom.io/en/docs/tools/api-reference/ or
// https://dev.soracom.io/jp/docs/api_guide/
type SoracomClient struct {
	APIKey   string // API key
	Token    string // API token
	Client   *http.Client
	Endpoint string
}

type apiParams struct {
	method string
	path   string
	body   string
}

// NewSoracomClient returns new SoracomClient for caller
func NewSoracomClient(coverageType, profileName string) (*SoracomClient, error) {
	akid, ak, ct, err := getAuthInfoFromProfile(profileName)
	if err != nil {
		return nil, err
	}

	if coverageType == "" {
		coverageType = ct
	}

	endpoint, err := getEndpoint(coverageType)
	if err != nil {
		return nil, err
	}

	c := SoracomClient{
		Client:   http.DefaultClient,
		Endpoint: endpoint,
		APIKey:   "",
		Token:    "",
	}

	body, err := json.Marshal(struct {
		AuthKeyID           string `json:"authKeyId"`
		AuthKey             string `json:"authKey"`
		TokenTimeoutSeconds int    `json:"tokenTimeoutSeconds"`
	}{
		AuthKeyID:           akid,
		AuthKey:             ak,
		TokenTimeoutSeconds: 24 * 60 * 60,
	})
	if err != nil {
		return nil, err
	}

	res, err := c.callAPI(&apiParams{
		method: "POST",
		path:   "auth",
		body:   string(body),
	})
	if err != nil {
		return nil, err
	}

	ar := struct {
		APIKey string `json:"apiKey"`
		Token  string `json:"token"`
	}{}
	if err := json.NewDecoder(res.Body).Decode(&ar); err != nil {
		return nil, fmt.Errorf("failed to decode auth response: %w", err)
	}

	c.APIKey = ar.APIKey
	c.Token = ar.Token
	return &c, nil
}

// FindSubscribersByName finds subscribers which has the specified name
func (c *SoracomClient) FindSubscribersByName(name string) ([]models.Subscriber, error) {
	res, err := c.callAPI(&apiParams{
		method: "GET",
		path:   fmt.Sprintf("subscribers?tag_name=name&tag_value=%s", url.QueryEscape(name)),
		body:   "",
	})
	if err != nil {
		return nil, err
	}

	var Subscribers []models.Subscriber
	err = json.NewDecoder(res.Body).Decode(&Subscribers)
	return Subscribers, err
}

// FindSIMsByName finds SIMs which has the specified name
func (c *SoracomClient) FindSIMsByName(name string) ([]models.SIM, error) {
	res, err := c.callAPI(&apiParams{
		method: "GET",
		path:   fmt.Sprintf("query/sims?name=%s", url.QueryEscape(name)),
		body:   "",
	})
	if err != nil {
		return nil, err
	}

	var sims []models.SIM
	err = json.NewDecoder(res.Body).Decode(&sims)
	return sims, err
}

// FindOnlineSubscribers finds online subscribers
func (c *SoracomClient) FindOnlineSubscribers() ([]models.Subscriber, error) {
	res, err := c.callAPI(&apiParams{
		method: "GET",
		path:   "query/subscribers?session_status=ONLINE",
		body:   "",
	})
	if err != nil {
		return nil, err
	}

	var Subscribers []models.Subscriber
	err = json.NewDecoder(res.Body).Decode(&Subscribers)
	return Subscribers, err
}

// FindOnlineSubscribersByName finds online subscribers which has the specified name
func (c *SoracomClient) FindOnlineSubscribersByName(name string) ([]models.Subscriber, error) {
	subscribers, err := c.FindSubscribersByName(name)
	if err != nil {
		return nil, err
	}

	var onlineSubscribers []models.Subscriber
	for _, s := range subscribers {
		if s.SessionStatus.Online {
			onlineSubscribers = append(onlineSubscribers, s)
		}
	}
	return onlineSubscribers, nil
}

// FindOnlineSIMsByName finds online SIMs which has the specified name
func (c *SoracomClient) FindOnlineSIMsByName(name string) ([]models.SIM, error) {
	sims, err := c.FindSIMsByName(name)
	if err != nil {
		return nil, err
	}

	var onlineSIMs []models.SIM
	for _, s := range sims {
		if s.SessionStatus.Online {
			onlineSIMs = append(onlineSIMs, s)
		}
	}
	return onlineSIMs, nil
}

// GetSubscriber gets subscriber information for specified IMSI
func (c *SoracomClient) GetSubscriber(imsi string) (*models.Subscriber, error) {
	res, err := c.callAPI(&apiParams{
		method: "GET",
		path:   fmt.Sprintf("subscribers/%s", imsi),
		body:   "",
	})
	if err != nil {
		return nil, err
	}

	var subscriber models.Subscriber
	err = json.NewDecoder(res.Body).Decode(&subscriber)

	return &subscriber, err
}

// ListPortMappings finds all port mappings
func (c *SoracomClient) ListPortMappings() ([]models.PortMapping, error) {
	res, err := c.callAPI(&apiParams{
		method: "GET",
		path:   "port_mappings",
		body:   "",
	})
	if err != nil {
		return nil, err
	}

	var portMapping []models.PortMapping
	err = json.NewDecoder(res.Body).Decode(&portMapping)
	return portMapping, err
}

// FindPortMappingsForSubscriber finds port mappings for specified subscriber
func (c *SoracomClient) FindPortMappingsForSubscriber(subscriber models.Subscriber) ([]models.PortMapping, error) {
	res, err := c.callAPI(&apiParams{
		method: "GET",
		path:   fmt.Sprintf("port_mappings/subscribers/%s", subscriber.Imsi),
		body:   "",
	})
	if err != nil {
		return nil, err
	}

	var portMapping []models.PortMapping
	err = json.NewDecoder(res.Body).Decode(&portMapping)
	return portMapping, err
}

// FindPortMappingsForSIM finds port mappings for specified SIM
func (c *SoracomClient) FindPortMappingsForSIM(sim models.SIM) ([]models.PortMapping, error) {
	res, err := c.callAPI(&apiParams{
		method: "GET",
		path:   fmt.Sprintf("port_mappings/sims/%s", sim.SimID),
		body:   "",
	})
	if err != nil {
		return nil, err
	}

	var portMapping []models.PortMapping
	err = json.NewDecoder(res.Body).Decode(&portMapping)
	return portMapping, err
}

// FindAvailablePortMappingsForSubscriber finds available port mappings for specified subscriber and port
func (c *SoracomClient) FindAvailablePortMappingsForSubscriber(subscriber models.Subscriber, port int) ([]models.PortMapping, error) {
	portMappings, err := c.FindPortMappingsForSubscriber(subscriber)
	if err != nil {
		return nil, err
	}

	var currentPortMappings []models.PortMapping
	var availablePortMappings []models.PortMapping

	for _, pm := range portMappings {
		if pm.Destination.Port == port {
			currentPortMappings = append(currentPortMappings, pm)
		}
	}

	if len(currentPortMappings) > 0 {
		fmt.Printf("nssh: → found %d port mapping(s) for %s:%d\n", len(currentPortMappings), subscriber.Imsi, port)
		ip, err := GetIP()

		// search port mappings which allows being connected from current IP address
		if err == nil { // ignore https://checkip.amazonaws.com/ error
			fmt.Printf("nssh: → check allowed CIDR for current IP address is %s\n", ip)
			for _, pm := range currentPortMappings {
				for _, r := range pm.Source.IPRanges {
					_, ipNet, err := net.ParseCIDR(r)
					if err == nil {
						if ipNet.Contains(ip) {
							availablePortMappings = append(availablePortMappings, pm)
						}
					}
				}
			}
		}
	}
	return availablePortMappings, nil
}

// FindAvailablePortMappingsForSIM finds available port mappings for specified SIM and port
func (c *SoracomClient) FindAvailablePortMappingsForSIM(sim models.SIM, port int) ([]models.PortMapping, error) {
	portMappings, err := c.FindPortMappingsForSIM(sim)
	if err != nil {
		return nil, err
	}

	var currentPortMappings []models.PortMapping
	var availablePortMappings []models.PortMapping

	for _, pm := range portMappings {
		if pm.Destination.Port == port {
			currentPortMappings = append(currentPortMappings, pm)
		}
	}

	if len(currentPortMappings) > 0 {
		fmt.Printf("nssh: → found %d port mapping(s) for %s:%d\n", len(currentPortMappings), sim.SimID, port)
		ip, err := GetIP()

		// search port mappings which allows being connected from current IP address
		if err == nil { // ignore https://checkip.amazonaws.com/ error
			fmt.Printf("nssh: → check allowed CIDR for current IP address is %s\n", ip)
			for _, pm := range currentPortMappings {
				for _, r := range pm.Source.IPRanges {
					_, ipNet, err := net.ParseCIDR(r)
					if err == nil {
						if ipNet.Contains(ip) {
							availablePortMappings = append(availablePortMappings, pm)
						}
					}
				}
			}
		}
	}
	return availablePortMappings, nil
}

// CreatePortMappingForSubscriber creates port mappings for specified
// subscriber, port, and duration
func (c *SoracomClient) CreatePortMappingForSubscriber(subscriber models.Subscriber, port, duration int) (*models.PortMapping, error) {
	body, err := json.Marshal(struct {
		Duration    int  `json:"duration"`
		TLSRequired bool `json:"tlsRequired"`
		Destination struct {
			Imsi string `json:"imsi"`
			Port int    `json:"port"`
		} `json:"destination"`
	}{
		Duration:    duration * 60,
		TLSRequired: false,
		Destination: struct {
			Imsi string `json:"imsi"`
			Port int    `json:"port"`
		}{
			Imsi: subscriber.Imsi,
			Port: port,
		},
	})
	if err != nil {
		return nil, err
	}

	res, err := c.callAPI(&apiParams{
		method: "POST",
		path:   "port_mappings",
		body:   string(body),
	})
	if err != nil {
		return nil, err
	}

	var portMapping models.PortMapping
	err = json.NewDecoder(res.Body).Decode(&portMapping)
	return &portMapping, err
}

// CreatePortMappingForSIM creates port mappings for specified
// subscriber, port, and duration
func (c *SoracomClient) CreatePortMappingForSIM(sim models.SIM, port, duration int) (*models.PortMapping, error) {
	body, err := json.Marshal(struct {
		Duration    int  `json:"duration"`
		TLSRequired bool `json:"tlsRequired"`
		Destination struct {
			SimID string `json:"simId"`
			Port  int    `json:"port"`
		} `json:"destination"`
	}{
		Duration:    duration * 60,
		TLSRequired: false,
		Destination: struct {
			SimID string `json:"simId"`
			Port  int    `json:"port"`
		}{
			SimID: sim.SimID,
			Port:  port,
		},
	})
	if err != nil {
		return nil, err
	}

	res, err := c.callAPI(&apiParams{
		method: "POST",
		path:   "port_mappings",
		body:   string(body),
	})
	if err != nil {
		return nil, err
	}

	var portMapping models.PortMapping
	err = json.NewDecoder(res.Body).Decode(&portMapping)
	return &portMapping, err
}

// Connect connects to specified port mapping with login name and identity. If
// identity is specified, use it for public key authentication. If not, use
// password authentication instead.
func (c *SoracomClient) Connect(login, identity string, portMapping *models.PortMapping) error {
	sshConfig, err := newSSHClientConfig(login, identity)
	if err != nil {
		return err
	}

	client, err := ssh.Dial("tcp", portMapping.Endpoint, sshConfig)
	if err != nil {
		return err
	}

	session, err := client.NewSession()
	if err != nil {
		return err
	}

	defer func() {
		err := session.Close()
		if err != nil {
			// do nothing
		}
	}()

	fd := int(os.Stdin.Fd())
	state, err := terminal.MakeRaw(fd)
	if err != nil {
		return err
	}

	defer func() {
		err := terminal.Restore(fd, state)
		if err != nil {
			fmt.Println("failed to restore terminal", err)
		}
	}()

	w, h, err := terminal.GetSize(fd)
	if err != nil {
		fmt.Println("failed to get terminal size, using default values", err)
		w = 80
		h = 24
	}

	err = session.RequestPty("xterm", h, w, ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	})
	if err != nil {
		return err
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to setup stdin for session: %v", err)
	}
	go dup(stdin, os.Stdin)

	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to setup stdout for session: %v", err)
	}
	go dup(os.Stdout, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to setup stderr for session: %v", err)
	}
	go dup(os.Stderr, stderr)

	err = session.Shell()
	if err != nil {
		fmt.Println(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, SIGWINCH)
	go func() {
		for {
			s := <-ch
			switch s {
			case SIGWINCH:
				fd := int(os.Stdout.Fd())
				w, h, _ = terminal.GetSize(fd)
				err := session.WindowChange(h, w)
				if err != nil {
					fmt.Println("failed to change window size", err)
				}
			}
		}
	}()

	err = session.Wait()
	return err
}

func readPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	return string(password), err
}

func dup(dst io.Writer, src io.Reader) {
	_, err := io.Copy(dst, src)
	if err != nil {
		fmt.Println("failed to copy stdin", err)
	}
}

func getEndpoint(coverageType string) (string, error) {
	if strings.HasPrefix(coverageType, "j") {
		return "https://api.soracom.io", nil
	} else if strings.HasPrefix(coverageType, "g") {
		return "https://g.api.soracom.io", nil
	} else {
		return "", fmt.Errorf("invalid coverage type: %s", coverageType)
	}
}

func getAuthInfoFromProfile(profileName string) (string, string, string, error) {
	dir, err := getProfileDir()
	if err != nil {
		return "", "", "", err
	}
	path := filepath.Join(dir, profileName+".json")

	b, err := os.ReadFile(path)
	if err != nil {
		return "", "", "", err
	}

	p := struct {
		AuthKeyID    *string `json:"authKeyId"`
		AuthKey      *string `json:"authKey"`
		CoverageType *string `json:"coverageType"`
	}{}
	err = json.Unmarshal(b, &p)
	return *p.AuthKeyID, *p.AuthKey, *p.CoverageType, err
}

func getProfileDir() (string, error) {
	profileDir := os.Getenv("SORACOM_PROFILE_DIR")

	if profileDir == "" {
		dir, err := homedir.Dir()
		if err != nil {
			return "", err
		}
		profileDir = filepath.Join(dir, ".soracom")
	}

	return profileDir, nil
}

func newSSHClientConfig(login string, identity string) (*ssh.ClientConfig, error) {
	var am ssh.AuthMethod

	if identity == "" {
		password, err := readPassword("nssh: password: ")
		if err != nil {
			return nil, err
		}
		am = ssh.Password(password)
		fmt.Println("")
	} else {
		_, err := os.Stat(identity)
		if err != nil {
			return nil, err
		}

		buf, err := os.ReadFile(identity)
		if err != nil {
			return nil, err
		}

		key, err := ssh.ParsePrivateKey(buf)
		if err != nil {
			return nil, err
		}
		am = ssh.PublicKeys(key)
	}

	return &ssh.ClientConfig{
		User:            login,
		Auth:            []ssh.AuthMethod{am},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}, nil
}

func (c *SoracomClient) callAPI(params *apiParams) (*http.Response, error) {
	req, err := c.makeRequest(params)
	if err != nil {
		return nil, err
	}
	res, err := c.doRequest(req)
	return res, err
}

func (c *SoracomClient) makeRequest(params *apiParams) (*http.Request, error) {
	var body io.Reader
	if params.body != "" {
		body = strings.NewReader(params.body)
	}

	req, err := http.NewRequest(params.method,
		fmt.Sprintf("%s/v1/%s", c.Endpoint, params.path),
		body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Soracom-Lang", "en")
	if c.APIKey != "" {
		req.Header.Set("X-Soracom-Api-Key", c.APIKey)
	}
	if c.Token != "" {
		req.Header.Set("X-Soracom-Token", c.Token)
	}
	return req, nil
}

func (c *SoracomClient) doRequest(req *http.Request) (*http.Response, error) {
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		defer func() {
			err := res.Body.Close()
			if err != nil {
				fmt.Println("failed to close response", err)
			}
		}()
		return nil, fmt.Errorf("%s: %s %s", res.Status, req.Method, req.URL)
	}
	return res, nil
}
