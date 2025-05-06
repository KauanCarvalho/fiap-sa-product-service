Feature: Product management

  Scenario: Create a new product successfully
    Given a product with SKU "hamburger" exists
    When I send a POST request to create a product with payload:
      """
      {
        "name": "Coxinha",
        "price": 7.50,
        "description": "Coxinha de frango",
        "category": {"name": "lanche"},
        "images": [{"url": "https://example.com/coxinha.jpg"}]
      }
      """
    Then the response status should be 201
    And the response should contain "name" with value "Coxinha"

  Scenario: Invalid product payload
    When I send a POST request to create a product with payload:
      """
      {invalid-json}
      """
    Then the response status should be 400

  Scenario: Update an existing product successfully
    Given a product with SKU "hamburger" exists
    When I send a PUT request to update product "hamburger" with payload:
      """
      {
        "name": "Coxinha",
        "price": 7.50,
        "description": "Coxinha de frango",
        "category": {"name": "lanche"},
        "images": [{"url": "https://example.com/coxinha.jpg"}]
      }
      """
    Then the response status should be 200
    And the response should contain "name" with value "Coxinha"

  Scenario: Product not found for update
    When I send a PUT request to update product "nonexistentsku" with payload:
      """
      {
        "name": "Coxinha",
        "price": 7.50,
        "description": "Coxinha de frango",
        "category": {"name": "lanche"},
        "images": [{"url": "https://example.com/coxinha.jpg"}]
      }
      """
    Then the response status should be 404

  Scenario: Delete an existing product successfully
    Given a product with SKU "soda" exists
    When I send a DELETE request to delete product "hamburger"
    Then the response status should be 204

  Scenario: Product not found for deletion
    When I send a DELETE request to delete product "nonexistentsku"
    Then the response status should be 204
