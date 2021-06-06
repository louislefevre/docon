package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var configFile string

var rootCmd = &cobra.Command{
	Use:   "docon",
	Short: "A brief description of your application",
	Long:  `Docon is a command line tool used for maintaining Linux dotfiles.`,
	//Run: func(cmd *cobra.Command, args []string) {},
}

type Configuration map[string]ConfigGroup

type ConfigGroup struct {
	Path    string   `mapstructure:"path"`
	Include []string `mapstructure:"include"`
	Exclude []string `mapstructure:"exclude"`
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.config/docon/config.yaml)")
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		viper.AddConfigPath(fmt.Sprintf("%s/.config/docon/", home))
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("missing config file: %s", err))
		} else {
			panic(fmt.Errorf("an error occurred: %s", err))
		}
	}

	config := make(Configuration)
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("cannot parse config file: %s", err))
	}

	for name, group := range config {
		if group.Path == "" {
			panic(fmt.Errorf("cannot parse config file: %s has no defined path", name))
		}
		fmt.Println(name)
		fmt.Println(group.Path)
		fmt.Println(group.Include)
		fmt.Println(group.Exclude)
		fmt.Println()
	}
}
