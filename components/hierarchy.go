package components

import (
	"context"
	"fmt"
	"net/http"

	"github.com/SKF/go-rest-utility/client"
	"github.com/SKF/go-rest-utility/client/auth"
	"github.com/SKF/go-utility/v2/log"
	"github.com/pkg/errors"
)

const (
	hierarchyBaseURL = "https://api.%s.hierarchy.enlight.skf.com"
)

func httpClient(stage, identityToken string) *client.Client {
	return client.NewClient(
		client.WithBaseURL(fmt.Sprintf(hierarchyBaseURL, stage)),
		client.WithDatadogTracing(),
		client.WithTokenProvider(auth.RawToken(identityToken)),
	)
}

func Create(identityToken, stage, parentNodeID, componenttype string) (Component, error) {
	return CreateComponentWithContext(context.Background(), identityToken, stage, parentNodeID, Component{
		Type: componenttype,
	})
}

func CreateShaft(identityToken, stage, parentNodeID string, fixedSpeed int) (Component, error) {
	return CreateComponentWithContext(context.Background(), identityToken, stage, parentNodeID, Component{
		Type:       "shaft",
		FixedSpeed: &fixedSpeed,
	})
}

// Deprecated: Please use CreateComponentWithContext instead
func CreateWithContext(ctx context.Context, identityToken, stage, parentNodeID, componenttype string, fixedSpeed *int) (Component, error) {
	component := Component{
		Type:       componenttype,
		Position:   1,
		FixedSpeed: fixedSpeed,
	}

	return CreateComponentWithContext(ctx, identityToken, stage, parentNodeID, component)
}

func CreateComponentWithContext(ctx context.Context, identityToken, stage, parentNodeID string, component Component) (Component, error) {
	if component.Position == 0 {
		component.Position = 1
	}

	log.WithTracing(ctx).
		WithField("body", component).
		WithField("assetID", parentNodeID).
		Debugf("creating component")

	req := client.Post(fmt.Sprintf("/assets/%s/components", parentNodeID)).
		WithJSONPayload(component)

	restClient := httpClient(stage, identityToken)
	resp, err := restClient.Do(ctx, req)
	if err != nil {
		err = errors.Wrap(err, "failed to execute request")
		return Component{}, err
	}

	var responseBody struct {
		Component Component `json:"component"`
	}
	if err = resp.Unmarshal(&responseBody); err != nil {
		err = errors.Wrap(err, "failed to unmarshal response")
		return Component{}, err
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.Errorf("wrong response status: %q", resp.Status)
		return Component{}, err
	}

	return responseBody.Component, nil
}

type Component struct {
	ID                  string  `json:"id"`
	Type                string  `json:"type"`
	AttachedTo          string  `json:"attachedTo,omitempty"`
	Position            int     `json:"position"`
	Designation         *string `json:"designation"`
	FixedSpeed          *int    `json:"fixedSpeed"`
	Manufacturer        *string `json:"manufacturer"`
	PositionDescription *string `json:"positionDescription"`
	RotatingRing        *string `json:"rotatingRing"`
	SerialNumber        *string `json:"serialNumber"`
	ShaftSide           *string `json:"shaftSide"`
}
