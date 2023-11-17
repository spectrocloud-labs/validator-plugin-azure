package test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/spectrocloud-labs/validator/pkg/types"
)

func CheckTestCase(t *testing.T, res *types.ValidationResult, expectedResult types.ValidationResult, err, expectedError error) {
	if !reflect.DeepEqual(res.State, expectedResult.State) {
		t.Errorf("expected state (%+v), got (%+v)", expectedResult.State, res.State)
	}
	if !reflect.DeepEqual(res.Condition.ValidationType, expectedResult.Condition.ValidationType) {
		t.Errorf("expected validation type (%s), got (%s)", expectedResult.Condition.ValidationType, res.Condition.ValidationType)
	}
	if !reflect.DeepEqual(res.Condition.ValidationRule, expectedResult.Condition.ValidationRule) {
		t.Errorf("expected validation rule (%s), got (%s)", expectedResult.Condition.ValidationRule, res.Condition.ValidationRule)
	}
	if !reflect.DeepEqual(res.Condition.Message, expectedResult.Condition.Message) {
		t.Errorf("expected message (%s), got (%s)", expectedResult.Condition.Message, res.Condition.Message)
	}
	if !reflect.DeepEqual(res.Condition.Details, expectedResult.Condition.Details) {
		t.Errorf("expected details (%s), got (%s)", expectedResult.Condition.Details, res.Condition.Details)
	}
	if !reflect.DeepEqual(res.Condition.Failures, expectedResult.Condition.Failures) {
		t.Errorf("expected failures (%s), got (%s)", expectedResult.Condition.Failures, res.Condition.Failures)
	}
	if !reflect.DeepEqual(res.Condition.Status, expectedResult.Condition.Status) {
		t.Errorf("expected status (%s), got (%s)", expectedResult.Condition.Status, res.Condition.Status)
	}
	if !errors.Is(err, expectedError) {
		t.Errorf("expected error (%v), got (%v)", expectedError, err)
	}
}
