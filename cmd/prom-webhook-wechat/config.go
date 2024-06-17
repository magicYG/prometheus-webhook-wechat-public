package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
	"text/template"
	"time"
	"unicode"
)

var cfg = struct {
	fs *flag.FlagSet

	listenAddress        string
	WechatProfiles       wechatProfilesFlag
	WechatAPIUrlProfiles string
	requestTimeout       time.Duration
	corpid               string
	corpsecret           string
	configdir            string
	templateid           string
}{}

type Config struct {
	TemplateID string         `yaml:"templateid"`
	Appid      string         `yaml:"appid"`
	Secret     string         `yaml:"secret"`
	Chatids    []chatgroupids `yaml:"chatgroups"`
}

type chatgroupids struct {
	Name    string   `yaml:"name"`
	Chatids []string `yaml:"chatids"`
}

func init() {
	cfg.fs = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	cfg.fs.Usage = usage

	cfg.fs.StringVar(&cfg.listenAddress, "web.listen-address", ":8060",
		"Address to listen on for web interface.",
	)
	// cfg.fs.StringVar(&cfg.WechatAPIUrlProfiles, "wechat.apiurl", "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=",
	// 	"Custom wechat api url ",
	// )
	cfg.fs.StringVar(&cfg.WechatAPIUrlProfiles, "wechat.apiurl", "api.weixin.qq.com",
		"Custom wechat api url ",
	)
	cfg.fs.DurationVar(&cfg.requestTimeout, "wechat.timeout", 5*time.Second,
		"Timeout for invoking wechat webhook.",
	)
	cfg.fs.StringVar(&cfg.configdir, "config.file", "/config/config.yaml",
		"config yaml path",
	)

	configBytes, err := os.ReadFile(cfg.configdir)
	if err != nil {
		fmt.Print("msg: ", "Load config file error: ", err)
		os.Exit(10)
	}
	var config Config
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		fmt.Print("msg: ", "Unmarshal config file error: ", err)
		os.Exit(10)
	}
	cfg.templateid = config.TemplateID
	cfg.corpid = config.Appid
	cfg.corpsecret = config.Secret
	cfg.WechatProfiles.Set(config.Chatids)

}

func parse(args []string) error {
	err := cfg.fs.Parse(args)
	if err != nil || len(cfg.fs.Args()) != 0 {
		if err != flag.ErrHelp {
			fmt.Fprintf(os.Stderr, "Invalid command line arguments. Help: %s -h", os.Args[0])
		}
		if err == nil {
			err = fmt.Errorf("Non-flag argument on command line: %q", cfg.fs.Args()[0])
		}
		return err
	}

	return nil
}

var helpTmpl = strings.TrimSpace(`
usage: prom-webhook-wechat [<args>]
{{ range $cat, $flags := . }}{{ if ne $cat "." }} == {{ $cat | upper }} =={{ end }}
  {{ range $flags }}
   -{{ .Name }} {{ .DefValue | quote }}
      {{ .Usage | wrap 80 6 }}
  {{ end }}
{{ end }}
`)

func usage() {
	t := template.New("usage")
	t = t.Funcs(template.FuncMap{
		"wrap": func(width, indent int, s string) (ns string) {
			width = width - indent
			length := indent
			for _, w := range strings.SplitAfter(s, " ") {
				if length+len(w) > width {
					ns += "\n" + strings.Repeat(" ", indent)
					length = 0
				}
				ns += w
				length += len(w)
			}
			return strings.TrimSpace(ns)
		},
		"quote": func(s string) string {
			if len(s) == 0 || s == "false" || s == "true" || unicode.IsDigit(rune(s[0])) {
				return s
			}
			return fmt.Sprintf("%q", s)
		},
		"upper": strings.ToUpper,
	})
	t = template.Must(t.Parse(helpTmpl))

	groups := make(map[string][]*flag.Flag)

	// Bucket flags into groups based on the first of their dot-separated levels.
	cfg.fs.VisitAll(func(fl *flag.Flag) {
		parts := strings.SplitN(fl.Name, ".", 2)
		if len(parts) == 1 {
			groups["."] = append(groups["."], fl)
		} else {
			name := parts[0]
			groups[name] = append(groups[name], fl)
		}
	})
	for cat, fl := range groups {
		if len(fl) < 2 && cat != "." {
			groups["."] = append(groups["."], fl...)
			delete(groups, cat)
		}
	}

	if err := t.Execute(os.Stdout, groups); err != nil {
		panic(fmt.Errorf("error executing usage template: %s", err))
	}
}

type wechatProfilesFlag struct {
	chatids map[string][]string
}

func (c *wechatProfilesFlag) Set(opt []chatgroupids) error {
	if c.chatids == nil {
		c.chatids = make(map[string][]string)
	}
	for _, value := range opt {
		for _, id := range value.Chatids {
			c.chatids[value.Name] = append(c.chatids[value.Name], id)
		}
	}
	return nil
}
