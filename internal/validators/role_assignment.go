package validators

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v2"
	"github.com/go-logr/logr"
	"github.com/spectrocloud-labs/validator-plugin-azure/api/v1alpha1"
	"github.com/spectrocloud-labs/validator-plugin-azure/internal/constants"
	azure_utils "github.com/spectrocloud-labs/validator-plugin-azure/internal/utils/azure"
	vapi "github.com/spectrocloud-labs/validator/api/v1alpha1"
	vapiconstants "github.com/spectrocloud-labs/validator/pkg/constants"
	"github.com/spectrocloud-labs/validator/pkg/types"
	"github.com/spectrocloud-labs/validator/pkg/util/ptr"
	corev1 "k8s.io/api/core/v1"
)

// Describes the data needed to validate a rule.
type roleAssignmentRule interface {
	GetRole() v1alpha1.Role
	GetServicePrincipalID() string
	GetSubscriptionID() string
}

// Describes the operations against Azure's role assignments API that this validator will be
// performing.
type roleAssignmentAPI interface {
	NewListForSubscriptionPager(options *armauthorization.RoleAssignmentsClientListForSubscriptionOptions) *runtime.Pager[armauthorization.RoleAssignmentsClientListForSubscriptionResponse]
}

type RoleAssignmentRuleService struct {
	log logr.Logger
	api roleAssignmentAPI
}

// NewRoleAssignmentRuleService creates a new instance of the service. Must be provided with an
// implementation of the role assignment API.
func NewRoleAssignmentRuleService(log logr.Logger, api roleAssignmentAPI) *RoleAssignmentRuleService {
	return &RoleAssignmentRuleService{
		log: log,
		api: api,
	}
}

// ReconcileRoleAssignmentRule reconciles a role assignment rule from a
// validation config.
func (s *RoleAssignmentRuleService) ReconcileRoleAssignmentRule(rule roleAssignmentRule) (*types.ValidationResult, error) {

	// Build the default ValidationResult for this role assignment rule
	vr := buildValidationResult(rule, constants.ValidationTypeRoleAssignment)

	// From Azure, get a list of all role assignments that have been created in
	// the subscription. These will implicitly have scopes of this subscription.
	// Therefore, there's no need to include a scope filter in the request.
	// TODO: Look into using filter style API request.
	filter := url.QueryEscape(fmt.Sprintf("principalId eq '%s'", rule.GetServicePrincipalID()))
	opts := &armauthorization.RoleAssignmentsClientListForSubscriptionOptions{
		Filter: &filter,
	}
	pager := s.api.NewListForSubscriptionPager(opts)
	var roleAssignmentsFound []*armauthorization.RoleAssignment
	for pager.More() {
		nextResult, err := pager.NextPage(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve next page of role assignment results from pager: %w", err)
		}
		if nextResult.RoleAssignmentListResult.Value != nil {
			roleAssignmentsFound = append(roleAssignmentsFound, nextResult.RoleAssignmentListResult.Value...)
		}
	}

	if len(roleAssignmentsFound) == 0 {
		// Failed validation result if no role assignments found.
		state := vapi.ValidationFailed
		condition := vapi.DefaultValidationCondition()
		condition.Status = corev1.ConditionFalse
		condition.Message = "No role assignments found."
		return &types.ValidationResult{
			Condition: &condition,
			State:     &state,
		}, nil
	}

	// Note: This is the name of the role, not the role name of the role.
	var roleName string

	// From the remaining role assignments, check whether there is one with the desired role. First,
	// find out whether we need to look the role up by its role name.
	if rule.GetRole().Name != nil {
		// The user has provided the name of the role directly, so we'll use it.
		roleName = *rule.GetRole().Name
	} else if rule.GetRole().RoleName != nil {
		// The user has not provided the name of the role, but they have provided the role name of
		// the role, so we'll try to use it to look up the name of the role.
		rolelookupMap, err := azure_utils.BuiltInRoleLookupMap(rule.GetSubscriptionID())
		if err != nil {
			return nil, fmt.Errorf("failed to get role lookup map: %w", err)
		}
		specifiedRoleName := *rule.GetRole().RoleName
		foundName, ok := rolelookupMap[specifiedRoleName]
		if !ok {
			failValidationResult(vr, rule, constants.ValidationTypeRoleAssignment, "Role name specified does not exist. Cannot validate.", []string{
				fmt.Sprintf("provided role name: %q", specifiedRoleName),
			})
			return vr, nil
		}
		roleName = foundName
	} else {
		// The user has provided neither, so we can't continue with validation.
		failValidationResult(vr, rule, constants.ValidationTypeRoleAssignment, "Neither role name nor name specified. Cannot validate.", []string{})
		return vr, nil
	}

	// We have the name of the role, so we can continue with validation. As soon as a match is
	// found, consider the validation successful.
	for _, ra := range roleAssignmentsFound {
		if ra.Properties == nil || ra.Properties.RoleDefinitionID == nil {
			return nil, errors.New("data from Azure API response malformed; missing Properties.RoleDefinitionID")
		}
		if azure_utils.RoleNameFromRoleDefinitionID(*ra.Properties.RoleDefinitionID) == roleName {
			return vr, nil
		}
	}

	failValidationResult(vr, rule, constants.ValidationTypeRoleAssignment, "Desired role assignment not found.", []string{
		fmt.Sprintf("specified role name of role? %t", rule.GetRole().RoleName != nil),
		fmt.Sprintf("specified name of role? %t", rule.GetRole().Name != nil),
		fmt.Sprintf("specified service principal ID = %q", rule.GetServicePrincipalID()),
		fmt.Sprintf("specified subscription ID = %q", rule.GetSubscriptionID()),
	})
	return vr, nil
}

// Creates a validation result that indicates validation succeeding. Parts of it
// can be overriden to indicate failure or provide more detail.
func buildValidationResult(rule roleAssignmentRule, validationType string) *types.ValidationResult {
	latestCondition := vapi.DefaultValidationCondition()

	latestCondition.Message = "Required role assignment was found."

	// Provided spec might not have enough data to make a good identifier here.
	identifier := "invalid-config"
	if rule.GetRole().RoleName != nil {
		identifier = *rule.GetRole().RoleName
	} else if rule.GetRole().Name != nil {
		identifier = *rule.GetRole().Name
	}

	latestCondition.ValidationRule = fmt.Sprintf("%s-%s", vapiconstants.ValidationRulePrefix, identifier)
	latestCondition.ValidationType = validationType

	return &types.ValidationResult{
		Condition: &latestCondition,
		State:     ptr.Ptr(vapi.ValidationSucceeded),
	}
}

// Modifies a validation result so that it is set to a failed state and the desired error message is
// used as its condition message.
func failValidationResult(result *types.ValidationResult, rule roleAssignmentRule, validationType, message string, details []string) {
	// Provided spec might not have enough data to make a good identifier here.
	identifier := "invalid-config"
	if rule.GetRole().RoleName != nil {
		identifier = *rule.GetRole().RoleName
	} else if rule.GetRole().Name != nil {
		identifier = *rule.GetRole().Name
	}

	result.Condition.Details = details
	result.Condition.Message = message
	result.Condition.Status = corev1.ConditionFalse
	result.Condition.ValidationRule = fmt.Sprintf("%s-%s", vapiconstants.ValidationRulePrefix, identifier)
	result.Condition.ValidationType = validationType

	result.State = ptr.Ptr(vapi.ValidationFailed)
}