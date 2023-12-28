package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/crossplane/crossplane-runtime/pkg/fieldpath"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/function-sdk-go/errors"
	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/request"
	"github.com/crossplane/function-sdk-go/resource"
	"github.com/crossplane/function-sdk-go/response"
	"github.com/crossplane/function-template-go/input/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Function returns whatever response you ask it to.
type Function struct {
	fnv1beta1.UnimplementedFunctionRunnerServiceServer

	log logging.Logger
}

// RunFunction runs the Function.
func (f *Function) RunFunction(_ context.Context, req *fnv1beta1.RunFunctionRequest) (*fnv1beta1.RunFunctionResponse, error) {
	f.log.Info("Running function", "tag", req.GetMeta().GetTag())

	rsp := response.To(req, response.DefaultTTL)

	in := &v1beta1.RandString{}
	if err := request.GetInput(req, in); err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "Can not get the function input from %T", req))
		return rsp, nil
	}

	desired, err := request.GetDesiredComposedResources(req)
	f.log.Info("Desired resource", "Object", desired, "ns", resource.Name("test"))
	if err != nil {
		return nil, err
	}
	observed, err := request.GetObservedComposedResources(req)
	f.log.Info("Observed Resource", "res", observed, "ns", observed)
	cmp, err := request.GetObservedCompositeResource(req)
	f.log.Info("Observed composite  Resource", "res", cmp, "ns", observed)
	if err != nil {
		return nil, err
	}
	b := make([]byte, in.Cfg.RandStr.Length/2)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	randString := fmt.Sprintf("%x", b)

	for _, obj := range in.Cfg.Objs {
		f.log.Info("Name of the", "object", obj)
		if observed[resource.Name(obj.Name)].Resource != nil {
			observedPaved, err := fieldpath.PaveObject(observed[resource.Name(obj.Name)].Resource)
			if err != nil {
				return nil, err
			}
			randString, err = observedPaved.GetString(obj.FieldPath)
			if err != nil {
				return nil, err
			}
		}
		if observed[resource.Name(obj.Name)].Resource == nil && obj.Prefix != "" {
			err = patchFieldValueToObject(obj.FieldPath, obj.Prefix+randString, desired[resource.Name(obj.Name)].Resource)
		} else {
			f.log.Info("here on before error", "resname", desired[resource.Name(obj.Name)].Resource, "path", obj.FieldPath)
			err = patchFieldValueToObject(obj.FieldPath, randString, desired[resource.Name(obj.Name)].Resource)
			if err != nil {
				f.log.Debug("Error", "is", err, "obj", obj.FieldPath)
			}
		}
	}
	err = response.SetDesiredComposedResources(rsp, desired)
	return rsp, nil
}

// actual does the job

func patchFieldValueToObject(fieldPath string, value any, to runtime.Object) error {
	paved, err := fieldpath.PaveObject(to)
	if err != nil {
		return err
	}

	if err := paved.SetValue(fieldPath, value); err != nil {
		return err
	}
	return runtime.DefaultUnstructuredConverter.FromUnstructured(paved.UnstructuredContent(), to)
}
