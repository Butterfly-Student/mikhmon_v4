package mikrotik

import (
	"context"

	"github.com/irhabi89/mikhmon/internal/domain/entity"
)

const (
	expireMonitorName = "Mikhmon-Expire-Monitor"
)

// EnsureExpireMonitor ensures scheduler "Mikhmon-Expire-Monitor" exists and enabled.
// Returns status: "created", "enabled", or "existing".
func (c *Client) EnsureExpireMonitor(ctx context.Context, router *entity.Router, script string) (string, error) {
	client, err := c.getClient(router)
	if err != nil {
		return "", err
	}

	reply, err := client.RunContext(ctx, "/system/scheduler/print", "?name="+expireMonitorName)
	if err != nil {
		return "", err
	}

	if len(reply.Re) == 0 {
		_, err := client.RunContext(ctx,
			"/system/scheduler/add",
			"=name="+expireMonitorName,
			"=start-time=00:00:00",
			"=interval=00:01:00",
			"=on-event="+script,
			"=disabled=no",
			"=comment=Mikhmon Expire Monitor",
		)
		if err != nil {
			return "", err
		}
		return "created", nil
	}

	entry := reply.Re[0].Map
	if parseBool(entry["disabled"]) {
		_, err := client.RunContext(ctx,
			"/system/scheduler/set",
			"=.id="+entry[".id"],
			"=interval=00:01:00",
			"=on-event="+script,
			"=disabled=no",
		)
		if err != nil {
			return "", err
		}
		return "enabled", nil
	}

	return "existing", nil
}
