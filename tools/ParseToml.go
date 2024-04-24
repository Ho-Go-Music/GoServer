package tools

import (
	"fmt"
	"github.com/BurntSushi/toml"
	diylog "github.com/Ho-Go-Music/GoServer/log"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
)

type Config struct {
	MySQLServer struct {
		Host      string `toml:"host"`
		Port      string `toml:"port"`
		User      string `toml:"user"`
		Password  string `toml:"password"`
		Database  string `toml:"database"`
		Charset   string `toml:"charset"`
		ParseTime string `toml:"parseTime"`
		Loc       string `toml:"loc"`
		Dev       struct {
			Host      string `toml:"host"`
			Port      string `toml:"port"`
			User      string `toml:"user"`
			Password  string `toml:"password"`
			Database  string `toml:"database"`
			Charset   string `toml:"charset"`
			ParseTime string `toml:"parseTime"`
			Loc       string `toml:"loc"`
		} `toml:"dev"`
	} `toml:"mysql-server"`
	RedisServer struct {
		Host           string `toml:"host"`
		Port           string `toml:"port"`
		Password       string `toml:"password"`
		Database       int    `toml:"database"`
		MaxActiveConns int    `toml:"maxActiveConns"`
		MaxIdleConns   int    `toml:"maxIdleConns"`
	} `toml:"redis"`
}

var (
	once sync.Once
	Conf *Config
)

func init() {
	once.Do(ParseToml)
}
func ParseToml() {
	// 获取当前运行堆栈的帧信息
	// Retrieve stack frame information for the current execution stack
	//_, filename, _, _ := runtime.Caller(0)
	absPath, err := filepath.Abs("config.toml")
	if err != nil {
		diylog.Sugar.Errorln("Retrieving configuration file path error\n")
	}
	//Retrieve file status information
	if _, err := os.Stat(absPath); err != nil {
		panic(err)
	}
	// decode toml file
	var config Config
	meta, err := toml.DecodeFile(absPath, &config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%#v", err)
		os.Exit(1)
	} else {
		Conf = &config
	}
	// display toml file's structure
	indent := strings.Repeat(" ", 14)
	fmt.Print("Decoded")
	typ, val := reflect.TypeOf(config), reflect.ValueOf(config)
	for i := 0; i < typ.NumField(); i++ {
		indent := indent
		if i == 0 {
			indent = strings.Repeat(" ", 7)
		}
		fmt.Printf("%s%-11s → %v\n", indent, typ.Field(i).Name, val.Field(i))
	}
	fmt.Print("\nKeys")
	keys := meta.Keys()
	sort.Slice(keys, func(i, j int) bool { return keys[i].String() < keys[j].String() })
	for i, k := range keys {
		indent := indent
		if i == 0 {
			indent = strings.Repeat(" ", 10)
		}
		log.Printf("%s%-10s %s\n", indent, meta.Type(k...), k)
	}
	fmt.Print("\nUndecoded")
	keys = meta.Undecoded()
	sort.Slice(keys, func(i, j int) bool { return keys[i].String() < keys[j].String() })
	for i, k := range keys {
		indent := indent
		if i == 0 {
			indent = strings.Repeat(" ", 5)
		}
		log.Printf("%s%-10s %s\n", indent, meta.Type(k...), k)
	}
}
