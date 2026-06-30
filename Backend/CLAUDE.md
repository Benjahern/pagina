# Backend — Tienda Online

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
- **Idiomas**: código, comentarios, logs en **inglés**. Mensajes al usuario final en **español** (cuando se agregue UI).
- **Errores**: services mapean `repository.ErrNotFound → service.ErrNotFound` y agregan específicos (`ErrEmailTaken`, `ErrInvalidCredentials`). El middleware de error en handler convierte a código HTTP.
- **DTOs siempre**: nunca serializar un `model` directo a JSON — filtra `password_hash` y columnas internas. Mapear a DTO.
- **Slug generation**: pendiente (hook `BeforeCreate` para auto-generar desde `name` si está vacío).
- **Timestamps**: `timestamptz` siempre (nunca `date`).

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

- [ ] Inicializar `go.mod` y descargar deps
- [ ] `config/` + `.env.example` + carga de env
- [ ] `database/postgres.go` + conexión + AutoMigrate
- [ ] Models GORM que reflejen el esquema
- [ ] Slice vertical **Auth** completo (register/login/me + JWT middleware)
- [ ] Lectura pública de productos y categorías (con paginación + filtro + búsqueda)
- [ ] Scaffolds (stubs 501) para admin CRUD, cart, orders
- [ ] Hook para auto-generar slug desde `name`
- [ ] Tests básicos (patrón con `services/auth_service_test.go`)
- [ ] `Makefile` con targets run/build/test/tidy
- [ ] `Dockerfile` multi-stage + docker-compose con postgres
- [ ] Seeder opcional (admin user + categorías de ejemplo)
- [ ] Decidir si fusionar `Admins` con `Users + role`

## Recordatorios importantes

- `database.sql` **es DBML**, no SQL. Para regenerar el SQL: dbdiagram.io o `dbml2sql`.
- `Go` está disponible en el sistema. Las dependencias se instalan vía `go get` / `go mod tidy`.
- El usuario es estudiante de Ing. de Software aprendiendo Go: prefiere respuestas explicativas (el por qué), no solo el qué. Si hay varias formas, mostrarlas y justificar la recomendación.
- Si aparece un error en cualquier código/compilación/ejecución, delegar al agente `error-diagnosis-specialist` antes de aplicar fix directo.
