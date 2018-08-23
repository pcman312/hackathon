// +build component

package component_tests

import (
	"context"
	"flag"
	"testing"

	"google.golang.org/grpc"

	"github.com/pcman312/hackathon/protos"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	ServerHost = flag.String("host", "", "Host of the server to connect to")
)

func TestServer(t *testing.T) {
	if *ServerHost == "" {
		t.Fatalf("Server must be specified with '--host'")
	}
	Convey("Testing save and retrieve", t, func() {
		Println("Dialing server...")
		cc, err := grpc.Dial(*ServerHost, grpc.WithInsecure())
		So(err, ShouldBeNil)

		client := hackathon.NewHackathonServiceClient(cc)

		Println("Saving message...")
		saveReq := &hackathon.SaveMessageRequest{
			Messages: []*hackathon.Message{
				&hackathon.Message{
					Text: "some message",
					User: "user",
				},
			},
		}

		_, err = client.SaveMessage(context.Background(), saveReq)
		So(err, ShouldBeNil)

		Println("Successfully saved")

		Println("Getting messages...")
		getReq := &hackathon.GetMessagesRequest{}
		getResp, err := client.GetMessages(context.Background(), getReq)
		So(err, ShouldBeNil)

		So(getResp.Messages, ShouldHaveLength, len(saveReq.Messages))
		for _, savedMsg := range getResp.Messages {
			So(savedMsg, ShouldBeIn, saveReq.Messages)
		}
	})
}
