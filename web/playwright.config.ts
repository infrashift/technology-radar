import { defineConfig } from '@playwright/test';

// BASE_PATH must match the one the site was built with (see astro.config.mjs).
const basePath = `/${(process.env.BASE_PATH ?? '').replace(/^\/+|\/+$/g, '')}/`.replace('//', '/');
const baseURL = `http://localhost:4321${basePath}`;

export default defineConfig({
  testDir: './e2e',
  timeout: 30_000,
  retries: 0,
  use: {
    baseURL,
    headless: true,
  },
  projects: [
    {
      name: 'chromium',
      use: { browserName: 'chromium' },
    },
  ],
  webServer: {
    // --ignore-lock keeps `astro preview` in the foreground even when Astro
    // detects a coding agent, which would otherwise background it.
    command: 'bun run preview --ignore-lock',
    url: baseURL,
    reuseExistingServer: !process.env.CI,
  },
});
