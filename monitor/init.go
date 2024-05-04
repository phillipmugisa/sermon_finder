package monitor

import (
	"errors"
	"net/http"
	"time"
)

const MAX_REQUEST_PER_MINUTE int = 50
const BLACKLIST_PERIOD time.Duration = time.Minute * 5 // suspend for 5 minutes

type RequestData struct {
	Ip                 string
	Last_reqest_time   time.Time
	Request_per_minute int
	Allow_access_at    time.Time
	BlackListed        bool
}

type RequestMonitor struct {
	requests []*RequestData
}

func NewRequestMethod() *RequestMonitor {
	return &RequestMonitor{
		requests: []*RequestData{},
	}
}

func (m *RequestMonitor) Monitor(r *http.Request) error {
	// analysis request for suspicious activiy
	for _, req := range m.requests {
		if req.Ip == r.RemoteAddr {
			if req.BlackListed {
				// check if blacklist time got done
				if req.Allow_access_at.Before(time.Now()) {
					req.BlackListed = false
					req.Request_per_minute = 1
					return nil
				}
				m.BlackList(req)
			}

			if req.Request_per_minute >= MAX_REQUEST_PER_MINUTE {
				return m.BlackList(req)
			}

			// valid request
			req.Last_reqest_time = time.Now()
			req.Request_per_minute = req.Request_per_minute + 1
			return nil
		}
	}

	// new request ip
	req := &RequestData{
		Ip:                 r.RemoteAddr,
		Last_reqest_time:   time.Now(),
		Request_per_minute: 1,
	}
	m.requests = append(m.requests, req)

	return nil
}

func (m *RequestMonitor) BlackList(r *RequestData) error {
	r.BlackListed = true
	r.Allow_access_at = time.Now().Add(BLACKLIST_PERIOD)

	return errors.New("blacklisted request")
}
