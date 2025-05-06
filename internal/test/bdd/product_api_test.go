package bdd_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/cucumber/godog"
)

var mockProducts = map[string]map[string]interface{}{
	"hamburger": {
		"sku":   "hamburger",
		"price": 10.00,
	},
}

var responseBody string

func givenAProductWithSKUExists(sku string) error {
	mockProducts[sku] = map[string]interface{}{
		"sku":   sku,
		"price": 10.00,
	}
	return nil
}

func iSendAGETRequestToProductBySKU(sku string) error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/products/%s", sku), nil)
	if err != nil {
		return err
	}
	recorder = httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)

	responseBody = recorder.Body.String()
	return nil
}

func iSendAGETRequestToProductsWithPageFilter(page string) error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/products/?page=%s", page), nil)
	if err != nil {
		return err
	}
	recorder = httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)

	responseBody = recorder.Body.String()

	return nil
}

func iSendAGETRequestToProductsWithPageSizeFilter(page string) error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/products/?pageSize=%s", page), nil)
	if err != nil {
		return err
	}
	recorder = httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)

	responseBody = recorder.Body.String()

	return nil
}

func iSendAGETRequestToProductsWithFilter(page, pageSize, category string) error {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/api/v1/products/?page=%s&pageSize=%s&category=%s",
			page,
			pageSize,
			category),
		nil)
	if err != nil {
		return err
	}
	recorder = httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)

	responseBody = recorder.Body.String()

	return nil
}

func theResponseStatusShouldBe(status int) error {
	if recorder.Code != status {
		return fmt.Errorf("expected status %d but got %d", status, recorder.Code)
	}
	return nil
}

func theResponseShouldContainFieldWithValue(field, expected string) error {
	var body map[string]interface{}
	if responseBody == "" {
		return errors.New("response body is empty, cannot parse")
	}

	if err := json.Unmarshal([]byte(responseBody), &body); err != nil {
		return fmt.Errorf("failed to parse response body: %v", err)
	}
	if value, exists := body[field]; exists {
		if fmt.Sprintf("%v", value) != expected {
			return fmt.Errorf("expected field %q to have value %q, but got %v", field, expected, value)
		}
	} else {
		return fmt.Errorf("field %q not found in response", field)
	}
	return nil
}

func InitializeScenarioProductAPI(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		ResetAndLoadFixtures()
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		ResetAndLoadFixtures()
		return ctx, nil
	})

	ctx.Step(`^a product with SKU "([^"]*)" exists$`, givenAProductWithSKUExists)
	ctx.Step(`^I send a GET request to path "/api/v1/products/([^"]+)"$`, iSendAGETRequestToProductBySKU)
	ctx.Step(`^I send a GET request to "/api/v1/products/\?page=([^&]+)&pageSize=([^&]+)&category=([^"]+)"$`,
		iSendAGETRequestToProductsWithFilter)
	ctx.Step(`^I send a GET request to "/api/v1/products/\?page=([^&]+)$`, iSendAGETRequestToProductsWithPageFilter)
	ctx.Step(`^I send a GET request to "/api/v1/products/\?pageSize=([^&]+)$`,
		iSendAGETRequestToProductsWithPageSizeFilter)
	ctx.Step(`^the response status should be (\d+)$`, theResponseStatusShouldBe)
	ctx.Step(`^the response should contain "([^"]*)" with value "([^"]*)"$`, theResponseShouldContainFieldWithValue)
}
