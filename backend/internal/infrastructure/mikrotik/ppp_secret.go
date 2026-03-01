package mikrotik

import (
	"context"

	"github.com/irhabi89/mikhmon/internal/domain/dto"
)

// parsePPPSecret maps a RouterOS sentence map to a PPPSecret DTO.
func parsePPPSecret(m map[string]string) *dto.PPPSecret {
	return &dto.PPPSecret{
		ID:                   m[".id"],
		Name:                 m["name"],
		Password:             m["password"],
		Profile:              m["profile"],
		Service:              m["service"],
		Disabled:             parseBool(m["disabled"]),
		CallerID:             m["caller-id"],
		LocalAddress:         m["local-address"],
		RemoteAddress:        m["remote-address"],
		Routes:               m["routes"],
		Comment:              m["comment"],
		LimitBytesIn:         parseInt(m["limit-bytes-in"]),
		LimitBytesOut:        parseInt(m["limit-bytes-out"]),
		LastLoggedOut:        m["last-logged-out"],
		LastCallerID:         m["last-caller-id"],
		LastDisconnectReason: m["last-disconnect-reason"],
	}
}

// GetPPPSecrets retrieves all PPP secrets, optionally filtered by profile.
func (c *Client) GetPPPSecrets(ctx context.Context, profile string) ([]*dto.PPPSecret, error) {
	args := []string{"/ppp/secret/print"}
	if profile != "" {
		args = append(args, "?profile="+profile)
	}

	reply, err := c.RunArgsContext(ctx, args)
	if err != nil {
		return nil, err
	}

	secrets := make([]*dto.PPPSecret, 0, len(reply.Re))
	for _, re := range reply.Re {
		secrets = append(secrets, parsePPPSecret(re.Map))
	}

	return secrets, nil
}

// GetPPPSecretByID retrieves a PPP secret by ID.
func (c *Client) GetPPPSecretByID(ctx context.Context, id string) (*dto.PPPSecret, error) {
	reply, err := c.RunContext(ctx, "/ppp/secret/print", "?.id="+id)
	if err != nil {
		return nil, err
	}

	if len(reply.Re) == 0 {
		return nil, nil
	}

	return parsePPPSecret(reply.Re[0].Map), nil
}

// GetPPPSecretByName retrieves a PPP secret by name.
func (c *Client) GetPPPSecretByName(ctx context.Context, name string) (*dto.PPPSecret, error) {
	reply, err := c.RunContext(ctx, "/ppp/secret/print", "?name="+name)
	if err != nil {
		return nil, err
	}

	if len(reply.Re) == 0 {
		return nil, nil
	}

	return parsePPPSecret(reply.Re[0].Map), nil
}

// AddPPPSecret adds a new PPP secret.
func (c *Client) AddPPPSecret(ctx context.Context, secret *dto.PPPSecret) error {
	args := []string{
		"/ppp/secret/add",
		"=name=" + secret.Name,
	}

	if secret.Password != "" {
		args = append(args, "=password="+secret.Password)
	}
	if secret.Profile != "" {
		args = append(args, "=profile="+secret.Profile)
	}
	if secret.Service != "" {
		args = append(args, "=service="+secret.Service)
	}
	if secret.CallerID != "" {
		args = append(args, "=caller-id="+secret.CallerID)
	}
	if secret.LocalAddress != "" {
		args = append(args, "=local-address="+secret.LocalAddress)
	}
	if secret.RemoteAddress != "" {
		args = append(args, "=remote-address="+secret.RemoteAddress)
	}
	if secret.Routes != "" {
		args = append(args, "=routes="+secret.Routes)
	}
	if secret.Comment != "" {
		args = append(args, "=comment="+secret.Comment)
	}
	if secret.LimitBytesIn > 0 {
		args = append(args, "=limit-bytes-in="+formatInt(secret.LimitBytesIn))
	}
	if secret.LimitBytesOut > 0 {
		args = append(args, "=limit-bytes-out="+formatInt(secret.LimitBytesOut))
	}

	_, err := c.RunArgsContext(ctx, args)
	return err
}

// UpdatePPPSecret updates an existing PPP secret.
func (c *Client) UpdatePPPSecret(ctx context.Context, id string, secret *dto.PPPSecret) error {
	args := []string{
		"/ppp/secret/set",
		"=.id=" + id,
	}

	if secret.Name != "" {
		args = append(args, "=name="+secret.Name)
	}
	if secret.Password != "" {
		args = append(args, "=password="+secret.Password)
	}
	if secret.Profile != "" {
		args = append(args, "=profile="+secret.Profile)
	}
	if secret.Service != "" {
		args = append(args, "=service="+secret.Service)
	}
	if secret.CallerID != "" {
		args = append(args, "=caller-id="+secret.CallerID)
	}
	if secret.LocalAddress != "" {
		args = append(args, "=local-address="+secret.LocalAddress)
	}
	if secret.RemoteAddress != "" {
		args = append(args, "=remote-address="+secret.RemoteAddress)
	}
	if secret.Routes != "" {
		args = append(args, "=routes="+secret.Routes)
	}
	if secret.Comment != "" {
		args = append(args, "=comment="+secret.Comment)
	}
	if secret.LimitBytesIn > 0 {
		args = append(args, "=limit-bytes-in="+formatInt(secret.LimitBytesIn))
	}
	if secret.LimitBytesOut > 0 {
		args = append(args, "=limit-bytes-out="+formatInt(secret.LimitBytesOut))
	}
	if secret.Disabled {
		args = append(args, "=disabled=yes")
	} else {
		args = append(args, "=disabled=no")
	}

	_, err := c.RunArgsContext(ctx, args)
	return err
}

// RemovePPPSecret removes a PPP secret.
func (c *Client) RemovePPPSecret(ctx context.Context, id string) error {
	_, err := c.RunContext(ctx, "/ppp/secret/remove", "=.id="+id)
	return err
}

// DisablePPPSecret disables a PPP secret.
func (c *Client) DisablePPPSecret(ctx context.Context, id string) error {
	_, err := c.RunContext(ctx, "/ppp/secret/disable", "=.id="+id)
	return err
}

// EnablePPPSecret enables a PPP secret.
func (c *Client) EnablePPPSecret(ctx context.Context, id string) error {
	_, err := c.RunContext(ctx, "/ppp/secret/enable", "=.id="+id)
	return err
}
