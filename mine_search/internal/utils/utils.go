package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"mine/internal"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"golang.org/x/text/unicode/norm"
)

func IsEmpty(value string) bool {
	return len(strings.TrimSpace(value)) == 0
}

func StructToMap(myStruct interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	jsonEnc, err := json.Marshal(myStruct)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(jsonEnc, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func GetDefaultStr(trying []string, defaultValue string) string {
	for idx := range trying {
		if !IsEmpty(trying[idx]) {
			return trying[idx]
		}
	}
	return defaultValue
}

func GetDefaultEnv(variable_name, defaultValue string) (string, bool) {
	value, ok := os.LookupEnv(variable_name)
	if ok {
		return value, ok
	}
	fmt.Println("Không có biến môi trường:", variable_name)
	return defaultValue, ok
}

func GetDefaultInt(trying []int, defaultValue int) int {
	for idx := range trying {
		if trying[idx] > 0 {
			return trying[idx]
		}
	}
	return defaultValue
}

func IsValidStar(answer []int) bool {
	if len(answer) != 1 {
		return false
	}
	return answer[0] <= 5 && answer[0] >= 1
}

func IsSubset(first, second []int) bool {
	set := make(map[int]int)
	for _, value := range second {
		set[value] += 1
	}

	for _, value := range first {
		if count, found := set[value]; !found {
			return false
		} else if count < 1 {
			return false
		} else {
			set[value] = count - 1
		}
	}

	return true
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func GetHmacSha256(data, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

func GetSha256(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

func CreateLocalToken(EcomClientKey, EcomSecretKey string) string {
	timestr := GetStringTimeUTC7("Y-D-M")
	return GetMD5Hash(EcomClientKey + "::" + EcomSecretKey + timestr)
}

func CheckRegexFrType(input string, regexType string) bool {
	var sampleRegexp *regexp.Regexp
	sampleRegexp = regexp.MustCompile(regexType)
	return sampleRegexp.MatchString(input)
}

func CheckValid(input string, mapCheckRegex map[string]string, typeValue string) (bool, string) {
	value, ok := mapCheckRegex[typeValue]
	if !ok {
		return false, "Check regex not have typeValue:" + typeValue
	}
	if !CheckRegexFrType(input, value) {
		return false, "Unvalid: " + typeValue + " - " + input
	}
	return true, "ok"
}

func GetTimeUTC7() time.Time {
	now := time.Now()
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	return now.In(loc)
}

func GetTimeUTC7WithAddedDays(daysToAdd int) time.Time {
	return GetTimeUTC7().Add(time.Hour * time.Duration(24*daysToAdd))
}

func GetEndOfDayUTC7() time.Time {
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	timeWithOffset := GetTimeUTC7()
	year, month, day := timeWithOffset.Date()
	endOfDay := time.Date(year, month, day, 0, 0, 0, 0, loc).Add(time.Hour*24 - time.Second)
	return endOfDay
}

func GetTimeUTC7FrTime(input time.Time) time.Time {
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	return input.In(loc)
}

func GetTimeUTC7FrTimeV2(input time.Time) time.Time {
	if input.Location().String() == "UTC" {
		input = input.Add(time.Hour * -7)
	}
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	return input.In(loc)
}

func GetStringTimeUTC7(stringtype string) string {
	t := GetTimeUTC7()
	switch stringtype {
	case "Y-D-M":
		return t.Format("2006-02-01")
	case "Y-M-D":
		return t.Format("2006-01-02")
	case "D-M-Y":
		return t.Format("02-01-2006")
	case "M-D-Y":
		return t.Format("01-02-2006")
	case "Y-D-M H:M:S":
		return t.Format("2006-02-01 15:04:05")
	case "Y-M-D H:M:S":
		return t.Format("2006-01-02 15:04:05")
	case "Y-M-D H:M:S -0700":
		return t.Format("2006-01-02 15:04:05 -0700")
	case "D/M/Y H:M:S":
		return t.Format("02/01/2006 15:04:05")
	}
	return ""
}

func GetStringTime(t time.Time, stringtype string) string {
	switch stringtype {
	case "Y-D-M":
		return t.Format("2006-02-01")
	case "Y-M-D":
		return t.Format("2006-01-02")
	case "D-M-Y":
		return t.Format("02-01-2006")
	case "D/M/Y":
		return t.Format("02/01/2006")
	case "D/M/Y H:M:S":
		return t.Format("02/01/2006 15:04:05")
	case "H:M:S D/M/Y":
		return t.Format("15:04:05 02/01/2006")
	case "M-D-Y":
		return t.Format("01-02-2006")
	case "Y-D-M H:M:S":
		return t.Format("2006-02-01 15:04:05")
	case "Y-M-D H:M:S":
		return t.Format("2006-01-02 15:04:05")
	case "Y-M-D H:M:S -0700":
		return t.Format("2006-01-02 15:04:05 -0700")
	}
	return ""
}

func GetBlackListContain(input string) string {
	blacklist := [...]string{
		"drop", "delete", "select", "update", "or", "and", "insert", "all",
		"=", "<>", "!=", ">", ">=", "<", "<=", "*", ";", "--",
	}
	for _, value := range blacklist {
		if strings.Contains(input, value) {
			return value
		}
	}
	return ""
}

func EncodeBase64(text string) string {
	return base64.StdEncoding.EncodeToString([]byte(text))
}

func DecodeBase64(textEncoded string) (string, error) {
	enc, err := base64.StdEncoding.DecodeString(textEncoded)
	return string(enc), err
}

func ParseTimeFrString(stringtype string, timeinput string) (time.Time, error) {
	layout := ""
	switch stringtype {
	case "Y-D-M":
		layout = "2006-02-01"
	case "Y-M-D":
		layout = "2006-01-02"
	case "D-M-Y":
		layout = "02-01-2006"
	case "D-M-Y H:M:S":
		layout = "02-01-2006 15:04:05"
	case "M-D-Y":
		layout = "01-02-2006"
	case "Y-D-M H:M:S":
		layout = "2006-02-01 15:04:05"
	case "Y-M-D H:M:S":
		layout = "2006-01-02 15:04:05"
	case "Y-M-D H:M:S -0700":
		layout = "2006-01-02 15:04:05 -0700"
	case "Y-M-DTH:M:S +0700":
		layout = "2006-01-02T15:04:05 +0700"
	case "Y-M-DTH:M:S+07:00":
		layout = "2006-01-02T15:04:05+07:00"
	case "Y-M-DTH:M:S.000":
		layout = "2006-01-02T15:04:05.999999999"
	case "D/M/Y":
		layout = "02/01/2006"
	case "D/M/Y H:M:S":
		layout = "02/01/2006 15:04:05"
	case "M/D/Y H:M:S":
		layout = "01/02/2006 15:04:05"
	case "H:M:S D/M/Y":
		layout = "15:04:05 02/01/2006"
	}
	t, err := time.Parse(layout, timeinput)
	return t, err
}

func ParseTimeFrStringV2(stringtype string, timeinput string) (time.Time, error) {
	t, err := ParseTimeFrString(stringtype, timeinput)
	if err != nil {
		return t, err
	}
	return GetTimeUTC7FrTimeV2(t), nil
}

func ConvertInterfaceToStruct(a interface{}) {
	type Discount []struct {
		Percentage      int `json:"percentage"`
		NumberUserStart int `json:"number_user_start"`
		NumberUserEnd   int `json:"number_user_end"`
	}
	jsonbody, err := json.Marshal(a)
	if err != nil {
		fmt.Println(err)
	}
	student := Discount{}
	if err := json.Unmarshal(jsonbody, &student); err != nil {
		// do error check
		fmt.Println(err)
	}
	fmt.Printf("student: %v\n", student)
}

func FloatToTime(input float64) time.Time {
	integ, decim := math.Modf(input)
	return time.Unix(int64(integ), int64(decim*(1e9)))
}

func StringToFloat64(input string) (float64, error) {
	result, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func StringToInt(input string) (int, error) {
	result, err := strconv.Atoi(input)
	if err != nil {
		return 0, err
	}
	return result, nil
}

var Weekdays = map[string]string{
	"Sunday":    "Chủ nhật",
	"Monday":    "Thứ hai",
	"Tuesday":   "Thứ ba",
	"Wednesday": "Thứ tư",
	"Thursday":  "Thứ năm",
	"Friday":    "Thứ sáu",
	"Saturday":  "Thứ bảy",
}

var WeekdayV2 = map[int]string{
	0: "Chủ nhật",
	1: "Thứ 2",
	2: "Thứ 3",
	3: "Thứ 4",
	4: "Thứ 5",
	5: "Thứ 6",
	6: "Thứ 7",
}

func DayInWeek(day time.Weekday) string {
	// fmt.Printf("day.String(): %v\n", day.String())
	return Weekdays[day.String()]
}

func DayInWeekV2(input time.Time) string {
	return WeekdayV2[int(input.Weekday())]
}

func TimeToTimeStr(start, end time.Time) string {
	format1 := "02 tháng 01, 2006 03:04 PM"
	format2 := "02 tháng 01, 2006 (03:04 PM - "
	format3 := "03:04 PM)"
	if start.Equal(end) {
		return DayInWeek(start.Weekday()) + ", " + start.Format(format1)
	} else {
		if start.Day() == end.Day() && start.Month() == end.Month() && start.Year() == end.Year() {
			return DayInWeek(start.Weekday()) + ", " + start.Format(format2) + end.Format(format3)
		} else {
			return DayInWeek(start.Weekday()) + ", " + start.Format(format1) + " - " + DayInWeek(end.Weekday()) + ", " + end.Format(format1)
		}
	}
}

func ConvertListToStringVN(input []string) string {
	lenStr := len(input)
	if lenStr == 1 {
		return input[0]
	}
	result := ""
	for index, value := range input {
		if index == 0 {
			result = value
		} else if index == lenStr-1 {
			result = result + " và " + value
		} else {
			result = result + ", " + value
		}
	}
	return result
}

func ConvertListToString(input []string) string {
	lenStr := len(input)
	if lenStr == 1 {
		return input[0]
	}
	result := ""
	for index, value := range input {
		if index == 0 {
			result = value
		} else {
			result = result + ";" + value
		}
	}
	return result
}

func ConvertStringToList(input string, key string) []string {
	if input == "" {
		return []string{}
	}
	return strings.Split(input, key)
}

func FindListInList(input []string, input2 []string) []string {
	result := []string{}
	mapInput2 := map[string]string{}
	for _, v := range input2 {
		mapInput2[v] = v
	}
	for _, v := range input {
		if value, ok := mapInput2[v]; ok {
			result = append(result, value)
		}
	}
	return result
}

func FormatIntToVND(n int64) string {
	in := []byte(strconv.FormatInt(n, 10))
	var out []byte
	if i := len(in) % 3; i != 0 {
		if out, in = append(out, in[:i]...), in[i:]; len(in) > 0 {
			out = append(out, '.')
		}
	}
	for len(in) > 0 {
		if out, in = append(out, in[:3]...), in[3:]; len(in) > 0 {
			out = append(out, '.')
		}
	}
	return string(out)
}

func NoAccentVietnamese(myStr string) string {
	normalized := norm.NFD.String(myStr)
	regex := regexp.MustCompile("[^\\p{L}\\p{N}\\s]+")
	clean := regex.ReplaceAllString(normalized, "")
	return clean
}

func RemoveAllWhiteSpace(myStr string) string {
	return strings.ReplaceAll(myStr, " ", "")
}

func ConvertStrToUpperWithoutSpacing(myStr string) string {
	myStrFilter := NoAccentVietnamese(myStr)
	myStrFilterUpper := strings.ToUpper(myStrFilter)
	return RemoveAllWhiteSpace(myStrFilterUpper)
}

func AddingDateFromTime(seconds int) time.Time {
	timeNow := GetTimeUTC7()
	duration := time.Duration(seconds) * time.Second
	timeAfterSec := timeNow.Add(duration)
	return timeAfterSec
}

func InterfaceToStruct(input interface{}, output interface{}) error {
	jsonBytes, err := json.Marshal(input)
	if err != nil {
		internal.Log.Error("InterfaceToStruct", zap.Any("Error Marshal", err))
		return err
	}
	err = json.Unmarshal(jsonBytes, &output)
	if err != nil {
		internal.Log.Error("InterfaceToStruct", zap.Any("Error Unmarshal", err))
		return err
	}
	return nil
}

func CreateParseURL(baseURL string, queryParams map[string]string) string {
	output := ""
	// Parse the base URL
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		fmt.Println("Error parsing base URL:", err)
		return output
	}
	// Add query parameters to the URL
	q := parsedURL.Query()
	for key, value := range queryParams {
		q.Add(key, value)
	}
	parsedURL.RawQuery = q.Encode()
	output = parsedURL.String()
	return output
}

func FormatMoneyVND(amount int, dot, cover string) string {
	// Convert the integer to a string
	amountStr := strconv.Itoa(amount)
	// Determine the length of the integer part
	length := len(amountStr)
	var formattedParts []string
	for i := 0; i < length; i += 3 {
		end := length - i
		if end > 3 {
			end = 3
		}
		formattedParts = append([]string{amountStr[length-i-end : length-i]}, formattedParts...)
	}
	// Join the parts with commas
	formattedAmount := ""
	if len(formattedParts) > 0 {
		formattedAmount = formattedParts[0]
		for i := 1; i < len(formattedParts); i++ {
			formattedAmount += dot + formattedParts[i]
		}
	}
	return formattedAmount + cover
}

func MaskString(s string, start, end int, key int32) (string, error) {
	if s == "" {
		return s, nil
	}
	rs := []rune(s)
	if len(s) < start || len(s) < end {
		return "", errors.New("index not valid")
	}
	for i := start; i <= end; i++ {
		rs[i] = key
	}
	return string(rs), nil
}

func MaskStringV2(s string, start, end int, key int32) string {
	if s == "" {
		return s
	}
	rs := []rune(s)
	if len(s) < start || len(s) < end {
		return ""
	}
	for i := start; i <= end; i++ {
		rs[i] = key
	}
	return string(rs)
}

func StringToFormatTime(stringtype string) string {
	switch stringtype {
	case "Y-D-M":
		return "2006-02-01"
	case "Y-M-D":
		return "2006-01-02"
	case "D-M-Y":
		return "02-01-2006"
	case "M-D-Y":
		return "01-02-2006"
	case "Y-D-M H:M:S":
		return "2006-02-01 15:04:05"
	case "Y-M-D H:M":
		return "2006-01-02 15:04"
	case "Y-M-D H:M:S":
		return "2006-01-02 15:04:05"
	case "Y-M-D H:M:S -0700":
		return "2006-01-02 15:04:05 -0700"
	case "D/M/Y":
		return "02/01/2006"
	case "D/M/Y H:M":
		return "02/01/2006 15:04"
	case "D-M-Y H:M":
		return "02-01-2006 15:04"
	case "H:M D/M/Y":
		return "15:04 02/01/2006"
	case "Y-M-DTH:M:S.000":
		return "2006-01-02T15:04:05.999999999"
	case "D/M/Y H:M:S.000":
		return "02-01-2006 15:04:05.999999999"
	case "D/M/Y H:M:S":
		return "02/01/2006 15:04:05"
	case "D/M/Y - H:M":
		return "02/01/2006 - 15:04"
	}
	return ""
}

func ConvertFormatTimeAtoB(stringtime, formatA, formatB string) (string, error) {
	time_input, err := time.Parse(StringToFormatTime(formatA), stringtime)
	if err != nil {
		return "", err
	}
	return time_input.Format(StringToFormatTime(formatB)), nil
}

func GetRelativeTime(tBegin, tEnd time.Time) string {
	// return 1 year ago, 1 month ago, 1 week ago, 1 day ago, 2 hour ago, 1 minute ago
	timeDiff := tEnd.Sub(tBegin).Seconds()
	// Get base on Milliseconds unit
	year, yearUnit := timeDiff/(60*60*24*30*12), " năm"
	month, monthUnit := timeDiff/(60*60*24*30), " tháng"
	week, weekUnit := timeDiff/(60*60*24*7), " tuần"
	day, dayUnit := timeDiff/(60*60*24), " ngày"
	hour, hourUnit := timeDiff/(60*60), " giờ"
	minute, minuteUnit := timeDiff/60, " phút"

	if int(year) > 0 {
		return fmt.Sprintf("%v%s", int(year), yearUnit)
	}
	if int(month) > 0 {
		return fmt.Sprintf("%v%s", int(month), monthUnit)
	}
	if int(week) > 0 {
		return fmt.Sprintf("%v%s", int(week), weekUnit)
	}
	if int(day) > 0 {
		return fmt.Sprintf("%v%s", int(day), dayUnit)
	}
	if int(hour) > 0 {
		return fmt.Sprintf("%v%s", int(hour), hourUnit)
	}
	if int(minute) > 0 {
		return fmt.Sprintf("%v%s", int(minute), minuteUnit)
	}
	return "vài giây"
}

func ParseJwt(token string, key string) (jwt.MapClaims, *internal.SystemStatus) {
	clams := jwt.MapClaims{}
	decodeToken, err := jwt.ParseWithClaims(token, clams, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		jwtErr := err.(*jwt.ValidationError).Errors
		if jwtErr == jwt.ValidationErrorExpired {
			internal.Log.Error("ValidationErrorExpired", zap.Any("input", token), zap.Any("tokenParse", decodeToken), zap.Error(err))
			return nil, internal.SysStatus.TokenExpired
		}
	}
	if decodeToken == nil || !decodeToken.Valid {
		internal.Log.Error("Token Valid", zap.Any("input", token), zap.Any("tokenParse", decodeToken))
		return nil, internal.SysStatus.InvalidToken
	}
	return clams, nil
}

// GetCurrentFuncName returns the name of the current function
func GetCurrentFuncName() string {
	pc, _, _, ok := runtime.Caller(1) // 1 refers to the current function
	if !ok {
		return "unknown"
	}

	// Get the function name using the program counter (PC)
	funcName := runtime.FuncForPC(pc).Name()

	// Optional: Extract the short function name (without the package path)
	return extractFuncName(funcName)
}

// extractFuncName shortens the full function name to its base name
func extractFuncName(fullName string) string {
	// Split the function name by "/" and get the last part, which is the function name
	parts := strings.Split(fullName, "/")
	return parts[len(parts)-1]
}
