package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/code-to-go/safepool.lib/api"
	"github.com/code-to-go/safepool.lib/api/chat"
	"github.com/code-to-go/safepool.lib/api/library"
	"github.com/code-to-go/safepool.lib/core"
	pool "github.com/code-to-go/safepool.lib/pool"
	"github.com/code-to-go/safepool.lib/security"
	"github.com/code-to-go/safepool.lib/transport"
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

func (a *App) CreatePool(config string) (string, error) {
	var c pool.Config
	err := yaml.Unmarshal([]byte(config), &c)
	if err != nil {
		return "", err
	}

	if len(c.Configs) == 0 || c.Name == "" {
		return "", pool.ErrInvalidConfig
	} else {
		core.Info("valid config for pool '%s' with %d exchanges", c.Name, len(c.Configs))
	}
	err = pool.Define(c)
	if core.IsErr(err, "cannot define pool '%s': %v", c.Name) {
		return "", err
	}

	p, err := pool.Create(api.Self, c.Name)
	if core.IsErr(err, "cannot create pool '%s': %v", c.Name) {
		return "", err
	}
	pools[p.Name] = p

	token, err := pool.EncodeToken(pool.Token{
		Config: c,
		Host:   api.Self,
	}, "")
	core.IsErr(err, "cannot encode universal token: %v")

	return token, err
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

func (a *App) GetMessages(poolName string, afterIdS string, beforeIdS string, limit int) ([]chat.Message, error) {
	beforeId, err := strconv.ParseUint(beforeIdS, 10, 64)
	if core.IsErr(err, "invalid beforeId parameter '%s': %v", beforeIdS) {
		return nil, err
	}
	afterId, err := strconv.ParseUint(afterIdS, 10, 64)
	if core.IsErr(err, "invalid afterId parameter '%s': %v", afterIdS) {
		return nil, err
	}

	if p, ok := pools[poolName]; ok {
		p.Sync()
		c := chat.Get(p)
		return c.GetMessages(afterId, beforeId, limit)
	}
	return nil, fmt.Errorf("invalid pool '%s'", poolName)
}

func (a *App) PostMessage(poolName string, content string, contentType string, attachements [][]byte) (string, error) {
	if p, ok := pools[poolName]; ok {
		c := chat.Get(p)
		id, err := c.SendMessage(content, contentType, attachements)
		if core.IsErr(err, "cannot post chat message: %v") {
			return "", err
		}
		return strconv.FormatUint(id, 10), nil
	}
	return "", fmt.Errorf("invalid pool '%s'", poolName)
}

func (a *App) GetToken(poolName string, guestId string) (string, error) {
	if p, ok := pools[poolName]; ok {
		c, err := pool.GetConfig(poolName)
		if core.IsErr(err, "cannot get pool config: %v") {
			return "", err
		}

		t := pool.Token{
			Config: c,
			Host:   api.Self,
		}

		if guestId != "" {
			err = p.SetAccess(guestId, pool.Active)
			if core.IsErr(err, "cannot set access for id '%s' in pool '%s': %v", guestId, p.Name) {
				return "", err
			}
		}
		return pool.EncodeToken(t, guestId)
	}
	return "", fmt.Errorf("invalid pool '%s'", poolName)
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
		{
			Id:          1,
			Name:        "test2.doc",
			Size:        192922,
			Author:      api.Self,
			ContentType: "application/word",
			LocalPath:   "~/Documents/pools/safepool.ch/test.doc",
			HasChanged:  false,
		},
	}
}

func (a *App) GetIdentities(poolName string) ([]security.Identity, error) {
	if p, ok := pools[poolName]; ok {
		return p.Identities()
	}
	return nil, fmt.Errorf("invalid pool '%s'", poolName)
}

func (a *App) UpdateIdentity(identity security.Identity) error {
	return security.SetIdentity(identity)
}

func (a *App) GetSelf() security.Identity {
	return api.Self.Public()
}
