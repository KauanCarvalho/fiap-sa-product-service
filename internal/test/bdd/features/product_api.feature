Feature: Product API Handling

  Scenario: Get product by SKU successfully
    Given a product with SKU "hamburger" exists
    When I send a GET request to path "/api/v1/products/hamburger"
    Then the response status should be 200
    And the response should contain "sku" with value "hamburger"
    And the response should contain "price" with value "10"

  Scenario: Product not found by SKU
    When I send a GET request to path "/api/v1/products/invalid-sku"
    Then the response status should be 404

  Scenario: Get products with valid filter parameters
    When I send a GET request to "/api/v1/products/?page=1&pageSize=10&category=food"
    Then the response status should be 200
    And the response should contain "currentPage" with value "1"
    And the response should contain "pageSize" with value "10"

  Scenario: Invalid page parameter
    When I send a GET request to "/api/v1/products/?page=abc"
    Then the response status should be 400

  Scenario: Invalid pageSize parameter
    When I send a GET request to "/api/v1/products/?pageSize=abc"
    Then the response status should be 400
