package nssh

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

// GetIP gets current global IP address using https://checkip.amazonaws.com/
func GetIP() (net.IP, error) {
	client := http.DefaultClient
	req, err := http.NewRequest("GET", "https://checkip.amazonaws.com/", nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
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

	ip, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response from https://checkip.amazonaws.com: %w", err)
	}

	return net.ParseIP(strings.TrimSpace(string(ip))), nil
}
