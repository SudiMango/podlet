# Podlet web

The public landing page for [podlet.dev](https://podlet.dev), built with Next.js.

## Development

```bash
npm install
npm run dev
```

Open [http://localhost:3000](http://localhost:3000).

## Structure

- `src/app` contains the page shell, global styles, and metadata.
- `src/components/layout` contains shared site chrome.
- `src/components/sections` contains page-level sections.
- `src/components/ui` contains small interactive components.

The landing page is intentionally lightweight. Only command-copy controls are
client components; the remaining content renders on the server.

## Checks

```bash
npm run lint
npm run build
```
