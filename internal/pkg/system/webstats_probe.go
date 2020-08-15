package system

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type WebStatsProbe struct {
	statsUrl         string
	samplingInterval time.Duration
	observers        []TextObserver
	initialDelay     time.Duration
}

func NewWebStatsProbe(
	statsUrl string,
	samplingInterval time.Duration,
	initialDelay time.Duration,
) WebStatsProbe {
	return WebStatsProbe{
		statsUrl:         statsUrl,
		samplingInterval: samplingInterval,
		initialDelay:     initialDelay,
	}
}

func (w *WebStatsProbe) Observe(observer TextObserver) {
	w.observers = append(w.observers, observer)

	go func() {
		// wait for initial delay to lapse before polling to give time
		// to the server to start
		time.Sleep(w.initialDelay)

		for {
			resp, err := http.Get(w.statsUrl)
			if err != nil {
				for _, o := range w.observers {
					o.Update(fmt.Sprintf("%s\n", err))
				}
			} else {
				// Response can only be read once
				responseBody, err := read(resp)
				if err != nil {
					break
				}

				for _, o := range w.observers {
					o.Update(responseBody)
				}
			}

			time.Sleep(w.samplingInterval)
		}
	}()
}

func read(resp *http.Response) (string, error) {
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(body), nil
}
