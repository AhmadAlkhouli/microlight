package main

import logger "microlight/logger"

var mylogger = logger.Create()

func main() {

	i := make(chan string)

	go func() {
		i <- "hello"
		close(i)
	}()
	for {
		select {
		case v, ok := <-i:
			if !ok {
				return
			}
			mylogger.Info.Println(v)

		}
	}
}
