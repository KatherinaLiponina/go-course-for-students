package grpc

import (
	context "context"
	"errors"
	"homework9/internal/ads"
	"homework9/internal/app"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type AdUserService struct {
	App app.App
}

func (serv *AdUserService) CreateAd(ctx context.Context, r *CreateAdRequest) (*AdResponse, error) {
	ad, err := serv.App.CreateAd(r.Title, r.Text, r.UserId)
	if err != nil {
		return &AdResponse{}, err
	}
	return &AdResponse{Id: ad.ID, Title: ad.Title,
		Text: ad.Text, AuthorId: ad.AuthorID, Published: ad.Published,
		CreationDate: timestamppb.New(ad.CreationDate), UpdateTime: timestamppb.New(ad.UpdateTime)}, nil
}

func (serv *AdUserService) ChangeAdStatus(ctx context.Context, r *ChangeAdStatusRequest) (*AdResponse, error) {
	ad, err := serv.App.ChangeAdStatus(r.AdId, r.UserId, r.Published)
	if err != nil {
		return &AdResponse{}, err
	}
	return &AdResponse{Id: ad.ID, Title: ad.Title,
		Text: ad.Text, AuthorId: ad.AuthorID, Published: ad.Published,
		CreationDate: timestamppb.New(ad.CreationDate), UpdateTime: timestamppb.New(ad.UpdateTime)}, nil
}

func (serv *AdUserService) UpdateAd(ctx context.Context, r *UpdateAdRequest) (*AdResponse, error) {
	ad, err := serv.App.UpdateAd(r.AdId, r.UserId, r.Title, r.Text)
	if err != nil {
		return &AdResponse{}, err
	}
	return &AdResponse{Id: ad.ID, Title: ad.Title,
		Text: ad.Text, AuthorId: ad.AuthorID, Published: ad.Published,
		CreationDate: timestamppb.New(ad.CreationDate), UpdateTime: timestamppb.New(ad.UpdateTime)}, nil
}

func createListAdResponse(a []ads.Ad) *ListAdResponse {
	var arr []*AdResponse
	for _, ad := range a {
		arr = append(arr, &AdResponse{Id: ad.ID, Title: ad.Title,
			Text: ad.Text, AuthorId: ad.AuthorID, Published: ad.Published,
			CreationDate: timestamppb.New(ad.CreationDate), UpdateTime: timestamppb.New(ad.UpdateTime)})
	}
	return &ListAdResponse{List: arr}
}

func (serv *AdUserService) ListAds(ctx context.Context, m *Mode) (*ListAdResponse, error) {
	var arr []ads.Ad
	mode := ModeType_name[int32(m.Mode)]

	var err error
	if mode == "ByAuthor" {
		data, ok := m.Data.(*Mode_AuthorId)
		if !ok {
			return &ListAdResponse{}, errors.New("wrong parameters")
		}
		arr, err = serv.App.SelectByAuthor(data.AuthorId)
	} else if mode == "ByCreation" {
		data, ok := m.Data.(*Mode_Time)
		if !ok {
			return &ListAdResponse{}, errors.New("wrong parameters")
		}
		arr = serv.App.SelectByCreation(data.Time.AsTime())
	} else if mode == "All" {
		arr = serv.App.SelectAll()
	} else if mode == "ByTitle" {
		data, ok := m.Data.(*Mode_Title)
		if !ok {
			return &ListAdResponse{}, errors.New("wrong parameters")
		}
		arr = serv.App.FindByTitle(data.Title)
	} else {
		arr = serv.App.Select()
	}
	if err != nil {
		return &ListAdResponse{}, err
	}
	return createListAdResponse(arr), nil
}

func (serv *AdUserService) CreateUser(ctx context.Context, r *CreateUserRequest) (*UserResponse, error) {
	usr := serv.App.CreateUser(r.Name, r.Email)
	return &UserResponse{Id: usr.ID, Name: usr.Nickname, Email: usr.Email}, nil
}

func (serv *AdUserService) GetUser(ctx context.Context, r *GetUserRequest) (*UserResponse, error) {
	usr, err := serv.App.GetUserByID(r.Id)
	if err != nil {
		return &UserResponse{}, err
	}
	return &UserResponse{Id: usr.ID, Name: usr.Nickname, Email: usr.Email}, nil
}

func (serv *AdUserService) DeleteUser(ctx context.Context, r *DeleteUserRequest) (*UserResponse, error) {
	usr, err := serv.App.DeleteUser(r.Id)
	if err != nil {
		return &UserResponse{}, err
	}
	return &UserResponse{Id: usr.ID, Name: usr.Nickname, Email: usr.Email}, nil
}

func (serv *AdUserService) DeleteAd(ctx context.Context, r *DeleteAdRequest) (*AdResponse, error) {
	ad, err := serv.App.DeleteAd(r.AdId, r.AuthorId)
	if err != nil {
		return &AdResponse{}, err
	}
	return &AdResponse{Id: ad.ID, Title: ad.Title,
		Text: ad.Text, AuthorId: ad.AuthorID, Published: ad.Published,
		CreationDate: timestamppb.New(ad.CreationDate), UpdateTime: timestamppb.New(ad.UpdateTime)}, nil
}

func (serv *AdUserService) mustEmbedUnimplementedAdServiceServer() {}
