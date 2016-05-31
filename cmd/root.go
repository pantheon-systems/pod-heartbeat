// Copyright © 2016 Pantheon Systems
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/pantheon-systems/pod-heartbeat/pkg/heartbeat"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type flags struct {
	connect      string
	url          *url.URL
	timeout      time.Duration
	timeoutFlag  string
	retries      int
	interval     time.Duration
	intervalFlag string

	configFile string
}

var cfg flags

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "pod-heartbeat",
	Short: "Connect Or Die",
	Long: `Sometimes you want your kube pod to die if it can't get to something.
That could be another container in your pod that has deadlocked.

This program runs connects, and  if it can't connect or times out it will exit.
When ran inside a kube pod the container exit event will cause the pod to be destroyed.`,

	PreRunE: validateRun,
	Run:     runHeartBeat,
}

// vlaidateRun does all the arg checks before we invoke the Run command
func validateRun(cmd *cobra.Command, args []string) error {
	// parse the timeout into a proper Duration
	t, err := time.ParseDuration(cfg.timeoutFlag)
	if err != nil {
		return fmt.Errorf("Unable to parse duration argument: %s", err.Error())
	}
	cfg.timeout = t

	// parse the interval into proper Duration
	t, err = time.ParseDuration(cfg.intervalFlag)
	if err != nil {
		return fmt.Errorf("Unable to parse interval argument: %s", err.Error())
	}
	cfg.interval = t

	// ensure the connect string is a proper uri
	u, err := url.Parse(cfg.connect)
	if err != nil {
		return fmt.Errorf("Unable to parse connect string: %s ", err.Error())
	}

	if u.Scheme == "" {
		return errors.New("connect string can not have an empty scheme, use tcp:// or http://")
	}

	if u.Host == "" {
		return errors.New("connect string can not have an empty host")
	}
	cfg.url = u

	return nil
}

// Functional Main for the application since we only have a root command
func runHeartBeat(cmd *cobra.Command, args []string) {
	c := heartbeat.Check{
		URL:      cfg.url,
		Timeout:  cfg.timeout,
		Retries:  cfg.retries,
		Interval: cfg.interval,
	}

	err := c.Beat()
	if err != nil {
		log.Fatalf("Heartbeat failed: %s\n", err.Error())
	}
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
	RootCmd.Flags().StringVarP(
		&cfg.configFile,
		"config-file",
		"f",
		"",
		"Config file (default is $HOME/.pod-heartbeat.yaml)")

	RootCmd.Flags().StringVarP(
		&cfg.connect,
		"connect",
		"c",
		"tcp://127.0.0.1:4000",
		"Connection URI, valid protocols are  tcp:// and http:// for now",
	)

	RootCmd.Flags().StringVarP(
		&cfg.timeoutFlag,
		"timeout",
		"t",
		"1s",
		"Timeout before considering the connection failed. Valid qualifiers: ns,ms,s,m,h,d",
	)

	RootCmd.Flags().IntVarP(
		&cfg.retries,
		"retries",
		"r",
		3,
		"How many times to retry before exiting.",
	)

	RootCmd.Flags().StringVarP(
		&cfg.intervalFlag,
		"interval",
		"i",
		"5s",
		"Interval for the Heartbeat action.",
	)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfg.configFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfg.configFile)
	}

	viper.SetConfigName(".pod-heartbeat") // name of config file (without extension)
	viper.AddConfigPath("$HOME")          // adding home directory as first search path
	viper.AutomaticEnv()                  // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
