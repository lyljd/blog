package common

import (
	"blog/config"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func DateDay(date time.Time) string {
	return date.Format("2006-01-02 15:04:05")
}

func TrimHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")
	return strings.TrimSpace(src)
}

func GetRequestJsonParam(r *http.Request) map[string]interface{} {
	var params map[string]interface{}
	body, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(body, &params)
	return params
}

func Sha256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

func NewToken() string {
	exp := strconv.FormatInt(time.Now().Add(time.Hour*24*7).Unix(), 10)
	return exp + "." + Sha256(exp+config.Cfg.System.Username+config.Cfg.System.Password)
}

func CheckToken(token string) (bool, string) {
	index := strings.Index(token, ".")
	if index < 1 {
		return false, "请登录"
	}

	exp, sig := token[:index], token[index+1:]
	trueSig := Sha256(exp + config.Cfg.System.Username + config.Cfg.System.Password)
	if sig != trueSig {
		return false, "身份异常，请重新登录"
	}

	timeExp, _ := strconv.ParseInt(exp, 10, 64)
	timeNow := time.Now().Unix()
	if timeExp <= timeNow {
		return false, "身份已过期，请重新登录"
	}

	return true, ""
}
