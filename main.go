package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/zhangyiming748/GetFileInfo"
	"github.com/zhangyiming748/goini"
	"golang.org/x/exp/slog"
	"io"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

const (
	configPath = "./settings.ini"
)

type RB struct {
	Msgtype string `json:"msgtype"`
	Image   struct {
		Base64 string `json:"base64"`
		Md5    string `json:"md5"`
	} `json:"image"`
}

var (
	conf   *goini.Config
	logger *slog.Logger
)

func setLevel(level string) {
	var opt slog.HandlerOptions
	switch level {
	case "Debug":
		opt = slog.HandlerOptions{ // 自定义option
			AddSource: true,
			Level:     slog.LevelDebug, // slog 默认日志级别是 info
		}
	case "Info":
		opt = slog.HandlerOptions{ // 自定义option
			AddSource: true,
			Level:     slog.LevelInfo, // slog 默认日志级别是 info
		}
	case "Warn":
		opt = slog.HandlerOptions{ // 自定义option
			AddSource: true,
			Level:     slog.LevelWarn, // slog 默认日志级别是 info
		}
	case "Err":
		opt = slog.HandlerOptions{ // 自定义option
			AddSource: true,
			Level:     slog.LevelError, // slog 默认日志级别是 info
		}
	default:
		slog.Warn("需要正确设置环境变量 Debug,Info,Warn or Err")
		slog.Info("默认使用Debug等级")
		opt = slog.HandlerOptions{ // 自定义option
			AddSource: true,
			Level:     slog.LevelDebug, // slog 默认日志级别是 info
		}
	}
	file := "pic2base64.log"
	logf, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}
	logger = slog.New(opt.NewJSONHandler(io.MultiWriter(logf, os.Stdout)))
}
func main() {
	conf = goini.SetConfig(configPath)
	level, not := conf.GetValue("log", "level")
	if not != nil {
		panic(not)
	}
	setLevel(level)
	src := getRoot()
	pattern := "jpg;png"
	files := GetFileInfo.GetAllFileInfo(src, pattern, level)
	for _, file := range files {
		MD5, _ := getMD5(file)
		BASE64, _ := getBase64(file)
		basedir := strings.Trim(file.FullPath, path.Ext(file.FullPath)) //带最后一个分隔符
		logger.Info("全路径-扩展名", slog.String("basedir", basedir))
		logger.Info("单个文件", slog.String("MD5", MD5), slog.String("base64", BASE64))
		md5File, err := os.OpenFile(strings.Join([]string{basedir, "md5"}, "."), 8|512|2, 0777)
		if err != nil {
			return
		}
		defer md5File.Close()
		md5File.WriteString(MD5)
		base64File, err := os.OpenFile(strings.Join([]string{basedir, "base64"}, "."), 8|512|2, 0777)
		if err != nil {
			return
		}
		defer base64File.Close()
		base64File.WriteString(BASE64)
		var j = &RB{
			Msgtype: "image",
		}
		j.setMd5(MD5)
		j.setBase64(BASE64)
		marshal, err := json.Marshal(j)
		if err != nil {
			return
		}
		jsonFile, err := os.OpenFile(strings.Join([]string{basedir, "json"}, "."), 8|512|2, 0777)
		if err != nil {
			return
		}
		defer jsonFile.Close()
		jsonFile.WriteString(string(marshal))

	}
}
func (r *RB) setMd5(s string) *RB {
	r.Image.Md5 = s
	return r
}
func (r *RB) setBase64(s string) *RB {
	r.Image.Base64 = s
	return r
}
func getRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Dir(filename)
}
func getMD5(file GetFileInfo.Info) (string, error) {
	f, err := os.Open(file.FullPath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	hash := md5.New()
	_, err = io.Copy(hash, f)
	if err != nil {
		return "", err
	}
	MD5 := fmt.Sprintf("%x", hash.Sum(nil))
	//logger.Info("文件计算MD5", slog.String("文件名", f.Name()), slog.String("MD5", MD5))
	return MD5, nil
}
func getBase64(file GetFileInfo.Info) (string, error) {
	var cmd *exec.Cmd
	//var output []byte
	if runtime.GOOS == "darwin" {
		cmd = exec.Command("base64", "-i", file.FullPath)
	} else if runtime.GOOS == "linux" {
		cmd = exec.Command("base64", "-w", "0", file.FullPath)
	} else {
		logger.Warn("垃圾Windows不支持")
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Warn("编码base64出错", slog.String("文件名", file.FullPath), slog.Any("错误", err))
		return "", err
	}
	base := strings.Trim(string(output), "\n")
	return base, nil
}
