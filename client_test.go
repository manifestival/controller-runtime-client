package client

import (
	"reflect"
	"testing"

	mf "github.com/manifestival/manifestival"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const owner = "manifestival"

var gracePeriodSeconds = int64(5)
var orphan = metav1.DeletePropagationOrphan

func TestDeleteOptions(t *testing.T) {
	uid := types.UID("test")
	tests := []struct {
		name     string
		options  []mf.DeleteOption
		expected client.DeleteOptions
	}{{
		name:     "GracePeriodSeconds",
		options:  []mf.DeleteOption{mf.DryRunAll, mf.GracePeriodSeconds(5), mf.IgnoreNotFound(false)},
		expected: client.DeleteOptions{GracePeriodSeconds: &gracePeriodSeconds},
	}, {
		name:     "Preconditions",
		options:  []mf.DeleteOption{mf.Preconditions(metav1.Preconditions{UID: &uid})},
		expected: client.DeleteOptions{Preconditions: &metav1.Preconditions{UID: &uid}},
	}, {
		name:     "PropagationPolicy",
		options:  []mf.DeleteOption{mf.PropagationPolicy(metav1.DeletePropagationOrphan)},
		expected: client.DeleteOptions{PropagationPolicy: &orphan},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := client.DeleteOptions{}
			actual.ApplyOptions(deleteWith(test.options))
			if !reflect.DeepEqual(actual, test.expected) {
				t.Errorf("wanted %v, got %v", test.expected, actual)
			}
		})
	}
}
