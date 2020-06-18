package main

import (
	"bufio"
	"fmt"
	"github.com/axgle/mahonia"
	ipdb "github.com/ipipdotnet/ipdb-go"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

const (
	version = "NaLi 1.5.0\n"
	helper  = "Usage: %s <command> [options] \n" +
		"\nOptions:" +
		"\n  -v, --version" +
		"\n  -h, --help\n" +
		"\n  -c, --color\n" +
		"\nCommands:" +
		"\n  IP Address\n"
)

var (
	help  = []string{"help", "--help", "-h", "h"}
	ver   = []string{"version", "--version", "-v", "v"}
	color = []string{"color", "--color", "-c", "c"}
)

type AppCfg struct {
	color bool
}

func main() {
	cfg := &AppCfg{}
	cfg.color = false

	if args := os.Args; len(args) > 1 {
		cmd(args[1], cfg)
		args = args[1:]
		for _, v := range args {
			if contains(help, v) || contains(ver, v) || contains(color, v) {
				continue
			}
			fmt.Println(Analyse(v, cfg))
		}
		os.Exit(0)
	}

	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if (info.Mode() & os.ModeCharDevice) != 0 {
		self := os.Args[0]
		fmt.Printf(helper, self)
		os.Exit(0)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		line = ConvertToString(line, "gbk", "utf-8")
		fmt.Printf("%s", Analyse(line, cfg))
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
}

func Analyse(item string, cfg *AppCfg) string {
	re4 := regexp.MustCompile(`((25[0-5]|(2[0-4]|1?[0-9])?[0-9])\.){3}(25[0-5]|(2[0-4]|1?[0-9])?[0-9])`)
	if ip := re4.FindStringSubmatch(item); len(ip) != 0 {
		res := findIpV4(ip[0])
		res = strings.Trim(res, " ")
		result := ""
		if cfg.color {
			result = ip[0] + " " + "\x1b[0;0;36m[" + res + "]\x1b[0m"
		} else {
			result = ip[0] + " " + "<" + res + ">"
		}

		return strings.ReplaceAll(item, ip[0], result)
	}

	return item
}

func findIpV4(ip string) string {
	addr := ""
	db, err := loadDb()
	if err != nil {
		return addr
	}

	arr, err := db.Find(ip, "CN")
	if err != nil {
		return addr
	}

	return addr + strings.Join(arr, " ")
}

func loadDb() (*ipdb.City, error) {
	ipdbFile := "D:\\test\\ipipfree_ipdb.dat"
	if _, err := os.Stat(ipdbFile); err == nil || os.IsExist(err) {
		return ipdb.NewCity(ipdbFile)
	}

	ipdbFile = "C:\\ipipfree_ipdb.dat"
	if _, err := os.Stat(ipdbFile); err == nil || os.IsExist(err) {
		return ipdb.NewCity(ipdbFile)
	}

	ipdbFile = "ipipfree_ipdb.dat"
	return ipdb.NewCity(ipdbFile)
}

func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

func contains(array []string, flag string) bool {
	for i := 0; i < len(array); i++ {
		if array[i] == flag {
			return true
		}
	}
	return false
}

func cmd(opt string, cfg *AppCfg) {
	if contains(help, opt) {
		fmt.Printf(helper, os.Args[0])
		os.Exit(0)
	} else if contains(ver, opt) {
		_version()
		os.Exit(0)
	} else if contains(color, opt) {
		cfg.color = true
	}
}

func _version() {
	fmt.Println(version)

	db, err := loadDb()
	if err != nil {
		fmt.Printf("IPv4 Version： Database Not Found.\n")
		log.Fatal(err)
		return
	}

	if db.IsIPv4() {
		fmt.Printf("ipdb BuildTime: %s\n", db.BuildTime())
		fmt.Printf("ipdb Fields: %v\n", db.Fields())
	}
}
