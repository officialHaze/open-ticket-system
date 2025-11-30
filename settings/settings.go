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
	Initial_admins            []*model.Admin           `json:"initial_admins"`
	Password_hash_rounds      int                      `json:"password_hash_rounds"`
	Pipeline_size             int                      `json:"pipeline_size"`
	Server_port               int                      `json:"server_port"`
	Ticket_assign_timeout_min int                      `json:"ticket_assign_timeout_min"`
	Reservoir_size            int                      `json:"reservoir_size"`
	Token_footer              string                   `json:"token_footer"`
	Access_token_exp_min      int                      `json:"access_token_exp_min"`
}

func ReadConfig() (*settingsConf, error) {
	settingsconf := &settingsConf{}
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

	if err = json.Unmarshal(b, settingsconf); err != nil {

		return nil, fmt.Errorf("error decoding settings: %v", err)
	}
	// log.Printf("Settings: %v", settings)

	return settingsconf, nil
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
		initial_admins:            conf.Initial_admins,
		password_hash_rounds:      conf.Password_hash_rounds,
		pipeline_size:             conf.Pipeline_size,
		server_port:               conf.Server_port,
		ticket_assign_timeout_min: time.Duration(conf.Ticket_assign_timeout_min) * time.Minute,
		reservoir_size:            conf.Reservoir_size,
		token_footer:              conf.Token_footer,
		access_token_exp_min:      time.Duration(conf.Access_token_exp_min) * time.Minute,
	}
}

type Settings struct {
	mongo_url                 string
	db_name                   string
	use_env                   string
	default_ticket_milestones []*model.TicketMilestone
	ctx_with_timeout          context.Context
	ctx_cancel                context.CancelFunc
	initial_admins            []*model.Admin
	password_hash_rounds      int
	pipeline_size             int
	server_port               int
	ticket_assign_timeout_min time.Duration
	reservoir_size            int
	token_footer              string
	access_token_exp_min      time.Duration
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

func (s *Settings) Get_InitialAdmins() []*model.Admin {
	return s.initial_admins
}

func (s *Settings) Get_PasswdHashRounds() int {
	return s.password_hash_rounds
}

func (s *Settings) Get_PipelineSize() int {
	return s.pipeline_size
}

func (s *Settings) Get_ServerPort() int {
	return s.server_port
}

func (s *Settings) Get_TicketAssignTimeoutMin() time.Duration {
	return s.ticket_assign_timeout_min
}

func (s *Settings) Get_ReservoirSize() int {
	return s.reservoir_size
}

func (s *Settings) Get_TokenFooter() string {
	return s.token_footer
}

func (s *Settings) Get_AccessTokenExpMin() time.Duration {
	return s.access_token_exp_min
}
