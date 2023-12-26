package usecase

import (
	"Dur4nC2/misc/protobuf/clientpb"
	"Dur4nC2/misc/protobuf/commonpb"
	"Dur4nC2/server/db"
	hostRepository "Dur4nC2/server/domain/host/repository/postgres"
	"Dur4nC2/server/domain/models"
	"github.com/gofrs/uuid"
)

type hostUsecase struct {
	repository models.HostRepository
}

func NewHostUsecase() models.HostUsecase {
	r := hostRepository.NewPostgresHostRepository(db.Session())
	return &hostUsecase{repository: r}
}

func (u *hostUsecase) List() (*clientpb.Hosts, error) {
	hosts, err := u.repository.List()
	if err != nil {
		return &clientpb.Hosts{}, nil
	}
	hosts_pb := &clientpb.Hosts{Hosts: []*clientpb.Host{}}
	for _, host := range hosts {
		var iocs_pb []*clientpb.IOC
		for _, ioc := range host.IOCs {
			iocs_pb = append(iocs_pb, &clientpb.IOC{
				ID:       ioc.ID.String(),
				Path:     ioc.Path,
				FileHash: ioc.FileHash,
			})
		}
		hosts_pb.Hosts = append(hosts_pb.Hosts, &clientpb.Host{
			HostUUID:     host.ID.String(),
			Hostname:     host.Hostname,
			OSVersion:    host.OSVersion,
			IOCs:         iocs_pb,
			Locale:       host.Locale,
			FirstContact: host.CreatedAt.Unix(),
		})
	}
	return hosts_pb, nil
}
func (u *hostUsecase) Get(host_pb *clientpb.Host) (*clientpb.Host, error) {
	id, _ := uuid.FromString(host_pb.HostUUID)
	host, err := u.repository.Read(id)
	if err != nil {
		return nil, err
	}
	var iocs_pb []*clientpb.IOC
	for _, ioc := range host.IOCs {
		iocs_pb = append(iocs_pb, &clientpb.IOC{
			ID:       ioc.ID.String(),
			Path:     ioc.Path,
			FileHash: ioc.FileHash,
		})
	}
	host_pb = &clientpb.Host{
		HostUUID:     host.ID.String(),
		Hostname:     host.Hostname,
		OSVersion:    host.OSVersion,
		IOCs:         iocs_pb,
		Locale:       host.Locale,
		FirstContact: host.CreatedAt.Unix(),
	}
	return host_pb, nil
}
func (u *hostUsecase) Delete(host_pb *clientpb.Host) (*commonpb.Empty, error) {
	id, _ := uuid.FromString(host_pb.HostUUID)
	err := u.repository.Delete(id)
	if err != nil {
		return &commonpb.Empty{}, err
	}
	return &commonpb.Empty{}, nil
}

func (u *hostUsecase) ListHostIOC(host_pb clientpb.Host) ([]*clientpb.IOC, error) {
	hostID, _ := uuid.FromString(host_pb.HostUUID)
	host := &models.Host{
		ID:       hostID,
		Hostname: host_pb.Hostname,
	}
	iocs, err := u.repository.ListHostIOC(*host)
	if err != nil {
		return nil, err
	}
	var iocs_pb []*clientpb.IOC
	for _, ioc := range iocs {
		iocs_pb = append(iocs_pb, &clientpb.IOC{
			ID:       ioc.ID.String(),
			FileHash: ioc.FileHash,
			Path:     ioc.Path,
		})
	}
	return iocs_pb, nil
}
func (u *hostUsecase) ListHostLoot(loot_pb *clientpb.Loot) (*clientpb.Loots, error) {
	return nil, nil
}
