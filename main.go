package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"os"

	"github.com/apex/log"

	"github.com/pyama86/wazuh-tailer/wazuh_notifier"
)

type Notifier interface {
	Notify(*wazuh_notifier.Alert) error
}

func newNotifier(c *wazuh_notifier.Config) Notifier {
	switch c.Notifier {
	case "slack":
		return wazuh_notifier.NewSlack(c)
	}
	return nil
}
func main() {
	path := flag.String("config", "/var/ossec/etc/wazuh-notifier.toml", "config file path")
	flag.Parse()
	config, err := wazuh_notifier.NewConfig(*path)
	if err != nil {
		log.Fatal(err.Error())
	}

	notifier := newNotifier(config)
	stdin := bufio.NewScanner(os.Stdin)
	if stdin.Scan() {
		a := wazuh_notifier.Alert{}
		err := json.Unmarshal(stdin.Bytes(), &a)
		if err != nil {
			log.Fatal(err.Error())
		}
		err = notifier.Notify(&a)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
