package main

import (
	_implantCrypto "Dur4nC2/implant/cryptography"
	_serverCrypto "Dur4nC2/misc/crypto"
	"Dur4nC2/server/db"
	_implantRepository "Dur4nC2/server/domain/implant/respository/postgres"
	"Dur4nC2/server/listener"
	"bufio"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
	"time"
)

var _ = Describe("Implant manual", func() {
	BeforeEach(func() {

	})
	Context("Manual connections", func() {
		It("Testing comipled implant", func() {
			//Create server
			server_keypair := _serverCrypto.ECCServerKeyPair()
			implant_keypair, err := _implantRepository.NewPosgresImplantRepository(db.Session()).Create()
			//fmt.Println("digest creaded = ", implant_keypair.ECCPublicKeyDigest)
			_implantCrypto.SetSecrets(implant_keypair.ECCPublicKey, implant_keypair.ECCPrivateKey, implant_keypair.ECCPublicKeyDigest, server_keypair.PublicBase64(), "", "", server_keypair.PrivateBase64())
			server, err := listener.StartHTTPListener(&listener.HTTPListenerConfig{
				Addr:   "127.0.0.1:8000",
				Domain: "127.0.0.1",
			})
			if err != nil {
				fmt.Println("Listener failed to start ", err)
				return
			}
			Ω(err).To(Succeed())
			go server.HTTPServer.ListenAndServe()
			if err != nil {
				fmt.Println("Listener failed to serve ", err)
				return
			}
			time.Sleep(1 * time.Second)

			fmt.Println("Waiting to launch and implant manually... Press enter")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
		})
	})
})
