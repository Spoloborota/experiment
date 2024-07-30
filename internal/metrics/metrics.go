package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RegisterCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "user_register_total",
			Help: "Total number of user registrations",
		},
	)
	LoginCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "user_login_total",
			Help: "Total number of user logins",
		},
	)
)

func Init() {
	prometheus.MustRegister(RegisterCounter)
	prometheus.MustRegister(LoginCounter)
}
