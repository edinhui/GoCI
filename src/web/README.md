# Configuration Management Component - Frontend

This is the frontend implementation for the Configuration Management Component, built with Vue 3 and Element Plus.

## Project Structure

```
src/web/
├── src/
│   ├── components/         # Reusable Vue components
│   │   ├── __tests__/      # Component unit tests
│   │   └── ...
│   ├── views/              # Page components
│   ├── router/             # Vue Router configuration
│   ├── App.vue             # Root component
│   └── main.js             # Application entry point
├── index.html              # HTML entry point
├── package.json            # Dependencies and scripts
└── vite.config.js          # Vite configuration
```

## Features Implemented

### JSON Schema Editor
- Basic layout with navigation and editor area
- Tree-based property visualization
- Property addition, editing, and deletion
- Schema preview and export functionality

## Setup Instructions

### Prerequisites
- Node.js (v14+)
- npm or yarn

### Installation

```bash
# Navigate to the project directory
cd src/web

# Install dependencies
npm install
# or
yarn install
```

### Development

```bash
# Start development server
npm run dev
# or
yarn dev
```

### Building for Production

```bash
# Build for production
npm run build
# or
yarn build
```

### Running Tests

```bash
# Run unit tests
npm run test:unit
# or
yarn test:unit
```

## Next Steps

1. Implement hierarchical structure creation (up to 3 levels)
2. Add validation rules configuration
3. Integrate with backend API for schema storage
4. Implement configuration editing based on schemas
