/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
    "github.com/mitchellh/go-homedir"
    "github.com/spf13/viper"
    "log"
    "github.com/DaoYoung/gen-model/handler"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
    Use:   "gen-model",
    Short: "A brief description of your application",
    Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
    // Uncomment the following line if your bare application
    // has an action associated with it:
    // 	Run: func(cmd *cobra.Command, args []string) { },
}
var genRequest handler.GenRequest

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func init() {
    cobra.OnInitialize(initConfig)
    dir, _ := os.Getwd()
    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is "+dir+"/"+handler.Yamlfile+".yaml)")

    rootCmd.PersistentFlags().StringVar(&genRequest.DbConfig.Host, "host", "localhost", "set DB host")
    rootCmd.PersistentFlags().StringVarP(&genRequest.DbConfig.Database, "database", "d", "foo", "set your database")
    rootCmd.PersistentFlags().IntVarP(&genRequest.DbConfig.Port, "port", "p", 3306, "set DB port")
    rootCmd.PersistentFlags().StringVarP(&genRequest.DbConfig.Username, "username", "u", "root", "set DB login username")
    rootCmd.PersistentFlags().StringVarP(&genRequest.DbConfig.Password, "password", "w", "", "set DB login password")
    viper.BindPFlag("mysql.host", rootCmd.PersistentFlags().Lookup("host"))
    viper.BindPFlag("mysql.database", rootCmd.PersistentFlags().Lookup("database"))
    viper.BindPFlag("mysql.port", rootCmd.PersistentFlags().Lookup("port"))
    viper.BindPFlag("mysql.username", rootCmd.PersistentFlags().Lookup("username"))
    viper.BindPFlag("mysql.password", rootCmd.PersistentFlags().Lookup("password"))
    rootCmd.PersistentFlags().StringVarP(&genRequest.OutPutPath, "outPutPath", "o", ".", "set your OutPutPath")
    viper.BindPFlag("outPutPath", rootCmd.PersistentFlags().Lookup("outPutPath"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
    if cfgFile != "" {
        // Use config file from the flag.
        viper.SetConfigFile(cfgFile)
    } else {
        // Find home directory.
        home, err := homedir.Dir()
        log.Println("home dir ", home)
        dir, _ := os.Getwd()
        log.Println("project dir ", dir)
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        // Search config in home directory with name ".gen-model" (without extension).
        viper.AddConfigPath(dir)
        viper.SetConfigName(handler.Yamlfile)
    }

    viper.AutomaticEnv() // read in environment variables that match

    // If a config file is found, read it in.
    if err := viper.ReadInConfig(); err == nil {
        fmt.Println("Using config file:", viper.ConfigFileUsed())
    }
}
