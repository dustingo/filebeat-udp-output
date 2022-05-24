package udpout

import (
	"context"
	"fmt"
	"net"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/outputs"
	"github.com/elastic/beats/v7/libbeat/outputs/codec"
	"github.com/elastic/beats/v7/libbeat/publisher"
	"github.com/elastic/elastic-agent-libs/logp"
)

func init() {
	fmt.Println("register putput udp")
	outputs.RegisterType("udp", makeUdpOutput)
}

type udpOutput struct {
	log           *logp.Logger
	connection    *net.UDPConn
	remoteAddress *net.UDPAddr
	beat          beat.Info
	observer      outputs.Observer
	codec         codec.Codec
	bulkMaxSize   int
	bulkSendDelay int
	onlyMessage   bool
}

func makeUdpOutput(
	_ outputs.IndexManager,
	beat beat.Info,
	observer outputs.Observer,
	//cfg *c.C,
	cfg *common.Config,
) (outputs.Group, error) {
	config := defaultConfig()
	if err := cfg.Unpack(&config); err != nil {
		fmt.Println("error============>", err)
		return outputs.Fail(err)
	}
	// // disable bulk support in publisher pipeline
	// _ = cfg.SetInt("bulk_max_size", -1, -1)
	uo := &udpOutput{
		log:      logp.NewLogger("terminal"),
		beat:     beat,
		observer: observer,
	}
	if err := uo.init(beat, config); err != nil {
		return outputs.Fail(err)
	}
	_ = cfg.SetInt("bulk_max_size", -1, int64(uo.bulkMaxSize))
	return outputs.Success(-1, 2, uo)
}

func (out *udpOutput) init(beat beat.Info, c udpConfig) error {
	var err error
	out.bulkMaxSize = c.BulkMaxSize
	out.bulkSendDelay = c.BulkSendDelay
	out.onlyMessage = c.OnlyMessage
	out.codec, err = codec.CreateEncoder(beat, c.Codec)
	if err != nil {
		return err
	}
	server, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", c.Host, c.Port))
	if err != nil {
		return err
	}
	out.remoteAddress = server
	out.log.Info("udp initialized :", server)
	return nil

}
func (out *udpOutput) Publish(_ context.Context, batch publisher.Batch) error {
	defer batch.ACK()
	st := out.observer
	events := batch.Events()
	st.NewBatch(len(events))
	dropped := 0
	conn, err := net.DialUDP("udp", nil, out.remoteAddress)
	if err != nil {
		return err
	}
	out.connection = conn
	defer out.Close()
	for i := range events {
		event := &events[i]
		serializedEvent, err := out.codec.Encode(out.beat.Beat, &event.Content)
		if err != nil {
			if event.Guaranteed() {
				out.log.Errorf("Failed to serialize the event: %+v", err)
			} else {
				out.log.Warnf("Failed to serialize the event: %+v", err)
			}
			out.log.Debugf("Failed event: %v", event)

			dropped++
			continue
		}

		if _, err = out.connection.Write(append(serializedEvent, '\n')); err != nil {
			st.WriteError(err)
			if event.Guaranteed() {
				out.log.Errorf("Writing event to file failed with: %+v", err)
			} else {
				out.log.Warnf("Writing event to file failed with: %+v", err)
			}

			dropped++
			continue
		}

		st.WriteBytes(len(serializedEvent) + 1)
	}

	st.Dropped(dropped)
	st.Acked(len(events) - dropped)

	return nil
}

func (out *udpOutput) Close() error {
	return out.connection.Close()
}

func (out *udpOutput) String() string {
	return "udp(" + out.remoteAddress.String() + ")"
}
