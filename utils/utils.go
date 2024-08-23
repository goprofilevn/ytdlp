package utils

import (
	"context"
	"errors"
	"fmt"
	"github.com/oklog/ulid/v2"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

func CheckOrCreateDir(dir string) error {
	_, errStat := os.Stat(dir)
	if os.IsNotExist(errStat) {
		if errMkdir := os.MkdirAll(dir, 0755); errMkdir != nil {
			return errMkdir
		}
	}
	return nil
}

func CheckOrDeleteFile(pathFile string) error {
	if _, err := os.Stat(pathFile); err == nil {
		if err := os.Remove(pathFile); err != nil {
			return err
		}
	}
	return nil
}

func CheckOrDeleteDir(dir string) error {
	if _, err := os.Stat(dir); err == nil {
		if err := os.RemoveAll(dir); err != nil {
			return err
		}
	}
	return nil
}

func CheckIsExistDir(dir string) bool {
	_, err := os.Stat(dir)
	return !os.IsNotExist(err)
}

func GetPathProfile(id string) string {
	homeDir := GetHomeDir()
	return filepath.Join(homeDir, "user-data-dir", id)
}

func GetZipProfilePath(id string) string {
	homeDir := GetHomeDir()
	return filepath.Join(homeDir, "user-data-dir", fmt.Sprintf("%s.zip", id))
}

func GetCwd() string {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return currentDir
}

func GetHomeDir() string {
	userHomeDir, _ := os.UserHomeDir()
	return filepath.Join(userHomeDir, ".ytdlp")
}

func GetResourceDir() string {
	homeDir := GetHomeDir()
	return filepath.Join(homeDir, "resources")
}

func GetTempDir() string {
	homeDir := GetHomeDir()
	return filepath.Join(homeDir, "temp")
}

func GetDownloadDir() string {
	homeDir := GetHomeDir()
	return filepath.Join(homeDir, "download")
}

func GetExtensionDir() string {
	homeDir := GetHomeDir()
	return filepath.Join(homeDir, "extensions")
}

func GetFFmpegPath() string {
	resourceDir := GetResourceDir()
	if runtime.GOOS == "windows" {
		return filepath.Join(resourceDir, "ffmpeg", "ffmpeg-N-116752-g507c2a5774-win64-gpl", "bin", "ffmpeg.exe")
	}
	return "ffmpeg"
}

func GetYtDlpPath() string {
	resourceDir := GetResourceDir()
	if runtime.GOOS == "windows" {
		return filepath.Join(resourceDir, "yt-dlp", "yt-dlp.exe")
	}
	return "youtube-dl"
}

func getCurrentVersion() (string, error) {
	browserFolder := filepath.Join(GetResourceDir(), "browser")
	folders, errFolders := os.ReadDir(browserFolder)
	if errFolders != nil {
		return "", errFolders
	}
	versions := make([]string, 0)
	for _, folder := range folders {
		if folder.Name() != ".DS_Store" {
			versions = append(versions, folder.Name())
		}
	}
	if len(versions) == 0 {
		return "", nil
	}
	sort.Slice(versions, func(i, j int) bool {
		ver := strings.Split(versions[i], ".")
		ver2 := strings.Split(versions[j], ".")
		for index, value := range ver {
			if index < len(ver2) {
				if value < ver2[index] {
					return true
				}
			}
		}
		return false
	})
	return versions[len(versions)-1], nil
}

func GetPlatform() string {
	switch platform := runtime.GOOS; platform {
	case "darwin":
		if runtime.GOARCH == "arm64" {
			return "darwin_x64"
		} else {
			return "darwin_x64"
		}
	case "linux":
		return "linux"
	case "windows":
		return "windows"
	default:
		return "linux"
	}
}

func CopyDir(src, dest string) error {
	srcInfo, errSrcInfo := os.Stat(src)
	if errSrcInfo != nil {
		return errSrcInfo
	}
	if errMkdir := os.MkdirAll(dest, srcInfo.Mode()); errMkdir != nil {
		return errMkdir
	}
	directory, errReadDir := os.ReadDir(src)
	if errReadDir != nil {
		return errReadDir
	}
	for _, item := range directory {
		srcPath := filepath.Join(src, item.Name())
		destPath := filepath.Join(dest, item.Name())
		if item.IsDir() {
			if errCopy := CopyDir(srcPath, destPath); errCopy != nil {
				return errCopy
			}
		} else {
			if errCopy := CopyFile(srcPath, destPath); errCopy != nil {
				return errCopy
			}
		}
	}
	return nil
}

func CopyFile(src, dest string) error {
	srcFile, errSrcFile := os.Open(src)
	if errSrcFile != nil {
		return errSrcFile
	}
	defer func(srcFile *os.File) {
		if errSrcFileClose := srcFile.Close(); errSrcFileClose != nil {
			log.Println(errSrcFileClose.Error())
		}
	}(srcFile)
	destFile, errDestFile := os.Create(dest)
	if errDestFile != nil {
		return errDestFile
	}
	defer func(destFile *os.File) {
		if errDestFileClose := destFile.Close(); errDestFileClose != nil {
			log.Println(errDestFileClose.Error())
		}
	}(destFile)
	if _, errCopy := io.Copy(destFile, srcFile); errCopy != nil {
		return errCopy
	}
	return nil
}

type BoudingBox struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Left   float64 `json:"left"`
	Right  float64 `json:"right"`
	Top    float64 `json:"top"`
	Bottom float64 `json:"bottom"`
}

