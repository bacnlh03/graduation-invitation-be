package guest

// Status constants
const (
	StatusNotSent   = 0 // Chưa gửi
	StatusSent      = 1 // Đã gửi
	StatusConfirmed = 2 // Xác nhận tham dự
)

type Guest struct {
	ID        string `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Status    int    `db:"status"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

func (g *Guest) FullName() string {
	if g.LastName == "" {
		return g.FirstName
	}
	return g.LastName + " " + g.FirstName
}
