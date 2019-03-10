package service

import (
	"fmt"
	"net"

	log "github.com/cohix/simplog"
	"github.com/pkg/errors"
	"github.com/taask/taask-server/auth"
	"github.com/taask/taask-server/brain"
	"github.com/taask/taask-server/model"
	"github.com/taask/taask-server/partner"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
)

const partnerServicePort = ":3690"

// StartPartnerService starts the partner service
func StartPartnerService(brain *brain.Manager, errChan chan error) {
	lis, err := net.Listen("tcp", partnerServicePort)
	if err != nil {
		errChan <- err
		return
	}

	grpcServer := grpc.NewServer()

	partner.RegisterPartnerServiceServer(grpcServer, &PartnerService{Manager: brain})

	log.LogInfo("starting taask-server partner service on :3690")
	if err := grpcServer.Serve(lis); err != nil {
		errChan <- err
	}
}

// PartnerService describes the service available to taask managing partners
type PartnerService struct {
	Manager *brain.Manager
}

// AuthPartner allows a partner to authenticate and get a session
func (ps *PartnerService) AuthPartner(ctx context.Context, attempt *auth.Attempt) (*auth.AttemptResponse, error) {
	defer log.LogTrace("AuthPartner")()

	log.LogInfo("received partner auth request")
	encPartnerChallenge, err := ps.Manager.AttemptPartnerAuth(attempt)
	if err != nil {
		log.LogWarn(errors.Wrap(err, "partner auth request denied").Error())
		return nil, err
	}

	log.LogInfo("partner auth request accepted")
	resp := &auth.AttemptResponse{
		EncChallenge: encPartnerChallenge.EncSessionChallenge,
		MasterPubKey: ps.Manager.GetMasterPartnerPubKey(),
	}

	return resp, nil
}

// StreamUpdates enables sync between partners
func (ps *PartnerService) StreamUpdates(stream partner.PartnerService_StreamUpdatesServer) error {
	if err := ps.Manager.RunPartnerManagerWithServer(stream); err != nil {
		return errors.Wrap(err, "StreamUpdates failed to StartWithServer")
	}

	return nil
}

func (ps *PartnerService) checkTaskVersion(update *model.TaskUpdate) error {
	task, err := ps.Manager.GetTask(update.UUID)
	if err != nil {
		return errors.Wrap(err, "failed to storage.Get")
	}

	if update.Version != task.Meta.Version+1 {
		return fmt.Errorf("runner tried to apply update with version %d to task %s with version %d", update.Version, task.UUID, task.Meta.Version)
	}

	return nil
}
