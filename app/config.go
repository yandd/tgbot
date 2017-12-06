package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var (
	ErrCfgIsNil               = fmt.Errorf("<cfg> is nil")
	ErrCfgTelebotTokenInvalid = fmt.Errorf("<cfg> telebot_token invalid")
	ErrCfgSQLiteDBFileInvalid = fmt.Errorf("<cfg> sqlite_dbfile invalid")
)

type Config struct {
	TelebotToken  string `json:"telebot_token"`
	FetchInterval uint64 `json:"fetch_interval"`
	SQLiteDBFile  string `json:"sqlite_dbfile"`
}

func (cfg *Config) validate() error {
	if cfg == nil {
		return ErrCfgIsNil
	}

	if len(cfg.TelebotToken) == 0 {
		return ErrCfgTelebotTokenInvalid
	}

	if len(cfg.SQLiteDBFile) == 0 {
		return ErrCfgSQLiteDBFileInvalid
	}

	if cfg.FetchInterval < 5*60 {
		cfg.FetchInterval = 5 * 60
	}

	return nil
}

func (cfg *Config) FromJsonFile(filename string) error {
	if cfg == nil {
		return ErrCfgIsNil
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, cfg)
	if err != nil {
		return err
	}

	err = cfg.validate()
	if err != nil {
		return err
	}

	return nil
}
