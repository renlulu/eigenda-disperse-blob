## EigenDA Disperse Blob Function

### Why

To enable payment mechanism, EigenDa introduces a new function type for submitting blob to the network: DisperseBlobAuthenticated, see following interface definition (as well as the comments/explanation)
```go
// DisperserClient is the client API for Disperser service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DisperserClient interface {
	// This API accepts blob to disperse from clients.
	// This executes the dispersal async, i.e. it returns once the request
	// is accepted. The client could use GetBlobStatus() API to poll the the
	// processing status of the blob.
	DisperseBlob(ctx context.Context, in *DisperseBlobRequest, opts ...grpc.CallOption) (*DisperseBlobReply, error)
	// DisperseBlobAuthenticated is similar to DisperseBlob, except that it requires the
	// client to authenticate itself via the AuthenticationData message. The protoco is as follows:
	//  1. The client sends a DisperseBlobAuthenticated request with the DisperseBlobRequest message
	//  2. The Disperser sends back a BlobAuthHeader message containing information for the client to
	//     verify and sign.
	//  3. The client verifies the BlobAuthHeader and sends back the signed BlobAuthHeader in an
	//     AuthenticationData message.
	//  4. The Disperser verifies the signature and returns a DisperseBlobReply message.
	DisperseBlobAuthenticated(ctx context.Context, opts ...grpc.CallOption) (Disperser_DisperseBlobAuthenticatedClient, error)
	// This API is meant to be polled for the blob status.
	GetBlobStatus(ctx context.Context, in *BlobStatusRequest, opts ...grpc.CallOption) (*BlobStatusReply, error)
	// This retrieves the requested blob from the Disperser's backend.
	// This is a more efficient way to retrieve blobs than directly retrieving
	// from the DA Nodes (see detail about this approach in
	// api/proto/retriever/retriever.proto).
	// The blob should have been initially dispersed via this Disperser service
	// for this API to work.
	RetrieveBlob(ctx context.Context, in *RetrieveBlobRequest, opts ...grpc.CallOption) (*RetrieveBlobReply, error)
}
```

It is basically using a simple two steps protocol to identify the sender address, thus EigenDA can track users' TPS and throughput more easily, as well as bandwidth payments. This repo mainly mostly for 
showing how the authentication process works, or anyone want get rid of the tedious code can use the function straightaway, although it is quite SIMPLE.