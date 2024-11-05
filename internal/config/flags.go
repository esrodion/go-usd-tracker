package config

import (
	"flag"
	"os"
)

var (
	readers map[string]func(cfg *Config, flag *flag.Flag)

	// default startup
	def bool

	// postgres settings
	dbconn, dbmigurl string
	amigup, amigdown bool

	// webservice settings
	grpcaddr, httpaddr string

	// logging setup
	ljson, lext                           bool
	llevel, lfile, ltimeformat, llifetime string
)

func init() {
	readers = make(map[string]func(cfg *Config, flag *flag.Flag))

	flag.BoolVar(&def, "default", false, "Start app with default options")

	flag.StringVar(&dbconn, "dbconn", "", "Database connection URL")
	readers["dbconn"] = func(cfg *Config, flag *flag.Flag) {
		cfg.Conn = dbconn
	}

	flag.StringVar(&dbmigurl, "dbmigurl", "", "Database migrations URL")
	readers["dbmigurl"] = func(cfg *Config, flag *flag.Flag) {
		cfg.MigrationsURL = dbmigurl
	}

	flag.BoolVar(&amigup, "amigup", true, "Auto migrate up on startup")
	readers["amigup"] = func(cfg *Config, flag *flag.Flag) {
		cfg.AutoMigrateUp = boolToText(amigup)
	}

	flag.BoolVar(&amigdown, "amigdown", false, "Auto migrate down on shutdown")
	readers["amigdown"] = func(cfg *Config, flag *flag.Flag) {
		cfg.AutoMigrateUp = boolToText(amigdown)
	}

	flag.StringVar(&grpcaddr, "grpcaddr", "", "address:port of grpc server (rates service)")
	readers["grpcaddr"] = func(cfg *Config, flag *flag.Flag) {
		cfg.GrpcAddress = grpcaddr
	}

	flag.StringVar(&httpaddr, "httpaddr", "", "address:port of http server (metrics and health check)")
	readers["httpaddr"] = func(cfg *Config, flag *flag.Flag) {
		cfg.HttpAddress = httpaddr
	}

	flag.BoolVar(&ljson, "ljson", false, "Write logs in json format")
	readers["ljson"] = func(cfg *Config, flag *flag.Flag) {
		os.Setenv("LOG_TO_JSON", boolToText(ljson))
	}

	flag.BoolVar(&lext, "lext", false, "Add extra inforamtion to log entries")
	readers["lext"] = func(cfg *Config, flag *flag.Flag) {
		os.Setenv("LOG_EXTENDED", boolToText(lext))
	}

	flag.StringVar(&llevel, "llevel", "DEBUG", "Log level: DEBUG (default), INFO, WARNING, ERROR")
	readers["llevel"] = func(cfg *Config, flag *flag.Flag) {
		os.Setenv("LOG_LEVEL", llevel)
	}

	flag.StringVar(&lfile, "lfile", "", "filepath to write logs to, if empty, logs go to stdout. May contain single '%s' token for timestamp")
	readers["lfile"] = func(cfg *Config, flag *flag.Flag) {
		os.Setenv("LOG_TO_FILE", "true")
		os.Setenv("LOG_FILE", lfile)
	}

	flag.StringVar(&ltimeformat, "ltimeformat", "", "format of logfile timestamp")
	readers["ltimeformat"] = func(cfg *Config, flag *flag.Flag) {
		os.Setenv("LOG_FILE_TIME_FORMAT", ltimeformat)
	}

	flag.StringVar(&llifetime, "llifetime", "0", "time (seconds) until current log file is closed and new is created")
	readers["llifetime"] = func(cfg *Config, flag *flag.Flag) {
		os.Setenv("LOG_FILE_LIFETIME", llifetime)
	}
}

func ReadFlags(cfg *Config) bool {
	if !flag.Parsed() {
		flag.Parse()
	}

	anySet := false
	flag.Visit(func(flag *flag.Flag) {
		anySet = true
		f, ok := readers[flag.Name]
		if ok {
			f(cfg, flag)
		}
	})

	return anySet
}

func PrintCLIUsage() {
	flag.Usage()
}

//// Service

func boolToText(v bool) string {
	if v {
		return "true"
	}
	return "false"
}
