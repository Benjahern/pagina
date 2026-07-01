# Frontend — Tienda Online

> Espejo de `Backend/CLAUDE.md` desde la óptica del frontend. Convenciones técnicas detalladas (modelos, endpoints, auth JWT, decisiones de arquitectura backend) **viven en ese archivo** y no se duplican aquí para evitar drift.

## Visión del proyecto (fuente: `../Idea.md`)

Mismo origen conceptual que el backend. **Lo que `Idea.md` exige funcionalmente al frontend:**

- Catálogo público mixto (productos físicos + items 3D del dueño).
- Registro de clientes con **aceptación obligatoria** del documento de privacidad y confidencialidad.
- Flujo de compra completo para clientes.
- **Toggle "vista admin → usuario normal"**: el admin debe poder navegar la tienda como cliente sin cerrar sesión.
- **Creación manual de pedidos** desde el panel admin (casos WhatsApp/Facebook).
- SEO decente (motivo por el que `Idea.md` eligió Nuxt).

Para entender el modelo de datos, los roles (dos admins: dueño + papá), la calculadora de costes 3D, y los pendientes completos, ver `Backend/CLAUDE.md` § "Visión del proyecto".

## Qué estamos construyendo (en este repo)

Sitio público (storefront) + panel admin. **Mismo código, dos modos** — el admin puede ver el storefront (toggle), y el panel admin vive detrás de un layout separado pero comparte componentes (catálogo, pedidos, etc.) con la vista pública.

## Stack

