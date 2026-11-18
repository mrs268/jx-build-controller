package tekton_test

import (
	"testing"

	"github.com/jenkins-x-plugins/jx-pipeline/pkg/pipelines"
	v1 "github.com/jenkins-x/jx-api/v4/pkg/apis/jenkins.io/v1"
	"github.com/stretchr/testify/require"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis"
)

func TestTaskRunMessagesWithActivityStepVariants(t *testing.T) {
	tests := []struct {
		name string
		step *v1.PipelineActivityStep
	}{
		{
			name: "Stage",
			step: &v1.PipelineActivityStep{
				Kind: v1.ActivityStepKindTypeStage,
				Stage: &v1.StageActivityStep{
					CoreActivityStep: v1.CoreActivityStep{Name: "build"},
				},
			},
		},
		{
			name: "Preview",
			step: &v1.PipelineActivityStep{
				Kind: v1.ActivityStepKindTypePreview,
				Preview: &v1.PreviewActivityStep{
					CoreActivityStep: v1.CoreActivityStep{Message: "preview unchanged"},
					ApplicationURL:   "https://preview.example.com",
				},
			},
		},
		{
			name: "Promote",
			step: &v1.PipelineActivityStep{
				Kind: v1.ActivityStepKindTypePromote,
				Promote: &v1.PromoteActivityStep{
					CoreActivityStep: v1.CoreActivityStep{Message: "promote unchanged"},
					Environment:      "staging",
				},
			},
		},
		{name: "No activity steps"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taskStatus := &v1beta1.TaskRunStatus{}
			taskStatus.Conditions = append(taskStatus.Conditions, apis.Condition{
				Type:    apis.ConditionSucceeded,
				Status:  corev1.ConditionUnknown,
				Message: "TaskRun is pending",
			})
			pr := &v1beta1.PipelineRun{}
			pr.Status.TaskRuns = map[string]*v1beta1.PipelineRunTaskRunStatus{
				"build": {PipelineTaskName: "build", Status: taskStatus},
			}
			pa := &v1.PipelineActivity{}
			if tt.step != nil {
				pa.Spec.Steps = []v1.PipelineActivityStep{*tt.step.DeepCopy()}
			}

			// A pending TaskRun has conditions but no generated stage steps yet.
			require.NotPanics(t, func() {
				pipelines.ToPipelineActivity(pr, pa, true)
			})

			if tt.step == nil {
				require.Empty(t, pa.Spec.Steps)
			} else if tt.step.Stage != nil {
				require.Len(t, pa.Spec.Steps, 1)
				require.Equal(t, v1.ActivityMessageType("Stage is pending"), pa.Spec.Steps[0].Stage.Message)
			} else {
				require.Equal(t, []v1.PipelineActivityStep{*tt.step}, pa.Spec.Steps)
			}
		})
	}
}
