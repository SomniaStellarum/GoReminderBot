package main

import (
	"context"
	"fmt"
	"log"

	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

func (s *Server) parseText(sessionID, text string) (reply string, queryResult *dialogflowpb.QueryResult, err error) {
	log.Printf("SessionID: %v\n Message: %v", sessionID, text)
	ctx := context.Background()
	sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", s.projectID, sessionID)
	textInput := dialogflowpb.TextInput{Text: text, LanguageCode: "en"}
	queryTextInput := dialogflowpb.QueryInput_Text{Text: &textInput}
	queryInput := dialogflowpb.QueryInput{Input: &queryTextInput}
	request := dialogflowpb.DetectIntentRequest{Session: sessionPath, QueryInput: &queryInput}

	response, err := s.df.DetectIntent(ctx, &request)
	if err != nil {
		log.Printf("Error Detecting Intent: %v", err)
		return "", nil, err
	}

	queryResult = response.GetQueryResult()
	//if queryResult.AllRequiredParamsPresent {
	reply = queryResult.GetFulfillmentText()
	log.Printf("Reply: %v", reply)
	return reply, queryResult, nil
	//}

	//return reply, nil
}
