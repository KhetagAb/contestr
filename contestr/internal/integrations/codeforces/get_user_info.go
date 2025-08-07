package codeforces

import (
	"context"
	"fmt"
	"github.com/togatoga/goforces"
)

func (c *Service) GetUserInfo(ctx context.Context, handle string) (*goforces.User, error) {
	users, err := c.client.GetUserInfo(ctx, []string{handle})
	if err != nil {
		return nil, fmt.Errorf("error getting user info: %w", err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return &users[0], nil
}
