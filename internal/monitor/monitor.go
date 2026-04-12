package monitor

import (
	"strconv"
	"time"

	"fyne.io/fyne/v2"

	"NekoSleep/internal/config"
	"NekoSleep/internal/locker"
)

var skipUntil time.Time
var sleepImgCache fyne.Resource 
var stopChan chan bool          

func Init(img fyne.Resource) {
	sleepImgCache = img
}

func ProlongSession() {
	skipUntil = time.Now().Add(time.Hour)
}

func isTimeInSleepWindow(nowHour, nowMin, sleepHour, sleepMin int) bool {
	nowTotal := nowHour*60 + nowMin
	startTotal := sleepHour*60 + sleepMin
	endTotal := (startTotal + 480) % 1440

	if startTotal < endTotal {
		return nowTotal >= startTotal && nowTotal < endTotal
	}
	return nowTotal >= startTotal || nowTotal < endTotal
}


func Start() {
	if stopChan != nil {
		return 
	}

	stopChan = make(chan bool)

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		isLockerOpen := false

		for {
			select {
			case <-stopChan:
				
				return 

			case <-ticker.C:
				
				if time.Now().Before(skipUntil) {
					isLockerOpen = false
					continue
				}

				data, err := config.Load()
				if err != nil {
					continue
				}

				sH, _ := strconv.Atoi(data.Hour)
				sM, _ := strconv.Atoi(data.Minute)
				now := time.Now()

				if isTimeInSleepWindow(now.Hour(), now.Minute(), sH, sM) {
					if !isLockerOpen {
						isLockerOpen = true
						fyne.Do(func() {locker.Show(sleepImgCache, ProlongSession)}) 
					}
				} else {
					isLockerOpen = false
				}
			}
		}
	}()
}

func Stop() {
	if stopChan != nil {
		close(stopChan) 
		stopChan = nil
	}
}