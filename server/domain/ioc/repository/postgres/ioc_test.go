package postgres

import (
	_hostRepo "Dur4nC2/server/domain/host/repository/postgres"
	"Dur4nC2/server/domain/models"
)

var _ = Describe("Repository", func() {
	var iocRepo models.IOCRepository
	var hosts []models.Host
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

		// *** Host ***
		h := _hostRepo.NewPostgresHostRepository(Db)
		err = h.Create(&models.Host{
			Hostname:  "DC-01",
			OSVersion: "Windows Server 2019",
			Locale:    "USA",
		})
		Ω(err).To(Succeed())
		hosts, err = h.List()
		Ω(hosts).To(HaveLen(1))

		iocRepo = NewPostgresIOCRepository(Db)

	})
	Context("Create", func() {
		It("Success", func() {
			iocMock := models.IOC{
				HostID:      hosts[0].ID,
				Path:        "C:\\Temp\\implant.exe",
				FileHash:    "asdfasdfasdfsadsf",
				Name:        "Loader",
				Description: "Initial loader",
				State:       models.OPEN,
			}
			err := iocRepo.Create(iocMock)
			Ω(err).To(Succeed())
			operators, err := iocRepo.List()
			Ω(err).To(Succeed())
			Ω(operators).To(HaveLen(1))
		})
	})
})
