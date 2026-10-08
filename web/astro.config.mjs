import { defineConfig } from 'astro/config';
import vue from '@astrojs/vue';
import mdx from '@astrojs/mdx';

// SITE_URL and BASE_PATH are set by the GitHub Pages workflow (see
// .github/workflows/pages.yml); locally the site is served from "/".
export default defineConfig({
  site: process.env.SITE_URL || undefined,
  base: process.env.BASE_PATH || '/',
  integrations: [vue(), mdx()],
});
