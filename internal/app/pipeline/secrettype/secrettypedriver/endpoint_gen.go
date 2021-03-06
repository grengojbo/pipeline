// Code generated by mga tool. DO NOT EDIT.
package secrettypedriver

import (
	"github.com/banzaicloud/pipeline/internal/app/pipeline/secrettype"
	"github.com/go-kit/kit/endpoint"
	kitoc "github.com/go-kit/kit/tracing/opencensus"
	kitxendpoint "github.com/sagikazarmark/kitx/endpoint"
)

// Endpoint name constants
const (
	GetSecretTypeEndpoint   = "secrettype.GetSecretType"
	ListSecretTypesEndpoint = "secrettype.ListSecretTypes"
)

// Endpoints collects all of the endpoints that compose the underlying service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
type Endpoints struct {
	GetSecretType   endpoint.Endpoint
	ListSecretTypes endpoint.Endpoint
}

// MakeEndpoints returns a(n) Endpoints struct where each endpoint invokes
// the corresponding method on the provided service.
func MakeEndpoints(service secrettype.TypeService, middleware ...endpoint.Middleware) Endpoints {
	mw := kitxendpoint.Combine(middleware...)

	return Endpoints{
		GetSecretType:   kitxendpoint.OperationNameMiddleware(GetSecretTypeEndpoint)(mw(MakeGetSecretTypeEndpoint(service))),
		ListSecretTypes: kitxendpoint.OperationNameMiddleware(ListSecretTypesEndpoint)(mw(MakeListSecretTypesEndpoint(service))),
	}
}

// TraceEndpoints returns a(n) Endpoints struct where each endpoint is wrapped with a tracing middleware.
func TraceEndpoints(endpoints Endpoints) Endpoints {
	return Endpoints{
		GetSecretType:   kitoc.TraceEndpoint("secrettype.GetSecretType")(endpoints.GetSecretType),
		ListSecretTypes: kitoc.TraceEndpoint("secrettype.ListSecretTypes")(endpoints.ListSecretTypes),
	}
}
