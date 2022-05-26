package udpout

import "github.com/elastic/beats/v7/libbeat/outputs/codec"

type udpConfig struct {
	Host          string       `config:"host"`
	Port          int          `config:"port"`
	Codec         codec.Config `config:"codec"`
	BulkMaxSize   int          `config:"bulk_max_size"`
	BulkSendDelay int          `config:"bulk_send_delay"`
	//OnlyMessage   bool         `config:"only_message"`
}

func defaultConfig() udpConfig {
	return udpConfig{
		Port:          514,
		BulkMaxSize:   2048,
		BulkSendDelay: 20,
		//OnlyMessage:   false,
	}
}

func (u *udpConfig) Validate() error {
	return nil
}
