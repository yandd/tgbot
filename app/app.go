package app

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/tucnak/telebot"
)

var (
	Cfg   *Config = &Config{}
	DB    *sqlx.DB
	Bot   *telebot.Bot
	Route *Router
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfgFile := flag.String("config", "tagbot.json", "config <tgbot.json>")

	flag.Parse()

	err := Cfg.FromJsonFile(*cfgFile)
	if err != nil {
		log.Fatalln("Error: ConfigFromJsonFile failed,", err)
		return
	}

	DB, err = OpenSQLite3(Cfg.SQLiteDBFile)
	if err != nil {
		log.Fatalln("Error: OpenSQLite3 failed,", err)
		return
	}

	Bot, err = telebot.NewBot(telebot.Settings{
		Token:  Cfg.TelebotToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatalln("Error: telebot.NewBot failed,", err)
		return
	}

	Bot.Handle("/start", func(m *telebot.Message) {
		Bot.Send(m.Chat, usage())
	})

	Bot.Handle("/help", func(m *telebot.Message) {
		Bot.Send(m.Chat, usage())
	})

	Route = &Router{
		Bot: Bot,
	}
}

func usage() string {
	return fmt.Sprintf(`
/start or /help : 显示此信息

%s`, Route.Usage())
}

func Run() {
	log.Println(usage())
	Bot.Start()
}
