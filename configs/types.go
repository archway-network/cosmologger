package configs

type Configuration struct {
	GRPC struct {
		Server       string `json:"server"`
		TLS          bool   `json:"tls"`
		APICallRetry int    `json:"api_call_retry"`
		CallTimeout  int    `json:"call_timeout"`
	} `json:"grpc"`

	TendermintClient struct {
		SubscriberName string `json:"subscriber_name"`
		ConnectRetry   int    `json:"connect_retry"`
	} `json:"tendermint_client"`

	Bech32Prefix struct {
		Account struct {
			Address string `json:"address"`
			PubKey  string `json:"pubkey"`
		} `json:"account"`

		Validator struct {
			Address string `json:"address"`
			PubKey  string `json:"pubkey"`
		} `json:"validator"`

		Consensus struct {
			Address string `json:"address"`
			PubKey  string `json:"pubkey"`
		} `json:"consensus"`
	} `json:"bech32_prefix"`
}
