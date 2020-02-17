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

func TestCreateOptions(t *testing.T) {
	tests := []struct {
		name     string
		options  []mf.ApplyOption
		expected []client.CreateOption
	}{{
		name:     "Manager w/DryRun",
		options:  []mf.ApplyOption{mf.DryRunAll, mf.FieldManager(owner)},
		expected: []client.CreateOption{client.DryRunAll, client.FieldOwner(owner)},
	}, {
		name:     "Manager",
		options:  []mf.ApplyOption{mf.FieldManager(owner)},
		expected: []client.CreateOption{client.FieldOwner(owner)},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := createWith(test.options)
			if !reflect.DeepEqual(actual, test.expected) {
				t.Errorf("wanted %v, got %v", test.expected, actual)
			}
		})
	}
}

func TestUpdateOptions(t *testing.T) {
	tests := []struct {
		name     string
		options  []mf.ApplyOption
		expected []client.UpdateOption
	}{{
		name:     "Manager w/DryRun",
		options:  []mf.ApplyOption{mf.DryRunAll, mf.FieldManager(owner)},
		expected: []client.UpdateOption{client.DryRunAll, client.FieldOwner(owner)},
	}, {
		name:     "Manager",
		options:  []mf.ApplyOption{mf.FieldManager(owner)},
		expected: []client.UpdateOption{client.FieldOwner(owner)},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := updateWith(test.options)
			if !reflect.DeepEqual(actual, test.expected) {
				t.Errorf("wanted %v, got %v", test.expected, actual)
			}
		})
	}
}

func TestDeleteOptions(t *testing.T) {
	uid := types.UID("test")
	tests := []struct {
		name     string
		options  []mf.DeleteOption
		expected []client.DeleteOption
	}{{
		name:     "GracePeriodSeconds",
		options:  []mf.DeleteOption{mf.DryRunAll, mf.GracePeriodSeconds(5), mf.IgnoreNotFound(false)},
		expected: []client.DeleteOption{client.DryRunAll, client.GracePeriodSeconds(5)},
	}, {
		name:     "Preconditions",
		options:  []mf.DeleteOption{mf.Preconditions(metav1.Preconditions{UID: &uid})},
		expected: []client.DeleteOption{client.Preconditions(metav1.Preconditions{UID: &uid})},
	}, {
		name:     "PropagationPolicy",
		options:  []mf.DeleteOption{mf.PropagationPolicy(metav1.DeletePropagationOrphan)},
		expected: []client.DeleteOption{client.PropagationPolicy(metav1.DeletePropagationOrphan)},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := deleteWith(test.options)
			if !reflect.DeepEqual(actual, test.expected) {
				t.Errorf("wanted %v, got %v", test.expected, actual)
			}
		})
	}
}
