package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

//Client represents a client
type Client struct {
	conn     net.Conn
	nickname string
	//id string
	//commenting chanel for client and writing on connection
	//ch       chan string
}

func main() {
	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//msgchan is used for all the messages any user types
	//the message is first populated by handleConnection and passed on to handleMessages
	msgchan := make(chan string)

	//A channel to keep track of Client connections, clients are added to this channel
	//handleMessages then iterates through clients and for each appends to its channel
	//the messages sent by other users, broadcast
	addchan := make(chan Client)

	go handleMessages(msgchan, addchan)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		//for each client, we do have a separate handleConnection goroutine
		go handleConnection(conn, msgchan, addchan)
	}
}

//ReadLinesInto is a method on Client type
//it keeps waiting for user to input a line, ch chan is the msgchannel
//it formats and writes the message to the channel
func (c Client) ReadLinesInto(ch chan<- string) {
	bufc := bufio.NewReader(c.conn)
	for {
		line, err := bufc.ReadString('\n')
		if err != nil {
			break
		}
		ch <- fmt.Sprintf("%s: %s", c.nickname, line)
	}
}

/*WriteLinesFrom is a method is not required as we are writing to connection in handlConnection
//each client routine is writing to channel
func (c Client) WriteLinesFrom(ch <-chan string) {
	for msg := range ch {
		_, err := io.WriteString(c.conn, msg)
		if err != nil {
			return
		}
	}
}*/

func promptNick(c net.Conn, bufc *bufio.Reader) string {
	io.WriteString(c, "\033[1;30;41mWelcome to the fancy demo chat!\033[0m\n")
	io.WriteString(c, "What is your nick? ")
	nick, _, _ := bufc.ReadLine()
	return string(nick)
}

//the core one
func handleConnection(c net.Conn, msgchan chan<- string, addchan chan<- Client) {
	bufc := bufio.NewReader(c)
	defer c.Close()

	//we first need to add current client to the channel
	//filling in the client structure
	client := Client{
		conn:     c,
		nickname: promptNick(c, bufc),
		//ch:       make(chan string),
	}
	if strings.TrimSpace(client.nickname) == "" {
		io.WriteString(c, "Invalid Username\n")
		return
	}

	// Register user, our messageHandler is waiting on this channel
	//it populates the map
	addchan <- client

	//just a welcome message
	io.WriteString(c, fmt.Sprintf("Welcome, %s!\n\n", client.nickname))

	//We are now populating the other channel now
	//our message handler is waiting on this channel as well
	//it reads this message and copies to the individual channel of each Client in map
	// effectively the broadcast
	msgchan <- fmt.Sprintf("New user %s has joined the chat room.\n", client.nickname)

	// another go routine whose purpose is to keep on waiting for user input
	//and write it with nick to the
	go client.ReadLinesInto(msgchan)

	//given a channel, writelines prints lines from it
	//we are giving here client.ch and this routine is for each client
	//so effectively each client is printitng its channel
	//to which our messagehandler has added messages for boroadcast
//here we need this infinite loop to keep the thread alive otherwise once ReadLinesInto ends the threads also ends with it. Writes lines into was keeping the thread alive while waiting for the something to be written.
for{

}
}

func handleMessages(msgchan <-chan string, addchan <-chan Client) {
//here instead of having a channel for the clients msg now a string key value pair would do.
	clients := make(map[net.Conn]string)

	for {
		select {
		case msg := <-msgchan:
			log.Printf("New message: %s", msg)
			for cli,_ := range clients {
				//go func(mch chan<- string) { mch <- "\033[1;33;40m" + msg + "\033[m" }(ch)
				//ch <- "\033[1;33;40m" + msg + "\033[m"
//here writing to connection the msgs written by a client instead of writing to the channel
		_, err := io.WriteString(cli, msg)
		if err != nil {
			return
	}

			}
		case client := <-addchan:
			log.Printf("New client: %v\n", client.conn)
			clients[client.conn] = client.nickname
//writing the name of the user on the connection to display

		}
	}
}

