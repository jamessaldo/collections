package view

import (
	"authorization/domain/model"
	"authorization/service"
	"regexp"
	"strings"

	uuid "github.com/satori/go.uuid"
)

var uuidPattern = regexp.MustCompile(`\b[0-9a-f]{8}\b-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-\b[0-9a-f]{12}\b`)

func Authorization(userID string, method string, path string, uow *service.UnitOfWork, endpoints map[string]*model.Endpoint) (bool, error) {
	rePath := uuidPattern.ReplaceAllString(path, ":id")

	if endpoint, ok := endpoints[rePath+"_"+method]; ok {
		if strings.Contains(path, "/v1/team") {
			teamID := strings.Split(path, "/")[4]
			grantedEndpoints, err := uow.Endpoint.ListFilteredBy(uuid.FromStringOrNil(teamID), uuid.FromStringOrNil(userID))
			if err != nil {
				return false, err
			}
			return containsEndpoint(grantedEndpoints, endpoint), nil
		}
		return true, nil
	}

	return true, nil
}

func containsEndpoint(endpoints []model.Endpoint, endpoint *model.Endpoint) bool {
	for _, e := range endpoints {
		if e == *endpoint {
			return true
		}
	}
	return false
}
