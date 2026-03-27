import { test, expect } from '@playwright/test';

test.describe('Radar page', () => {
  test('renders the main heading', async ({ page }) => {
    await page.goto('/radar');
    await expect(page.getByRole('heading', { name: 'Technology Radar', level: 1 })).toBeVisible();
  });

  test('renders the SVG radar visualization', async ({ page }) => {
    await page.goto('/radar');
    await expect(page.locator('.radar-page svg').first()).toBeVisible();
  });

  test('displays quadrant labels in the radar', async ({ page }) => {
    await page.goto('/radar');
    const svg = page.locator('.radar-page svg').first();
    await expect(svg).toBeVisible();

    for (const label of ['Techniques', 'Platforms', 'Tools', 'Languages & Frameworks']) {
      await expect(svg.getByText(label)).toBeVisible();
    }
  });

  test('displays ring labels in the radar', async ({ page }) => {
    await page.goto('/radar');
    const svg = page.locator('.radar-page svg').first();
    await expect(svg).toBeVisible();

    for (const label of ['Adopt', 'Trial', 'Assess', 'Hold']) {
      await expect(svg.getByText(label).first()).toBeVisible();
    }
  });

  test('blip circles are links to detail pages', async ({ page }) => {
    await page.goto('/radar');
    const blipLink = page.locator('svg a[href^="/radar/"]').first();
    await expect(blipLink).toBeVisible();
  });
});

test.describe('Navigation', () => {
  test('header contains navigation links', async ({ page }) => {
    await page.goto('/radar');
    const header = page.locator('header').first();
    await expect(header).toBeVisible();
    await expect(header.locator('a[href="/radar/search"]')).toBeVisible();
    await expect(header.locator('a[href="/radar/faq"]')).toBeVisible();
    await expect(header.locator('a[href="/radar/contribute"]')).toBeVisible();
  });

  test('home redirects to /radar', async ({ page }) => {
    await page.goto('/');
    await page.waitForURL('/radar');
    await expect(page.getByRole('heading', { name: 'Technology Radar', level: 1 })).toBeVisible();
  });
});

test.describe('Search page', () => {
  test('renders search input', async ({ page }) => {
    await page.goto('/radar/search');
    await expect(page.locator('input[type="text"]')).toBeVisible();
  });

  test('displays filter buttons', async ({ page }) => {
    await page.goto('/radar/search');

    // Quadrant filters
    for (const q of ['Techniques', 'Platforms', 'Tools', 'Languages & Frameworks']) {
      await expect(page.getByRole('button', { name: q })).toBeVisible();
    }

    // Ring filters
    for (const r of ['adopt', 'trial', 'assess', 'hold']) {
      await expect(page.getByRole('button', { name: r })).toBeVisible();
    }
  });

  test('filters results by search query', async ({ page }) => {
    await page.goto('/radar/search');
    const input = page.locator('input[type="text"]');
    await input.fill('Kubernetes');

    // Should show at least one result containing "Kubernetes"
    const results = page.locator('.search-results a, .results a, [class*="result"] a');
    await expect(results.first()).toBeVisible();
  });

  test('filters results by ring', async ({ page }) => {
    await page.goto('/radar/search');

    // Click adopt ring filter
    await page.getByRole('button', { name: 'adopt' }).click();

    // Results should be visible (at least one blip is in adopt)
    const resultLinks = page.locator('a[href^="/radar/"]').filter({ hasNotText: 'Search' });
    await expect(resultLinks.first()).toBeVisible();
  });
});

test.describe('Blip detail page', () => {
  test('renders blip title and metadata', async ({ page }) => {
    await page.goto('/radar/platforms/kubernetes');

    await expect(page.getByRole('heading', { name: 'Kubernetes', level: 1 })).toBeVisible();
    // Should show ring badge
    await expect(page.locator('.ring-badge').first()).toBeVisible();
  });

  test('renders breadcrumb navigation', async ({ page }) => {
    await page.goto('/radar/platforms/kubernetes');

    await expect(page.locator('.breadcrumb a[href="/radar"]')).toBeVisible();
  });

  test('renders sidebar with quadrant radar', async ({ page }) => {
    await page.goto('/radar/platforms/kubernetes');
    await expect(page.locator('.quadrant-radar-sidebar svg').first()).toBeVisible();
  });
});

test.describe('FAQ page', () => {
  test('renders FAQ heading', async ({ page }) => {
    await page.goto('/radar/faq');
    await expect(page.getByRole('heading', { name: 'Frequently Asked Questions', level: 1 })).toBeVisible();
  });

  test('FAQ items are expandable', async ({ page }) => {
    await page.goto('/radar/faq');
    const details = page.locator('details').first();
    await expect(details).toBeVisible();

    // Click to expand
    await details.locator('summary').click();
    await expect(details).toHaveAttribute('open', '');
  });
});

test.describe('Contribute page', () => {
  test('renders contribute heading', async ({ page }) => {
    await page.goto('/radar/contribute');
    await expect(page.getByRole('heading', { name: 'Contribute', level: 1 })).toBeVisible();
  });
});

test.describe('Theme toggle', () => {
  test('toggles dark mode', async ({ page }) => {
    await page.goto('/radar');
    const toggleBtn = page.locator('header button').first();
    await expect(toggleBtn).toBeVisible();

    // Click to toggle theme
    await toggleBtn.click();
    const htmlClass = await page.locator('html').getAttribute('class');
    const isDark = htmlClass?.includes('dark');

    // Click again to toggle back
    await toggleBtn.click();
    const htmlClass2 = await page.locator('html').getAttribute('class');
    const isDark2 = htmlClass2?.includes('dark');

    expect(isDark).not.toBe(isDark2);
  });
});
