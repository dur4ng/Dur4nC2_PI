package docker_test

import (
	"Dur4nC2/server/domain/models"
	_operatorRepository "Dur4nC2/server/domain/operator/repository/postgres"
)

var _ = Describe("Repository", func() {
	var repo models.OperatorRepository

	BeforeEach(func() {
		err := Db.AutoMigrate(
			&models.Operator{},
			&models.Host{},
			&models.Beacon{},
			&models.BeaconTask{},
			&models.IOC{},
			&models.Loot{},
		)
		Ω(err).To(Succeed())

		repo = _operatorRepository.NewPostgresOperatorRepository(Db)
	})
	Context("Create", func() {
		It("Success", func() {
			operatorMock := &models.Operator{Username: "bob"}
			err := repo.Create(operatorMock)
			Ω(err).To(Succeed())
			operators, err := repo.List()
			Ω(err).To(Succeed())
			Ω(operators).To(HaveLen(1))
		})
	})
})
