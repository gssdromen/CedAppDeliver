package utils

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"crypto/md5"

	"encoding/hex"

	"github.com/revel/revel"
)

// TrimQuote 去除字符串前后的单引号和双引号
func TrimQuote(str string) string {
	for strings.HasPrefix(str, "'") || strings.HasPrefix(str, "\"") {
		str = str[1:]
	}
	for strings.HasSuffix(str, "'") || strings.HasSuffix(str, "\"") {
		str = str[:len(str)-1]
	}
	return str
}

// GetBaseURL 得到基础的URL
func GetBaseURL() string {
	return revel.HttpAddr + ":" + strconv.Itoa(revel.HttpPort)
}

// GetLocalIP 得到本机的IP
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func getMD5(b []byte) string {
	h := md5.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

// RandSeq 得到随机字符串
func RandSeq(n int) string {
	letters := []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range b {
		rand := r.Intn(len(letters))
		b[i] = letters[rand]
	}
	return string(b)
}
