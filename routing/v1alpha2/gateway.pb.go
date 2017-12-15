// Code generated by protoc-gen-go. DO NOT EDIT.
// source: routing/v1alpha2/gateway.proto

/*
Package istio_routing_v1alpha2 is a generated protocol buffer package.

It is generated from these files:
	routing/v1alpha2/gateway.proto
	routing/v1alpha2/route_rule.proto

It has these top-level messages:
	Gateway
	Server
	RouteRule
	Destination
	HTTPRoute
	TCPRoute
	HTTPMatchRequest
	DestinationWeight
	L4MatchAttributes
	HTTPRedirect
	HTTPRewrite
	StringMatch
	HTTPRetry
	CorsPolicy
	HTTPFaultInjection
	PortSelector
*/
package istio_routing_v1alpha2

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// The location of this service with respect to the service mesh.
type Server_Location int32

const (
	// Services that are managed by Istio service mesh.
	Server_MESH Server_Location = 0
	// External services are hosted services (e.g. Google Maps) that may be
	// consumed by services in the mesh.
	Server_EXTERNAL Server_Location = 1
)

var Server_Location_name = map[int32]string{
	0: "MESH",
	1: "EXTERNAL",
}
var Server_Location_value = map[string]int32{
	"MESH":     0,
	"EXTERNAL": 1,
}

func (x Server_Location) String() string {
	return proto.EnumName(Server_Location_name, int32(x))
}
func (Server_Location) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 0} }

// TLS modes enforced by the gateway
type Server_TLSOptions_TLSmode int32

const (
	// If set to "passthrough", the gateway will forward the connection to the
	// upstream server as is.
	Server_TLSOptions_PASSTHROUGH Server_TLSOptions_TLSmode = 0
	// If set to "simple", the gateway will secure connections with
	// standard TLS semantics (server certs only).
	Server_TLSOptions_SIMPLE Server_TLSOptions_TLSmode = 1
	// If set to "mutual", the gateway will use standard mTLS authentication.
	Server_TLSOptions_MUTUAL Server_TLSOptions_TLSmode = 2
	// If set to originate, the gateway implementation will use TLS on
	// the connection to the service. Applicable only to services whose
	// Location is EXTERNAL to the mesh. TLS to services inside the mesh
	// is controlled by global mTLS settings.
	Server_TLSOptions_ORIGINATE Server_TLSOptions_TLSmode = 3
)

var Server_TLSOptions_TLSmode_name = map[int32]string{
	0: "PASSTHROUGH",
	1: "SIMPLE",
	2: "MUTUAL",
	3: "ORIGINATE",
}
var Server_TLSOptions_TLSmode_value = map[string]int32{
	"PASSTHROUGH": 0,
	"SIMPLE":      1,
	"MUTUAL":      2,
	"ORIGINATE":   3,
}

func (x Server_TLSOptions_TLSmode) String() string {
	return proto.EnumName(Server_TLSOptions_TLSmode_name, int32(x))
}
func (Server_TLSOptions_TLSmode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor0, []int{1, 1, 0}
}

