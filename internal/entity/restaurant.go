package entity

import "time"

type Restaurant struct {
	ID                 int64     `pg:"id,pk"`
	UserID             int64     `pg:"user_id"`
	RestaurantName     string    `pg:"restaurant_name,notnull"`
	RestaurantLogo     string    `pg:"restaurant_logo,notnull"`
	RestaurantFavicon  string    `pg:"restaurant_favicon"`
	ThumbnailDesktop   string    `pg:"thumbnail_desktop,notnull"`
	RestaurantPhone    string    `pg:"restaurant_phone,default:''"`
	RestaurantWhatsapp string    `pg:"restaurant_whatsapp,default:''"`
	RestaurantEmail    string    `pg:"restaurant_email,default:''"`
	RestaurantAddress  string    `pg:"restaurant_address"`
	RestaurantWebsite  string    `pg:"restaurant_website"`
	CreatedAt          time.Time `pg:"created_at,notnull"`
	UpdatedAt          time.Time `pg:"updated_at,notnull"`
	User               User      `pg:"rel:has-one,join:user_id"`
}
