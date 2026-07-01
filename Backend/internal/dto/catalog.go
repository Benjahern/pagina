package dto

import "time"

// ItemCreateRequest is the body for POST /admin/products (admin-only).
//
// Money fields (Price, Cost) are int64 representing CLP. CLP has no
// decimals so storing cents is unnecessary — the numeric(12,2) in DBML
// is safety margin for future currencies, not a precision requirement.
//
// Slug is optional here; the service auto-generates it from Name when
// empty (per Item.BeforeCreate in models/item.go).
//
// Cost appears here because this is an admin endpoint. It must never
// appear in ItemResponse — that's the leak boundary.
type ItemCreateRequest struct {
	SKU             string  `json:"sku" binding:"required,min=1,max=100"`
	Name            string  `json:"name" binding:"required,min=1,max=200"`
	Slug            string  `json:"slug,omitempty" binding:"omitempty,max=500"`
	Description     string  `json:"description,omitempty"`
	Price           int64   `json:"price" binding:"required,gte=0"`
	Cost            int64   `json:"cost" binding:"required,gte=0"`
	Stock           int     `json:"stock"`
	Backorder       bool    `json:"backorder"`
	Status          string  `json:"status" binding:"required,oneof=activo inactivo archivado"`
	CategoryID      uint64  `json:"category_id" binding:"required"`
	Brand           string  `json:"brand,omitempty" binding:"omitempty,max=200"`
	Color           string  `json:"color,omitempty" binding:"omitempty,max=100"`
	ImageURL        string  `json:"image_url,omitempty" binding:"omitempty,url"`
	Items3DID       *uint64 `json:"items3d_id,omitempty"`
	MetaTitle       string  `json:"meta_title,omitempty" binding:"omitempty,max=200"`
	MetaDescription string  `json:"meta_description,omitempty"`
	IsFeatured      bool    `json:"is_featured"`
}

// ItemUpdateRequest is the body for PUT /admin/products/:id.
//
// All fields are pointers so partial updates work: a nil field means
// "don't touch"; a non-nil pointer means "set to this value (even if
// it's the zero value)". The service distinguishes missing from
// explicit-empty for nullable strings (Description, ImageURL, etc.)
// — pointer to "" maps to NULL in DB.
//
// Known limitation: clearing a nullable FK (Items3DID) requires the
// same "pointer to nil" pattern, which is indistinguishable from
// "field absent". Defer a solution (sentinel value, JSON Merge Patch,
// or a separate clear flag) until we actually need to clear it.
type ItemUpdateRequest struct {
	SKU             *string `json:"sku,omitempty" binding:"omitempty,min=1,max=100"`
	Name            *string `json:"name,omitempty" binding:"omitempty,min=1,max=200"`
	Slug            *string `json:"slug,omitempty" binding:"omitempty,max=500"`
	Description     *string `json:"description,omitempty"`
	Price           *int64  `json:"price,omitempty" binding:"omitempty,gte=0"`
	Cost            *int64  `json:"cost,omitempty" binding:"omitempty,gte=0"`
	Stock           *int    `json:"stock,omitempty" binding:"omitempty,gte=0"`
	Backorder       *bool   `json:"backorder,omitempty"`
	Status          *string `json:"status,omitempty" binding:"omitempty,oneof=activo inactivo archivado"`
	CategoryID      *uint64 `json:"category_id,omitempty"`
	Brand           *string `json:"brand,omitempty" binding:"omitempty,max=200"`
	Color           *string `json:"color,omitempty" binding:"omitempty,max=100"`
	ImageURL        *string `json:"image_url,omitempty" binding:"omitempty,url"`
	Items3DID       *uint64 `json:"items3d_id,omitempty"`
	MetaTitle       *string `json:"meta_title,omitempty" binding:"omitempty,max=200"`
	MetaDescription *string `json:"meta_description,omitempty"`
	IsFeatured      *bool   `json:"is_featured,omitempty"`
}

