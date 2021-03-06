package service

import (
	"context"
	"log"
	"time"

	"github.com/influenzanet/go-utils/pkg/api_types"
	"github.com/influenzanet/go-utils/pkg/constants"
	"github.com/influenzanet/go-utils/pkg/token_checks"
	loggingAPI "github.com/influenzanet/logging-service/pkg/api"
	"github.com/influenzanet/study-service/pkg/api"
	"github.com/influenzanet/study-service/pkg/dbs/studydb"
	"github.com/influenzanet/study-service/pkg/types"
	"github.com/influenzanet/study-service/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *studyServiceServer) EnterStudy(ctx context.Context, req *api.EnterStudyRequest) (*api.AssignedSurveys, error) {
	if req == nil || token_checks.IsTokenEmpty(req.Token) || req.StudyKey == "" {
		return nil, status.Error(codes.InvalidArgument, "missing argument")
	}

	if err := utils.CheckIfProfileIDinToken(req.Token, req.ProfileId); err != nil {
		s.SaveLogEvent(req.Token.InstanceId, req.Token.Id, loggingAPI.LogEventType_SECURITY, constants.LOG_EVENT_WRONG_PROFILE_ID, "enter study:"+req.ProfileId)
		return nil, status.Error(codes.Internal, "permission denied")
	}

	// ParticipantID
	participantID, err := s.profileIDToParticipantID(req.Token.InstanceId, req.StudyKey, req.ProfileId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Exists already?
	exists := s.checkIfParticipantExists(req.Token.InstanceId, req.StudyKey, participantID, "active")
	if exists {
		log.Printf("error: participant (%s) already exists for this study", participantID)
		return nil, status.Error(codes.Internal, "participant already exists for this study")
	}

	// Init state and perform rules
	pState := types.ParticipantState{
		ParticipantID: participantID,
		EnteredAt:     time.Now().Unix(),
		StudyStatus:   "active",
	}

	// perform study rules/actions
	currentEvent := types.StudyEvent{
		Type:       "ENTER",
		InstanceID: req.Token.InstanceId,
		StudyKey:   req.StudyKey,
	}
	pState, err = s.getAndPerformStudyRules(req.Token.InstanceId, req.StudyKey, pState, currentEvent)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// save state back to DB
	pState, err = s.studyDBservice.SaveParticipantState(req.Token.InstanceId, req.StudyKey, pState)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Prepare response
	resp := api.AssignedSurveys{
		Surveys: []*api.AssignedSurvey{},
	}
	for _, as := range pState.AssignedSurveys {
		cs := as.ToAPI()
		cs.StudyKey = req.StudyKey
		resp.Surveys = append(resp.Surveys, cs)
	}
	return &resp, nil
}

func (s *studyServiceServer) GetAssignedSurveys(ctx context.Context, req *api_types.TokenInfos) (*api.AssignedSurveys, error) {
	if token_checks.IsTokenEmpty(req) {
		return nil, status.Error(codes.InvalidArgument, "missing argument")
	}

	studies, err := s.studyDBservice.GetStudiesByStatus(req.InstanceId, types.STUDY_STATUS_ACTIVE, true)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// for every profile form the token
	profileIDs := []string{req.ProfilId}
	profileIDs = append(profileIDs, req.OtherProfileIds...)

	resp := api.AssignedSurveys{
		Surveys:     []*api.AssignedSurvey{},
		SurveyInfos: []*api.SurveyInfo{},
	}
	for _, study := range studies {
		studySurveys, err := s.studyDBservice.FindAllSurveyDefsForStudy(req.InstanceId, study.Key, false)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		for _, profileID := range profileIDs {

			participantID, err := utils.ProfileIDtoParticipantID(profileID, s.StudyGlobalSecret, study.SecretKey)
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			pState, err := s.studyDBservice.FindParticipantState(req.InstanceId, study.Key, participantID)
			if err != nil || pState.StudyStatus != types.PARTICIPANT_STUDY_STATUS_ACTIVE {
				continue
			}

			for _, as := range pState.AssignedSurveys {
				cs := as.ToAPI()
				cs.StudyKey = study.Key
				cs.ProfileId = profileID
				resp.Surveys = append(resp.Surveys, cs)

				sDef := types.Survey{}
				for _, def := range studySurveys {
					if def.Current.SurveyDefinition.Key == cs.SurveyKey {
						sDef = def
						break
					}
				}

				found := false
				for _, info := range resp.SurveyInfos {
					if info.SurveyKey == sDef.Current.SurveyDefinition.Key && info.StudyKey == cs.StudyKey {
						found = true
						break
					}
				}
				if !found {
					apiS := sDef.ToAPI()
					resp.SurveyInfos = append(resp.SurveyInfos, &api.SurveyInfo{
						StudyKey:        cs.StudyKey,
						SurveyKey:       apiS.Current.SurveyDefinition.Key,
						Name:            apiS.Props.Name,
						Description:     apiS.Props.Description,
						TypicalDuration: apiS.Props.TypicalDuration,
					})
				}
			}
		}
	}

	return &resp, nil
}

func (s *studyServiceServer) GetAssignedSurvey(ctx context.Context, req *api.SurveyReferenceRequest) (*api.SurveyAndContext, error) {
	if req == nil || token_checks.IsTokenEmpty(req.Token) || req.StudyKey == "" || req.SurveyKey == "" {
		return nil, status.Error(codes.InvalidArgument, "missing argument")
	}

	if err := utils.CheckIfProfileIDinToken(req.Token, req.ProfileId); err != nil {
		s.SaveLogEvent(req.Token.InstanceId, req.Token.Id, loggingAPI.LogEventType_SECURITY, constants.LOG_EVENT_WRONG_PROFILE_ID, "get assigned survey:"+req.ProfileId)
		return nil, status.Error(codes.Internal, "permission denied")
	}

	// ParticipantID
	participantID, err := s.profileIDToParticipantID(req.Token.InstanceId, req.StudyKey, req.ProfileId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Get survey definition
	surveyDef, err := s.studyDBservice.FindSurveyDef(req.Token.InstanceId, req.StudyKey, req.SurveyKey)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	surveyContext, err := s.resolveContextRules(req.Token.InstanceId, req.StudyKey, participantID, surveyDef.ContextRules)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	prefill, err := s.resolvePrefillRules(req.Token.InstanceId, req.StudyKey, participantID, surveyDef.PrefillRules)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// empty irrelevant fields for this purpose
	surveyDef.ContextRules = nil
	surveyDef.PrefillRules = []types.Expression{}
	surveyDef.History = []types.SurveyVersion{}

	resp := api.SurveyAndContext{
		Survey:  surveyDef.ToAPI(),
		Context: surveyContext.ToAPI(),
	}
	if len(prefill.Responses) > 0 {
		resp.Prefill = prefill.ToAPI()
	}
	return &resp, nil
}

func (s *studyServiceServer) PostponeSurvey(ctx context.Context, req *api.PostponeSurveyRequest) (*api.AssignedSurveys, error) {
	if req == nil || token_checks.IsTokenEmpty(req.Token) || req.StudyKey == "" || req.SurveyKey == "" {
		return nil, status.Error(codes.InvalidArgument, "missing argument")
	}

	if err := utils.CheckIfProfileIDinToken(req.Token, req.ProfileId); err != nil {
		s.SaveLogEvent(req.Token.InstanceId, req.Token.Id, loggingAPI.LogEventType_SECURITY, constants.LOG_EVENT_WRONG_PROFILE_ID, "postpone survey:"+req.ProfileId)
		return nil, status.Error(codes.Internal, "permission denied")
	}

	// ParticipantID
	participantID, err := s.profileIDToParticipantID(req.Token.InstanceId, req.StudyKey, req.ProfileId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	pState, err := s.studyDBservice.FindParticipantState(req.Token.InstanceId, req.StudyKey, participantID)
	if err != nil {
		log.Println("PostponeSurvey: participant state not found")
		return nil, status.Error(codes.Internal, err.Error())
	}

	for i, as := range pState.AssignedSurveys {
		if as.SurveyKey == req.SurveyKey {
			newValidFrom := time.Now().Unix() + req.Delay

			if as.ValidUntil > 0 {
				if newValidFrom > as.ValidUntil-1800 {
					// submit survey as empty
					emptyResponse := types.SurveyResponse{
						Key:           req.SurveyKey,
						ParticipantID: participantID,
						SubmittedAt:   time.Now().Unix(),
						ArrivedAt:     time.Now().Unix(),
					}
					// perform study rules/actions
					currentEvent := types.StudyEvent{
						Type:       "SUBMIT",
						Response:   emptyResponse,
						InstanceID: req.Token.InstanceId,
						StudyKey:   req.StudyKey,
					}
					pState, err = s.getAndPerformStudyRules(req.Token.InstanceId, req.StudyKey, pState, currentEvent)
					if err != nil {
						return nil, status.Error(codes.Internal, err.Error())
					}
					break
				}
			}
			pState.AssignedSurveys[i].ValidFrom = newValidFrom
		}
	}

	// save state back to DB
	pState, err = s.studyDBservice.SaveParticipantState(req.Token.InstanceId, req.StudyKey, pState)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := api.AssignedSurveys{
		Surveys: []*api.AssignedSurvey{},
	}
	for _, as := range pState.AssignedSurveys {
		cs := as.ToAPI()
		cs.StudyKey = req.StudyKey
		resp.Surveys = append(resp.Surveys, cs)
	}

	return &resp, nil
}

func (s *studyServiceServer) SubmitStatusReport(ctx context.Context, req *api.StatusReportRequest) (*api.AssignedSurveys, error) {
	if req == nil || token_checks.IsTokenEmpty(req.Token) || req.StatusSurvey == nil {
		return nil, status.Error(codes.InvalidArgument, "missing argument")
	}

	if err := utils.CheckIfProfileIDinToken(req.Token, req.ProfileId); err != nil {
		s.SaveLogEvent(req.Token.InstanceId, req.Token.Id, loggingAPI.LogEventType_SECURITY, constants.LOG_EVENT_WRONG_PROFILE_ID, "submit status report:"+req.ProfileId)
		return nil, status.Error(codes.Internal, "permission denied")
	}

	studies, err := s.studyDBservice.GetStudiesByStatus(req.Token.InstanceId, "active", true)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp := api.AssignedSurveys{
		Surveys: []*api.AssignedSurvey{},
	}
	for _, study := range studies {
		participantID, err := utils.ProfileIDtoParticipantID(req.ProfileId, s.StudyGlobalSecret, study.SecretKey)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		pState, err := s.studyDBservice.FindParticipantState(req.Token.InstanceId, study.Key, participantID)
		if err != nil {
			// user not in the study - log.Println(err)
			continue
		}

		if pState.StudyStatus != types.PARTICIPANT_STUDY_STATUS_ACTIVE {
			continue
		}

		// Save responses
		response := types.SurveyResponseFromAPI(req.StatusSurvey)
		response.ParticipantID = participantID
		err = s.studyDBservice.AddSurveyResponse(req.Token.InstanceId, study.Key, response)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		// perform study rules/actions
		currentEvent := types.StudyEvent{
			Type:       "SUBMIT",
			Response:   response,
			InstanceID: req.Token.InstanceId,
			StudyKey:   study.Key,
		}
		pState, err = s.getAndPerformStudyRules(req.Token.InstanceId, study.Key, pState, currentEvent)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		// save state back to DB
		pState, err = s.studyDBservice.SaveParticipantState(req.Token.InstanceId, study.Key, pState)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		for _, as := range pState.AssignedSurveys {
			cs := as.ToAPI()
			cs.StudyKey = study.Key
			resp.Surveys = append(resp.Surveys, cs)
		}
	}
	return &resp, nil
}

func (s *studyServiceServer) SubmitResponse(ctx context.Context, req *api.SubmitResponseReq) (*api.AssignedSurveys, error) {
	if req == nil || token_checks.IsTokenEmpty(req.Token) || req.StudyKey == "" || req.Response == nil || len(req.Response.Responses) < 1 {
		return nil, status.Error(codes.InvalidArgument, "missing argument")
	}

	if err := utils.CheckIfProfileIDinToken(req.Token, req.ProfileId); err != nil {
		s.SaveLogEvent(req.Token.InstanceId, req.Token.Id, loggingAPI.LogEventType_SECURITY, constants.LOG_EVENT_WRONG_PROFILE_ID, "submit responses study:"+req.ProfileId)
		return nil, status.Error(codes.Internal, "permission denied")
	}

	// ParticipantID
	participantID, err := s.profileIDToParticipantID(req.Token.InstanceId, req.StudyKey, req.ProfileId)
	if err != nil {
		return nil, status.Error(codes.Internal, "could not compute participant id")
	}

	pState, err := s.studyDBservice.FindParticipantState(req.Token.InstanceId, req.StudyKey, participantID)
	if err != nil {
		return nil, status.Error(codes.Internal, "participant state not found")
	}
	if pState.StudyStatus != types.PARTICIPANT_STUDY_STATUS_ACTIVE {
		return nil, status.Error(codes.Internal, "user is not active in the current study")
	}

	// Save responses
	response := types.SurveyResponseFromAPI(req.Response)
	response.ParticipantID = participantID
	err = s.studyDBservice.AddSurveyResponse(req.Token.InstanceId, req.StudyKey, response)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// perform study rules/actions
	currentEvent := types.StudyEvent{
		Type:       "SUBMIT",
		Response:   response,
		InstanceID: req.Token.InstanceId,
		StudyKey:   req.StudyKey,
	}
	pState, err = s.getAndPerformStudyRules(req.Token.InstanceId, req.StudyKey, pState, currentEvent)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// save state back to DB
	pState, err = s.studyDBservice.SaveParticipantState(req.Token.InstanceId, req.StudyKey, pState)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Prepare response
	resp := api.AssignedSurveys{
		Surveys: []*api.AssignedSurvey{},
	}
	for _, as := range pState.AssignedSurveys {
		cs := as.ToAPI()
		cs.StudyKey = req.StudyKey
		resp.Surveys = append(resp.Surveys, cs)
	}
	return &resp, nil
}

func (s *studyServiceServer) LeaveStudy(ctx context.Context, req *api.LeaveStudyMsg) (*api.AssignedSurveys, error) {
	if req == nil || token_checks.IsTokenEmpty(req.Token) || req.StudyKey == "" {
		return nil, status.Error(codes.InvalidArgument, "missing argument")
	}

	if err := utils.CheckIfProfileIDinToken(req.Token, req.ProfileId); err != nil {
		s.SaveLogEvent(req.Token.InstanceId, req.Token.Id, loggingAPI.LogEventType_SECURITY, constants.LOG_EVENT_WRONG_PROFILE_ID, "leave study:"+req.ProfileId)
		return nil, status.Error(codes.Internal, "permission denied")
	}

	// ParticipantID
	participantID, err := s.profileIDToParticipantID(req.Token.InstanceId, req.StudyKey, req.ProfileId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	pState, err := s.studyDBservice.FindParticipantState(req.Token.InstanceId, req.StudyKey, participantID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if pState.StudyStatus != types.PARTICIPANT_STUDY_STATUS_ACTIVE {
		return nil, status.Error(codes.Internal, "not active in the study")
	}

	// Init state and perform rules
	pState = types.ParticipantState{
		ParticipantID: participantID,
		StudyStatus:   types.PARTICIPANT_STUDY_STATUS_EXITED,
	}
	// perform study rules/actions
	currentEvent := types.StudyEvent{
		Type:       "LEAVE",
		InstanceID: req.Token.InstanceId,
		StudyKey:   req.StudyKey,
	}
	pState, err = s.getAndPerformStudyRules(req.Token.InstanceId, req.StudyKey, pState, currentEvent)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	_, err = s.studyDBservice.SaveParticipantState(req.Token.InstanceId, req.StudyKey, pState)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Prepare response
	resp := api.AssignedSurveys{
		Surveys: []*api.AssignedSurvey{},
	}
	for _, as := range pState.AssignedSurveys {
		cs := as.ToAPI()
		cs.StudyKey = req.StudyKey
		resp.Surveys = append(resp.Surveys, cs)
	}
	return &resp, nil
}

func (s *studyServiceServer) DeleteParticipantData(ctx context.Context, req *api_types.TokenInfos) (*api.ServiceStatus, error) {
	if req == nil || token_checks.IsTokenEmpty(req) {
		return nil, status.Error(codes.InvalidArgument, "missing argument")
	}

	studies, err := s.studyDBservice.GetStudiesByStatus(req.InstanceId, "", true)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	profileIDs := []string{req.ProfilId}
	profileIDs = append(profileIDs, req.OtherProfileIds...)

	for _, study := range studies {
		for _, profileID := range profileIDs {
			// ParticipantID
			participantID, err := s.profileIDToParticipantID(req.InstanceId, study.Key, profileID)
			if err != nil {
				log.Printf("DeleteParticipantData: %v", err)
				continue
			}
			err = s.studyDBservice.DeleteParticipantState(req.InstanceId, study.Key, participantID)
			if err != nil {
				continue
			}
			_, err = s.studyDBservice.DeleteSurveyResponses(req.InstanceId, study.Key, studydb.ResponseQuery{ParticipantID: participantID})
			if err != nil {
				continue
			}
		}

	}
	return &api.ServiceStatus{
		Status: api.ServiceStatus_NORMAL,
		Msg:    "all responses deleted",
	}, nil
}
