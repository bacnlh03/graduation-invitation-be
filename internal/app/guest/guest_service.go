package guest

import (
	"context"
	"graduation-invitation/internal/domain/guest"
	"strings"

	"github.com/google/uuid"
)

type GuestService interface {
	BulkRegister(ctx context.Context, req BulkCreateGuestRequest) error
	UpdateGuest(ctx context.Context, id string, req UpdateGuestRequest) error
	UpdateStatus(ctx context.Context, id string, status int) error
	ConfirmAttendance(ctx context.Context, id string) error
	DeleteGuest(ctx context.Context, id string) error
	ListGuests(ctx context.Context, filter FilterGuestsRequest) ([]*GuestResponse, error)
	GetByID(ctx context.Context, id string) (*GuestResponse, error)
}

type guestService struct {
	repo guest.GuestRepo
}

func NewGuestService(repo guest.GuestRepo) GuestService {
	return &guestService{repo: repo}
}

func (s *guestService) BulkRegister(ctx context.Context, req BulkCreateGuestRequest) error {
	var domainGuests []*guest.Guest

	for _, r := range req.Guests {
		full := strings.TrimSpace(r.Name)
		firstName := ""
		lastName := ""

		lastSpaceIndex := strings.LastIndex(full, " ")

		if lastSpaceIndex == -1 {
			firstName = full
			lastName = ""
		} else {
			lastName = strings.TrimSpace(full[:lastSpaceIndex])
			firstName = strings.TrimSpace(full[lastSpaceIndex+1:])
		}

		domainGuests = append(domainGuests, &guest.Guest{
			ID:        uuid.Must(uuid.NewV7()).String(),
			FirstName: firstName,
			LastName:  lastName,
		})
	}

	return s.repo.CreateMany(ctx, domainGuests)
}

func (s *guestService) UpdateGuest(ctx context.Context, id string, req UpdateGuestRequest) error {
	guest, err := s.repo.FindByID(ctx, id)
	if err != nil || guest == nil {
		return err
	}

	full := strings.TrimSpace(req.Name)
	firstName := ""
	lastName := ""

	lastSpaceIndex := strings.LastIndex(full, " ")

	if lastSpaceIndex == -1 {
		firstName = full
		lastName = ""
	} else {
		lastName = strings.TrimSpace(full[:lastSpaceIndex])
		firstName = strings.TrimSpace(full[lastSpaceIndex+1:])
	}

	guest.FirstName = firstName
	guest.LastName = lastName

	return s.repo.Update(ctx, guest)
}

func (s *guestService) UpdateStatus(ctx context.Context, id string, status int) error {
	guest, err := s.repo.FindByID(ctx, id)
	if err != nil || guest == nil {
		return err
	}
	guest.Status = status
	return s.repo.Update(ctx, guest)
}

func (s *guestService) ConfirmAttendance(ctx context.Context, id string) error {
	// 1. Tìm khách bằng ID
	guest, err := s.repo.FindByID(ctx, id)
	if err != nil || guest == nil {
		return err
	}
	// 2. Cập nhật trạng thái -> Xác nhận tham dự
	guest.Status = 2 // StatusConfirmed
	// 3. Lưu lại
	return s.repo.Update(ctx, guest)
}

func (s *guestService) DeleteGuest(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *guestService) ListGuests(ctx context.Context, filter FilterGuestsRequest) ([]*GuestResponse, error) {
	guests, err := s.repo.FindMany(ctx, filter.Search, filter.Status)
	if err != nil {
		return nil, err
	}

	var res []*GuestResponse
	for _, g := range guests {
		res = append(res, &GuestResponse{
			ID:        g.ID,
			FullName:  g.FullName(),
			Status:    g.Status,
			UpdatedAt: g.UpdatedAt,
		})
	}

	return res, nil
}
func (s *guestService) GetByID(ctx context.Context, id string) (*GuestResponse, error) {
	g, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, nil
	}

	return &GuestResponse{
		ID:        g.ID,
		FullName:  g.FullName(),
		Status:    g.Status,
		UpdatedAt: g.UpdatedAt,
	}, nil
}
