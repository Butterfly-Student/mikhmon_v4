# Mikhmon v4 Frontend

Frontend modern untuk Mikhmon Hotspot Management System dengan React + TypeScript + Tailwind CSS.

## 🎨 Features

- **Modern UI**: Desain colorful dengan gradient accents
- **Responsive**: Mobile-first approach dengan bottom navigation
- **Dark Mode**: Toggle antara light/dark mode
- **Real-time**: Auto-refresh data dashboard setiap 5 detik
- **Type Safe**: Full TypeScript support
- **Animations**: Smooth transitions dengan Framer Motion

## 🚀 Quick Start

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build
```

## 📁 Project Structure

```
src/
├── api/           # API functions
├── components/    # React components
│   ├── ui/        # Base UI components
│   ├── common/    # Shared components
│   ├── layout/    # Layout components
│   └── charts/    # Chart components
├── hooks/         # Custom React hooks
├── pages/         # Page components
├── stores/        # Zustand stores
├── styles/        # Global styles
├── types/         # TypeScript types
└── utils/         # Utility functions
```

## 🎨 Color Palette

| Color | Hex | Usage |
|-------|-----|-------|
| Primary | #4F46E5 | Main actions, navigation |
| Secondary | #EC4899 | Accents, highlights |
| Success | #10B981 | Success states, income |
| Warning | #F59E0B | Warnings |
| Danger | #EF4444 | Errors, delete |
| Info | #06B6D4 | Info, active users |

## 📱 Responsive Breakpoints

- Mobile: < 768px (bottom navigation)
- Tablet: 768px - 1024px
- Desktop: > 1024px (sidebar navigation)

## 🔧 Tech Stack

- React 18
- TypeScript 5
- Vite 5
- Tailwind CSS 3
- TanStack Query 5
- Zustand
- React Router 6
- Framer Motion
- Recharts

## 📝 Environment Variables

```env
VITE_API_URL=http://localhost:8080/api/v1
```

## 📄 License

MIT License
