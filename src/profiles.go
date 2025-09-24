package alexchatapp

import (
	"alexchatapp/src/data"
	pb "alexchatapp/src/proto/profiles"
	"alexchatapp/src/utils"
	"context"
	"errors"
)

type ProfileServer struct {
	pb.UnimplementedProfileServiceServer
	profile_repo *data.ProfilesRepository
}

func NewProfilesServer(profile_repo *data.ProfilesRepository) *ProfileServer {
	return &ProfileServer{
		profile_repo: profile_repo,
	}
}

func (p *ProfileServer) RegisterProfile(ctx *context.Context, req *pb.CreateProfileRequest) (*pb.CreateProfileResponse, error) {
	if err := utils.ValidateProfileName(req.ProfileName); err != nil {
		return &pb.CreateProfileResponse{
			StatusCode: 400,
		}, err
	}

	if isExists := p.profile_repo.DoesProfileExist(uint(req.UserId)); isExists {
		return &pb.CreateProfileResponse{
			StatusCode: 400,
		}, errors.New("the user already has profile")
	}

	err := p.profile_repo.CreateProfile(uint(req.UserId), req.ProfileName)
	if err != nil {
		return &pb.CreateProfileResponse{
			StatusCode: 400,
		}, err
	}

	return &pb.CreateProfileResponse{
		StatusCode: 200,
	}, nil
}
