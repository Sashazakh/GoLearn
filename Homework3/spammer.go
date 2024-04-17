package main

import (
	"fmt"
	"log"
	"sort"
	"sync"
)

func RunPipeline(cmds ...cmd) {
	var waitGroup sync.WaitGroup
	in := make(chan interface{})

	for _, command := range cmds {
		out := make(chan interface{})
		waitGroup.Add(1)

		go func(command cmd, in, out chan interface{}) {
			defer close(out)
			command(in, out)
			waitGroup.Done()
		}(command, in, out)

		in = out
	}

	waitGroup.Wait()
}

func SelectUsers(in, out chan interface{}) {
	// 	in - string
	// 	out - User

	var waitGroup sync.WaitGroup
	var mutex sync.Mutex
	processedUsers := make(map[string]struct{})

	for email := range in {
		waitGroup.Add(1)
		var emailStr string
		var ok bool

		if emailStr, ok = email.(string); !ok {
			continue
		}
		go func(email string) {
			defer waitGroup.Done()
			user := GetUser(email)
			mutex.Lock()
			if _, isProcessed := processedUsers[user.Email]; !isProcessed {
				processedUsers[user.Email] = struct{}{}
				out <- user
			}
			mutex.Unlock()
		}(emailStr)
	}

	waitGroup.Wait()
}

func SelectMessages(in, out chan interface{}) {
	// 	in - User
	// 	out - MsgID

	var waitGroup sync.WaitGroup
	usersBatch := make([]User, 0, GetMessagesMaxUsersBatch)

	for user := range in {
		var currUser User
		var ok bool
		if currUser, ok = user.(User); !ok {
			continue
		}

		usersBatch = append(usersBatch, currUser)

		if len(usersBatch) == GetMessagesMaxUsersBatch {
			waitGroup.Add(1)

			go func(batch []User) {
				defer waitGroup.Done()
				messagesID, err := GetMessages(batch...)
				if err != nil {
					log.Printf("Error getting message: %s", err)
					return
				}
				for _, id := range messagesID {
					out <- id
				}
			}(usersBatch)

			usersBatch = make([]User, 0, GetMessagesMaxUsersBatch)
		}
	}

	// Остаток батчей
	if len(usersBatch) > 0 {
		messageIDs, err := GetMessages(usersBatch...)
		if err != nil {
			log.Printf("Error getting message: %s", err)
			return
		}
		for _, id := range messageIDs {
			out <- id
		}
	}

	waitGroup.Wait()
}

func CheckSpam(in, out chan interface{}) {
	// in - MsgID
	// out - MsgData

	var waitGroup sync.WaitGroup
	semaphore := make(chan struct{}, HasSpamMaxAsyncRequests)

	for msg := range in {
		waitGroup.Add(1)
		semaphore <- struct{}{}

		var msgID MsgID
		var ok bool
		if msgID, ok = msg.(MsgID); !ok {
			continue
		}

		go func(msg MsgID) {
			defer waitGroup.Done()
			hasSpam, err := HasSpam(msgID)
			if err != nil {
				log.Panicf("An error occurred while checking spam: %v", err)
				return
			}
			out <- MsgData{ID: msgID, HasSpam: hasSpam}
			<-semaphore
		}(msgID)
	}

	waitGroup.Wait()
}

func CombineResults(in, out chan interface{}) {
	// in - MsgData
	// out - string

	var results []MsgData
	for data := range in {
		var msgData MsgData
		var ok bool
		if msgData, ok = data.(MsgData); !ok {
			continue
		}

		results = append(results, msgData)
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].HasSpam != results[j].HasSpam {
			return results[i].HasSpam
		}
		return results[i].ID < results[j].ID
	})

	for _, result := range results {
		formatted := fmt.Sprintf("%t %v", result.HasSpam, result.ID)
		out <- formatted
	}
}
