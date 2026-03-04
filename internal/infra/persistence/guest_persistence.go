package persistence

import (
	"context"
	"errors"
	"graduation-invitation/internal/domain/guest"

	"github.com/jmoiron/sqlx"
)

type guestPersistence struct {
	db *sqlx.DB
}

func NewGuestPersistence(db *sqlx.DB) guest.GuestRepo {
	return &guestPersistence{db: db}
}

func (p *guestPersistence) CreateMany(ctx context.Context, guests []*guest.Guest) error {
	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO guests (id, first_name, last_name, status, created_at, updated_at)
		VALUES (:id, :first_name, :last_name, 0, NOW(), NOW())`

	rows, err := p.db.NamedQueryContext(ctx, query, guests)
	if err != nil {
		return err
	}
	defer rows.Close()

	return tx.Commit()
}

func (p *guestPersistence) FindMany(ctx context.Context, search string, status *int) ([]*guest.Guest, error) {
	guests := []*guest.Guest{}

	// Khởi tạo query cơ bản
	query := `SELECT * FROM guests WHERE 1=1`
	args := make(map[string]any)

	if search != "" {
		// Lọc cả FirstName hoặc LastName
		query += ` AND (first_name ILIKE :search OR last_name ILIKE :search)`
		args["search"] = "%" + search + "%"
	}

	// Kiểm tra nếu status không phải nil thì mới thêm điều kiện lọc
	if status != nil {
		query += ` AND status = :status`
		args["status"] = *status
	}

	query += ` ORDER BY first_name ASC, last_name ASC`

	// Sử dụng NamedQuery để map args an toàn
	rows, err := p.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		g := new(guest.Guest)
		if err := rows.StructScan(g); err != nil {
			return nil, err
		}
		guests = append(guests, g)
	}

	return guests, nil
}

func (p *guestPersistence) FindByID(ctx context.Context, id string) (*guest.Guest, error) {
	var guest guest.Guest
	query := `SELECT * FROM guests WHERE id = $1`

	err := p.db.GetContext(ctx, &guest, query, id)
	if err != nil {
		return nil, err
	}
	return &guest, nil
}

// Update: Cập nhật thông tin khách mời (ví dụ cập nhật trạng thái RSVP)
func (p *guestPersistence) Update(ctx context.Context, guest *guest.Guest) error {
	query := `
		UPDATE guests 
		SET first_name = :first_name,
		    last_name = :last_name,
		    status = :status,
		    updated_at = NOW()
		WHERE id = :id`

	result, err := p.db.NamedExecContext(ctx, query, guest)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("không có bản ghi nào được cập nhật")
	}
	return nil
}

func (p *guestPersistence) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM guests WHERE id = $1`
	_, err := p.db.ExecContext(ctx, query, id)
	return err
}