func GetPositionClick(boundingBox BoudingBox, paddingWidth float64, paddingHeight float64) (float64, float64) {
	centerX := boundingBox.X + boundingBox.Width/2
	centerY := boundingBox.Y + boundingBox.Height/2
	width := boundingBox.Width * paddingWidth
	height := boundingBox.Height * paddingHeight
	x := RandomFloat64(centerX-width/2, centerX+width/2)
	y := RandomFloat64(centerY-height/2, centerY+height/2)
	return x, y
}

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func RandomInt64(min, max int64) int64 {
	return min + rand.Int63n(max-min)
}

func RandomFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func RandomFloat64Min(min, max float64) float64 {
	if max <= min {
		return min
	}
	return min + rand.Float64()*(max-min)
}

func ByteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

func ParseFilename(keyString string) (filename string) {
	ss := strings.Split(keyString, "/")
	s := ss[len(ss)-1]
	return s
}

func Capitalize(str string) string {
	if str == "id" {
		return "ID"
	} else {
		return strings.ToUpper(str[0:1]) + str[1:]
	}
}

func InArray(array []string, value string) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

func GetTime(strTime string) int {
	infoTime := strings.TrimSpace(strTime)
	if strings.Contains(infoTime, "(") {
		infoTime = strings.Split(infoTime, "(")[1]
		infoTime = strings.Split(infoTime, ")")[0]
		infos := strings.Split(infoTime, ",")
		if len(infos) == 2 {
			from, errFrom := strconv.Atoi(infos[0])
			if errFrom != nil {
				from = 0
			}
			to, errTo := strconv.Atoi(infos[1])
			if errTo != nil {
				to = 0
			}
			seconds := RandomInt(from, to)
			return seconds
		} else {
			seconds, err := strconv.Atoi(infos[0])
			if err != nil {
				seconds = 0
			}
			return seconds
		}
	} else {
		seconds, err := strconv.Atoi(infoTime)
		if err != nil {
			seconds = 0
		}
		return seconds
	}
}

