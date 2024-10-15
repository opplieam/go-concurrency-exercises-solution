//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"sync"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

var userTimeTracking = make(map[int]int)
var mu sync.Mutex

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	if u.IsPremium {
		process()
		return true
	}
	mu.Lock()
	_, ok := userTimeTracking[u.ID]
	if !ok {
		userTimeTracking[u.ID] = 0
	}
	mu.Unlock()

	go process()

	done := make(chan bool)
	timeTracking(u.ID, done)
	<-done

	return false
}

func timeTracking(uID int, done chan bool) {
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				mu.Lock()
				if userTimeTracking[uID] == 10 {
					mu.Unlock()
					done <- true
					return
				}
				userTimeTracking[uID]++
				mu.Unlock()
			}
		}
	}()
}

func main() {
	RunMockServer()
}
