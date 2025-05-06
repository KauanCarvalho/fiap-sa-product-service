package bdd_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/cucumber/godog"
)

func iSendAPostRequestToCreateProduct(payload string) error {
	req, err := http.NewRequest(http.MethodPost, "/api/v1/admin/products/", bytes.NewReader([]byte(payload)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	recorder = httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)
	responseBody = recorder.Body.String()

	return nil
}

func iSendAPutRequestToUpdateProduct(sku, payload string) error {
	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("/api/v1/admin/products/%s", sku),
		bytes.NewReader([]byte(payload)),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	recorder = httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)
	responseBody = recorder.Body.String()

	return nil
}
func iSendADeleteRequestToDeleteProduct(sku string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/admin/products/%s", sku), nil)
	if err != nil {
		return err
	}
	recorder = httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)
	responseBody = recorder.Body.String()

	return nil
}

func InitializeScenarioProductAdmin(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		ResetAndLoadFixtures()
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		ResetAndLoadFixtures()
		return ctx, nil
	})

	ctx.Step(`^a product with SKU "([^"]*)" exists$`, givenAProductWithSKUExists)
	ctx.Step(`^I send a POST request to create a product with payload:$`, iSendAPostRequestToCreateProduct)
	ctx.Step(`^the response status should be (\d+)$`, theResponseStatusShouldBe)
	ctx.Step(`^the response should contain "([^"]*)" with value "([^"]*)"$`, theResponseShouldContainFieldWithValue)
	ctx.Step(`^I send a PUT request to update product "([^"]*)" with payload:$`, iSendAPutRequestToUpdateProduct)
	ctx.Step(`^I send a DELETE request to delete product "([^"]*)"$`, iSendADeleteRequestToDeleteProduct)
}
