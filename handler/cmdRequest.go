package handler

import (
    "github.com/spf13/viper"
    "strings"
    "path/filepath"
    "os"
    "sync"
    "io/ioutil"
    "path"
    "regexp"
)

type CmdRequest struct {
    Db  dbConfig
    Gen genConfig
    Wg  sync.WaitGroup
}

type dbConfig struct {
    Host     string
    Database string
    Username string
    Password string
    Port     int
}

type genConfig struct {
    SearchTableName      string
    OutPutPath           string
    IsLowerCamelCaseJson bool
    HasGormTag           bool   // gorm tag, `gorm:"column:name"`
    HasJsonTag           bool   // json tag, `json:"age"`
    HasGureguNullPackage bool   // have package: "gopkg.in/guregu/null.v3"
    ModelSuffix          string // model name suffix
    SourceType           string // self-table: struct create by connect mysql tables local: struct create by local mappers gen-table: struct create by table "gen_model_mapper"
    PersistType          string // persist struct mappers at local or db
    LocalMapperPath      string
}

const (
    sourceSelfTable = "self-table"
    sourceLocal     = "local-mapper"
    sourceGenTable  = "gen-table"
)

func (g *CmdRequest) getTables() []string {
    if strings.Contains(g.Gen.SearchTableName, "*") {
        return matchTables(g.Db.Database, g.Gen.SearchTableName)
    }
    return []string{g.Gen.SearchTableName}
}


func (g *CmdRequest) getOutPutPath() string {
    if g.Gen.OutPutPath == "" {
        g.Gen.OutPutPath = "model"
    }
    p, _ := filepath.Abs(g.Gen.OutPutPath)
    outDir := filepath.Dir(p)
    mkdir(outDir)
    return p
}

func (g *CmdRequest) getAbsPathAndPackageName() (absPath, packageName string) {
    if g.Gen.OutPutPath == "" {
        g.Gen.OutPutPath = "model"
    }
    var err error
    var appPath string
    if absPath, err = filepath.Abs(g.Gen.OutPutPath); err != nil {
        printErrorAndExit(err)
    }
    if !isExist(absPath) {
        printMessageAndExit("OutPutPath not exist: " + absPath)
    }
    if appPath, err = os.Getwd(); err != nil {
        printErrorAndExit(err)
    }
    if absPath == appPath {
        packageName = "main"
    } else {
        _, packageName = filepath.Split(absPath)
    }
    return absPath, packageName
}

func (g *CmdRequest) SetDataByViper() {
    g.Gen.SearchTableName = viper.GetString("gen.searchTableName")
    g.Gen.OutPutPath = viper.GetString("gen.outPutPath")
    g.Gen.IsLowerCamelCaseJson = viper.GetBool("gen.isLowerCamelCaseJson")
    g.Gen.ModelSuffix = viper.GetString("gen.modelSuffix")
    g.Gen.SourceType = viper.GetString("gen.sourceType")
    g.Gen.PersistType = viper.GetString("gen.persistType")
    g.Db.Host = viper.GetString("mysql.host")
    g.Db.Database = viper.GetString("mysql.database")
    g.Db.Port = viper.GetInt("mysql.port")
    g.Db.Username = viper.GetString("mysql.username")
    g.Db.Password = viper.GetString("mysql.password")
}
func (cmdRequest *CmdRequest) selfTable2Struct() {
    initSchemaDb()
    tables := cmdRequest.getTables()
    for _, tn := range tables {
        cmdRequest.Wg.Add(1)
        go mkStructFromSelfTable(tn, cmdRequest)
    }
    cmdRequest.Wg.Wait()
    os.Exit(0)
}
func (g *CmdRequest) localMap2Struct() {
    modelPath, packageName := g.getAbsPathAndPackageName()
    files, _ := ioutil.ReadDir(modelPath)
    for _, f := range files {
        fn := f.Name()
        suffix := path.Ext(fn)
        if suffix == YamlExt {
            fileName := strings.TrimSuffix(fn, suffix)
            if isFileNameMatch(g.Gen.SearchTableName, g.Gen.ModelSuffix, fileName) {
                g.Wg.Add(1)
                go mkStructFromYaml(g, fileName, packageName, modelPath)
            }
        }
    }
    g.Wg.Wait()
    os.Exit(0)
}
func (cmdRequest *CmdRequest) genTable2Struct() {
    initSchemaDb()
    initGenDb()
    tables := cmdRequest.getTables()
    for _, tn := range tables {
        cmdRequest.Wg.Add(1)
        go mkStructFromGenTable(tn, cmdRequest)
    }
    cmdRequest.Wg.Wait()
    os.Exit(0)
}
func (g *CmdRequest) CreateModelStruct() {
    switch g.Gen.SourceType {
    case sourceSelfTable:
        g.selfTable2Struct()
        break
    case sourceLocal:
        g.localMap2Struct()
        break
    case sourceGenTable:
        g.genTable2Struct()
        break
    default:
        printMessageAndExit("wrong sourceType, set value with \"" + sourceSelfTable + "\" or \"" + sourceLocal + "\" or \"" + sourceGenTable + "\"")
    }
}

func isFileNameMatch(pattern, suffix, fileName string) bool {
    fileName = strings.TrimSuffix(fileName, YamlMap)
    pattern = camelString(fileName)
    if strings.Contains(pattern, "*") {
        isMatch, _ := regexp.MatchString(strings.Replace(pattern, "*", "(.*)", -1), "Hello World!")
        return isMatch
    }
    return fileName == pattern
}
