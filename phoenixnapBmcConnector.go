package exporter

import (
	"context"

	"github.com/phoenixnap/go-sdk-bmc/billingapi"
	"github.com/phoenixnap/go-sdk-bmc/bmcapi"
	"golang.org/x/oauth2/clientcredentials"
)

func getContext() context.Context {
	ctx := context.Background()

	return ctx
}

func getBillingApiClient(config clientcredentials.Config) billingapi.APIClient {
	configuration := billingapi.NewConfiguration()
	ctx := getContext()
	configuration.HTTPClient = config.Client(ctx)

	apiClient := billingapi.NewAPIClient(configuration)

	return *apiClient
}

func getBmcApiClient(config clientcredentials.Config) bmcapi.APIClient {
	configuration := bmcapi.NewConfiguration()
	ctx := getContext()
	configuration.HTTPClient = config.Client(ctx)

	apiClient := bmcapi.NewAPIClient(configuration)

	return *apiClient
}

func getTagValue(tags []bmcapi.TagAssignment, name string) string {
	for _, tag := range tags {
		if tag.GetName() == name {
			return tag.GetValue()
		}
	}
	return ""
}
