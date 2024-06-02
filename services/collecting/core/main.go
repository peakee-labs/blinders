package collecting

import (
	"fmt"
	"log"

	"blinders/packages/db/collectingdb"
	"blinders/packages/transport"
	"blinders/packages/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	ExplainLogsRepo   *collectingdb.ExplainLogsRepo
	TranslateLogsRepo *collectingdb.TranslateLogsRepo
}

func NewService(
	explainLogsRepo *collectingdb.ExplainLogsRepo,
	translateLogsRepo *collectingdb.TranslateLogsRepo,
) *Service {
	return &Service{
		ExplainLogsRepo:   explainLogsRepo,
		TranslateLogsRepo: translateLogsRepo,
	}
}

func (s Service) HandlePushEvent(event transport.Event) error {
	switch event.Type {
	case transport.AddTranslateLog:
		event, err := utils.JSONConvert[transport.AddTranslateLogEvent](event)
		if err != nil {
			log.Printf("invalid AddTranslateLogEvent, err: %v\n", err)
			return fmt.Errorf("invalid AddTranslateLogEvent")
		}
		_, err = s.TranslateLogsRepo.InsertRaw(&event.Payload)
		if err != nil {
			log.Println("can not insert translate log", err)
			return fmt.Errorf("can not insert translate log")
		}

	case transport.AddExplainLog:
		event, err := utils.JSONConvert[transport.AddExplainLogEvent](event)
		if err != nil {
			log.Printf("invalid AddExplainLogEvent, err: %v\n", err)
			return fmt.Errorf("invalid AddExplainLogEvent")
		}
		_, err = s.ExplainLogsRepo.InsertRaw(&event.Payload)
		if err != nil {
			log.Println("can not insert explain log", err)
			return fmt.Errorf("can not insert explain log")
		}

	default:
		log.Printf("event type mismatch: %v\n", event.Type)
		return fmt.Errorf("event type mismatch")
	}

	return nil
}

func (s Service) HandleGetRequest(request transport.Request) (any, error) {
	switch request.Type {
	case transport.GetTranslateLog:
		request, err := utils.JSONConvert[transport.GetCollectingLogRequest](
			request,
		)
		if err != nil {
			log.Printf("invalid GetCollectingLogRequest, err: %v\n", err)
			return nil, fmt.Errorf("invalid GetCollectingLogRequest")
		}
		userID, err := primitive.ObjectIDFromHex(request.Payload.UserID)
		if err != nil {
			log.Printf("invalid user id, err: %v\n", err)
			return nil, fmt.Errorf("invalid user id")
		}

		return s.TranslateLogsRepo.GetLogWithSmallestGetCountByUserID(userID)
	case transport.GetExplainLog:
		request, err := utils.JSONConvert[transport.GetCollectingLogRequest](
			request,
		)
		if err != nil {
			log.Printf("invalid GetCollectingLogRequest, err: %v\n", err)
			return nil, fmt.Errorf("invalid GetCollectingLogRequest")
		}
		userID, err := primitive.ObjectIDFromHex(request.Payload.UserID)
		if err != nil {
			log.Printf("invalid user id, err: %v\n", err)
			return nil, fmt.Errorf("invalid user id")
		}

		return s.ExplainLogsRepo.GetLogWithSmallestGetCountByUserID(userID)

	case transport.GetExplainLogBatch:
		request, err := utils.JSONConvert[transport.GetCollectingLogRequest](
			request,
		)
		if err != nil {
			log.Println("invalid GetCollectingLogRequest", err)
			return nil, fmt.Errorf("invalid GetCollectingLogRequest, err: %v\n", err)
		}

		userID, err := primitive.ObjectIDFromHex(request.Payload.UserID)
		if err != nil {
			log.Println("invalid user id", err)
			return nil, fmt.Errorf("invalid user id")
		}

		logs, pagination, err := s.ExplainLogsRepo.GetLogWithPagination(
			userID,
			request.Payload.PagintionInfo,
		)
		if err != nil {
			log.Println("can not get explain log", err)
			return nil, fmt.Errorf("can not get explain log, err: %v", err)
		}

		return transport.GetExplainLogBatchResponse{
			Logs:          logs,
			PagintionInfo: *pagination,
		}, nil

	default:
		log.Printf("request type mismatch: %v\n", request.Type)
		return nil, fmt.Errorf("request type mismatch")
	}
}
