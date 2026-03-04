package guest

import "context"

type GuestRepo interface {
	FindMany(ctx context.Context, search string, status *int) ([]*Guest, error)
	FindByID(ctx context.Context, id string) (*Guest, error)
	CreateMany(ctx context.Context, guests []*Guest) error
	Update(ctx context.Context, guest *Guest) error
	Delete(ctx context.Context, id string) error
}
