// Esquema DBML — dbdiagram.io
// Para exportar a SQL: https://dbdiagram.io/d → "Export to PostgreSQL"

Table "public"."Orders" {
  "order_id" bigint [pk, not null, increment]
  "user_id" bigint [ref: > "public"."Users"."user_id"]
  "total" numeric(12,2) [not null, note: 'Total final del pedido']
  "status" varchar(50) [default: 'pendiente', note: 'pendiente|pagado|enviado|entregado|cancelado']
  "shipping_address_id" bigint [not null, ref: - "public"."shipping_address"."shipping_address_id"]
  "created_at" timestamptz [not null, default: `now()`]
  "updated_at" timestamptz [not null, default: `now()`]
  Indexes {
    user_id [name: 'idx_orders_user']
    (user_id, created_at) [name: 'idx_orders_user_date']
    status [name: 'idx_orders_status']
    shipping_address_id [name: 'idx_orders_shipping']
    (status, created_at) [name: 'idx_orders_status_date']
  }
}

Table "public"."Users" {
  "user_id" bigint [pk, not null, increment]
  "email" varchar(255) [unique, not null, note: 'Único por persona']
  "pass" varchar(255) [not null, note: 'Hasheado con bcrypt']
  "name" varchar(200) [not null]
  "phone" varchar(50) [note: 'Opcional; útil para envíos']
  "created_at" timestamptz [not null, default: `now()`]
  "updated_at" timestamptz [not null, default: `now()`]
  Indexes {
    email [name: 'idx_users_email']
  }
}

Table "public"."Admin_Permission" {
  "admin_permission_id" bigint [pk, not null, increment]
  "admin_id" bigint [not null, ref: > "public"."Admins"."admin_id"]
  "permission_id" bigint [not null, ref: > "public"."Permission"."permission_id"]
  "created_at" timestamptz [not null, default: `now()`]
  Indexes {
    (admin_id, permission_id) [unique, name: 'uq_admin_permission']
  }
}

Table "public"."Order_item" {
  "order_item_id" bigint [pk, not null, increment]
  "order_id" bigint [not null, ref: > "public"."Orders"."order_id"]
  "item_id" bigint [not null, ref: > "public"."Items"."item_id"]
  "quantity" int [not null, note: 'CHECK > 0']
  "unit_price" numeric(12,2) [not null, note: 'Precio al momento de comprar (snapshot)']
  "subtotal" numeric(12,2) [not null, note: '= unit_price * quantity']
  "created_at" timestamptz [not null, default: `now()`]
  Indexes {
    (order_id, item_id) [unique, name: 'uq_order_item']
    item_id [name: 'idx_order_item_item']
  }
}

Table "public"."Cart_Item" {
  "cart_item_id" bigint [pk, not null, increment]
  "cart_id" bigint [not null, ref: > "public"."Cart"."cart_id"]
  "item_id" bigint [not null, ref: > "public"."Items"."item_id"]
  "quantity" int [not null, default: 1, note: 'CHECK > 0; si agregas el mismo producto, sumar cantidad']
  "created_at" timestamptz [not null, default: `now()`]
  "updated_at" timestamptz [not null, default: `now()`]
  Indexes {
    (cart_id, item_id) [unique, name: 'uq_cart_item']
  }
}

Table "public"."Permission" {
  "permission_id" bigint [pk, not null, increment]
  "access" varchar(200) [unique, not null, note: 'Antes "acces" (typo)']
  "description" text
  "created_at" timestamptz [not null, default: `now()`]
}

Table "public"."Filament" {
  "filament_id" bigint [pk, not null]
  "name" varchar(200) [not null]
  "slug" varchar(200) [unique, note: 'URL-friendly; ej: pla-negro']
  "cost_kilogram" numeric(12,2)
  "color" varchar(100)
  "brand" varchar(200)
  "created_at" timestamptz [not null, default: `now()`]
}

Table "public"."Payment" {
  "payment_id" bigint [pk, not null, increment]
  "order_id" bigint [not null, ref: - "public"."Orders"."order_id"]
  "method" varchar(50) [not null, note: 'efectivo|tarjeta|transferencia|webpay']
  "amount" numeric(12,2) [not null]
  "status" varchar(50) [default: 'pendiente', note: 'pendiente|aprobado|rechazado|reembolsado']
  "transaction_id" varchar(200) [unique, note: 'ID externo del gateway']
  "paid_at" timestamptz
  "created_at" timestamptz [not null, default: `now()`]
  "updated_at" timestamptz [not null, default: `now()`]
  Indexes {
    order_id [name: 'idx_payment_order']
  }
}

Table "public"."Category" {
  "category_id" bigint [pk, not null, increment]
  "name" varchar(200) [not null]
  "slug" varchar(200) [unique, not null, note: 'URL-friendly; ej: electronica']
  "description" text
  "image_url" varchar(500)
  "meta_title" varchar(200)
  "meta_description" text
  "parent_id" bigint [ref: > "public"."Category"."category_id", note: 'Para subcategorías']
  "created_at" timestamptz [not null, default: `now()`]
  "updated_at" timestamptz [not null, default: `now()`]
  Indexes {
    parent_id [name: 'idx_category_parent']
  }
}

Table "public"."Cart" {
  "cart_id" bigint [pk, not null, increment]
  "user_id" bigint [unique, ref: - "public"."Users"."user_id", note: 'Un carrito activo por usuario']
  "status" varchar(50) [not null, default: 'activo', note: 'activo|abandonado|completado']
  "created_at" timestamptz [not null, default: `now()`]
  "updated_at" timestamptz [not null, default: `now()`]
}