// Gateway describes a load balancer operating at the edge of the mesh
// receiving incoming or outgoing HTTP/TCP connections. The specification
// describes a set of ports that should be exposed, the type of protocol to
// use, SNI configuration for the load balancer, etc.
//
// For example, the following gateway spec sets up a proxy to act as a load
// balancer exposing port 80 and 9080 (http), 443 (https), and port 2379
// (TCP) for ingress.  While Istio will configure the proxy to listen on
// these ports, it is the responsibility of the user to ensure that
// external traffic to these ports are allowed into the mesh.
//
//     apiVersion: config.istio.io/v1alpha2
//     kind: Gateway
//     metadata:
//       name: my-gateway
//     spec:
//       servers:
//       - port:
//           number: 80
//           name: http
//         domains:
//         - uk.bookinfo.com
//         - eu.bookinfo.com
//         tls:
//           httpsRedirect: true # sends 302 redirect for http requests
//       - port:
//           number: 443
//           name: https
//         domains:
//         - uk.bookinfo.com
//         - eu.bookinfo.com
//         tls:
//           mode: simple #enables HTTPS on this port
//           serverCert: server.crt
//           clientCABundle: client.ca-bundle
//       - port:
//           number: 9080
//           name: http-wildcard
//         # no domains implies wildcard match
//       - port:
//           number: 2379 #to expose internal service via external port 2379
//           name: Mongo
//           protocol: MONGO
//
// The gateway specification above describes the L4-L6 properties of a load
// balancer. Routing rules can then be bound to a gateway to control
// the forwarding of traffic arriving at a particular domain or gateway port.
//
// The following sample route rule splits traffic for
// https://uk.bookinfo.com/reviews, https://eu.bookinfo.com/reviews,
// http://uk.bookinfo.com:9080/reviews, http://eu.bookinfo.com:9080/reviews
// into two versions (prod and qa) of an internal reviews service on port
// 9080. In addition, requests containing the cookie user: dev-123 will be
// sent to special port 7777 in the qa version. The same rule is also
// applicable inside the mesh for requests to the reviews.prod
// service. This rule is applicable across ports 443, 9080. Note that
// http://uk.bookinfo.com gets redirected to https://uk.bookinfo.com
// (i.e. 80 redirects to 443).
//
//     apiVersion: config.istio.io/v1alpha2
//     kind: RouteRule
//     metadata:
//       name: bookinfo-rule
//     spec:
//       hosts:
//       - reviews.prod
//       - uk.bookinfo.com
//       - eu.bookinfo.com
//       gateways:
//       - my-gateway
//       - mesh # applies to all the sidecars in the mesh
//       http:
//       - match:
//         - headers:
//             cookie:
//               user: dev-123
//         route:
//         - destination:
//             port:
//               number: 7777
//             name: reviews.qa
//       - match:
//           uri:
//             prefix: /reviews/
//         route:
//         - destination:
//             port:
//               number: 9080 # can be omitted if its the only port for reviews
//             name: reviews.prod
//           weight: 80
//         - destination:
//             name: reviews.qa
//           weight: 20
//
// The following routing rule forwards traffic arriving at (external) port
// 2379 from 172.17.16.0/24 subnet to internal Mongo server on port 5555. This
// rule is not applicable internally in the mesh as the gateway list omits
// the reserved name "mesh".
//
//     apiVersion: config.istio.io/v1alpha2
//     kind: RouteRule
//     metadata:
//       name: bookinfo-Mongo
//     spec:
//       hosts:
//       - Mongosvr #name of Mongo service
//       gateways:
//       - my-gateway
//       tcp:
//       - match:
//         - port:
//             number: 2379
//           sourceSubnet: "172.17.16.0/24"
//         route:
//         - destination:
//             name: mongo.prod
//
type Gateway struct {
	// REQUIRED: A list of server specifications.
	Servers []*Server `protobuf:"bytes,1,rep,name=servers" json:"servers,omitempty"`
}

func (m *Gateway) Reset()                    { *m = Gateway{} }
func (m *Gateway) String() string            { return proto.CompactTextString(m) }
func (*Gateway) ProtoMessage()               {}
func (*Gateway) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Gateway) GetServers() []*Server {
	if m != nil {
		return m.Servers
	}
	return nil
}

