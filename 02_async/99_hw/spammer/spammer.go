package main

import (
	"fmt"
	"slices"
	"sort"
	"sync"
)

/*
SelectUsers() -> SelectMessages() -> CheckSpam() -> CombineResults()
Max exec time: 3s
Available API:
1. GetUser(),     1.0 sec,
2. GetMessages(), 1.0 sec, batch 2
3. HasSpam(),     0.1 sec, limit 5 connections

in, out -> out, new(out)
*/

func RunPipeline(cmds ...cmd) {
	if len(cmds) == 0 {
		return
	}

	var wg sync.WaitGroup
	var channels = make([]chan any, len(cmds)+1)

	for i := range channels {
		channels[i] = make(chan any, 100)
	}
	close(channels[0])

	for i, command := range cmds {
		in, out := channels[i], channels[i+1]

		wg.Add(1)
		go func(command cmd, in, out chan any) {
			defer wg.Done()
			defer close(out)
			command(in, out)
		}(command, in, out)
	}

	wg.Wait()
}

// in - string, out - User
func SelectUsers(in, out chan any) {
	ch := make(chan User, 100)

	var wg sync.WaitGroup
	userIds := make(map[uint64]struct{})

	for userEmail := range in {
		wg.Add(1)
		go func(userEmail any) {
			defer wg.Done()
			strEmail, ok := userEmail.(string)
			if !ok {
				return
			}
			ch <- GetUser(strEmail)
		}(userEmail)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for user := range ch {
		if _, exist := userIds[user.ID]; !exist {
			userIds[user.ID] = struct{}{}
			out <- user
		}
	}
}

// in - User, out - MsgID
func SelectMessages(in, out chan any) {
	buffer := make([]User, GetMessagesMaxUsersBatch)
	i := 0
	var wg sync.WaitGroup

	for userAny := range in {
		user, ok := userAny.(User)
		if !ok {
			continue
		}
		buffer[i] = user
		i++

		if i == GetMessagesMaxUsersBatch {
			wg.Add(1)
			bufferCopy := slices.Clone(buffer) // копия!
			go func(buffer []User) {
				defer wg.Done()
				messages, err := GetMessages(buffer...)
				if err != nil {
					fmt.Println("ERROR: ", err)
					return
				}
				for _, msg := range messages {
					out <- msg
				}
			}(bufferCopy)
			i = 0
		}
	}

	if i != 0 {
		wg.Add(1)
		go func(buffer []User) {
			defer wg.Done()
			messages, err := GetMessages(buffer...)
			if err != nil {
				fmt.Println("ERROR: ", err)
				return
			}
			for _, msg := range messages {
				out <- msg
			}
		}(buffer[:i])
		i = 0
	}
	wg.Wait()
}

// in - MsgID out - MsgData
func CheckSpam(in, out chan any) {
	sem := make(chan struct{}, HasSpamMaxAsyncRequests)
	var wg sync.WaitGroup

	for v := range in {
		msg, ok := v.(MsgID)
		if !ok {
			continue
		}
		sem <- struct{}{}
		wg.Add(1)
		go func(m MsgID) {
			defer wg.Done()
			defer func() { <-sem }()
			isSpam, err := HasSpam(m)
			if err != nil {
				return
			}
			out <- MsgData{ID: m, HasSpam: isSpam}
		}(msg)
	}
	wg.Wait()
}

// in - MsgData out - string
func CombineResults(in, out chan any) {
	results := []MsgData{}
	for msgDataAny := range in {
		msgData, ok := msgDataAny.(MsgData)
		if !ok {
			continue
		}
		results = append(results, msgData)
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].HasSpam != results[j].HasSpam {
			return results[i].HasSpam
		} else {
			return results[i].ID < results[j].ID
		}
	})

	for _, msgData := range results {
		out <- fmt.Sprintf("%t %d", msgData.HasSpam, msgData.ID)
	}
}

func PrintResults(in, out chan any) {
	// in - any
	// out - any
	for item := range in {
		strItem, ok := item.(string)
		if ok {
			fmt.Println(strItem)
		}
		out <- item
	}
}
