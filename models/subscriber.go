package models

import "fmt"

// A Subscriber represents a SORACOM IoT SIM
type Subscriber struct {
	Imsi          string `json:"imsi"`         // IMSI of the subscriber
	Subscription  string `json:"subscription"` // subscription e.g. plan01s, plan-D
	SpeedClass    string `json:"type"`         // speed class e.g. s1.4xfast
	SessionStatus struct {
		Online bool `json:"online"` // represents subscriber is online or not
	} `json:"sessionStatus"`
	Tags struct {
		Name string `json:"name,omitempty"` // name of the subscriber
	} `json:"tags"`
}

func (s Subscriber) String() string {
	name := s.Tags.Name
	if s.Tags.Name == "" {
		name = "Unknown"
	}
	return fmt.Sprintf("%v (%v / %v / %v)", name, s.Imsi, s.Subscription, s.SpeedClass)
}

// Title returns IMSI and name as its title of the subscriber, for interactive command
func (s Subscriber) Title() string {
	name := s.Tags.Name
	if s.Tags.Name == "" {
		name = "Unknown"
	}
	return fmt.Sprintf("%v %v", s.Imsi, name)
}

// Description returns subscription and type (speed class) as its description of the subscriber, for interactive command
func (s Subscriber) Description() string {
	return fmt.Sprintf("%s (%s)", s.Subscription, s.SpeedClass)
}

// FilterValue uses all fields as source of filter value of the subscriber, for interactive command
func (s Subscriber) FilterValue() string {
	return fmt.Sprintf("%s%s%s%s", s.Imsi, s.Subscription, s.Tags.Name, s.SpeedClass)
}