// Server describes the properties of the proxy on a given load balancer port.
// For example,
//
//     apiVersion: config.istio.io/v1alpha2
//     kind: Gateway
//     metadata:
//       name: my-ingress
//     spec:
//       servers:
//       - port:
//           number: 80
//           protocol: HTTP2
//         domains:
//         - uk.bookinfo.com
//         - eu.bookinfo.com
//
// Another example
//
//     apiVersion: config.istio.io/v1alpha2
//     kind: Gateway
//     metadata:
//       name: my-tcp-ingress
//     spec:
//       servers:
//       - port:
//           number: 27018
//           protocol: MONGO
//         domains:
//         - uk.bookinfo.com
//         - eu.bookinfo.com
//
// The following is an example of TLS configuration for port 443
//
//     apiVersion: config.istio.io/v1alpha2
//     kind: Gateway
//     metadata:
//       name: my-ingress
//     spec:
//       servers:
//       - port:
//           number: 443
//           protocol: HTTP
//         domains:
//         - uk.bookinfo.com
//         - eu.bookinfo.com
//         tls:
//           mode: simple
//           serverCertificate: server.crt
//
// The following is an example of a gateway abstracting an external
// service. The caller is expected to access the service using
// http://maps.google.com:443
//
//     apiVersion: config.istio.io/v1alpha2
//     kind: Gateway
//     metadata:
//       name: my-egress
//     spec:
//       servers:
//       - port:
//           number: 443
//           protocol: HTTP
//         domains:
//         - maps.google.com
//         tls:
//           mode: originate
//         location: external
//
type Server struct {
	// REQUIRED: The Port on which the proxy should listen for incoming
	// connections
	Port *Server_Port `protobuf:"bytes,1,opt,name=port" json:"port,omitempty"`
	// REQUIRED. A list of domains exposed by this gateway. While
	// typically applicable to HTTP services, it can also be used for TCP
	// services using TLS with SNI. Standard DNS wildcard prefix syntax
	// is permitted.
	//
	// RouteRules that are bound to a gateway must having a matching domain
	// in their default destination. Specifically one of the route rule
	// destination domains is a strict suffix of a gateway domain or
	// a gateway domain is a suffix of one of the route rule domains.
	Domains []string `protobuf:"bytes,2,rep,name=domains" json:"domains,omitempty"`
	// Set of TLS related options that govern the server's behavior. Use
	// these options to control if all http requests should be redirected to
	// https, and the TLS modes to use.
	Tls *Server_TLSOptions `protobuf:"bytes,3,opt,name=tls" json:"tls,omitempty"`
	// Optional: Indicates the location of the server with respect to the
	// other services in the mesh.
	Location Server_Location `protobuf:"varint,4,opt,name=location,enum=istio.routing.v1alpha2.Server_Location" json:"location,omitempty"`
}

func (m *Server) Reset()                    { *m = Server{} }
func (m *Server) String() string            { return proto.CompactTextString(m) }
func (*Server) ProtoMessage()               {}
func (*Server) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Server) GetPort() *Server_Port {
	if m != nil {
		return m.Port
	}
	return nil
}

func (m *Server) GetDomains() []string {
	if m != nil {
		return m.Domains
	}
	return nil
}

func (m *Server) GetTls() *Server_TLSOptions {
	if m != nil {
		return m.Tls
	}
	return nil
}

func (m *Server) GetLocation() Server_Location {
	if m != nil {
		return m.Location
	}
	return Server_MESH
}

// Port describes the properties of a specific port of a service.
type Server_Port struct {
	// REQUIRED: A valid non-negative integer port number.
	Number uint32 `protobuf:"varint,1,opt,name=number" json:"number,omitempty"`
	// The protocol exposed on the port.
	// MUST BE one of HTTP|HTTPS|GRPC|HTTP2|MONGO|TCP.
	Protocol string `protobuf:"bytes,2,opt,name=protocol" json:"protocol,omitempty"`
	// Label assigned to the port.
	Name string `protobuf:"bytes,3,opt,name=name" json:"name,omitempty"`
}

func (m *Server_Port) Reset()                    { *m = Server_Port{} }
func (m *Server_Port) String() string            { return proto.CompactTextString(m) }
func (*Server_Port) ProtoMessage()               {}
func (*Server_Port) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 0} }

func (m *Server_Port) GetNumber() uint32 {
	if m != nil {
		return m.Number
	}
	return 0
}

func (m *Server_Port) GetProtocol() string {
	if m != nil {
		return m.Protocol
	}
	return ""
}

func (m *Server_Port) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type Server_TLSOptions struct {
	// If set to true, the load balancer will send a 302 redirect for all
	// http connections, asking the clients to use HTTPS.
	HttpsRedirect bool `protobuf:"varint,1,opt,name=https_redirect,json=httpsRedirect" json:"https_redirect,omitempty"`
	// Optional: Indicates whether connections to this port should be
	// secured using TLS.  The value of this field determines how TLS is
	// enforced.
	Mode Server_TLSOptions_TLSmode `protobuf:"varint,2,opt,name=mode,enum=istio.routing.v1alpha2.Server_TLSOptions_TLSmode" json:"mode,omitempty"`
	// REQUIRED if mode == SIMPLE/MUTUAL. The name of the file holding the
	// server-side TLS certificate to use.  It is the responsibility of the
	// underlying platform to mount the certificate as a file under
	// /etc/istio/ingress-certs with the same name as the specified in this
	// field.
	ServerCertificate string `protobuf:"bytes,3,opt,name=server_certificate,json=serverCertificate" json:"server_certificate,omitempty"`
	// REQUIRED if mode == MUTUAL. To use mutual TLS for external clients,
	// specify the name of the file holding the CA certificate to validate
	// the client's certificate. It is the responsibility of the underlying
	// platform to mount the certificate as a file under
	// /etc/istio/ingress-certs with the same name as specified in this
	// field.
	ClientCaBundle string `protobuf:"bytes,4,opt,name=client_ca_bundle,json=clientCaBundle" json:"client_ca_bundle,omitempty"`
}

