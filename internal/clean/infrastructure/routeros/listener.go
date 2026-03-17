//go:build modern

package routeros

import (
	"context"
	"fmt"

	rosv3 "github.com/go-routeros/routeros/v3"
)

type Listener struct {
	address  string
	username string
	password string
}

func NewListener(host string, port int, username, password string) *Listener {
	return &Listener{
		address:  fmt.Sprintf("%s:%d", host, port),
		username: username,
		password: password,
	}
}

// StreamTraffic membuka stream realtime RouterOS menggunakan ListenArgs.
// Tidak menggunakan pooling agar koneksi stream dedicated per subscription.
func (l *Listener) StreamTraffic(ctx context.Context, args []string, onData func(map[string]string)) error {
	client, err := rosv3.Dial(l.address, l.username, l.password)
	if err != nil {
		return err
	}
	defer client.Close()

	listen, err := client.ListenArgs(args)
	if err != nil {
		return err
	}
	defer listen.Cancel()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case sentence, ok := <-listen.Chan():
			if !ok {
				return nil
			}
			if onData != nil {
				onData(sentence.Map)
			}
		}
	}
}
