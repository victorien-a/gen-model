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
	"errors"
	"github.com/DaoYoung/gen-model/handler"
	"github.com/DaoYoung/gen-model/manager/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create model struct",
	Long: `with mysql connect, generate model file`,
	Args:cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		generateModel()
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&genRequest.SearchTableName,"searchTableName","t","", "set your searchTableName, support patten with '*'")
	createCmd.Flags().BoolVarP(&genRequest.IsLowerCamelCaseJson,"isLowerCamelCaseJson","i",true, "set IsLowerCamelCaseJson true/false")
	flagBindviper(createCmd, false,"searchTableName","searchTableName")
	flagBindviper(createCmd, false,"isLowerCamelCaseJson","isLowerCamelCaseJson")
}

func validArgs() error {
	if viper.GetString("mysql.host") == ""{
		return errors.New("mysql.host is empty")
	}
	if viper.GetString("mysql.database") == ""{
		return errors.New("mysql.database is empty")
	}
	if viper.GetString("mysql.username") == ""{
		return errors.New("mysql.username is empty")
	}
	if viper.GetString("mysql.password") == ""{
		return errors.New("mysql.password is empty")
	}
	if genRequest.SearchTableName == ""{
		return errors.New("tableName is empty")
	}

	if genRequest.OutPutPath == ""{
		return errors.New("outPutPath is empty")
	}
	return  nil
}
func generateModel()  {
	genRequest.SetDataByViper()
	log.Printf("%+v", genRequest)
	if err := validArgs();err != nil{
		log.Println(err)
		os.Exit(1)
	}
	if err := db.InitDb();err != nil{
		log.Println(err)
		os.Exit(1)
	}
	handler.Table2struct(&genRequest)
}