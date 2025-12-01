package acl

import "fmt"

// SecurityLabel represents MAC security classification levels
type SecurityLabel string

const (
	LabelPublic       SecurityLabel = "public"
	LabelInternal     SecurityLabel = "internal"
	LabelConfidential SecurityLabel = "confidential"
)

// labelHierarchy defines security clearance levels (higher index = higher clearance)
var labelHierarchy = []SecurityLabel{LabelPublic, LabelInternal, LabelConfidential}

// getLabelLevel returns the clearance level of a security label
func getLabelLevel(label SecurityLabel) int {
	for i, l := range labelHierarchy {
		if l == label {
			return i
		}
	}
	return -1
}

// CanAccessResource checks if a user with userLabel can access a resource with resourceLabel
// MAC Rule: User clearance must be >= resource classification
func CanAccessResource(userLabel string, resourceLabel string) bool {
	userLevel := getLabelLevel(SecurityLabel(userLabel))
	resourceLevel := getLabelLevel(SecurityLabel(resourceLabel))

	// Invalid labels deny access
	if userLevel == -1 || resourceLevel == -1 {
		return false
	}

	return userLevel >= resourceLevel
}

// ValidateLabel checks if a security label is valid
func ValidateLabel(label string) error {
	for _, l := range labelHierarchy {
		if SecurityLabel(label) == l {
			return nil
		}
	}
	return fmt.Errorf("invalid security label: %s (must be public, internal, or confidential)", label)
}

// GetLabelName returns human-readable name for label
func GetLabelName(label SecurityLabel) string {
	names := map[SecurityLabel]string{
		LabelPublic:       "Public",
		LabelInternal:     "Internal",
		LabelConfidential: "Confidential",
	}
	return names[label]
}
