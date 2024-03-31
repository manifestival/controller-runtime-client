package client

import (
	"context"

	mf "github.com/manifestival/manifestival"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewManifest(pathname string, client client.Client, opts ...mf.Option) (mf.Manifest, error) {
	return mf.NewManifest(pathname, append(opts, mf.UseClient(NewClient(client)))...)
}

func NewClient(client client.Client) mf.Client {
	return &controllerRuntimeClient{client: client}
}

type controllerRuntimeClient struct {
	client client.Client
}

// verify implementation
var _ mf.Client = (*controllerRuntimeClient)(nil)

func (c *controllerRuntimeClient) Create(ctx context.Context, obj *unstructured.Unstructured, options ...mf.ApplyOption) error {
	return c.client.Create(ctx, obj, createWith(options)...)
}

func (c *controllerRuntimeClient) Update(ctx context.Context, obj *unstructured.Unstructured, options ...mf.ApplyOption) error {
	return c.client.Update(ctx, obj, updateWith(options)...)
}

func (c *controllerRuntimeClient) Delete(ctx context.Context, obj *unstructured.Unstructured, options ...mf.DeleteOption) error {
	err := c.client.Delete(ctx, obj, deleteWith(options)...)
	opts := mf.DeleteWith(options)
	if opts.IgnoreNotFound {
		return client.IgnoreNotFound(err)
	}
	return err
}

func (c *controllerRuntimeClient) Get(ctx context.Context, obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	key := client.ObjectKey{Namespace: obj.GetNamespace(), Name: obj.GetName()}
	result := &unstructured.Unstructured{}
	result.SetGroupVersionKind(obj.GroupVersionKind())
	err := c.client.Get(ctx, key, result)
	return result, err
}

func createWith(opts []mf.ApplyOption) []client.CreateOption {
	result := []client.CreateOption{}
	for _, opt := range opts {
		if opt == mf.DryRunAll {
			result = append(result, client.DryRunAll)
		}
		if t, ok := opt.(mf.FieldManager); ok {
			result = append(result, client.FieldOwner(string(t)))
		}
	}
	return result
}

func updateWith(opts []mf.ApplyOption) []client.UpdateOption {
	result := []client.UpdateOption{}
	for _, opt := range opts {
		if opt == mf.DryRunAll {
			result = append(result, client.DryRunAll)
		}
		if t, ok := opt.(mf.FieldManager); ok {
			result = append(result, client.FieldOwner(string(t)))
		}
	}
	return result
}

func deleteWith(opts []mf.DeleteOption) []client.DeleteOption {
	result := []client.DeleteOption{}
	for _, opt := range opts {
		if opt == mf.DryRunAll {
			result = append(result, client.DryRunAll)
			continue
		}
		switch v := opt.(type) {
		case mf.GracePeriodSeconds:
			result = append(result, client.GracePeriodSeconds(int64(v)))
		case mf.Preconditions:
			result = append(result, client.Preconditions(metav1.Preconditions(v)))
		case mf.PropagationPolicy:
			result = append(result, client.PropagationPolicy(metav1.DeletionPropagation(v)))
		}
	}
	return result
}
