# Romchek Asteriska Frontend

Vue 3 + TypeScript + Tailwind CSS Single Page Application

## Tech Stack

- **Framework:** Vue 3.4+ with Composition API
- **Language:** TypeScript 5.6+ (strict mode)
- **Bundler:** Vite 5.1+
- **Styling:** Tailwind CSS 3.4+
- **Build Tools:** vue-tsc, PostCSS, Autoprefixer

## Project Structure

```
frontend/
├── src/
│   ├── assets/          # Static assets and global styles
│   │   └── main.css     # Tailwind CSS directives
│   ├── components/      # Vue components
│   │   └── HelloWorld.vue
│   ├── App.vue          # Root component
│   ├── main.ts          # Application entry point
│   └── vite-env.d.ts    # Vite type definitions
├── public/              # Public static assets
├── dist/                # Production build output (generated)
├── index.html           # HTML entry point
├── vite.config.ts       # Vite configuration
├── tsconfig.json        # TypeScript configuration
├── tailwind.config.js   # Tailwind CSS configuration
├── postcss.config.js    # PostCSS configuration
└── package.json         # Project dependencies and scripts
```

## Prerequisites

- Node.js 18+ or 20+
- npm 9+ or yarn 1.22+

## Installation

```bash
# Navigate to frontend directory
cd frontend

# Install dependencies
npm install
```

## Development

### Start Development Server

```bash
npm run dev
```

The application will be available at:
- Local: http://localhost:3000/
- Network: http://[your-ip]:3000/

The dev server includes:
- Hot Module Replacement (HMR)
- Fast refresh for Vue components
- TypeScript type checking in background
- Automatic browser reload

### Type Checking

```bash
# Run TypeScript type checker
npm run type-check
```

This runs `vue-tsc --noEmit` to check for type errors without emitting output files.

## Production Build

```bash
# Build for production
npm run build
```

This command:
1. Runs TypeScript type checking (`vue-tsc --noEmit`)
2. Builds optimized production assets (`vite build`)
3. Generates output to `dist/` directory

Build output includes:
- Minified and tree-shaken JavaScript
- Optimized CSS with Tailwind utilities
- Asset hashing for cache busting
- Gzip size estimates

### Preview Production Build

```bash
# Preview production build locally
npm run preview
```

Serves the `dist/` directory on a local server for testing production build.

## Configuration

### Vite Configuration (`vite.config.ts`)

- Vue plugin enabled for `.vue` file support
- Path alias `@/` mapped to `./src/`
- Development server on port 3000
- Host mode enabled for network access

### TypeScript Configuration (`tsconfig.json`)

- Strict mode enabled for type safety
- ES2020 target with ESNext modules
- Vue 3 JSX support
- Path aliases for clean imports
- Unused locals/parameters detection

### Tailwind CSS Configuration (`tailwind.config.js`)

- Content paths configured for Vue files
- Scans `index.html` and `src/**/*.{vue,js,ts,jsx,tsx}`
- Default theme with extend option for customization

## Component Development

### Example: HelloWorld Component

```vue
<script setup lang="ts">
// TypeScript props definition
defineProps<{
  msg: string
}>()
</script>

<template>
  <div class="bg-white rounded-2xl shadow-2xl p-8">
    <h1 class="text-4xl font-bold text-indigo-600">
      {{ msg }}
    </h1>
  </div>
</template>
```

Best practices:
- Use `<script setup lang="ts">` for Composition API
- Define props with TypeScript interfaces
- Use Tailwind utility classes for styling
- Keep components focused and composable

## Troubleshooting

### Common Issues

1. **Port 3000 already in use**
   ```bash
   # Change port in vite.config.ts or use environment variable
   PORT=3001 npm run dev
   ```

2. **Type errors in IDE but build succeeds**
   - Restart TypeScript server in your IDE
   - Check that IDE is using workspace TypeScript version

3. **Tailwind classes not applying**
   - Verify `main.css` is imported in `main.ts`
   - Check that file paths are included in `tailwind.config.js` content array
   - Clear browser cache and rebuild

4. **Module resolution errors**
   - Ensure `@/` alias is configured in both `tsconfig.json` and `vite.config.ts`
   - Run `npm install` to ensure dependencies are current

## Browser Support

- Modern browsers with ES2020 support
- Chrome/Edge 88+
- Firefox 78+
- Safari 14+

## Scripts Reference

| Script | Description |
|--------|-------------|
| `npm run dev` | Start development server with HMR |
| `npm run build` | Build production bundle |
| `npm run preview` | Preview production build locally |
| `npm run type-check` | Run TypeScript type checking |

## Next Steps

- Add Vue Router for navigation
- Set up Pinia for state management
- Configure API client for backend integration
- Add component testing with Vitest
- Set up ESLint and Prettier for code quality

## License

Private project - all rights reserved
