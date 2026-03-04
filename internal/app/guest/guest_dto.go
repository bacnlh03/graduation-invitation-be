package guest

type GuestResponse struct {
	ID        string `json:"id"`
	FullName  string `json:"full_name"`
	Status    int    `json:"status"`
	UpdatedAt string `json:"updated_at"`
}

type FilterGuestsRequest struct {
	Status *int   `json:"status,omitempty" query:"status"`
	Search string `json:"search,omitempty" query:"search"`
}

type CreateGuestRequest struct {
	Name string `json:"name" validate:"required"`
}

type BulkCreateGuestRequest struct {
	Guests []CreateGuestRequest `json:"guests" validate:"required,dive"`
}

type UpdateGuestRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateStatusRequest struct {
	Status int `json:"status" validate:"min=0,max=2"`
}
