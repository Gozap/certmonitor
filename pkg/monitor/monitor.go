package monitor

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Gozap/certmonitor/pkg/alarm"

	"github.com/robfig/cron"

	"github.com/Sirupsen/logrus"

	"github.com/spf13/viper"

	"github.com/Gozap/certmonitor/pkg/utils"
)

type Config struct {
	WebSites   []string
	Cron       string
	BeforeTime time.Duration
}

func ExampleConfig() Config {
	return Config{
		WebSites: []string{
			"https://google.com",
			"https://mritd.me",
		},
		Cron:       "@every 1h",
		BeforeTime: 7 * 24 * time.Hour,
	}
}

func check(address string, beforeTime time.Duration) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(address)
	if !utils.CheckErr(err) {
		return err
	}
	defer resp.Body.Close()

	for _, cert := range resp.TLS.PeerCertificates {
		if !cert.NotAfter.After(time.Now()) {
			return errors.New(fmt.Sprintf("Website [%s] certificate has expired: %s", address, cert.NotAfter.Format("2006-01-02 15:04:05")))
		}

		if cert.NotAfter.Sub(time.Now()) < beforeTime {
			return errors.New(fmt.Sprintf("Website [%s] certificate will expire, remaining time: %sh", address, cert.NotAfter.Sub(time.Now()).Hours()))
		}
	}

	return nil
}

func Start() {
	var config Config
	err := viper.UnmarshalKey("monitor", &config)
	if err != nil {
		logrus.Fatalf("Can't parse server config: %s", err)
	}

	c := cron.New()

	for _, website := range config.WebSites {
		addr := website
		c.AddFunc(config.Cron, func() {
			err := check(addr, config.BeforeTime)
			if err != nil {
				alarm.Alarm(err.Error())
			}
		})
	}
}
