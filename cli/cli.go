package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"google.golang.org/grpc"

	"github.com/pcman312/hackathon/protos"
)

type Opts struct {
	Host string `short:"h" required:"true" description:"Host to connect to"`
	User string `short:"u" required:"true" description:"User that is sending the message"`
	Text string `short:"m" required:"true" description:"Text of the message"`
}

func main() {
	opts := Opts{}
	parser := flags.NewParser(&opts, flags.HelpFlag)
	_, err := parser.Parse()
	if err != nil {
		fmt.Printf("%s\n", err)

		e, ok := err.(*flags.Error)
		if ok && e.Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	fmt.Printf("Dialing server...\n")
	cc, err := grpc.Dial(opts.Host, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := hackathon.NewHackathonServiceClient(cc)

	fmt.Printf("Saving message...\n")
	req := &hackathon.SaveMessageRequest{
		Messages: []*hackathon.Message{
			&hackathon.Message{
				Text: opts.Text,
				User: opts.User,
			},
		},
	}

	_, err = client.SaveMessage(context.Background(), req)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Successfully saved\n")

	fmt.Printf("Getting messages...\n")
	getReq := &hackathon.GetMessagesRequest{}
	getResp, err := client.GetMessages(context.Background(), getReq)
	if err != nil {
		panic(err)
	}

	for _, msg := range getResp.Messages {
		fmt.Printf("Mesage: [%s] [%s]\n", msg.User, msg.Text)
	}

	fmt.Printf("Goodbye\n")
}
