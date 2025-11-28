package settings

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"ots/model"
	"path"
	"strings"
	"time"
)

var MySettings *Settings

type settingsConf struct {
	Mongo_url                 string                   `json:"mongo_url"`
	Db_name                   string                   `json:"db_name"`
	Use_env                   string                   `json:"use_env"`
	Default_ticket_milestones []*model.TicketMilestone `json:"default_ticket_milestones"`
	Ctx_timeout_min           int                      `json:"ctx_timeout_min"`
}

func ReadConfig() (*settingsConf, error) {
	settings := &settingsConf{}
	filepath := path.Join("settings", "settings.jsonc")

	f, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening %s: %v", filepath, err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %v", filepath, err)
	}

	if err = json.Unmarshal(b, settings); err != nil {

		return nil, fmt.Errorf("error decoding settings: %v", err)
	}
	// log.Printf("Settings: %v", settings)

	return settings, nil
}

func Generate() {
	conf, err := ReadConfig()
	if err != nil {
		log.Println(err)
	}

	ctxBase := context.TODO()
	ctx, cancel := context.WithTimeout(ctxBase, time.Duration(conf.Ctx_timeout_min)*time.Minute)

	MySettings = &Settings{
		mongo_url:                 conf.Mongo_url,
		db_name:                   conf.Db_name,
		use_env:                   conf.Use_env,
		default_ticket_milestones: conf.Default_ticket_milestones,
		ctx_with_timeout:          ctx,
		ctx_cancel:                cancel,
	}
}

type Settings struct {
	mongo_url                 string
	db_name                   string
	use_env                   string
	default_ticket_milestones []*model.TicketMilestone
	ctx_with_timeout          context.Context
	ctx_cancel                context.CancelFunc
}

// Getters
func (s *Settings) Get_UseEnv() string {
	return s.use_env
}

func (s *Settings) Get_MongoURL() string {
	urlWithourPass := s.mongo_url

	pass := os.Getenv("MONGO_PASS")

	normURL := strings.Replace(urlWithourPass, "<DB_PASSWORD>", pass, 1)
	// log.Printf("Mongo normalized URL: %s", normURL)

	return normURL
}

func (s *Settings) Get_DBName() string {
	return s.db_name
}

func (s *Settings) Get_DefaultTicketMilestones() []*model.TicketMilestone {
	return s.default_ticket_milestones
}

func (s *Settings) Get_CtxWithTimeout() context.Context {
	return s.ctx_with_timeout
}

func (s *Settings) Get_CtxCancel() context.CancelFunc {
	return s.ctx_cancel
}
