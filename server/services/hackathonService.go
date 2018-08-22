package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	log "github.com/cihub/seelog"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/olivere/elastic.v3"

	"github.com/pcman312/hackathon/protos"
)

type HackathonService struct {
	esClient *elastic.Client
	index    string
}

func NewHackathonService(esClient *elastic.Client, index string) (svc HackathonService, err error) {
	if esClient == nil {
		return svc, errors.New("missing elasticsearch client")
	}
	if index == "" {
		return svc, errors.New("missing index")
	}

	svc = HackathonService{
		esClient: esClient,
		index:    index,
	}
	return svc, nil
}

func (s HackathonService) SaveMessage(ctx context.Context, req *hackathon.SaveMessageRequest) (resp *hackathon.SaveMessageResponse, err error) {
	resp = &hackathon.SaveMessageResponse{}

	for _, msg := range req.Messages {
		msg.Timestamp = time.Now().UnixNano()

		_, err = s.esClient.Index().
			Index(s.index).
			Type("message").
			BodyJson(msg).
			Refresh(true).
			DoC(ctx)
		if err != nil {
			return resp, status.Error(codes.Internal, err.Error())
		}
	}

	return resp, nil
}

func (s HackathonService) GetMessages(ctx context.Context, req *hackathon.GetMessagesRequest) (resp *hackathon.GetMessagesResponse, err error) {
	resp = &hackathon.GetMessagesResponse{}

	scroll := s.esClient.
		Scroll(s.index).
		Type("message").
		Sort("timestamp", false).
		FetchSource(true)

	messages := []*hackathon.Message{}

	scrollId := ""
	for {
		result, err := scroll.Do()
		if err == io.EOF {
			cancelScroll(s.esClient, scrollId)
			break
		}
		if elastic.IsNotFound(err) {
			break
		}
		if err != nil {
			cancelScroll(s.esClient, scrollId)
			return resp, errors.Wrap(err, "unable to retrieve messages")
		}

		scrollId = result.ScrollId

		for _, hit := range result.Hits.Hits {
			msg := &hackathon.Message{}
			err = json.Unmarshal(*hit.Source, msg)
			if err != nil {
				log.Infof("ERROR: Unable to unmarshal message: %s\n", err)
				continue
			}
			messages = append(messages, msg)
		}
	}
	resp.Messages = messages

	return resp, nil
}

func cancelScroll(client *elastic.Client, id string) {
	if len(id) == 0 {
		fmt.Printf("attempting to cancel empty scroll ID\n")
		return
	}

	_, err := client.ClearScroll(id).Do()
	if err != nil {
		fmt.Printf("Error clearing the scroll id: %v", err)
	}
}
