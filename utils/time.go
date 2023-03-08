package utils

import (
	"database/sql/driver"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"strconv"
	"strings"
	"time"
)

// JSONTime format json time field by myself
type JSONTime struct {
	time.Time
}

func (t JSONTime) FormatYyyyMmDdHhMmSs() string {
	return t.Format("2006-01-02 15:04:05")
}

// MarshalJSON on JSONTime format Time field with %Y-%m-%d %H:%M:%S
func (t JSONTime) MarshalJSON() ([]byte, error) {
	//if t.UnixMilli() > 0 {
	//	return []byte(strconv.FormatInt(t.UnixMilli(), 10)), nil
	//}
	//return []byte(strconv.FormatInt(0, 10)), nil
	b, _ := jsoniter.Marshal(t.FormatYyyyMmDdHhMmSs())
	return b, nil
}

func (t *JSONTime) UnmarshalJSON(s []byte) error {
	str := string(s)
	if str == "null" {
		return nil
	}
	switch {
	case strings.Contains(str, "-"):
		// Fractional seconds are handled implicitly by Parse.
		var err error
		t.Time, err = time.Parse(`"`+time.RFC3339+`"`, str)
		if err != nil {
			return err
		}
	default:
		q, err := strconv.ParseInt(string(s), 10, 64)
		if err != nil {
			return err
		}
		t.Time = time.UnixMilli(q)
	}
	return nil
}

// Value insert timestamp into mysql need this function.
func (t JSONTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan value time.Time
func (t *JSONTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JSONTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
