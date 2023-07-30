package view

import (
	"authorization/domain"
	"authorization/repository"
	"context"
	"regexp"
	"strings"

	uuid "github.com/satori/go.uuid"
)

var uuidPattern = regexp.MustCompile(`\b[0-9a-f]{8}\b-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-\b[0-9a-f]{12}\b`)

func Authorization(ctx context.Context, userID string, method string, path string, endpoints map[string]domain.Endpoint) (bool, error) {
	rePath := uuidPattern.ReplaceAllString(path, ":id")

	if endpoint, ok := endpoints[rePath+"_"+method]; ok {
		if strings.Contains(path, "/v1/team") {
			teamID := strings.Split(path, "/")[4]
			access, err := repository.Role.GetAccess(ctx, uuid.FromStringOrNil(teamID), uuid.FromStringOrNil(userID), endpoint)
			if err != nil {
				return false, err
			}
			return access.IsAllowed, nil
		}
		return true, nil
	}

	return true, nil
}
