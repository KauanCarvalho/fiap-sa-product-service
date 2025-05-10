#!/bin/bash

# Check for required tools
command -v uuidgen >/dev/null 2>&1 || { echo >&2 "uuidgen is required but not installed. Aborting."; exit 1; }
command -v jq >/dev/null 2>&1 || { echo >&2 "jq is required but not installed. Aborting."; exit 1; }

base_url=$1
if [ -z "$base_url" ]; then
  echo "Usage: $0 <base_url>"
  exit 1
fi

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Store SKUs
skus=()

# Function to generate a random price
generate_random_price() {
  echo $(( RANDOM % 10000 + 1000 ))
}

# Function to format JSON or handle empty/invalid responses
print_pretty_body() {
  body="$1"
  if [ -z "$body" ]; then
    echo "no body"
  elif echo "$body" | jq . >/dev/null 2>&1; then
    echo "$body" | jq .
  else
    echo "$body"
  fi
}

# Function to color status code
print_status_code() {
  code="$1"
  if [ "$code" -ge 200 ] && [ "$code" -lt 300 ]; then
    echo -e "${GREEN}$code${NC}"
  elif [ "$code" -ge 400 ] && [ "$code" -lt 500 ]; then
    echo -e "${BLUE}$code${NC}"
  else
    echo -e "${RED}$code${NC}"
  fi
}

# Healthcheck before anything else
echo "Checking healthcheck..."
response=$(curl -s -w "\n%{http_code}" "$base_url/healthcheck")
body=$(echo "$response" | sed '$d')
status_code=$(echo "$response" | tail -n1)

print_status_code "$status_code"
print_pretty_body "$body"
echo "---------------------------------------------"

# Exit if healthcheck failed
if [ "$status_code" -lt 200 ] || [ "$status_code" -ge 300 ]; then
  echo "Healthcheck failed. Aborting."
  exit 1
fi

# Create 4 products
for i in {1..4}; do
  uuid=$(uuidgen)
  price=$(generate_random_price)

  echo "Creating product #$i"

  response=$(curl -s -w "\n%{http_code}" -X POST "$base_url/api/v1/admin/products/" \
    -H "Content-Type: application/json" \
    -d "{
      \"name\": \"name $uuid\",
      \"price\": $price,
      \"description\": \"description $uuid\",
      \"category\": {
          \"name\": \"lanche\"
      },
      \"images\": [
          {\"url\": \"https://placehold.co/$uuid\"},
          {\"url\": \"https://placehold.co/$uuid\"}
      ]
    }")

  body=$(echo "$response" | sed '$d')
  status_code=$(echo "$response" | tail -n1)

  sku=$(echo "$body" | jq -r '.sku // empty')
  skus+=("$sku")

  echo -n "Status: "
  print_status_code "$status_code"
  echo "Response body:"
  print_pretty_body "$body"
  echo "---------------------------------------------"
done

# Get first product
first_sku=${skus[0]}
echo "Fetching first created product (SKU: $first_sku)"
response=$(curl -s -w "\n%{http_code}" "$base_url/api/v1/products/$first_sku")
body=$(echo "$response" | sed '$d')
status_code=$(echo "$response" | tail -n1)

echo -n "Status: "
print_status_code "$status_code"
echo "Response body:"
print_pretty_body "$body"
echo "---------------------------------------------"

# Update penultimate product
penultimate_sku=${skus[2]}
new_uuid=$(uuidgen)
new_price=$(generate_random_price)

echo "Updating penultimate product (SKU: $penultimate_sku)"
response=$(curl -s -w "\n%{http_code}" -X PUT "$base_url/api/v1/admin/products/$penultimate_sku" \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"updated name $new_uuid\",
    \"price\": $new_price,
    \"description\": \"updated description $new_uuid\",
    \"category\": {
        \"name\": \"lanche\"
    },
    \"images\": [
        {\"url\": \"https://placehold.co/$new_uuid\"},
        {\"url\": \"https://placehold.co/$new_uuid\"}
    ]
  }")

body=$(echo "$response" | sed '$d')
status_code=$(echo "$response" | tail -n1)

echo -n "Status: "
print_status_code "$status_code"
echo "Response body:"
print_pretty_body "$body"
echo "---------------------------------------------"

# Delete last product
last_sku=${skus[3]}
echo "Deleting last product (SKU: $last_sku)"
response=$(curl -s -w "\n%{http_code}" -X DELETE "$base_url/api/v1/admin/products/$last_sku")
body=$(echo "$response" | sed '$d')
status_code=$(echo "$response" | tail -n1)

echo -n "Status: "
print_status_code "$status_code"
echo "Response body:"
print_pretty_body "$body"
echo "---------------------------------------------"

# Confirm deletion
echo "Confirming deletion (GET $last_sku)"
response=$(curl -s -w "\n%{http_code}" "$base_url/api/v1/products/$last_sku")
body=$(echo "$response" | sed '$d')
status_code=$(echo "$response" | tail -n1)

echo -n "Status: "
print_status_code "$status_code"
echo "Response body:"
print_pretty_body "$body"
echo "---------------------------------------------"

# List all products
echo "Listing products (page=1, pageSize=10)"
response=$(curl -s -w "\n%{http_code}" "$base_url/api/v1/products/?page=1&pageSize=10")
body=$(echo "$response" | sed '$d')
status_code=$(echo "$response" | tail -n1)

echo -n "Status: "
print_status_code "$status_code"
echo "Response body:"
print_pretty_body "$body"
echo "---------------------------------------------"
