package modele

import "time"

// Voucher is the modele for vouchers. It has a proprietary and an expiration date.
type Voucher struct {
	Id         int64
	Code       string
	Expiration time.Time
	Prop       int64
}