func (m *Server_TLSOptions) Reset()                    { *m = Server_TLSOptions{} }
func (m *Server_TLSOptions) String() string            { return proto.CompactTextString(m) }
func (*Server_TLSOptions) ProtoMessage()               {}
func (*Server_TLSOptions) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 1} }

func (m *Server_TLSOptions) GetHttpsRedirect() bool {
	if m != nil {
		return m.HttpsRedirect
	}
	return false
}

func (m *Server_TLSOptions) GetMode() Server_TLSOptions_TLSmode {
	if m != nil {
		return m.Mode
	}
	return Server_TLSOptions_PASSTHROUGH
}

func (m *Server_TLSOptions) GetServerCertificate() string {
	if m != nil {
		return m.ServerCertificate
	}
	return ""
}

func (m *Server_TLSOptions) GetClientCaBundle() string {
	if m != nil {
		return m.ClientCaBundle
	}
	return ""
}

func init() {
	proto.RegisterType((*Gateway)(nil), "istio.routing.v1alpha2.Gateway")
	proto.RegisterType((*Server)(nil), "istio.routing.v1alpha2.Server")
	proto.RegisterType((*Server_Port)(nil), "istio.routing.v1alpha2.Server.Port")
	proto.RegisterType((*Server_TLSOptions)(nil), "istio.routing.v1alpha2.Server.TLSOptions")
	proto.RegisterEnum("istio.routing.v1alpha2.Server_Location", Server_Location_name, Server_Location_value)
	proto.RegisterEnum("istio.routing.v1alpha2.Server_TLSOptions_TLSmode", Server_TLSOptions_TLSmode_name, Server_TLSOptions_TLSmode_value)
}

