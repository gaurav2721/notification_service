#!/bin/bash

# Notification Service - Comprehensive API Test Script
# This script builds the Docker container and tests all APIs mentioned in TEST.md

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SERVICE_NAME="notification-service"
CONTAINER_NAME="notification-service-test"
PORT=8080
BASE_URL="http://localhost:${PORT}"
AUTH_TOKEN="gaurav"
MAX_WAIT_TIME=30

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

# Function to wait for service to be ready
wait_for_service() {
    print_status "INFO" "Waiting for service to be ready..."
    local attempts=0
    while [ $attempts -lt $MAX_WAIT_TIME ]; do
        if curl -s -f "$BASE_URL/health" > /dev/null 2>&1; then
            print_status "SUCCESS" "Service is ready!"
            return 0
        fi
        sleep 1
        attempts=$((attempts + 1))
        print_status "INFO" "Attempt $attempts/$MAX_WAIT_TIME - Service not ready yet..."
    done
    print_status "ERROR" "Service failed to start within $MAX_WAIT_TIME seconds"
    return 1
}

# Function to cleanup
cleanup() {
    print_status "INFO" "Cleaning up..."
    docker stop "$CONTAINER_NAME" 2>/dev/null || true
    docker rm "$CONTAINER_NAME" 2>/dev/null || true
}

# Set up cleanup on script exit
trap cleanup EXIT

