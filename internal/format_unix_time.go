package internal

import "time"

func formatUnixTime(timeStamp int64) string {
	t := time.Unix(timeStamp, 0)
	return t.Format("2006-01-02 15:04:05 EST")
}
