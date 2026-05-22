package objectstorage

import "contestr/internal/configs"

// NewClientOptional returns a client when object_storage is configured, otherwise (nil, nil).
func NewClientOptional(cfg *configs.Config) (*Client, error) {
	if cfg == nil || !cfg.ObjectStorage.Enabled() {
		return nil, nil
	}
	return NewClient(cfg)
}
