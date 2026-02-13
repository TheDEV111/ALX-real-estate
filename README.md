
# 🏠 ALX Real Estate (Airbnb Clone)

A high-performance, full-stack real estate marketplace built with **Go (Backend)** and **React 19 (Frontend)**. This project follows a modern monorepo-style structure and is designed for scalability and type safety.

## 🎨 Design Reference

The UI/UX for this project is based on the following Figma design:
**[Figma Design Link](https://www.figma.com/design/TPbjDkOGfymtb52xZPB9dV/Project-Airbnb--Copy-?node-id=1-4&p=f&t=ZPjC5TXo0c8PVFPE-0)**

---

## 🛠 Tech Stack

### Frontend 

* **Framework:** React 19 (using New Hooks & Actions)
* **Build Tool:** Vite + pnpm
* **Styling:** Tailwind CSS v4
* **Linting/Formatting:** Biome
* **State Management:** React 19 `use()` API & Context

### Backend 

* **Language:** Go 1.22+
* **API Style:** REST / JSON
* **Database:** PostgreSQL (via Railway/Docker)
* **CI/CD:** GitHub Actions (Automated Linting & Build Checks)

---

## 📂 Project Structure

```text
.
├── .github/workflows/   # CI/CD Automation (Linting & Builds)
├── backend/             # Go source code (API)
└── frontend/            # React 19 source code (Vite + Tailwind)

```

---

## 🚀 Getting Started

### Prerequisites

* [pnpm](https://pnpm.io/) (v9+)
* [Go](https://go.dev/) (v1.22+)
* [Node.js](https://nodejs.org/) (v20+)

### Setup

1. **Clone the repository:**
```bash
git clone https://github.com/TheDEV111/ALX-real-estate.git
cd ALX-real-estate

```


2. **Frontend Setup:**
```bash
cd frontend
pnpm install
pnpm run dev

```


3. **Backend Setup:**
```bash
cd backend
go run main.go

```



---

## 🤝 Collaboration Workflow (For Interns)

To maintain code quality, please follow these steps for every task:

1. **Pick an Issue:** Find a task assigned to you in the [GitHub Issues](https://www.google.com/search?q=https://github.com/TheDEV111/ALX-real-estate/issues) tab.
2. **Branching:** Create a new branch for your feature: `git checkout -b feat/your-feature-name`.
3. **Coding Standards:** Before pushing, ensure your code is formatted and linted:
* Run `pnpm biome check --write` in the frontend directory.




---

## 📝 Roadmap (MVP)

* [ ] **Auth:** User signup/login (Frontend + Go JWT)
* [ ] **Listings:** View properties from the Figma design.
* [ ] **Search:** Filter by category (Icons in Figma).
* [ ] **Details:** Individual property page with booking simulation.

---




