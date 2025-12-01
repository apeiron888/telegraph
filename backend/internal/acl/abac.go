package acl

import "fmt"

// Policy represents an ABAC policy rule
type Policy struct {
	Attribute string      `json:"attribute"` // e.g., "mfa_enabled", "premium_status"
	Operator  string      `json:"operator"`  // "equals", "not_equals", "in", "contains"
	Value     interface{} `json:"value"`     // expected value
}

// EvaluatePolicy evaluates a single ABAC policy against user attributes
func EvaluatePolicy(userAttributes map[string]interface{}, policy Policy) (bool, error) {
	userValue, exists := userAttributes[policy.Attribute]
	if !exists {
		// Attribute doesn't exist - treat as false for security
		return false, nil
	}

	switch policy.Operator {
	case "equals":
		return userValue == policy.Value, nil

	case "not_equals":
		return userValue != policy.Value, nil

	case "exists":
		return exists, nil

	case "not_exists":
		return !exists, nil

	default:
		return false, fmt.Errorf("unknown operator: %s", policy.Operator)
	}
}

// EvaluatePolicies evaluates multiple policies with AND logic
// All policies must pass for access to be granted
func EvaluatePolicies(userAttributes map[string]interface{}, policies []Policy) (bool, error) {
	if len(policies) == 0 {
		return true, nil // No policies = allow
	}

	for _, policy := range policies {
		result, err := EvaluatePolicy(userAttributes, policy)
		if err != nil {
			return false, err
		}
		if !result {
			return false, nil // One failure = deny all
		}
	}

	return true, nil
}

// EvaluatePoliciesOR evaluates multiple policies with OR logic
// At least one policy must pass for access to be granted
func EvaluatePoliciesOR(userAttributes map[string]interface{}, policies []Policy) (bool, error) {
	if len(policies) == 0 {
		return true, nil
	}

	for _, policy := range policies {
		result, err := EvaluatePolicy(userAttributes, policy)
		if err != nil {
			return false, err
		}
		if result {
			return true, nil // One success = allow
		}
	}

	return false, nil
}

// Common ABAC policy helpers

// RequireMFA creates a policy requiring multi-factor authentication
func RequireMFA() Policy {
	return Policy{
		Attribute: "mfa_enabled",
		Operator:  "equals",
		Value:     true,
	}
}

// RequirePremium creates a policy requiring premium account
func RequirePremium() Policy {
	return Policy{
		Attribute: "premium_status",
		Operator:  "equals",
		Value:     true,
	}
}

// RequireRegion creates a policy for geographic restriction
func RequireRegion(region string) Policy {
	return Policy{
		Attribute: "geographic_region",
		Operator:  "equals",
		Value:     region,
	}
}
