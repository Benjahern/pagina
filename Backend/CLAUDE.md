# Backend — Tienda Online

## Visión del proyecto (fuente: `../Idea.md`)

> Resumen conceptual de la idea original. Las decisiones técnicas (stack, estructura, modelos) viven en las secciones siguientes — este apartado **no** las redefine, solo fija el alcance funcional que la idea original plantea.

### Roles y catálogo

- **Dos administradores** (el dueño + su papá). Ambos pueden agregar productos al catálogo y gestionar la tienda. Cada uno aporta sus propios productos pero comparten el mismo storefront público.
- **Clientes** se registran por su cuenta para comprar. Los admins pueden gestionarlos (ver pedidos, datos, historial).
- El catálogo es **mixto**: productos físicos + items 3D del dueño. La nota sobre `Items3D` como herramienta personal del dueño (no categoría pública) está en la siguiente sección.

### Calculadora de costes 3D

Para los items 3D propios del dueño, el precio se calcula como:

```
coste = (filamento_gastado × precio_kg_filamento) + (horas_impresión × costo_eléctrico_por_hora_impresora)
```

Los insumos de esta fórmula viven en las tablas `Filament` (precio/kg por marca/color) y `impresora3d` (watts → costo hora). Pendiente de feature, ver sección de Pendientes.

### Funcionalidades administrativas especiales

Estas son **exigencias explícitas de Idea.md** que el sistema debe soportar — no son opcionales:

1. **Vista "admin → usuario normal"**: cuando el dueño inicia sesión con cuenta admin, debe existir un toggle para cambiar a la vista de cliente y navegar la tienda como un usuario cualquiera. Sirve para QA, demos y soporte.
2. **Pedidos manuales**: clientes también contactan por WhatsApp o Facebook. El admin debe poder **crear el pedido manualmente** en el sistema, asignarlo al cliente correspondiente y tener visión centralizada de qué quiere cada usuario, aunque la venta no haya pasado por el checkout web. Esto convierte el panel admin en un mini-CRM de pedidos.
3. **Privacidad**: el registro de usuarios debe requerir la aceptación de un documento de acuerdo de privacidad y confidencialidad. El backend debe persistir el consentimiento (timestamp + versión del documento + flag) para auditoría.

### Pagos (pendiente de decidir)

La pasarela de pagos **aún no está elegida**. Candidatos:

- **Webpay** (Transbank, estándar chileno, alta adopción).
- **Flow** (alternativa chilena, multi-proveedor).

La decisión requiere investigar: comisiones, facilidad de integración con Go, soporte de devoluciones, webhooks, y manejo de pagos manuales para el caso de pedidos creados por admin. Ver Pendientes.

### Auth (recordatorio)

JWT emitido por el backend, **dos tipos de cuenta claramente diferenciados**: `User` (cliente) y `Admin` (administrador). Esto ya está implementado abajo en Stack/Convenciones; se reitera aquí porque Idea.md lo marca como preocupación de seguridad explícita.

## Qué estamos construyendo

Tienda online **general** (productos físicos al público). Polivalente — no es una tienda de impresión 3D. El dueño tiene impresoras y filamento 3D como **herramientas personales/herramientas de producción**, por eso el esquema tiene tablas auxiliares (`Items3D`, `impresora3d`, `Filament`) que registran su operación privada. Esas tablas no son catálogo público: si un producto sale de la impresora del dueño, aparece en la tienda como un `Item` normal, no como algo especial.

Cuando el dueño arranca algo y lo vende, ese algo **es** un `Item` común, con su `slug`, `sku`, precio, stock, etc. El link opcional `Items.items3d_id` permite rastrear qué configuración de impresión se usó para producir ese item, nada más.

## Stack

