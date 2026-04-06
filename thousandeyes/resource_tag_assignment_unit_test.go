package thousandeyes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tags"
)

func TestDiffTagAssignments(t *testing.T) {
	tests := []struct {
		name            string
		oldAssignments  []tags.Assignment
		newAssignments  []tags.Assignment
		expectedAdded   []string
		expectedRemoved []string
	}{
		{
			name: "add and remove assignments",
			oldAssignments: []tags.Assignment{
				testTagAssignment("test-1", tags.ASSIGNMENTTYPE_TEST),
				testTagAssignment("test-2", tags.ASSIGNMENTTYPE_TEST),
			},
			newAssignments: []tags.Assignment{
				testTagAssignment("test-2", tags.ASSIGNMENTTYPE_TEST),
				testTagAssignment("test-3", tags.ASSIGNMENTTYPE_TEST),
			},
			expectedAdded:   []string{"test|test-3"},
			expectedRemoved: []string{"test|test-1"},
		},
		{
			name: "keep identical assignments",
			oldAssignments: []tags.Assignment{
				testTagAssignment("test-1", tags.ASSIGNMENTTYPE_TEST),
			},
			newAssignments: []tags.Assignment{
				testTagAssignment("test-1", tags.ASSIGNMENTTYPE_TEST),
			},
			expectedAdded:   []string{},
			expectedRemoved: []string{},
		},
		{
			name: "different assignment type changes the key",
			oldAssignments: []tags.Assignment{
				testTagAssignment("shared-id", tags.ASSIGNMENTTYPE_TEST),
			},
			newAssignments: []tags.Assignment{
				testTagAssignment("shared-id", tags.AssignmentType("dashboard")),
			},
			expectedAdded:   []string{"dashboard|shared-id"},
			expectedRemoved: []string{"test|shared-id"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			addedAssignments, removedAssignments := diffTagAssignments(tc.oldAssignments, tc.newAssignments)

			assert.ElementsMatch(t, tc.expectedAdded, tagAssignmentKeys(addedAssignments))
			assert.ElementsMatch(t, tc.expectedRemoved, tagAssignmentKeys(removedAssignments))
		})
	}
}

func TestTagAssignmentKey(t *testing.T) {
	assert.Equal(t, "|", tagAssignmentKey(tags.Assignment{}))
	assert.Equal(t, "test|test-1", tagAssignmentKey(testTagAssignment("test-1", tags.ASSIGNMENTTYPE_TEST)))
}

func testTagAssignment(id string, assignmentType tags.AssignmentType) tags.Assignment {
	return tags.Assignment{
		Id:   &id,
		Type: &assignmentType,
	}
}

func tagAssignmentKeys(assignments []tags.Assignment) []string {
	keys := make([]string, 0, len(assignments))
	for _, assignment := range assignments {
		keys = append(keys, tagAssignmentKey(assignment))
	}
	return keys
}
