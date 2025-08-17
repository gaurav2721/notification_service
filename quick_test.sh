#!/bin/bash

# Notification Service - Quick Test Script
# This script tests APIs assuming the service is already running

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BASE_URL="http://localhost:8080"
AUTH_TOKEN="gaurav"

# Test counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Function to print colored output
print_status() {
    local status=$1
    local message=$2
    case $status in
        "INFO")
            echo -e "${BLUE}[INFO]${NC} $message"
            ;;
        "SUCCESS")
            echo -e "${GREEN}[SUCCESS]${NC} $message"
            ;;
        "WARNING")
            echo -e "${YELLOW}[WARNING]${NC} $message"
            ;;
        "ERROR")
            echo -e "${RED}[ERROR]${NC} $message"
            ;;
    esac
}

# Function to run a test
run_test() {
    local test_name="$1"
    local curl_command="$2"
    local expected_status="$3"
    local expected_pattern="$4"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    print_status "INFO" "Running test: $test_name"
    
    # Run the curl command and capture output
    local response
    local http_status
    local exit_code
    
    response=$(eval "$curl_command" 2>/dev/null) || exit_code=$?
    http_status=$(echo "$response" | tail -n1)
    response_body=$(echo "$response" | sed '$d')
    
    # Check if the command succeeded
    if [ $? -eq 0 ] || [ -n "$response" ]; then
        # Check HTTP status if expected
        if [ -n "$expected_status" ]; then
            if [[ "$http_status" == *"$expected_status"* ]]; then
                print_status "SUCCESS" "HTTP status check passed for: $test_name"
            else
                print_status "ERROR" "HTTP status check failed for: $test_name. Expected: $expected_status, Got: $http_status"
                FAILED_TESTS=$((FAILED_TESTS + 1))
                return 1
            fi
        fi
        
        # Check response pattern if expected
        if [ -n "$expected_pattern" ]; then
            if echo "$response_body" | grep -q "$expected_pattern"; then
                print_status "SUCCESS" "Response pattern check passed for: $test_name"
            else
                print_status "ERROR" "Response pattern check failed for: $test_name. Expected pattern: $expected_pattern"
                print_status "ERROR" "Response: $response_body"
                FAILED_TESTS=$((FAILED_TESTS + 1))
                return 1
            fi
        fi
        
        print_status "SUCCESS" "Test passed: $test_name"
        PASSED_TESTS=$((PASSED_TESTS + 1))
        return 0
    else
        print_status "ERROR" "Test failed: $test_name"
        print_status "ERROR" "Command: $curl_command"
        print_status "ERROR" "Exit code: $exit_code"
        FAILED_TESTS=$((FAILED_TESTS + 1))
        return 1
    fi
}

# Function to check if service is running
check_service() {
    if ! curl -s -f "$BASE_URL/health" > /dev/null 2>&1; then
        print_status "ERROR" "Service is not running at $BASE_URL"
        print_status "INFO" "Please start the service first using:"
        print_status "INFO" "  make docker-run"
        print_status "INFO" "  or"
        print_status "INFO" "  go run main.go"
        exit 1
    fi
    print_status "SUCCESS" "Service is running at $BASE_URL"
}

# Main execution
main() {
    print_status "INFO" "Starting Quick API Tests"
    print_status "INFO" "======================"
    
    # Check if service is running
    check_service
    
    print_status "INFO" "Running quick API tests..."
    print_status "INFO" "========================"
    
    # Test 1: Health Check
    run_test "Health Check" \
        "curl -s -w '\n%{http_code}' $BASE_URL/health" \
        "200" \
        '"status":"healthy"'
    
    # Test 2: Get Predefined Templates
    run_test "Get Predefined Templates" \
        "curl -s -w '\n%{http_code}' -H 'Authorization: Bearer $AUTH_TOKEN' $BASE_URL/api/v1/templates/predefined" \
        "200" \
        '"count":'
    
    # Test 3: Send Immediate Email Notification
    run_test "Send Immediate Email Notification" \
        "curl -s -w '\n%{http_code}' -X POST $BASE_URL/api/v1/notifications \
        -H 'Authorization: Bearer $AUTH_TOKEN' \
        -H 'Content-Type: application/json' \
        -d '{
            \"type\": \"email\",
            \"content\": {
                \"subject\": \"Quick Test Email\",
                \"email_body\": \"This is a quick test email.\"
            },
            \"recipients\": [\"user-001\"],
            \"from\": {
                \"email\": \"noreply@company.com\"
            }
        }'" \
        "200" \
        '"status":"sent"'
    
    # Test 4: Send Immediate Slack Notification
    run_test "Send Immediate Slack Notification" \
        "curl -s -w '\n%{http_code}' -X POST $BASE_URL/api/v1/notifications \
        -H 'Authorization: Bearer $AUTH_TOKEN' \
        -H 'Content-Type: application/json' \
        -d '{
            \"type\": \"slack\",
            \"content\": {
                \"text\": \"Quick test Slack message\"
            },
            \"recipients\": [\"user-001\"]
        }'" \
        "200" \
        '"status":"sent"'
    
    # Test 5: Send Immediate In-App Notification
    run_test "Send Immediate In-App Notification" \
        "curl -s -w '\n%{http_code}' -X POST $BASE_URL/api/v1/notifications \
        -H 'Authorization: Bearer $AUTH_TOKEN' \
        -H 'Content-Type: application/json' \
        -d '{
            \"type\": \"in_app\",
            \"content\": {
                \"title\": \"Quick Test\",
                \"body\": \"This is a quick test notification.\"
            },
            \"recipients\": [\"user-001\"]
        }'" \
        "200" \
        '"status":"sent"'
    
    # Test 6: Test missing authorization (should fail)
    run_test "Test Missing Authorization" \
        "curl -s -w '\n%{http_code}' $BASE_URL/api/v1/templates/predefined" \
        "401" \
        ""
    
    # Test 7: Test invalid user ID (should fail gracefully)
    run_test "Test Invalid User ID" \
        "curl -s -w '\n%{http_code}' -X POST $BASE_URL/api/v1/notifications \
        -H 'Authorization: Bearer $AUTH_TOKEN' \
        -H 'Content-Type: application/json' \
        -d '{
            \"type\": \"email\",
            \"content\": {
                \"subject\": \"Test\",
                \"email_body\": \"Test\"
            },
            \"recipients\": [\"invalid-user-id\"],
            \"from\": {
                \"email\": \"noreply@company.com\"
            }
        }'" \
        "500" \
        ""
    
    print_status "INFO" "========================"
    print_status "INFO" "Quick Test Summary:"
    print_status "INFO" "Total Tests: $TOTAL_TESTS"
    print_status "SUCCESS" "Passed: $PASSED_TESTS"
    if [ $FAILED_TESTS -gt 0 ]; then
        print_status "ERROR" "Failed: $FAILED_TESTS"
    else
        print_status "SUCCESS" "Failed: $FAILED_TESTS"
    fi
    
    # Calculate success rate
    if [ $TOTAL_TESTS -gt 0 ]; then
        success_rate=$((PASSED_TESTS * 100 / TOTAL_TESTS))
        print_status "INFO" "Success Rate: ${success_rate}%"
    fi
    
    if [ $FAILED_TESTS -eq 0 ]; then
        print_status "SUCCESS" "All quick tests passed! 🎉"
        exit 0
    else
        print_status "ERROR" "Some tests failed. Please check the output above."
        exit 1
    fi
}

# Run the main function
main "$@" 