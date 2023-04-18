package language

var supportedLanguages = map[string]map[string]string{
	"en": {
		"coin_flip":               "Coin flip",
		"coin_congrats":           "Congratulations, you won.",
		"coin_sorry":              "Sorry, you lost.",
		"coin_head":               "Head",
		"coin_tail":               "Tail",
		"coin_lhead":              "head",
		"coin_ltail":              "tail",
		"game_started":            "Game engine started",
		"game_gen_lb_msg":         "Time's up! Here's todays points:\n",
		"game_lb_broke_streak":    "%s broke their streak of %d days",
		"main_cant_init_discord":  "Unable to initialize discord session,",
		"main_config_load":        "Loading config",
		"main_starting_bot":       "Starting bot",
		"main_config":             "Config:",
		"main_no_streak_backfill": "No streaks found, backfilling",
		"main_cant_conn_discord":  "Can't connect to discord: ",
		"leet_lb_alltime":         "\n\n**Leaderboard all time:**\n",
		"leet_lb_week":            "\n\n**Leaderboard this week:**\n",
		"leet_lb_nopoints":        "No points yet!",
		"leet_lb_nostreak":        "No active streaks :(",
		"leet_days":               "%s: %d days\n",
		"streak_new":              "New streak for %s",
		"streak_continue":         "Continuing streak for %s",
		"streak_new_because":      "New streak for %s, because %s != %s",
	},
	"no": {
		"coin_flip":               "Kron eller Mynt",
		"coin_congrats":           "Gratulerer, du vant.",
		"coin_sorry":              "Desverre, du tapte.",
		"coin_head":               "Mynt",
		"coin_tail":               "Krone",
		"coin_lhead":              "mynt",
		"coin_ltail":              "krone",
		"game_started":            "Spill motor startet",
		"game_gen_lb_msg":         "Tiden er ute! Her er dagens poeng:\n",
		"game_lb_broke_streak":    "%s brøt sin streak på %d dager",
		"main_cant_init_discord":  "Kan ikke initialisere discord,",
		"main_config_load":        "Laster konfigurasjon",
		"main_starting_bot":       "Starter bot",
		"main_config":             "Konfigurasjon:",
		"main_no_streak_backfill": "Ingen streaks funnet, fyller på bakover",
		"main_cant_conn_discord":  "Kan ikke koble til discord: ",
		"leet_lb_alltime":         "\n\n**Ledertavle gjennom tidene:**\n",
		"leet_lb_week":            "\n\n**Ledertavle denne uken:**\n",
		"leet_lb_nopoints":        "Ingen poeng enda!",
		"leet_lb_nostreak":        "Ingen aktive streaks :(",
		"leet_days":               "%s: %d dager\n",
		"streak_new":              "Ny streak for %s",
		"streak_continue":         "Fortsetter streak for %s",
		"streak_new_because":      "Ny streak for %s, fordi %s != %s",
	},
}

var Config struct {
	lang string
}

func SetLanguage(language string) {
	Config.lang = language
}

func GetTranslation(key string) string {
	if langTranslations, ok := supportedLanguages[Config.lang]; ok {
		if translation, ok := langTranslations[key]; ok {
			return translation
		}
	}
	return key
}
