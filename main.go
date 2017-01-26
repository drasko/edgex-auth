package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/drasko/go-auth/api"
	"github.com/drasko/go-auth/config"
	"github.com/drasko/go-auth/domain"
	"github.com/drasko/go-auth/services"
)

const (
	defaultConfig string = "/src/github.com/drasko/go-auth/config/config.toml"
	httpPort      string = ":8180"
	help          string = `
		Usage: mainflux-auth [options]
		Options:
			-c, --config <file>         Configuration file
			-h, --help                  Prints this message end exits`
)

func main() {
	opts := struct {
		Config string
		Help   bool
	}{}

	flag.StringVar(&opts.Config, "c", "", "Configuration file.")
	flag.StringVar(&opts.Config, "config", "", "Configuration file.")
	flag.BoolVar(&opts.Help, "h", false, "Show help.")
	flag.BoolVar(&opts.Help, "help", false, "Show help.")

	flag.Parse()

	if opts.Help {
		fmt.Printf("%s\n", help)
		os.Exit(0)
	}

	if opts.Config == "" {
		opts.Config = os.Getenv("GOPATH") + defaultConfig
	}

	cfg := config.Config{}
	cfg.Load(opts.Config)

	if cfg.SecretKey != "" {
		domain.SetSecretKey(cfg.SecretKey)
	}

	services.StartCaching(cfg.RedisURL)
	defer services.StopCaching()

	fmt.Println(banner)
	http.ListenAndServe(httpPort, api.Server())
}

var banner = `
 _______    ______       ________   __  __   _________  ___   ___     
/______/\  /_____/\     /_______/\ /_/\/_/\ /________/\/__/\ /__/\    
\::::__\/__\:::_ \ \    \::: _  \ \\:\ \:\ \\__.::.__\/\::\ \\  \ \   
 \:\ /____/\\:\ \ \ \    \::(_)  \ \\:\ \:\ \  \::\ \   \::\/_\ .\ \  
  \:\\_  _\/ \:\ \ \ \    \:: __  \ \\:\ \:\ \  \::\ \   \:: ___::\ \ 
   \:\_\ \ \  \:\_\ \ \    \:.\ \  \ \\:\_\:\ \  \::\ \   \: \ \\::\ \
    \_____\/   \_____\/     \__\/\__\/ \_____\/   \__\/    \__\/ \::\/
                                                                      

                == Sleep well, everything's locked ==
`
