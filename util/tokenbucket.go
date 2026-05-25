package util

import (
	// "log"
	"ots/settings"
	"time"
)

var Limmiter chan struct{}

func StartTokenBucket(){

	// log.Println("Request Limit : ",settings.MySettings.Get_RateLimit())
	setLimit(settings.MySettings.Get_RateLimit())
	//refill tokens
	ticker := time.NewTicker(1*time.Second)

	go func() {
		for range ticker.C{
			for i:=0 ;i<cap(Limmiter);i++{
				select {
				case Limmiter <- struct{}{}:
				default :
					//bucket full
				}
			}
		}
	}()
}


func setLimit(limit int){
	Limmiter = make(chan struct{},limit)
}