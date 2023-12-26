package docker_test

import (
	"Dur4nC2/misc/protobuf/rpcpb"
	"Dur4nC2/server/rpc/docker"
	"fmt"
	"net"
	"runtime/debug"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	"Dur4nC2/server/db"
)

const (
	dbName = "test"
	passwd = "test"
)

var (
	kb                          = 1024
	mb                          = kb * 1024
	gb                          = mb * 1024
	bufSize                     = 2 * mb
	ServerMaxMessageSize        = 2 * gb
	ClientMaxReceiveMessageSize = 2 * gb
	host                        = "localhost"
	port                        = 7878
)

func TestDocker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Docker Suite")
}

var Db *gorm.DB
var cleanupDocker func()
var RPCClient rpcpb.TeamServerRPCClient

var _ = BeforeSuite(func() {
	// setup *gorm.Db with docker
	Db, cleanupDocker = db.SetupGormWithDocker()
	_, _, err := localListener(host, port)
	if err != nil {
		panic(err)
	}
	RPCClient, _ = localClient()

})

var _ = AfterSuite(func() {
	// cleanup resource
	cleanupDocker()
})

var _ = BeforeEach(func() {
	// clear db tables before each test
	err := Db.Exec(`DROP SCHEMA public CASCADE;CREATE SCHEMA public;`).Error
	Ω(err).To(Succeed())
})

func localClient() (rpcpb.TeamServerRPCClient, error) {
	fmt.Printf("Local client ...\n")

	//ctxDialer := grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
	//	return ln.Dial()
	//})
	//
	//options := []grpc.DialOption{
	//	ctxDialer,
	//	grpc.WithInsecure(), // This is an in-memory listener, no need for secure transport
	//	grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(ClientMaxReceiveMessageSize)),
	//}
	//conn, err := grpc.DialContext(context.Background(), "bufnet", options...)
	//if err != nil {
	//	fmt.Printf("Failed to dial bufnet: %s\n", err)
	//	return nil, err
	//}
	//defer conn.Close()
	//options := []grpc.DialOption{
	//	grpc.WithInsecure(), // This is an in-memory listener, no need for secure transport
	//	grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(ClientMaxReceiveMessageSize)),
	//}
	conn, err := grpc.Dial("localhost:7878", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	//defer conn.Close() grpc: the client connection is closing

	localRPC := rpcpb.NewTeamServerRPCClient(conn)
	return localRPC, nil
}

// LocalListener - Bind gRPC server to an in-memory listener, which is
//                 typically used for unit testing, but ... it should be fine
func localListener(host string, port int) (*grpc.Server, net.Listener, error) {
	fmt.Printf("Binding gRPC to listener ...\n")
	//ln := bufconn.Listen(bufSize) //in-memory listener returns -> *bufconn.Listener
	// normal tcp listener returns -> net.Listener
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port)) //
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	options := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(ServerMaxMessageSize),
		grpc.MaxSendMsgSize(ServerMaxMessageSize),
	}
	grpcServer := grpc.NewServer(options...)
	rpcpb.RegisterTeamServerRPCServer(grpcServer, docker.NewServer(Db))
	go func() {
		panicked := true
		defer func() {
			if panicked {
				fmt.Printf("stacktrace from panic: %s", string(debug.Stack()))
			}
		}()
		if err := grpcServer.Serve(ln); err != nil {
			fmt.Printf("gRPC local listener error: %v", err)
		} else {
			panicked = false
		}
	}()
	return grpcServer, ln, nil
}