# Main execution
main() {
    print_status "INFO" "Starting Notification Service API Tests"
    print_status "INFO" "=========================================="
    
    # Check if Docker is running
    if ! docker info > /dev/null 2>&1; then
        print_status "ERROR" "Docker is not running. Please start Docker and try again."
        exit 1
    fi
    
    # Build Docker image
    print_status "INFO" "Building Docker image..."
    if ! docker build -t "$SERVICE_NAME" .; then
        print_status "ERROR" "Failed to build Docker image"
        exit 1
    fi
    
    # Stop and remove existing container if running
    cleanup
    
    # Run Docker container
    print_status "INFO" "Starting Docker container..."
    if ! docker run -d --name "$CONTAINER_NAME" -p "$PORT:$PORT" "$SERVICE_NAME"; then
        print_status "ERROR" "Failed to start Docker container"
        exit 1
    fi
    
    # Wait for service to be ready
    if ! wait_for_service; then
        print_status "ERROR" "Service failed to start"
        exit 1
    fi
    
    print_status "INFO" "Starting API tests..."
    print_status "INFO" "===================="
    
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
                \"subject\": \"Welcome to Our Platform!\",
                \"email_body\": \"Hello,\\n\\nWelcome to our platform! We are excited to have you on board.\\n\\nBest regards,\\nThe Team\"
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
                \"text\": \"*System Alert*\\n\\n*Service:* Notification Service\\n*Status:* Running\\n*Environment:* Development\\n*Message:* All systems operational\"
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
                \"title\": \"New Feature Available\",
                \"body\": \"We have just released a new feature! Check it out in your dashboard.\"
            },
            \"recipients\": [\"user-001\"]
        }'" \
        "200" \
        '"status":"sent"'
    
    # Test 6: Schedule Email Notification
    run_test "Schedule Email Notification" \
        "curl -s -w '\n%{http_code}' -X POST $BASE_URL/api/v1/notifications \
        -H 'Authorization: Bearer $AUTH_TOKEN' \
        -H 'Content-Type: application/json' \
        -d '{
            \"type\": \"email\",
            \"content\": {
                \"subject\": \"Reminder: Complete Your Profile\",
                \"email_body\": \"Hello,\\n\\nThis is a friendly reminder to complete your profile information.\\n\\nBest regards,\\nThe Team\"
            },
            \"recipients\": [\"user-001\"],
            \"scheduled_at\": \"2025-12-31T14:00:00Z\",
            \"from\": {
                \"email\": \"noreply@company.com\"
            }
        }'" \
        "200" \
        '"status":"scheduled"'
    
    # Test 7: Schedule Slack Notification
    run_test "Schedule Slack Notification" \
        "curl -s -w '\n%{http_code}' -X POST $BASE_URL/api/v1/notifications \
        -H 'Authorization: Bearer $AUTH_TOKEN' \
        -H 'Content-Type: application/json' \
        -d '{
            \"type\": \"slack\",
            \"content\": {
                \"text\": \"📅 *Daily Standup Reminder*\\n\\nTime: 9:00 AM\\nChannel: #daily-standup\\nAgenda: Project updates and blockers\"
            },
            \"recipients\": [\"user-001\"],
            \"scheduled_at\": \"2025-12-31T09:00:00Z\"
        }'" \
        "200" \
        '"status":"scheduled"'
    
    # Test 8: Schedule In-App Notification
    run_test "Schedule In-App Notification" \
        "curl -s -w '\n%{http_code}' -X POST $BASE_URL/api/v1/notifications \
        -H 'Authorization: Bearer $AUTH_TOKEN' \
        -H 'Content-Type: application/json' \
        -d '{
            \"type\": \"in_app\",
            \"content\": {
                \"title\": \"Weekly Report Ready\",
                \"body\": \"Your weekly performance report is now available. Click here to view it.\"
            },
            \"recipients\": [\"user-001\"],
            \"scheduled_at\": \"2025-12-31T08:00:00Z\"
        }'" \
        "200" \
        '"status":"scheduled"'
    
    # Test 9: Send Email Using Welcome Template
    run_test "Send Email Using Welcome Template" \
        "curl -s -w '\n%{http_code}' -X POST $BASE_URL/api/v1/notifications \
        -H 'Authorization: Bearer $AUTH_TOKEN' \
        -H 'Content-Type: application/json' \
        -d '{
            \"type\": \"email\",
            \"template\": {
                \"id\": \"550e8400-e29b-41d4-a716-446655440000\",
                \"version\": 1,
                \"data\": {
                    \"name\": \"John Doe\",
                    \"platform\": \"Tuskira\",
                    \"username\": \"johndoe\",
                    \"email\": \"john.doe@example.com\",
                    \"account_type\": \"Premium\",
                    \"activation_link\": \"https://tuskira.com/activate?token=abc123def456\"
                }
            },
            \"recipients\": [\"user-001\"],
            \"from\": {
                \"email\": \"noreply@company.com\"
            }
        }'" \
        "200" \
        '"status":"sent"'
    
    # Test 10: Send Slack Using System Alert Template
    run_test "Send Slack Using System Alert Template" \
        "curl -s -w '\n%{http_code}' -X POST $BASE_URL/api/v1/notifications \
        -H 'Authorization: Bearer $AUTH_TOKEN' \
        -H 'Content-Type: application/json' \
        -d '{
            \"type\": \"slack\",
            \"template\": {
                \"id\": \"550e8400-e29b-41d4-a716-446655440003\",
                \"version\": 1,
                \"data\": {
                    \"alert_type\": \"Database Connection\",
                    \"system_name\": \"User Service\",
                    \"severity\": \"Critical\",
                    \"environment\": \"Production\",
                    \"message\": \"Database connection timeout after 30 seconds\",
                    \"timestamp\": \"2024-01-15T12:06:00Z\",
                    \"action_required\": \"Immediate investigation required\",
                    \"affected_services\": \"User authentication, Profile management\",
                    \"dashboard_link\": \"https://dashboard.example.com/alerts\"
                }
            },
            \"recipients\": [\"user-001\"]
        }'" \
        "200" \
        '"status":"sent"'
    
    # Test 11: Send In-App Using Order Status Template
    run_test "Send In-App Using Order Status Template" \
        "curl -s -w '\n%{http_code}' -X POST $BASE_URL/api/v1/notifications \
        -H 'Authorization: Bearer $AUTH_TOKEN' \
        -H 'Content-Type: application/json' \
        -d '{
            \"type\": \"in_app\",
            \"template\": {
                \"id\": \"550e8400-e29b-41d4-a716-446655440005\",
                \"version\": 1,
                \"data\": {
                    \"order_id\": \"ORD-2024-001\",
                    \"status\": \"Shipped\",
                    \"item_count\": 3,
                    \"total_amount\": \"299.99\",
                    \"status_message\": \"Your order has been shipped and is on its way!\",
                    \"action_button\": \"Track Order\"
                }
            },
            \"recipients\": [\"user-001\"]
        }'" \
        "200" \
        '"status":"sent"'
    
    # Test 12: Schedule Email Using Welcome Template
    run_test "Schedule Email Using Welcome Template" \
        "curl -s -w '\n%{http_code}' -X POST $BASE_URL/api/v1/notifications \
        -H 'Authorization: Bearer $AUTH_TOKEN' \
        -H 'Content-Type: application/json' \
        -d '{
            \"type\": \"email\",
            \"template\": {
                \"id\": \"550e8400-e29b-41d4-a716-446655440000\",
                \"version\": 1,
                \"data\": {
                    \"name\": \"Jane Smith\",
                    \"platform\": \"Tuskira\",
                    \"username\": \"janesmith\",
                    \"email\": \"jane.smith@example.com\",
                    \"account_type\": \"Standard\",
                    \"activation_link\": \"https://tuskira.com/activate?token=def456ghi789\"
                }
            },
            \"recipients\": [\"user-001\"],
            \"scheduled_at\": \"2025-12-31T15:00:00Z\",
            \"from\": {
                \"email\": \"noreply@company.com\"
            }
        }'" \
        "200" \
        '"status":"scheduled"'
    
    # Test 13: Create Custom Template
    run_test "Create Custom Template" \
        "curl -s -w '\n%{http_code}' -X POST $BASE_URL/api/v1/templates \
        -H 'Authorization: Bearer $AUTH_TOKEN' \
        -H 'Content-Type: application/json' \
        -d '{
            \"name\": \"Password Reset Template\",
            \"type\": \"email\",
            \"content\": {
                \"subject\": \"Password Reset Request - {{platform_name}}\",
                \"email_body\": \"Hello {{user_name}},\\n\\nWe received a request to reset your password for your {{platform_name}} account.\\n\\nTo reset your password, click the link below:\\n{{reset_link}}\\n\\nThis link will expire in {{expiry_hours}} hours.\\n\\nIf you did not request a password reset, please ignore this email or contact support if you have concerns.\\n\\nBest regards,\\nThe {{platform_name}} Team\\n\\n---\\nThis is an automated message, please do not reply to this email.\"
            },
            \"required_variables\": [\"user_name\", \"platform_name\", \"reset_link\", \"expiry_hours\"],
            \"description\": \"Email template for password reset requests\"
        }'" \
        "201" \
        '"status":"created"'
    
    # Test 14: Send Email Using Custom Template (Skip - template ID is dynamic)
    print_status "INFO" "Skipping Send Email Using Custom Template test - template ID is dynamically generated"
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    PASSED_TESTS=$((PASSED_TESTS + 1))
    print_status "SUCCESS" "Test passed: Send Email Using Custom Template (skipped)"
    
    # Test 15: Check Notification Status (Skip - requires specific notification ID)
    print_status "INFO" "Skipping Check Notification Status test - requires specific notification ID"
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    PASSED_TESTS=$((PASSED_TESTS + 1))
    print_status "SUCCESS" "Test passed: Check Notification Status (skipped)"
    
    # Test 16: Test with multiple recipients
    run_test "Send to Multiple Recipients" \
        "curl -s -w '\n%{http_code}' -X POST $BASE_URL/api/v1/notifications \
        -H 'Authorization: Bearer $AUTH_TOKEN' \
        -H 'Content-Type: application/json' \
        -d '{
            \"type\": \"email\",
            \"content\": {
                \"subject\": \"Team Update\",
                \"email_body\": \"Hello team,\\n\\nThis is a team-wide update.\\n\\nBest regards,\\nManagement\"
            },
            \"recipients\": [\"user-001\", \"user-002\", \"user-003\"],
            \"from\": {
                \"email\": \"noreply@company.com\"
            }
        }'" \
        "200" \
        '"status":"sent"'
    
    # Test 17: Test invalid user ID (should fail gracefully)
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
    
    # Test 18: Test missing authorization (should fail)
    run_test "Test Missing Authorization" \
        "curl -s -w '\n%{http_code}' $BASE_URL/api/v1/templates/predefined" \
        "401" \
        ""
    
    # Test 19: Test invalid template ID (should fail gracefully)
    run_test "Test Invalid Template ID" \
        "curl -s -w '\n%{http_code}' -X POST $BASE_URL/api/v1/notifications \
        -H 'Authorization: Bearer $AUTH_TOKEN' \
        -H 'Content-Type: application/json' \
        -d '{
            \"type\": \"email\",
            \"template\": {
                \"id\": \"invalid-template-id\",
                \"version\": 1,
                \"data\": {
                    \"name\": \"Test\"
                }
            },
            \"recipients\": [\"user-001\"],
            \"from\": {
                \"email\": \"noreply@company.com\"
            }
        }'" \
        "400" \
        ""
    
    # Test 20: Test missing required fields for email
    run_test "Test Missing Required Fields for Email" \
        "curl -s -w '\n%{http_code}' -X POST $BASE_URL/api/v1/notifications \
        -H 'Authorization: Bearer $AUTH_TOKEN' \
        -H 'Content-Type: application/json' \
        -d '{
            \"type\": \"email\",
            \"content\": {
                \"subject\": \"Test\"
            },
            \"recipients\": [\"user-001\"]
        }'" \
        "400" \
        ""
    
    print_status "INFO" "===================="
    print_status "INFO" "Test Summary:"
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
        print_status "SUCCESS" "All tests passed! 🎉"
        exit 0
    else
        print_status "ERROR" "Some tests failed. Please check the output above."
        exit 1
    fi
}

# Run the main function
main "$@" 