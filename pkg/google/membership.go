package google

import (
	"fmt"
	"strings"

	"google.golang.org/api/cloudidentity/v1"
)

// CheckGroupMembershipForEmailIDs determines if the given emailIDs are already
// a member in the google-group with given groupID and divides them into an existing and noexisting list
// accordingly.
func CheckGroupMembershipForEmailIDs(groupID string, emailID string) (found *cloudidentity.Membership, hErr error) {
	// Get Group by groupID
	g, hErr := GetGroupByID(groupID)
	if hErr != nil {
		return nil, hErr
	}

	svc := instance.CloudIdentityService
	var nextPageToken string
	memberships := []*cloudidentity.Membership{}

	for {
		call := svc.Groups.Memberships.List(g.Name).PageSize(100)
		if nextPageToken != "" {
			call.PageToken(nextPageToken)
		}

		resp, err := call.Do()
		if err != nil {
			return nil, err
		}

		memberships = append(memberships, resp.Memberships...)

		if resp.NextPageToken == "" {
			break
		}
		nextPageToken = resp.NextPageToken
	}

	for _, m := range memberships {
		if strings.EqualFold(m.PreferredMemberKey.Id, strings.ToLower(emailID)) {
			return m, nil
		}
	}

	return nil, fmt.Errorf("member not found")
}
