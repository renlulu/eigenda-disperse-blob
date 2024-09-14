## EigenDA Disperse Blob Function

### Why

To enable payment mechanism, EigenDa introduces a new function type for submitting blob to the network: DisperseBlobAuthenticated, see following interface definition (as well as the comments/explanation)
```go
type DisperserClient interface {
	// DisperseBlobAuthenticated is similar to DisperseBlob, except that it requires the
	// client to authenticate itself via the AuthenticationData message. The protoco is as follows:
	//  1. The client sends a DisperseBlobAuthenticated request with the DisperseBlobRequest message
	//  2. The Disperser sends back a BlobAuthHeader message containing information for the client to
	//     verify and sign.
	//  3. The client verifies the BlobAuthHeader and sends back the signed BlobAuthHeader in an
	//     AuthenticationData message.
	//  4. The Disperser verifies the signature and returns a DisperseBlobReply message.
	DisperseBlobAuthenticated(ctx context.Context, opts ...grpc.CallOption) (Disperser_DisperseBlobAuthenticatedClient, error)
```

It is basically using a simple two steps protocol to identify the sender address, thus EigenDA can track users' TPS and throughput more easily, as well as bandwidth payments. This repo mainly mostly for 
showing how the authentication process works, or anyone want get rid of the tedious code can use the function straightaway, although it is quite SIMPLE.