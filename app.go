package main

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/code-to-go/safepool/api"
	"github.com/code-to-go/safepool/api/chat"
	"github.com/code-to-go/safepool/api/library"
	"github.com/code-to-go/safepool/core"
	pool "github.com/code-to-go/safepool/pool"
	"github.com/code-to-go/safepool/transport"
	"gopkg.in/yaml.v3"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	api.Start()
	for _, n := range pool.List() {
		p, err := pool.Open(api.Self, n)
		if err == nil {
			pools[n] = p
		}
	}
	return &App{}
}

var pools = map[string]*pool.Pool{}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) GetNick() string {
	return api.Self.Nick
}

func (a *App) GetPoolList() []string {
	var names []string
	for n := range pools {
		names = append(names, n)
	}
	return names
}

func (a *App) GetConfigTemplate() (string, error) {
	c := pool.Config{
		Name:    "pool name",
		Configs: []transport.Config{transport.SampleConfig},
	}

	data, err := yaml.Marshal(c)
	return string(data), err
}

func (a *App) CreatePool(config string) error {
	var c pool.Config
	err := yaml.Unmarshal([]byte(config), &c)
	if err != nil {
		return err
	}

	if len(c.Configs) == 0 || c.Name == "" {
		return pool.ErrInvalidConfig
	} else {
		core.Info("valid config for pool '%s' with %d exchanges", c.Name, len(c.Configs))
	}
	err = pool.Define(c)
	if core.IsErr(err, "cannot define pool '%s': %v", c.Name) {
		return err
	}

	p, err := pool.Create(api.Self, c.Name)
	if core.IsErr(err, "cannot create pool '%s': %v", c.Name) {
		return err
	}
	pools[p.Name] = p
	return nil
}

func (a *App) AddPool(token string) error {
	var c pool.Config

	bs, err := base64.StdEncoding.DecodeString(token)
	if core.IsErr(err, "cannot decode token: %v") {
		return err
	}
	err = yaml.Unmarshal(bs, &c)
	if core.IsErr(err, "cannot unmarshal config from token: %v") {
		return err
	}

	if c.Name == "" || len(c.Configs) == 0 {
		core.IsErr(pool.ErrInvalidToken, "invalid config '%s': %v", string(bs))
		return pool.ErrInvalidToken
	} else {
		core.Info("valid token for pool '%s' with %d exchanges", c.Name, len(c.Configs))
	}

	err = pool.Define(c)
	if core.IsErr(err, "cannot define pool '%s': %v", c.Name) {
		return err
	}
	p, err := pool.Open(api.Self, c.Name)
	if core.IsErr(err, "cannot open pool '%s': %v", c.Name) {
		return err
	}
	pools[p.Name] = p
	return nil
}

func (a *App) GetMessages(poolName string, beforeId uint64, limit int) ([]chat.Message, error) {
	if p, ok := pools[poolName]; ok {
		p.Sync()
		c := chat.Get(p)
		return c.Pull(beforeId, limit)
	}
	return nil, fmt.Errorf("invalid pool '%s'", poolName)
}

func (a *App) PostMessage(poolName string, m chat.Message) error {
	if p, ok := pools[poolName]; ok {
		c := chat.Get(p)
		return c.Post(m)
	}
	return fmt.Errorf("invalid pool '%s'", poolName)
}

func (a *App) ListLibrary(poolName string) []library.Document {
	return []library.Document{
		{
			Id:          1,
			Name:        "test.doc",
			Size:        192922,
			Author:      api.Self,
			ContentType: "application/word",
			LocalPath:   "~/Documents/pools/safepool.ch/test.doc",
			HasChanged:  false,
		},
	}
}
