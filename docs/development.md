# Development

We welcome contributions! **homectl** is built using a modern stack focused on performance.

## Tech Stack

- **Backend:** Go 1.24 (Fiber)
- **Frontend:** React 18 (Vite, TypeScript)
- **State Management:** Zustand & React Query
- **Styling:** Vanilla CSS (Sharp edges, monochrome)

## Local Setup

1. **Clone the repo:**
   ```bash
   git clone https://github.com/palta-dev/homectl.git
   ```

2. **Run Backend (Go):**
   ```bash
   cd apps/server
   go run ./cmd
   ```

3. **Run Frontend (React):**
   ```bash
   cd apps/web
   npm install
   npm run dev
   ```

The frontend dev server proxies API requests to `localhost:8080` by default.
