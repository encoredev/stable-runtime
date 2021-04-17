package runtime

import (
	"context"
	"log"
	"os"

	runtimepb "encore.dev/internal/proto/encore/engine/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	runtimeAddr string
	procID      string
)

var traceMeta = metadata.Pairs("version", "3")

func RecordTrace(ctx context.Context, traceID [16]byte, data []byte) error {
	ctx = metadata.NewOutgoingContext(ctx, traceMeta)
	_, err := runtime.RecordTrace(ctx, &runtimepb.RecordTraceRequest{
		TraceId: traceID[:],
		Data:    data,
	})
	return err
}

func fetchSecrets(ctx context.Context) (map[string]string, error) {
	secrets, err := runtime.Secrets(ctx, &runtimepb.SecretsRequest{})
	if err != nil {
		return nil, err
	}
	return secrets.Secrets, nil
}

var runtime = func() runtimepb.RuntimeClient {
	const env = "ENCORE_RUNTIME_ADDRESS"
	addr := os.Getenv(env)
	os.Unsetenv(env)
	if addr == "" {
		log.Fatalln("encore: internal error: no runtime address given")
	}
	cc, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalln("encore: internal error: could not dial runtime:", err)
	}
	return runtimepb.NewRuntimeClient(cc)
}()