func init() { proto.RegisterFile("routing/v1alpha2/gateway.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 439 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x51, 0x4d, 0x6f, 0x9b, 0x40,
	0x14, 0x2c, 0x36, 0xb2, 0xe1, 0xa5, 0x76, 0xe9, 0x3b, 0x44, 0xc8, 0x87, 0x08, 0x51, 0x55, 0xa5,
	0x87, 0x12, 0x85, 0x1e, 0x5a, 0xa9, 0x27, 0x6a, 0x21, 0xc7, 0x12, 0xfe, 0xd0, 0x82, 0xa5, 0xde,
	0xac, 0x35, 0xde, 0x26, 0x48, 0x98, 0x45, 0xcb, 0x3a, 0x55, 0xff, 0x48, 0x7f, 0x64, 0x7f, 0x45,
	0xc5, 0x02, 0xc9, 0xa5, 0x6a, 0x7a, 0x7b, 0x6f, 0x76, 0x66, 0xde, 0x30, 0xc0, 0x95, 0xe0, 0x67,
	0x99, 0x97, 0x77, 0xd7, 0x0f, 0x37, 0xb4, 0xa8, 0xee, 0x69, 0x70, 0x7d, 0x47, 0x25, 0xfb, 0x41,
	0x7f, 0xfa, 0x95, 0xe0, 0x92, 0xe3, 0x65, 0x5e, 0xcb, 0x9c, 0xfb, 0x1d, 0xcb, 0xef, 0x59, 0xee,
	0x1c, 0xc6, 0x8b, 0x96, 0x88, 0x9f, 0x61, 0x5c, 0x33, 0xf1, 0xc0, 0x44, 0x6d, 0x6b, 0xce, 0xd0,
	0xbb, 0x08, 0xae, 0xfc, 0xbf, 0x8b, 0xfc, 0x44, 0xd1, 0x48, 0x4f, 0x77, 0x7f, 0xeb, 0x30, 0x6a,
	0x31, 0xfc, 0x04, 0x7a, 0xc5, 0x85, 0xb4, 0x35, 0x47, 0xf3, 0x2e, 0x82, 0x37, 0xff, 0x76, 0xf0,
	0xb7, 0x5c, 0x48, 0xa2, 0x04, 0x68, 0xc3, 0xf8, 0xc8, 0x4f, 0x34, 0x2f, 0x6b, 0x7b, 0xe0, 0x0c,
	0x3d, 0x93, 0xf4, 0x2b, 0x7e, 0x81, 0xa1, 0x2c, 0x6a, 0x7b, 0xa8, 0x1c, 0xdf, 0x3f, 0xe3, 0x98,
	0xc6, 0xc9, 0xa6, 0x92, 0x39, 0x2f, 0x6b, 0xd2, 0xa8, 0x70, 0x0e, 0x46, 0xc1, 0x33, 0xda, 0x20,
	0xb6, 0xee, 0x68, 0xde, 0x34, 0x78, 0xf7, 0x8c, 0x43, 0xdc, 0xd1, 0xc9, 0xa3, 0x70, 0xb6, 0x06,
	0xbd, 0x49, 0x8a, 0x97, 0x30, 0x2a, 0xcf, 0xa7, 0x03, 0x13, 0xea, 0xf3, 0x26, 0xa4, 0xdb, 0x70,
	0x06, 0x86, 0x6a, 0x39, 0xe3, 0x85, 0x3d, 0x70, 0x34, 0xcf, 0x24, 0x8f, 0x3b, 0x22, 0xe8, 0x25,
	0x3d, 0x31, 0x15, 0xdf, 0x24, 0x6a, 0x9e, 0xfd, 0x1a, 0x00, 0x3c, 0x05, 0xc5, 0xb7, 0x30, 0xbd,
	0x97, 0xb2, 0xaa, 0xf7, 0x82, 0x1d, 0x73, 0xc1, 0xb2, 0xb6, 0x3d, 0x83, 0x4c, 0x14, 0x4a, 0x3a,
	0x10, 0x23, 0xd0, 0x4f, 0xfc, 0xc8, 0xd4, 0x85, 0x69, 0x70, 0xf3, 0xdf, 0x45, 0x34, 0x63, 0x23,
	0x24, 0x4a, 0x8e, 0x1f, 0x00, 0xdb, 0xff, 0xb6, 0xcf, 0x98, 0x90, 0xf9, 0xf7, 0x3c, 0xa3, 0xb2,
	0x8f, 0xf7, 0xba, 0x7d, 0x99, 0x3f, 0x3d, 0xa0, 0x07, 0x56, 0x56, 0xe4, 0xac, 0x94, 0xfb, 0x8c,
	0xee, 0x0f, 0xe7, 0xf2, 0x58, 0x30, 0x55, 0xa4, 0x49, 0xa6, 0x2d, 0x3e, 0xa7, 0x5f, 0x15, 0xea,
	0x86, 0x30, 0xee, 0x2e, 0xe1, 0x2b, 0xb8, 0xd8, 0x86, 0x49, 0x92, 0xde, 0x92, 0xcd, 0x6e, 0x71,
	0x6b, 0xbd, 0x40, 0x80, 0x51, 0xb2, 0x5c, 0x6d, 0xe3, 0xc8, 0xd2, 0x9a, 0x79, 0xb5, 0x4b, 0x77,
	0x61, 0x6c, 0x0d, 0x70, 0x02, 0xe6, 0x86, 0x2c, 0x17, 0xcb, 0x75, 0x98, 0x46, 0xd6, 0xd0, 0x75,
	0xc1, 0xe8, 0xeb, 0x47, 0x03, 0xf4, 0x55, 0x94, 0x34, 0xe2, 0x97, 0x60, 0x44, 0xdf, 0xd2, 0x88,
	0xac, 0xc3, 0xd8, 0xd2, 0x0e, 0x23, 0x55, 0xed, 0xc7, 0x3f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xbc,
	0xb5, 0xfe, 0xd1, 0xf2, 0x02, 0x00, 0x00,
}