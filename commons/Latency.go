package commons

import "github.com/prometheus/client_golang/prometheus"

var (
	Latency = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  "user_data_management",
			Name:       "latency",
			Help:       "latency at different intercepts",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"component", "tag"},
	)
)

const (
	ConnPoolComponent = "connection_pool"
	GetConnLabel      = "get_connection"
	SendMsgLabel      = "send_message"
	GetResponseLabel  = "get_response"
	FreeConnLabel     = "free_conn"
	InitConnLabel     = "init_conn"

	ServerComponent    = "server"
	LoginLabel         = "handle_login"
	ProcessPacketLabel = "process_packet"
	TokenLabel         = "GenerateToken"

	DBComponent       = "database"
	GetUserDBLabel    = "get_user_DB"
	GetUserCacheLabel = "get_user_cache"
	UpdateTokenLabel  = "update_token"

	CacheComponent  = "cache"
	CacheStoreLabel = "cache_store"
	CacheGetLabel   = "cache_get"

	StressTestComponent = "stress_test"
	RoundTripLabel      = "round_trip"
)