| Capa | Tecnología | Notas |
|---|---|---|
| Framework | [Nuxt 3](https://nuxt.com) | `Idea.md` lo nombra. SSR para SEO, file-based routing. |
| Lenguaje | TypeScript | Default moderno. **Confirmar preferencia.** |
| Estado | Pinia | Default de Nuxt 3. |
| UI library | **pendiente** | Decisión bloqueante antes de empezar. |
| HTTP | `ofetch` (nativo Nuxt) | Wrapper sobre `fetch`. |
| Auth storage | httpOnly cookie (recomendado) | Detalle abajo. |

### Decisión de seguridad: dónde guardar el JWT

`Idea.md` dice "JWT mediante el backend" — el backend emite, pero **dónde lo guarda el frontend** es decisión de UI.

| Opción | Seguridad | UX |
|---|---|---|
| `localStorage` | Malo: vulnerable a XSS | Persistente fácil |
| **cookie httpOnly** (recomendado) | Bueno: JS no puede leer | Requiere CSRF protection |
| Memoria (ref/composable) | Bueno: | Se pierde en refresh |

Recomendación: **httpOnly cookie** + CSRF token. Tienda con datos de pago = un solo XSS roba sesión. Implica alinear con backend (cambio de `Authorization: Bearer` a cookie reader). **Decisión pendiente — bloqueante para empezar el auth flow.**

## Estructura de carpetas (a crear)

```
Frontend/
├── pages/
│   ├── index.vue              # home / catálogo
│   ├── producto/[slug].vue    # detalle público
│   ├── categoria/[slug].vue
│   ├── login.vue
│   ├── registro.vue           # registro + privacidad
│   ├── carrito.vue
│   ├── checkout.vue
│   ├── cuenta/                # área cliente (auth)
│   └── admin/                 # panel admin (RequireAdmin)
├── components/                # CardProducto, etc.
├── layouts/
│   ├── default.vue            # storefront
│   └── admin.vue              # layout admin (nav lateral)
├── composables/               # useAuth, useCart, useApi
├── stores/                    # Pinia (auth, cart, ui)
├── middleware/                # auth, admin, redirect-logged-in
├── plugins/                   # cliente API, etc.
├── public/
│   └── legal/                 # documento de privacidad versionado
├── server/                    # (si necesitamos endpoints proxy)
├── assets/
├── nuxt.config.ts
└── package.json
```

## Funcionalidades específicas del frontend

### 1. Registro con privacidad

El registro del cliente exige **scroll-to-bottom** del documento de privacidad antes de habilitar el botón "Aceptar y registrarme". El payload a enviar al backend:

```ts
{
  email: string
  password: string
  phone?: string
  privacyAccepted: true              // solo si el checkbox está activo
  privacyVersion: string             // ej "v1.0-2026-07-01"
}
```

`privacyVersion` debe leerlo el frontend desde `public/legal/privacidad-<version>.md` para evitar drift con el backend.

### 2. Toggle "admin → usuario normal"

Tres approaches reales:

1. **Botón en layout** que setea `useState('viewMode')` y cambia el contenido mostrado. Simple, no toca red.
2. **Ruta paralela** `/vista-usuario` que monta el storefront enviando `X-View-Mode: user` al backend para que responses no incluyan campos admin. Más complejo, más fiel a la experiencia real.
3. **Dos sesiones paralelas** (admin + test-user simulado). Innecesario para QA.

Recomendación: **opción 1** si el admin solo necesita ver la tienda vacía; **opción 2** si quiere ver la experiencia del cliente real con sus datos. **Decisión pendiente.**

### 3. Pedidos manuales (panel admin)

UI para que el admin cree pedidos en nombre de un cliente (WhatsApp/Facebook). Componentes:

- Buscar cliente por email/phone (autocomplete contra `/api/v1/admin/users/search`)
- Si no existe → crear cliente rápido en el mismo flujo
- Agregar items desde catálogo (reutilizar componentes del carrito del cliente)
- Marcar como `pagado` + método (`efectivo` / `transferencia`)
- Notas internas (visibles solo para admins)

### 4. SEO (la razón de Nuxt)

`Idea.md` menciona Nuxt específicamente por SEO. Cosas a configurar:

- `useSeoMeta()` por página (title, description, og:image, twitter card)
- `<NuxtImage>` para optimización automática
- Sitemap dinámico (`@nuxtjs/sitemap`)
- `robots.txt` favorable
- Schema.org `Product` en páginas de detalle (rich results en Google)

## Convenciones

- **Idiomas**: código y comentarios en **inglés**. **UI strings en español** (la tienda es para público chileno, implícito en `Idea.md`). Si el alcance crece, evaluar `i18n` o composable `t()`.
- **Naming**: componentes `PascalCase.vue`, composables `useXxx.ts`, stores `useXxxStore.ts`, páginas `kebab-case.vue`.
- **HTTP**: nunca `fetch` raw a la API. Wrapper central en composable `useApi()`.
- **Errores**: si la API responde error, vista o toast **humanizado**. Nunca JSON crudo al usuario.
- **Auth state**: Pinia store `useAuthStore` con `user`, `isAuthenticated`, `isAdmin`, `viewMode`. Mutaciones solo desde `useAuth()`.

## Endpoints que el frontend consume

Pendiente generar `docs/api.md` con detalle. Resumen actual (sincronizar con `Backend/CLAUDE.md` cuando crezca):

- Auth: `/api/v1/auth/register`, `/login`, `/me`
- Catálogo público: `/api/v1/products`, `/products/:slug`, `/categories`
- Cliente: `/cart/*`, `/orders/*`
- Admin: `/api/v1/admin/products`, `/admin/users`, `/admin/orders/manual` (pendiente)

## Pendientes

### Bootstrap (antes de cualquier feature)

- [ ] Inicializar proyecto Nuxt (decidir entre scaffolding en `Frontend/` vacío o reestructurar)
- [ ] Elegir **librería de UI** — decisión bloqueante
- [ ] Confirmar TypeScript vs JavaScript
- [ ] Configurar ESLint + Prettier
- [ ] Crear `nuxt.config.ts` con módulos SEO + image

### Funcionalidades (de `Idea.md`)

- [ ] Composable `useApi()` + manejo de errores centralizado
- [ ] Auth flow (login/register/logout) sobre la decisión de storage
- [ ] Registro con **scroll-to-bottom** del doc de privacidad + persistencia de versión
- [ ] Catálogo (grid, filtros básicos, paginación)
- [ ] Detalle de producto con `useSeoMeta` + Schema.org `Product`
- [ ] Carrito (Pinia + persist en `localStorage` para invitados)
- [ ] Checkout (UI mínima; pago real delegado al backend)
- [ ] Layout y rutas del panel admin
- [ ] Toggle "vista admin → usuario normal"
- [ ] Creación manual de pedidos desde admin (buscar cliente + armar pedido)
- [ ] Página del documento de privacidad versionada en `public/legal/`
- [ ] Sitemap, robots.txt, meta tags base

### Diferidos (post-MVP)

- Búsqueda full-text, wishlist, reviews UI, notificaciones email, modo oscuro, i18n real.

## Estado del proyecto (2026-07-01)

**Carpeta `Frontend/` vacía.** No hay scaffolding aún.

### Bloqueantes para arrancar

1. **Librería de UI** — sin esto no se puede crear el primer componente.
   - Sugerencia: **Nuxt UI** (oficial del equipo Nuxt, integración nativa) o **PrimeVue** (más componentes out-of-the-box, mejor para paneles admin densos).
2. **Storage del JWT** (cookie httpOnly vs header Bearer) — bloquea `useAuth()`.
3. **Stack language** (TS confirmado vs JS).

Recomendación de orden: **(1) → alinear con backend sobre (2) → confirmar (3) → scaffolding.**

## Recordatorios importantes

- Package manager: **`pnpm`**. Nunca `npm install` ni `yarn`.
- Si aparece un error técnico en cualquier parte del frontend, delegar al agente `error-diagnosis-specialist`.
- Cualquier cambio en el modelo de dominio / endpoints se sincroniza con `Backend/CLAUDE.md` primero. Este archivo **no** redefine lo que el backend ya decidió.