| Capa | Tecnología |
|---|---|
| Lenguaje | Go |
| HTTP framework | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [GORM](https://gorm.io) |
| Base de datos | PostgreSQL |
| Auth | JWT (`golang-jwt/jwt v5`) + bcrypt |
| Validación | `go-playground/validator` (via `binding:` de Gin) |
| Config | `godotenv` + variables de entorno |
| DB driver | `gorm.io/driver/postgres` |

## Arquitectura

Capas en orden de dependencia: **handler → service → repository → models**. Cada capa expone interfaces que la capa superior consume (no structs concretas), salvo en `cmd/server/main.go` donde se hace el cableado.

| Capa | Carpeta | Qué vive ahí |
|---|---|---|
| HTTP | `handlers/` | Bind de request, formato de response, código HTTP. Sin lógica. |
| Negocio | `services/` | Hashing de passwords, JWT, reglas, transacciones. Errores de dominio. |
| Datos | `repository/` | Queries GORM. Errores `ErrNotFound`. Sin reglas. |
| Modelos | `models/` | Entidades GORM, relaciones, constantes enum. |
| Wire shape | `dto/` | Request/response con tags `binding:`. Separado de models. |
| Auth/transversal | `middleware/` | `RequireAuth`, `RequireAdmin`, error handler. |
| Rutas | `routes/` | Registro central con grupos public/auth/admin. |
| Helpers | `utils/` | JWT, bcrypt, response helpers, errores sentinels. |
| Bootstrap | `config/`, `database/`, `seed/` | Env, conexión, datos iniciales. |
| Entry | `cmd/server/` | `main.go` — único punto de entrada. |

`internal/` deja todo privado al módulo: otros proyectos no pueden importar estos paquetes.

## Estructura de carpetas

```
Backend/
├── cmd/server/                  # main.go
├── internal/
│   ├── config/                  # carga de env vars + validación
│   ├── database/                # postgres.go: Connect + AutoMigrate
│   ├── models/                  # User, Item, Category, Order, etc.
│   ├── dto/                     # shapes req/res con binding tags
│   ├── repository/              # interfaces + impls GORM
│   ├── services/                # lógica de negocio
│   ├── handlers/                # HTTP handlers de Gin
│   ├── middleware/              # RequireAuth, RequireAdmin, error
│   ├── routes/                  # registerRoutes() central
│   ├── utils/                   # jwt, password, response, errors
│   └── seed/                    # seeder (admin + samples)
└── docs/                        # documentación de la API
```

## Archivos importantes

- `internal/database/database.sql` — **es DBML (dbdiagram.io), NO SQL crudo**. Para SQL: importar en [dbdiagram.io](https://dbdiagram.io/d) y exportar, o usar `dbml2sql`. Es la fuente de verdad del modelo de datos.
- `internal/utils/errors.go` — sentinels: `ErrNotFound`, `ErrConflict`, `ErrUnauthorized`, `ErrValidation`, `ErrNotImplemented`. Usar estos, no `errors.New()` ad-hoc.
- `internal/utils/response.go` — `Success`, `Created`, `Error`, `Paginated`. Todas las respuestas JSON deben pasar por aquí para mantener el formato `{ data, meta, error }` consistente.

## Convenciones

- **Naming Go**: packages `snake_case`, structs `PascalCase`, exported `PascalCase`, interno `camelCase`.
- **Naming DB**: columnas `snake_case`. Tablas mantienen la mezcla original del esquema (`Orders`, `Users`, `Cart_Item`, `shipping_address`). No renombrar sin migración planeada.
- **Naming models vs tablas**: struct en singular Go (`User`, `Item`), `TableName()` retorna la forma plural/original del DBML (`Users`, `Items`). Snake_case columnas vía `NamingStrategy` de GORM — sin taggear cada campo.
- **Organización de archivos en `internal/models/`**: **un archivo por modelo**. Extras (enums, constantes, helpers) del modelo en archivo separado del mismo paquete con prefijo: `item.go` + `item_status.go`, `payment.go` + `payment_method.go` + `payment_status.go`, etc.
- **Organización de archivos en `internal/dto/`**: por dominio funcional, no por modelo. Ej: `auth_user.go` agrupa todos los DTOs del flujo auth de usuarios; `auth_admin.go` los de admin. Evitar un DTO por archivo.
- **Idiomas**: código, comentarios, logs en **inglés**. Mensajes al usuario final en **español** (cuando se agregue UI).
- **Errores**: services mapean `repository.ErrNotFound → service.ErrNotFound` y agregan específicos (`ErrEmailTaken`, `ErrInvalidCredentials`). El middleware de error en handler convierte a código HTTP.
- **DTOs siempre**: nunca serializar un `model` directo a JSON — filtra `password_hash` y columnas internas. Mapear a DTO.
- **Slug generation**: pendiente (hook `BeforeCreate` para auto-generar desde `name` si está vacío).
- **Timestamps**: `timestamptz` siempre (nunca `date`).
- **Seguridad**: GORM parametriza queries por default. Nunca `fmt.Sprintf` con user input para construir SQL. Ver comentario en `postgres.go`.

## Comandos comunes

```bash
# Cuando exista docker-compose.yml:
docker compose up -d postgres

go mod tidy                       # descargar deps
go run ./cmd/server               # arrancar
go build -o bin/server ./cmd/server
go test ./...
```

(Targets de Makefile vendrán cuando se cree.)

## Variables de entorno

Requeridas (sin default, fallan si están vacías): `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `JWT_SECRET` (mínimo 32 bytes).

Opcionales con default: `PORT=8080`, `GIN_MODE=debug`, `JWT_EXPIRATION_HOURS=24`, `BCRYPT_COST=12`, `SEED=false`, `CORS_ALLOWED_ORIGINS=http://localhost:3000`, `LOG_LEVEL=info`.

## Modelo de dominio (resumen)

**Tienda (este es el producto):**
- `Users` — clientes. Email único, password hasheado, phone opcional.
- `Category` — categorías de productos, con `slug` para SEO y `parent_id` opcional para subcategorías.
- `Items` — productos a la venta: `sku`, `slug`, `name`, precio, costo, stock, estado, marca, color, `category_id`, meta SEO (`meta_title`, `meta_description`), `is_featured`, `view_count`, `items3d_id` opcional (link a la config de impresión si fue hecho por el dueño).
- `Cart` + `Cart_Item` — un carrito activo por usuario (unique en `user_id`), UNIQUE(`cart_id`, `item_id`).
- `Orders` + `Order_item` — órdenes con `shipping_address_id` (NOT NULL), UNIQUE(`order_id`, `item_id`).
- `Payment` — pagos por orden con `transaction_id` único (idempotencia con gateway).
- `shipping_address` — varias direcciones por usuario, una `is_default`.
- `Stock_movement` — auditoría de cambios de stock (positivo/negativo, `reason` ∈ venta|devolucion|ajuste|reposicion|danado).
- `Review` — reseñas con rating 1-5, `approved` (moderación), UNIQUE implícito por user+item (pendiente agregar).
- `Admins` + `Admin_Permission` + `Permission` — administradores y permisos (la fusión con `Users + role` es un pendiente).

**Personal/herramientas del dueño (no es producto público):**
- `impresora3d` — impresoras 3D registradas, costos operativos.
- `Filament` — inventario de filamentos (marca, color, costo/kg, slug).
- `Items3D` — configuraciones de impresión (filamento, horas, costo), referenciada opcionalmente desde `Items`.

## Endpoints (a implementar)

Públicos: `POST /api/v1/auth/register`, `POST /api/v1/auth/login`, `GET /api/v1/products`, `GET /api/v1/products/:slug`, `GET /api/v1/categories`, `GET /api/v1/categories/:slug`.

Auth: `GET /api/v1/auth/me`, `GET /api/v1/cart`, `POST /api/v1/cart/items`, `PUT /api/v1/cart/items/:id`, `DELETE /api/v1/cart/items/:id`, `POST /api/v1/orders`, `GET /api/v1/orders`, `GET /api/v1/orders/:id`.

Admin: `POST/PUT/DELETE /api/v1/admin/products`, `POST/PUT/DELETE /api/v1/admin/categories`.

Health: `GET /health`.

(Pendiente mover a `docs/api.md` cuando se cree.)

## Pendientes conocidos (v1)

- [x] Inicializar `go.mod` y descargar deps
- [x] `config/` + `.env.example` + carga de env
- [x] `database/postgres.go` — conexión GORM + pool + ping (sin AutoMigrate aún)
- [x] **DBML: añadir índices faltantes para query patterns reales** (Orders.shipping_address_id, Payment.order_id, Order_item.item_id, Review UNIQUE, Items.created_at, Stock_movement.order_id, Orders(status, created_at))
- [x] **Models GORM** que reflejen el esquema (ver lista completa abajo) — **cada índice nuevo en el DBML debe tener su tag `gorm:"index:idx_xxx"` en el struct, o AutoMigrate no los creará**
- [ ] DTOs de auth (`auth_user.go`, `auth_admin.go`) — primer set de DTOs
- [x] `db.AutoMigrate(...)` en `Connect()` o en un `Migrate()` separado
- [ ] Slice vertical **Auth** completo (register/login/me + JWT middleware)
- [ ] Lectura pública de productos y categorías (con paginación + filtro + búsqueda)
- [ ] Scaffolds (stubs 501) para admin CRUD, cart, orders
- [ ] Hook para auto-generar slug desde `name` (en `Item.BeforeCreate`)
- [ ] Tests básicos (patrón con `services/auth_service_test.go`)
- [ ] `Makefile` con targets run/build/test/tidy
- [ ] `Dockerfile` multi-stage + docker-compose con postgres
- [ ] Seeder opcional (admin user + categorías de ejemplo)
- [ ] Decidir si fusionar `Admins` con `Users + role`

### Pendientes funcionales (de `Idea.md`)

Estos vienen de la visión del proyecto, no del desglose técnico. **No son opcionales**.

- [ ] **Calculadora de costes 3D** — endpoint que reciba `item3d_id` (o parámetros) y devuelva precio calculado: filamento + horas × costo eléctrico. Depende de `Filament` y `impresora3d` pobladas con costos reales.
- [ ] **Vista "admin → usuario normal"** — mecanismo para que un admin navegue la tienda como cliente sin cerrar sesión. Probablemente un header `X-View-Mode: user` o un token secundario; el frontend es quien más lo usa, pero el backend debe respetarlo (no revelar campos admin en respuestas cuando está en modo user).
- [ ] **Pedidos manuales desde panel admin** — endpoints para que un admin cree un pedido en nombre de un cliente (caso WhatsApp/Facebook). Acepta cliente nuevo o existente, agrega items, fija total, opcionalmente marca como `pagado` con método `efectivo`/`transferencia`.
- [ ] **Gestión de usuarios por admin** — CRUD/lectura de clientes: listar, buscar, ver historial de pedidos, ver direcciones. No necesariamente edición libre (probablemente restricción a soft-ban + notas internas).
- [ ] **Aceptación de acuerdo de privacidad** — al registrarse el cliente debe aceptar el documento. Persistir: `user_id`, `document_version`, `accepted_at`, `ip`. Campo nuevo en `Users` o tabla aparte (`PrivacyConsent`).
- [ ] **Integración con pasarela de pagos** — investigar Webpay vs Flow vs otros. Decisión basada en comisiones, integración con Go (SDK oficial vs REST), webhooks, devoluciones. Implementar contra el ganador.
- [ ] **Documento de privacidad como archivo versionado** — almacenar las versiones del documento (markdown/pdf) en repo o storage; el backend expone la versión vigente que el frontend debe mostrar al cliente.

## Models planificados (no implementados aún)

Plan acordado el 2026-06-30: 17 structs + 6 archivos extra de enums/helpers = **23 archivos** en `internal/models/`.

### Capa pública (catálogo/operación)

| Archivo | Struct Go | Tabla | Notas |
|---|---|---|---|
| `user.go` | `User` | `Users` | Customer del storefront |
| `category.go` | `Category` | `Category` | Categorías con jerarquía |
| `item.go` | `Item` | `Items` | Productos (con hook BeforeCreate para slug) |
| `item_status.go` | enum `ItemStatus` | — | `activo\|inactivo\|archivado` |
| `cart.go` | `Cart` | `Cart` | 1 activo por usuario |
| `cart_status.go` | enum `CartStatus` | — | `activo\|abandonado\|completado` |
| `cart_item.go` | `CartItem` | `Cart_Item` | Línea de carrito |
| `order.go` | `Order` | `Orders` | Pedidos |
| `order_status.go` | enum `OrderStatus` | — | `pendiente\|pagado\|enviado\|entregado\|cancelado` |
| `order_item.go` | `OrderItem` | `Order_item` | Línea de pedido con snapshot de precio |
| `payment.go` | `Payment` | `Payment` | Pagos por orden |
| `payment_method.go` | enum `PaymentMethod` | — | `efectivo\|tarjeta\|transferencia\|webpay` |
| `payment_status.go` | enum `PaymentStatus` | — | `pendiente\|aprobado\|rechazado\|reembolsado` |
| `shipping_address.go` | `ShippingAddress` | `shipping_address` | Direcciones por usuario |
| `stock_movement.go` | `StockMovement` | `Stock_movement` | Auditoría de stock |
| `stock_movement_reason.go` | enum `StockMovementReason` | — | `venta\|devolucion\|ajuste\|reposicion\|danado` |
| `review.go` | `Review` | `Review` | Reseñas con moderación |
| `review_rating.go` | helpers rating | — | `IsValidRating(int) bool` (1–5) |

### Capa admin (backend)

| Archivo | Struct | Tabla | Notas |
|---|---|---|---|
| `admin.go` | `Admin` | `Admins` | Login admin |
| `permission.go` | `Permission` | `Permission` | Permisos granulares |
| `admin_permission.go` | `AdminPermission` | `Admin_Permission` | Join table |

### Capa personal/herramientas (no público)

| Archivo | Struct | Tabla | Notas |
|---|---|---|---|
| `impresora3d.go` | `Impresora3d` | `impresora3d` | Impresoras del dueño |
| `filament.go` | `Filament` | `Filament` | Inventario filamentos |
| `item3d.go` | `Item3D` | `Items3D` | Configs de impresión |

## DTOs planificados (no implementados aún)

### Auth (primer set, acordado el 2026-06-30)

```
internal/dto/
├── auth_user.go         # UserRegisterRequest, UserLoginRequest, UserLoginResponse, UserResponse
└── auth_admin.go        # AdminLoginRequest, AdminLoginResponse, AdminResponse
```

Separados en archivos distintos porque `User` y `Admin` viven en dominios distintos (auth público vs auth admin con middleware separado), aunque el request de login comparta forma. Mantenerlos separados evita cruces accidentales.

## Estado del proyecto (2026-06-30)

**Capa `config/` + `database/postgres.go` + índices del DBML completadas y validadas** (smoke test live OK con `deployment-postgres-1`).

- `go.mod` + `godotenv` + `gorm.io/gorm` + `gorm.io/driver/postgres` vía `go mod tidy`
- `internal/config/config.go` con `Load()`, validación en boot (`JWT_SECRET ≥ 32 bytes`, `BCRYPT_COST ∈ [4,31]`)
- `Backend/.env.example` (plantilla) + `Backend/.env` (real, gitignored) con DB host `dbtienda`
- `internal/database/postgres.go` con `Connect(cfg)`: DSN seguro, GORM abierto, pool (25/5, 5min/10min), ping con timeout 5s. **Comentario de seguridad explícito** sobre queries parametrizadas vs `fmt.Sprintf(userInput)` — defensa contra SQL injection desde el día 1.
- `cmd/server/main.go` cablea `config.Load()` → `database.Connect()` y muestra resumen + cierre limpio con `defer Close` (vía `db.DB().Close()` — *gorm.DB* no expone Close directo).
- **DBML actualizado con 7 índices nuevos** cubriendo FKs faltantes y sort patterns reales (ver sección "Pendientes" arriba).

**Verificación live:** contenedor `deployment-postgres-1` arrancado en `localhost:5432`, conexión con override de env (`DB_USER=postgres DB_PASSWORD=lo DB_NAME=postgres`) retorna "✓ open, pool configured, ping ok". El host `dbtienda` (nombre de servicio docker) solo resuelve desde la red `deployment_template-network` — para el docker-compose venidero.

**Sin commitear:** `cmd/server/main.go`, `internal/config/config.go`, `Backend/.env.example`, `internal/database/postgres.go`, `internal/database/database.sql` (con los 7 índices nuevos), `go.mod`, `go.sum`. Commit recomendado antes de seguir.

## Estado del proyecto (2026-07-01)

**17 modelos GORM + AutoMigrate completados y validados live.**

- `internal/models/` — los 17 archivos del plan, cada uno con `TableName()` y los tags `gorm:"index:...,uniqueIndex:...,not null,size:...,type:numeric(...)"` correspondientes. Money como `int64` + `numeric(12,2)` (CLP sin decimales). Nullable como punteros. `Item.BeforeCreate` para auto-slug sigue pendiente (en el plan).
- `internal/database/migrate.go` — `Migrate(db)` lista los 17 structs en orden de dependencia (raíz → joins) y delega en `db.AutoMigrate(...)`. Documenta que **NO** aplica FKs `REFERENCES`, políticas `ON DELETE`, ni `CHECK` constraints del DBML — eso vive en el SQL que se genere desde `database.sql` (dbdiagram.io) y debe correrse por separado cuando se necesite el FK topológico completo.
- `cmd/server/main.go` invoca `database.Migrate(db)` después de `Connect`. Imprime "✓ 17 models reconciled" en boot.

**Validación live contra `deployment-postgres-1`**:
- 17 tablas creadas con nombres exactos del DBML (`Users`, `Admins`, `Items`, `Items3D`, `Cart_Item`, `Admin_Permission`, `Order_item`, `Stock_movement`, `shipping_address`, `impresora3d`, etc.).
- 27 índices verificados: únicos simples, compuestos (`uq_admin_permission`, `uq_cart_item`, `uq_order_item`, `uq_review_user_item`, `uq_payment_transaction`) y multi-columna (`idx_review_item_approved`, `idx_orders_status_date`, `idx_orders_user_date`).
- Defaults y tipos numéricos respetados (`Items.status` → `'activo'`, `Payment.status` → `'pendiente'`, `impresora3d.error_margin` → `numeric(5,2) DEFAULT 10.00`).

**Bug detectado y corregido**:
- La `NamingStrategy` de GORM partía `Items3DID` / `Item3DID` en `items3_d_id` / `item3_d_id` (interpreta `3D` como transición camelCase). El DBML espera `items3d_id` / `item3d_id`. Sin esto, los `FOREIGN KEY ... REFERENCES "Items3D"` del SQL generado desde el DBML habrían fallado en el momento de aplicarlos.
- **Fix**: tag explícito `gorm:"column:item3d_id"` en los 2 structs afectados (`item.go`, `item3d.go`). Tablas re-creadas, columnas verificadas.
- **Lección**: cuando un identificador Go termina en `3D` (o cualquier `#<letra>` donde `#` es dígito), GORM añade `_` delante de la letra. Mejor tag explícito que renombrar el campo Go y romper el dominio.

**Sin commitear:** los 17 archivos en `internal/models/`, `internal/database/migrate.go`, cambios en `cmd/server/main.go`, y los 2 tags nuevos en `item.go`/`item3d.go`. Commit recomendado antes de seguir con DTOs de auth.

**Próximo paso (sesión siguiente):** DTOs de auth — `internal/dto/auth_user.go` (`UserRegisterRequest`, `UserLoginRequest`, `UserLoginResponse`, `UserResponse`) y `internal/dto/auth_admin.go` (`AdminLoginRequest`, `AdminLoginResponse`, `AdminResponse`). Después: slice vertical Auth (register/login/me + JWT middleware) o 7 archivos de enums/helpers primero. Decisión del usuario.

## Recordatorios importantes

- `database.sql` **es DBML**, no SQL. Para regenerar el SQL: dbdiagram.io o `dbml2sql`.
- `Go` está disponible en el sistema. Las dependencias se instalan vía `go get` / `go mod tidy`.
- El usuario es estudiante de Ing. de Software aprendiendo Go: prefiere respuestas explicativas (el por qué), no solo el qué. Si hay varias formas, mostrarlas y justificar la recomendación.
- Si aparece un error en cualquier código/compilación/ejecución, delegar al agente `error-diagnosis-specialist` antes de aplicar fix directo.
