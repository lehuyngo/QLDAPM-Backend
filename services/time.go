package services

import "time"

func FormatDate(t time.Time) string {
	return t.Format(Config.Format.DateFormat)
}

func FormatTime(t time.Time) string {
	return t.Format(Config.Format.TimeFormat)
}

func FormatTimestamp(miliSecond int64) string {
	return time.Unix(miliSecond/1000, 0).Format(Config.Format.DateFormat)
}

func TimestampToDate(miliSecond int64) time.Time {
	return time.Unix(miliSecond/1000, 0)
}

func StringToTimestamp(str string) (int64, error) {
	result, err := time.Parse(Config.Format.TimeFormat, str)
	if err == nil {
		return result.UnixMilli(), nil
	}

	result, err = time.Parse(Config.Format.DateFormat, str)
	if err != nil {
		return 0, err
	}

	return result.Add(12 * time.Hour).UnixMilli(), nil
}

func StringToTime(str string) (time.Time, error) {
	return time.Parse(Config.Format.DateFormat, str)
}