// ItemResponse is the public shape for GET /products and
// GET /products/:slug.
//
// Omitted fields:
//   - Cost: internal business data; would let competitors see margins.
//   - Items3DID: the owner's private printing config; not buyer-relevant.
//   - ViewCount: internal analytics; not exposed in product listing.
//
// InStock is a derived boolean: true if Stock > 0 OR Backorder is true.
// Saves the frontend from computing the OR client-side on every render.
type ItemResponse struct {
	ItemID          uint64    `json:"item_id"`
	SKU             string    `json:"sku"`
	Slug            string    `json:"slug"`
	Name            string    `json:"name"`
	Description     *string   `json:"description,omitempty"`
	Price           int64     `json:"price"`
	Stock           int       `json:"stock"`
	Backorder       bool      `json:"backorder"`
	InStock         bool      `json:"in_stock"`
	Status          string    `json:"status"`
	CategoryID      uint64    `json:"category_id"`
	Brand           *string   `json:"brand,omitempty"`
	Color           *string   `json:"color,omitempty"`
	ImageURL        *string   `json:"image_url,omitempty"`
	MetaTitle       *string   `json:"meta_title,omitempty"`
	MetaDescription *string   `json:"meta_description,omitempty"`
	IsFeatured      bool      `json:"is_featured"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ItemFilter is the query-string shape for GET /products (public list).
//
// form: tags (not json:) because Gin binds query strings via these
// tags. Status defaults to "activo" in the service if empty — public
// listing must never expose inactive/archived items even if the
// caller tampers with the query.
//
// CategoryID and CategorySlug are alternatives: clients choose
// whichever is more convenient. The service resolves slug → id if
// only slug is provided.
type ItemFilter struct {
	CategoryID   *uint64 `form:"category_id"`
	CategorySlug string  `form:"category_slug" binding:"omitempty,max=200"`
	MinPrice     *int64  `form:"min_price" binding:"omitempty,gte=0"`
	MaxPrice     *int64  `form:"max_price" binding:"omitempty,gte=0"`
	Search       string  `form:"search" binding:"omitempty,max=200"`
	Featured     *bool   `form:"featured"`
	Status       string  `form:"status" binding:"omitempty,oneof=activo inactivo archivado"`
	Sort         string  `form:"sort" binding:"omitempty,oneof=newest oldest price_asc price_desc popular"`
	Page         int     `form:"page" binding:"omitempty,min=1"`
	PerPage      int     `form:"per_page" binding:"omitempty,min=1,max=100"`
}

// CategoryCreateRequest is the body for POST /admin/categories.
// Slug is optional — the service auto-generates from Name.
type CategoryCreateRequest struct {
	Name            string  `json:"name" binding:"required,min=1,max=200"`
	Slug            string  `json:"slug,omitempty" binding:"omitempty,max=200"`
	Description     string  `json:"description,omitempty"`
	ParentID        *uint64 `json:"parent_id,omitempty"`
	MetaTitle       string  `json:"meta_title,omitempty" binding:"omitempty,max=200"`
	MetaDescription string  `json:"meta_description,omitempty"`
}

// CategoryUpdateRequest is the body for PUT /admin/categories/:id.
// Same partial-update pattern as ItemUpdateRequest. Same clearing
// limitation for ParentID.
type CategoryUpdateRequest struct {
	Name            *string `json:"name,omitempty" binding:"omitempty,min=1,max=200"`
	Slug            *string `json:"slug,omitempty" binding:"omitempty,max=200"`
	Description     *string `json:"description,omitempty"`
	ParentID        *uint64 `json:"parent_id,omitempty"`
	MetaTitle       *string `json:"meta_title,omitempty" binding:"omitempty,max=200"`
	MetaDescription *string `json:"meta_description,omitempty"`
}

// CategoryResponse is the public shape for GET /categories and
// GET /categories/:slug.
//
// ItemCount is an aggregate: the service can populate it via a
// separate COUNT(*) query when requested (e.g. admin listing) or
// leave it at 0 for public-facing responses where the cost of
// counting every category on every page load is too high.
type CategoryResponse struct {
	CategoryID      uint64    `json:"category_id"`
	Name            string    `json:"name"`
	Slug            string    `json:"slug"`
	Description     *string   `json:"description,omitempty"`
	ParentID        *uint64   `json:"parent_id,omitempty"`
	MetaTitle       *string   `json:"meta_title,omitempty"`
	MetaDescription *string   `json:"meta_description,omitempty"`
	ItemCount       int       `json:"item_count,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}