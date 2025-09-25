package alexchatapp

import (
	"alexchatapp/src/data"
	"alexchatapp/src/models"
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

func (p *ProfileServer) RegisterProfile(ctx context.Context, req *pb.CreateProfileRequest) (*pb.CreateProfileResponse, error) {
	var new_profile models.Profile
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

	new_profile = models.Profile{
		User_id:      uint(req.UserId),
		Profile_name: req.ProfileName,
		Bio:          *req.Bio,
		Avatar_url:   *req.AvatarUrl,
		Status:       *req.Status,
	}

	err := p.profile_repo.CreateProfileByModel(new_profile)
	if err != nil {
		return &pb.CreateProfileResponse{
			StatusCode: 400,
		}, err
	}

	return &pb.CreateProfileResponse{
		StatusCode: 200,
	}, nil
}

func (p *ProfileServer) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	var updated_profile models.Profile

	if err := utils.ValidateProfileName(*req.ProfileName); err != nil {
		return &pb.UpdateProfileResponse{
			StatusCode: 400,
		}, err
	}

	var profile, err = p.profile_repo.GetProfileByID(uint(req.UserId))
	if err != nil {
		return &pb.UpdateProfileResponse{
			StatusCode: 400,
		}, err
	}

	updated_profile = models.Profile{
		User_id:      profile.User_id,
		Profile_name: *req.ProfileName,
		Bio:          *req.Bio,
		Avatar_url:   *req.AvatarUrl,
		Status:       *req.Status,
	}

	if err := p.profile_repo.UpdateProfile(&updated_profile); err != nil {
		return &pb.UpdateProfileResponse{
			StatusCode: 400,
		}, err
	}

	return &pb.UpdateProfileResponse{
		StatusCode: 200,
	}, nil
}

func (p *ProfileServer) GetProfileByID(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	var profile *pb.Profile

	var profileModel, err = p.profile_repo.GetProfileByID(uint(req.UserId))
	if err != nil {
		return &pb.GetProfileResponse{
			Profile: nil,
		}, err
	}

	profile = &pb.Profile{
		UserId:      uint64(profileModel.User_id),
		ProfileName: profileModel.Profile_name,
		Bio:         &profileModel.Bio,
		AvatarUrl:   &profileModel.Avatar_url,
		Status:      &profileModel.Status,
	}

	return &pb.GetProfileResponse{
		Profile: profile,
	}, nil
}

func (p *ProfileServer) UpdateOnlineStatus(ctx context.Context, req *pb.UpdateOnlineStatusRequest) (pb.UpdateOnlineStatusResponse, error) {
	var profile, err = p.profile_repo.GetProfileByID(uint(req.UserId))
	if err != nil {
		return pb.UpdateOnlineStatusResponse{
			StatusCode: 400,
		}, err
	}

	profile.Last_seen = req.LastSeen.AsTime()

	return pb.UpdateOnlineStatusResponse{
		StatusCode: 200,
	}, nil
}
