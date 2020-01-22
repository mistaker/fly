package common

import "log"

func GoSafe(fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()
		fn()
	}()
}
