package exchange

import "time"

func IsoToUnix(date string) (int64, error) {
	ts, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return 0, err
	}
	return ts.UnixNano() / int64(time.Millisecond), nil
}

func UnixToISO(date int64) string {
	dt := time.Unix(date/1000, 0)
	return dt.UTC().Format(time.RFC3339)
}
