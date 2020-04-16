package util

import (
	"time"
)

func Retry(f func() error, retryTimes int, interval time.Duration) error {
	var err error
	for i := 0; i < retryTimes; i++ {
		if err = f(); err != nil {
			if i == retryTimes {
				return err
			}
			time.Sleep(interval)
			continue
		}
		return nil
	}
	return err
}
