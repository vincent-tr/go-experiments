package main

import (
	"os"
	"time"
)

// https://ha101-1.overkiz.com/enduser-mobile-web
// const bvouest = "io://0220-6975-2311/14430852"
//  deviceRefreshInterval = 30 * MINUTE, eventPollInterval = 2 * SECOND, stateRefreshInterval = 1 * MINUTE

func main() {
	client, err := MakeClient(os.Getenv("KIZ_USERNAME"), os.Getenv("KIZ_PASSWORD"))
	if err != nil {
		panic(err)
	}

	for {
		time.Sleep(time.Hour)
	}

	client.Terminate()
}
