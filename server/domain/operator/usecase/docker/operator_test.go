package docker_test

import (
	"Dur4nC2/server/domain/models"
	_operatorRepository "Dur4nC2/server/domain/operator/repository/postgres"
	_operatorUsecase "Dur4nC2/server/domain/operator/usecase/docker"
)

var _ = Describe("Repository", func() {
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

		repository = _operatorRepository.NewPostgresOperatorRepository(Db)
		usecase = _operatorUsecase.NewOperatorUsecase(repository)
	})
	Context("List", func() {
		It("Success", func() {
			operatorMock := &models.Operator{Username: "bob"}
			err := repository.Create(operatorMock)
			Ω(err).To(Succeed())
			operators, err := repository.List()
			Ω(err).To(Succeed())
			Ω(operators).To(HaveLen(1))

			operators_pb, err := usecase.List()
			Ω(err).To(Succeed())
			Ω(operators_pb.Operators).To(HaveLen(1))
		})
	})
})