Table "public"."Items" {
  "item_id" bigint [pk, not null, increment]
  "sku" varchar(100) [unique, not null, note: 'Stock-keeping unit interno']
  "slug" varchar(500) [unique, not null, note: 'URL-friendly; ej: mando-xbox-series-x']
  "name" varchar(200) [not null]
  "description" text
  "price" numeric(12,2) [not null, note: 'Precio de venta']
  "cost" numeric(12,2) [not null, note: 'Costo de producción/adquisición']
  "stock" int [not null, default: 0]
  "backorder" boolean [not null, default: false, note: '¿Se puede vender sin stock?']
  "status" varchar(50) [not null, default: 'activo', note: 'activo|inactivo|archivado']
  "category_id" bigint [not null, ref: > "public"."Category"."category_id"]
  "brand" varchar(200)
  "color" varchar(100)
  "image_url" varchar(500)
  "items3d_id" bigint [ref: > "public"."Items3D"."items3d_id", note: 'Config de impresión 3D asociada (opcional)']
  "meta_title" varchar(200)
  "meta_description" text
  "view_count" int [not null, default: 0]
  "is_featured" boolean [not null, default: false, note: 'Destacado en homepage']
  "created_at" timestamptz [not null, default: `now()`]
  "updated_at" timestamptz [not null, default: `now()`]
  Indexes {
    category_id [name: 'idx_items_category']
    status [name: 'idx_items_status']
    (status, category_id) [name: 'idx_items_status_category']
    is_featured [name: 'idx_items_featured']
    created_at [name: 'idx_items_created']
  }
}

Table "public"."Admins" {
  "admin_id" bigint [pk, not null, increment]
  "email" varchar(200) [unique, not null]
  "nombre" varchar(200) [not null]
  "phone" varchar(50)
  "pass" varchar(255) [note: 'Hasheado; considera fusionar con Users+role']
  "is_active" boolean [not null, default: true]
  "created_at" timestamptz [not null, default: `now()`]
  "updated_at" timestamptz [not null, default: `now()`]
}

Table "public"."impresora3d" {
  "impresora3d_id" bigint [pk, not null, increment]
  "name" varchar(200) [note: 'Nombre identificador; ej: Prusa-01']
  "electricity_cost_per_hour" numeric(12,4)
  "cost_reparation" numeric(12,2)
  "error_margin" numeric(5,2) [default: 10.00, note: '% de margen de error (antes: margén_error como bigint)']
  "useful_life_hours" int [note: 'Vida útil en horas']
  "is_active" boolean [not null, default: true]
  "created_at" timestamptz [not null, default: `now()`]
  "updated_at" timestamptz [not null, default: `now()`]
}

Table "public"."Items3D" {
  "items3d_id" bigint [pk, not null, increment]
  "name" varchar(200) [not null, note: 'Nombre de la configuración de impresión']
  "impresora3d_id" bigint [ref: > "public"."impresora3d"."impresora3d_id"]
  "filament_grams" double precision [note: 'Gramos de filamento usados']
  "hours" int [not null, default: 0]
  "minutes" int [not null, default: 0]
  "extra_cost" numeric(12,2) [note: 'Costos extra (post-procesado, etc.)']
  "cost" numeric(12,2) [note: 'Costo total calculado']
  "filament_id" bigint [ref: > "public"."Filament"."filament_id"]
  "created_at" timestamptz [not null, default: `now()`]
  "updated_at" timestamptz [not null, default: `now()`]
  Indexes {
    impresora3d_id [name: 'idx_items3d_printer']
  }
}

Table "public"."shipping_address" {
  "shipping_address_id" bigint [pk, not null, increment]
  "user_id" bigint [not null, ref: > "public"."Users"."user_id"]
  "address_line" varchar(500) [not null]
  "city" varchar(200) [not null]
  "postal_code" varchar(50) [not null]
  "commune" varchar(200) [not null]
  "is_default" boolean [not null, default: false]
  "created_at" timestamptz [not null, default: `now()`]
  "updated_at" timestamptz [not null, default: `now()`]
  Indexes {
    user_id [name: 'idx_shipping_user']
  }
}

Table "public"."Stock_movement" {
  "movement_id" bigint [pk, not null, increment]
  "item_id" bigint [not null, ref: > "public"."Items"."item_id"]
  "change" int [not null, note: 'Positivo=entrada, negativo=salida; ej: +10, -5']
  "reason" varchar(50) [not null, note: 'venta|devolucion|ajuste|reposicion|danado']
  "order_id" bigint [ref: > "public"."Orders"."order_id", note: 'NULL si no связа a venta (reposición, ajuste)']
  "user_id" bigint [ref: > "public"."Users"."user_id", note: 'Quién hizo el movimiento']
  "created_at" timestamptz [not null, default: `now()`]
  Indexes {
    item_id [name: 'idx_stock_movement_item']
    created_at [name: 'idx_stock_movement_date']
    order_id [name: 'idx_stock_movement_order']
  }
}

Table "public"."Review" {
  "review_id" bigint [pk, not null, increment]
  "user_id" bigint [not null, ref: > "public"."Users"."user_id"]
  "item_id" bigint [not null, ref: > "public"."Items"."item_id"]
  "rating" int [not null, note: '1-5 estrellas']
  "comment" text
  "approved" boolean [not null, default: false, note: 'Pasa por moderación']
  "created_at" timestamptz [not null, default: `now()`]
  Indexes {
    item_id [name: 'idx_review_item']
    (item_id, approved) [name: 'idx_review_item_approved']
    (user_id, item_id) [unique, name: 'uq_review_user_item']
  }
}
