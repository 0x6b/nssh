package models

import "fmt"

// A SIM represents a SORACOM IoT SIM
type SIM struct {
	ActiveProfileID string `json:"activeProfileId"`
	SimID           string `json:"simId"`      // IMSI of the subscriber
	SpeedClass      string `json:"speedClass"` // speed class e.g. s1.4xfast

	Profiles map[string]struct {
		PrimaryImsi string `json:"primaryImsi"`
		Subscribers map[string]struct {
			Imsi         string `json:"imsi"`
			Subscription string `json:"subscription"` // subscription e.g. plan01s, plan-D
		} `json:"subscribers"`
	} `json:"profiles"`
	SessionStatus struct {
		Online bool   `json:"online"` // represents subscriber is online or not
		Imsi   string `json:"imsi"`
	} `json:"sessionStatus"`
	Tags struct {
		Name string `json:"name,omitempty"` // name of the subscriber
	} `json:"tags"`
}

func (s SIM) String() string {
	name := s.Tags.Name
	if s.Tags.Name == "" {
		name = "Unknown"
	}

	return fmt.Sprintf("%v (%v / %v / %v)", name, s.SimID, s.getActiveSubscription(), s.SpeedClass)
}

// Title returns SIM ID and name as its title of the SIM, for interactive command
func (s SIM) Title() string {
	name := s.Tags.Name
	if s.Tags.Name == "" {
		name = "Unknown"
	}
	return fmt.Sprintf("%v %v", s.SimID, name)
}

// Description returns subscription and type (speed class) as its description of the SIM, for interactive command
func (s SIM) Description() string {
	return fmt.Sprintf("%s (%s)", s.getActiveSubscription(), s.SpeedClass)
}

// FilterValue uses all fields as source of filter value of the SIM, for interactive command
func (s SIM) FilterValue() string {
	return fmt.Sprintf("%s%s%s%s", s.SimID, s.getActiveSubscription(), s.Tags.Name, s.SpeedClass)
}

func (s SIM) getActiveSubscription() string {
	activeProfile := s.Profiles[s.ActiveProfileID]
	primaryImsi := activeProfile.PrimaryImsi
	return activeProfile.Subscribers[primaryImsi].Subscription
}
