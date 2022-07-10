/*
Copyright © 2022 Remi Ferrand

Contributor(s): Remi Ferrand <riton.github_at_gmail(dot)com>, 2022

This software is governed by the CeCILL-B license under French law and
abiding by the rules of distribution of free software.  You can  use,
modify and/ or redistribute the software under the terms of the CeCILL-B
license as circulated by CEA, CNRS and INRIA at the following URL
"http://www.cecill.info".

As a counterpart to the access to the source code and  rights to copy,
modify and redistribute granted by the license, users are provided only
with a limited warranty  and the software's author,  the holder of the
economic rights,  and the successive licensors  have only  limited
liability.

In this respect, the user's attention is drawn to the risks associated
with loading,  using,  modifying and/or developing or reproducing the
software by the user in light of its specific status of free software,
that may mean  that it is complicated to manipulate,  and  that  also
therefore means  that it is reserved for developers  and  experienced
professionals having in-depth computer knowledge. Users are therefore
encouraged to load and test the software's suitability as regards their
requirements in conditions enabling the security of their systems and/or
data to be ensured and,  more generally, to use and operate it in the
same conditions as regards security.

The fact that you are presently reading this means that you have had
knowledge of the CeCILL-B license and that you accept its terms.

*/
package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	apiV1 "github.com/toutdoux-app/toutdoux-cli/api/v1"
	"github.com/toutdoux-app/toutdoux-cli/config"
)

var cfgFile string
var cfg config.Config
var apiClient apiV1.Client

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "toutdoux-cli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) {},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		client, err := apiV1.NewClientWithOptions(apiV1.ClientOptions{
			Endpoint: cfg.API.Endpoint,
			Username: cfg.API.Username,
			Password: cfg.API.Password,
		})
		if err != nil {
			return errors.Wrap(err, "creating api client")
		}

		if err := client.Initialize(); err != nil {
			return errors.Wrap(err, "initializing api client")
		}

		apiClient = client

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.toutdoux/cli.yaml)")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "enable debug mode")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	debug, _ := rootCmd.Flags().GetBool("debug")
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in ~/.toutdoux/ with name "cli" (without extension).
		viper.AddConfigPath(filepath.Join(home, ".toutdoux"))
		viper.SetConfigType("yaml")
		viper.SetConfigName("cli")
	}

	viper.SetEnvPrefix("toutdoux")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logrus.Debugf("Using config file %s", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":       err,
			"config_file": viper.ConfigFileUsed(),
		}).Fatal("fail to unmarshal configuration")
	}
}