func ParseVideoId(url string) string {
	if strings.Contains(url, "youtube.com") {
		url = strings.Trim(url, " ")
		url = strings.ReplaceAll(url, "https://www.youtube.com/", "")
		url = strings.ReplaceAll(url, "https://youtu.be/", "")
		url = strings.ReplaceAll(url, "https://m.youtube.com/", "")
		url = strings.ReplaceAll(url, "playlist?list=", "")
		url = strings.ReplaceAll(url, "watch?v=", "")
		if strings.Contains(url, "&") {
			url = strings.Split(url, "&")[0]
		}
		url = strings.ReplaceAll(url, "/", "")
	} else if strings.Contains(url, "drive.google.com") {
		url = strings.Trim(url, " ")
		url = strings.ReplaceAll(url, "https://drive.google.com/file/d/", "")
		url = strings.ReplaceAll(url, "https://drive.google.com/open?id=", "")
		url = strings.Split(url, "/")[0]
	} else if strings.Contains(url, "twitch.tv") {
		url = strings.Trim(url, " ")
		url = strings.ReplaceAll(url, "https://www.twitch.tv/", "")
		url = strings.ReplaceAll(url, "https://m.twitch.tv/", "")
		url = strings.ReplaceAll(url, "https://twitch.tv/", "")
		url = strings.ReplaceAll(url, "videos/", "")
	}
	return url
}

func FormatTime(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

func FormatSize(fileSize uint64) string {
	if fileSize == 0 {
		return "Unknown"
	} else {
		i := uint64(math.Floor(math.Log(float64(fileSize)) / math.Log(1024)))
		size := (float64(fileSize) / math.Pow(1024, float64(i))) * 1
		return fmt.Sprintf("%f %s", size, []string{"B", "kB", "MB", "GB", "TB"}[i])
	}
}

func ValidTwofaSecret(secret string) string {
	secret = strings.ReplaceAll(secret, " ", "")
	secret = strings.ToUpper(secret)
	return secret
}

func GenerateBirthday() (year, month, day int) {
	year = RandomInt(1970, 2000)
	month = RandomInt(1, 12)
	day = RandomInt(1, 28)
	return year, month, day
}

func GetLocalIP() string {
	url := "https://api.ipify.org?format=text" // we are using a pulib IP API, we're using ipify here, below are some others
	// https://www.ipify.org
	// http://myexternalip.com
	// http://api.ident.me
	// http://whatismyipaddress.com/api
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	//defer resp.Body.Close()
	defer func(Body io.ReadCloser) {
		if errClose := Body.Close(); errClose != nil {
			log.Println(errClose.Error())
		}
	}(resp.Body)
	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("My IP is:%s\n", ip)
	return string(ip)
}

func GenerateSessionID() string {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())
	randomUUID, err := ulid.New(ms, entropy)
	if err != nil {
		randomUUID = ulid.Make()
	}
	return randomUUID.String()
}

func GetFileSize(localPath string) (int64, error) {
	file, err := os.Open(localPath)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	if err != nil {
		return 0, err
	}
	stat, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return stat.Size(), nil
}

const dataCustomKey = "DataCustom"

func SetValueToContext(ctx *context.Context, key string, value interface{}) (*context.Context, error) {
	if ctx == nil || key == "" || value == nil {
		return nil, errors.New("Valid SetValueToContext")
	}

	dataCustom, ok := (*ctx).Value(dataCustomKey).(map[string]interface{})

	if !ok {
		dataCustom = make(map[string]interface{})
		*ctx = context.WithValue(*ctx, dataCustomKey, dataCustom)
	}
	dataCustom[key] = value
	return ctx, nil
}

func GetValueFormContext[T any](ctx *context.Context, key string) (data T, err error) {
	if ctx == nil || key == "" {
		return data, errors.New("Not valid input GetValueFormContext")
	}
	dataCustom, ok := (*ctx).Value(dataCustomKey).(map[string]interface{})

	if !ok {
		return data, errors.New("DataCustom not exist")
	}

	if value, ok := dataCustom[key]; ok {
		return value.(T), nil
	}
	return data, fmt.Errorf("Key %s not found", key)
}

func IsInRange(value, min, max float64) bool {
	return value >= min && value <= max
}

func UpperFirstLetter(s string) string {
	if len(s) == 0 {
		return ""
	}
	return string(s[0]-32) + s[1:]
}
