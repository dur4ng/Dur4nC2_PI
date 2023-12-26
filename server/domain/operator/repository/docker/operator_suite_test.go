package docker_test

import (
	"Dur4nC2/server/db"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

func TestDocker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Docker Suite")
}

var Db *gorm.DB
var cleanupDocker func()

var _ = BeforeSuite(func() {
	// setup *gorm.Db with docker
	Db, cleanupDocker = db.SetupGormWithDocker()
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
