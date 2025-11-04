package plaza

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	jsoniter "github.com/json-iterator/go"
	uuid "github.com/satori/go.uuid"
)

func MakeSureDir(folder string) {
	if !IsDirExist(folder) {
		if err := MakeDir(folder); err != nil {
			panic(err)
		}
	}
}

func PrintBytes(data []byte) {
	for i := 0; i < len(data); i++ {
		if i%16 == 0 && i != 0 {
			fmt.Printf("\r\n")
		}
		fmt.Printf("%02x ", data[i])
	}
	fmt.Println("\r\n--------------------------------")
}

func Md5(data []byte) []byte {
	md5Ctx := md5.New()
	md5Ctx.Write(data)
	return md5Ctx.Sum(nil)
}

func Md5ToString(data []byte) string {
	return hex.EncodeToString(Md5(data))
}

func IsFileExist(filepath string) bool {
	finfo, err := os.Stat(filepath)
	if err != nil {
		return false
	}
	return !finfo.IsDir()
}

func IsDirExist(filepath string) bool {
	finfo, err := os.Stat(filepath)
	if err != nil {
		return false
	}
	return finfo.IsDir()
}

func MakeDir(dir string) error {
	return os.MkdirAll(dir, os.ModePerm) // 等价 0777，受系统默认权限/umask 影响
}

func IteratorFiles(dir string, ext string) []string {
	if !IsDirExist(dir) {
		return nil
	}

	var paths []string
	ext = strings.ToLower("." + ext)
	filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if ext == ".*" {
			paths = append(paths, p)
		} else {
			e := path.Ext(info.Name())
			if strings.ToLower(e) == ext {
				paths = append(paths, info.Name())
			}
		}
		return nil
	})

	return paths
}

func UuidWithoutDash() string {
	return strings.ReplaceAll(uuid.Must(uuid.NewV4(), nil).String(), "-", "")
}

// SlotResponse 是替换 has/core.SlotResponse 的本地 DTO。
type SlotResponse struct {
	Data  interface{} `json:"Data"`
	Error struct {
		Code    int    `json:"Code"`
		Message string `json:"Message"`
	} `json:"Error"`
}

func RequestHttpPost(ip string, port int, api string, data interface{}) (interface{}, error) {
	defer func() {
		err := recover()
		if err != nil {
			logger.Error(err)
		}
	}()
	bytesData, _ := jsoniter.Marshal(data)
	resp, err := http.Post(fmt.Sprintf("http://%s:%d/v1/%s", ip, port, api), "application/json", bytes.NewReader(bytesData))
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var res SlotResponse
	jsoniter.Unmarshal(body, &res)

	if res.Error.Code == 0 {
		return res.Data, nil
	} else {
		return nil, errors.New("发送请求失败")
	}
}

func TodayMidnight() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

func YesterdayMidnight() time.Time {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	return time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, now.Location())
}

func LastDaysMidnight(days int) time.Time {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -days)
	return time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, now.Location())
}

func TomorrowMidnight() time.Time {
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)
	return time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, now.Location())
}

func LastMondayMidnight() time.Time {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	lastMonday := now.AddDate(0, 0, -6-weekday)
	return time.Date(lastMonday.Year(), lastMonday.Month(), lastMonday.Day(), 0, 0, 0, 0, now.Location())
}

func ThisMondayMidnight() time.Time {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	thisMonday := now.AddDate(0, 0, 1-weekday)
	return time.Date(thisMonday.Year(), thisMonday.Month(), thisMonday.Day(), 0, 0, 0, 0, now.Location())
}

func cloneMembers(in []*GroupMember) []*GroupMember {
	if in == nil {
		return nil
	}
	out := make([]*GroupMember, len(in))
	copy(out, in)
	return out
}

func IsChienseString(txt string) bool {
	for _, r := range txt {
		if !unicode.Is(unicode.Han, r) {
			return false
		}
	}
	return true
}
func token(group int, pk string, st string) string {
	s := fmt.Sprintf("groupid=%d&privatekey=%s&servertime=%s", group, pk, st)
	return md5ToString([]byte(s))
}
func md5ToString(data []byte) string {
	return strings.ToUpper(hex.EncodeToString(_md5(data)))
}
func _md5(data []byte) []byte {
	md5Ctx := md5.New()
	md5Ctx.Write(data)
	return md5Ctx.Sum(nil)
}
