package client

import (
	"fmt"
	"github.com/livekit/livekit-server/pkg/config"
	"github.com/livekit/protocol/logger"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"os"
	"strings"
	"syscall"
)

func CreateNATSClient(conf *config.Config) (*nats.Conn, error) {
	if conf.NATS.Enabled {
		if len(conf.NATS.ClusterAddresses) == 0 {
			panic("NATS cluster addresses not configured")
		}
		nc, err := nats.Connect(strings.Join(conf.NATS.ClusterAddresses, ","),
			nats.ClosedHandler(func(_ *nats.Conn) {
				pid := os.Getpid()
				// Send SIGKILL to the current process
				err := syscall.Kill(pid, syscall.SIGTERM)
				if err != nil {
					fmt.Printf("Error sending SIGKILL: %v\n", err)
				} else {
					fmt.Println("SIGKILL sent successfully")
				}
				logger.Errorw("NATS connection closed", nil)
			}),
			nats.ReconnectHandler(func(_ *nats.Conn) {
				logger.Infow("NATS reconnected")
			}),
			nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
				logger.Errorw("NATS disconnected", err)
			}),
			nats.ConnectHandler(func(_ *nats.Conn) {
				logger.Infow("NATS connected")
			}),
			nats.MaxReconnects(10))
		if err != nil {
			return nil, errors.Wrap(err, "failed to connect to NATS")
		}
		return nc, nil
	}
	return nil, nil
}
