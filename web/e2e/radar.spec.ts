import { test, expect, type Page } from '@playwright/test';
import radarConfig from '../src/data/radar-config.json' with { type: 'json' };
import { QUADRANT_ORDER, quadrantBounds, toViewBox, FULL_BOUNDS } from '../src/config/radar';

// The suite runs against the site at "/" and, with BASE_PATH set, at a
// GitHub Pages-style sub-path. Navigate with relative paths (resolved against
// baseURL) and build expected hrefs with href().
const BASE = (process.env.BASE_PATH ?? '/').replace(/\/+$/, '');
const href = (path: string) => `${BASE}${path}`;

const quadrantLabels = radarConfig.quadrants.map((q) => q.label);
const ringLabels = radarConfig.rings.map((r) => r.label);

const radarSvg = (page: Page) => page.locator('.radar-page svg.radar-svg');

test.describe('Radar page', () => {
  test('renders the main heading and version', async ({ page }) => {
    await page.goto('radar');
    await expect(page.getByRole('heading', { name: 'Technology Radar', level: 1 })).toBeVisible();
    await expect(page.getByText(`Volume ${radarConfig.version}`)).toBeVisible();
  });

  test('renders the SVG radar visualization', async ({ page }) => {
    await page.goto('radar');
    await expect(radarSvg(page)).toBeVisible();
  });

  test('displays quadrant labels in the radar', async ({ page }) => {
    await page.goto('radar');
    for (const label of quadrantLabels) {
      await expect(radarSvg(page).getByText(label)).toBeVisible();
    }
  });

  test('displays ring labels in the radar', async ({ page }) => {
    await page.goto('radar');
    for (const label of ringLabels) {
      await expect(radarSvg(page).getByText(label).first()).toBeVisible();
    }
  });

  test('blip circles are links to detail pages', async ({ page }) => {
    await page.goto('radar');
    await expect(radarSvg(page).locator(`a[href^="${href('/radar/')}"]`).first()).toBeVisible();
  });

  test('hover tooltip is drawn above every blip', async ({ page }) => {
    await page.goto('radar');
    const blip = radarSvg(page).locator('a.blip-link').first();
    const title = (await blip.getAttribute('aria-label'))!.replace(/^\d+\.\s*/, '');
    await blip.hover();

    const label = radarSvg(page).locator('.radar-label');
    await expect(label).toHaveText(title);
    const isLast = await radarSvg(page).evaluate((svg) => svg.lastElementChild?.classList.contains('radar-label'));
    expect(isLast).toBe(true);
  });

  test('hovering a sidebar entry highlights its blip', async ({ page }) => {
    await page.goto('radar');
    await page.locator('.blip-lists li').first().hover();
    await expect(radarSvg(page).locator('.radar-label')).toBeVisible();
  });
});

test.describe('Quadrant zoom', () => {
  const quadrant = QUADRANT_ORDER[2];
  const label = quadrantLabels[2];

  test('clicking a quadrant title zooms in and updates the hash', async ({ page }) => {
    await page.goto('radar');
    await page.getByRole('button', { name: `Zoom to ${label}` }).click();
    await expect(radarSvg(page)).toHaveAttribute('viewBox', toViewBox(quadrantBounds(quadrant)));
    await expect(page).toHaveURL(new RegExp(`#${quadrant}$`));

    await page.getByRole('button', { name: /Back to full view/ }).click();
    await expect(radarSvg(page)).toHaveAttribute('viewBox', toViewBox(FULL_BOUNDS));
    await expect(page).not.toHaveURL(/#/);
  });

  test('a #quadrant link opens the radar zoomed in', async ({ page }) => {
    await page.goto(`radar#${quadrant}`);
    await expect(radarSvg(page)).toHaveAttribute('viewBox', toViewBox(quadrantBounds(quadrant)));
    await expect(page.getByRole('button', { name: /Back to full view/ })).toBeVisible();
    await expect(page.locator('.quadrant-section.dimmed')).toHaveCount(QUADRANT_ORDER.length - 1);
  });

  test('the zoomed radar does not paint over the sidebar', async ({ page }) => {
    await page.setViewportSize({ width: 1440, height: 900 });
    await page.goto(`radar#${quadrant}`);
    await expect(radarSvg(page)).toHaveCSS('overflow', 'hidden');

    // A square element means no letterbox space where neighbouring quadrants could show.
    const svgBox = (await radarSvg(page).boundingBox())!;
    expect(Math.abs(svgBox.width - svgBox.height)).toBeLessThan(1);

    // The sidebar's first heading must be the topmost element at its own position.
    const heading = page.locator('.quadrant-heading').first();
    const box = (await heading.boundingBox())!;
    const topmostIsHeading = await page.evaluate(
      ({ x, y }) => !!document.elementFromPoint(x, y)?.closest('.quadrant-heading'),
      { x: box.x + box.width / 2, y: box.y + box.height / 2 },
    );
    expect(topmostIsHeading).toBe(true);
  });

  test('quadrant titles work from the keyboard', async ({ page }) => {
    await page.goto('radar');
    await page.getByRole('button', { name: `Zoom to ${label}` }).focus();
    await page.keyboard.press('Enter');
    await expect(radarSvg(page)).toHaveAttribute('viewBox', toViewBox(quadrantBounds(quadrant)));
  });
});

test.describe('Radar fits the viewport', () => {
  for (const viewport of [
    { width: 1280, height: 720 },
    { width: 1440, height: 900 },
    { width: 1920, height: 1080 },
  ]) {
    test(`no page scroll at ${viewport.width}x${viewport.height}`, async ({ page }) => {
      await page.setViewportSize(viewport);
      await page.goto('radar');
      await expect(radarSvg(page)).toBeVisible();

      const { scrollHeight, innerHeight } = await page.evaluate(() => ({
        scrollHeight: document.documentElement.scrollHeight,
        innerHeight: window.innerHeight,
      }));
      expect(scrollHeight).toBeLessThanOrEqual(innerHeight);

      const box = (await radarSvg(page).boundingBox())!;
      expect(box.y).toBeGreaterThanOrEqual(0);
      expect(box.y + box.height).toBeLessThanOrEqual(viewport.height);
      // The radar should use the space, not shrink to a thumbnail.
      expect(Math.min(box.width, box.height)).toBeGreaterThan(viewport.height * 0.6);
    });
  }

  test('stacks and scrolls normally on a phone', async ({ page }) => {
    await page.setViewportSize({ width: 390, height: 844 });
    await page.goto('radar');
    const box = (await radarSvg(page).boundingBox())!;
    expect(box.width).toBeLessThanOrEqual(390);
    expect(box.width).toBeGreaterThan(300);

    const { scrollWidth, innerWidth, scrollHeight, innerHeight } = await page.evaluate(() => ({
      scrollWidth: document.documentElement.scrollWidth,
      innerWidth: window.innerWidth,
      scrollHeight: document.documentElement.scrollHeight,
      innerHeight: window.innerHeight,
    }));
    expect(scrollWidth).toBeLessThanOrEqual(innerWidth);
    expect(scrollHeight).toBeGreaterThan(innerHeight);
    await expect(page.locator('.blip-lists')).toBeAttached();
  });
});

test.describe('Navigation', () => {
  test('header contains navigation links', async ({ page }) => {
    await page.goto('radar');
    const header = page.locator('header.site-header');
    await expect(header.locator(`a[href="${href('/radar/search')}"]`)).toBeVisible();
    await expect(header.locator(`a[href="${href('/radar/faq')}"]`)).toBeVisible();
    await expect(header.locator(`a[href="${href('/radar/contribute')}"]`)).toBeVisible();
  });

  test('home redirects to the radar', async ({ page }) => {
    await page.goto('./');
    await page.waitForURL(/\/radar\/?$/);
    await expect(page.getByRole('heading', { name: 'Technology Radar', level: 1 })).toBeVisible();
  });

  test('unknown pages show the 404 page', async ({ page }) => {
    const response = await page.goto('radar/no-such-blip');
    expect(response?.status()).toBe(404);
    await expect(page.getByRole('heading', { name: 'Page not found' })).toBeVisible();
  });

  for (const path of ['radar', 'radar/search', 'radar/faq', 'radar/contribute', 'radar/platforms/hashicorp-nomad']) {
    test(`internal links on /${path} include the base path`, async ({ page }) => {
      await page.goto(path);
      const hrefs = await page.locator('a[href^="/"]').evaluateAll((els) => els.map((el) => el.getAttribute('href')!));
      expect(hrefs.length).toBeGreaterThan(0);
      for (const h of hrefs) {
        expect(h === BASE || h.startsWith(`${BASE}/`), `${h} should start with "${BASE}/"`).toBe(true);
      }
    });
  }
});

test.describe('Search page', () => {
  test('renders search input', async ({ page }) => {
    await page.goto('radar/search');
    await expect(page.getByRole('textbox', { name: 'Search blips' })).toBeVisible();
  });

  test('displays filter buttons', async ({ page }) => {
    await page.goto('radar/search');
    for (const q of quadrantLabels) {
      await expect(page.getByRole('button', { name: q })).toBeVisible();
    }
    for (const r of ringLabels) {
      await expect(page.getByRole('button', { name: r })).toBeVisible();
    }
  });

  test('filters results by search query', async ({ page }) => {
    await page.goto('radar/search');
    await page.getByRole('textbox', { name: 'Search blips' }).fill('Nomad');
    await expect(page.locator('.result-link').first()).toContainText('Nomad');
  });

  test('filters results by ring', async ({ page }) => {
    await page.goto('radar/search');
    const ringButton = page.getByRole('button', { name: ringLabels[0] });
    await ringButton.click();
    await expect(ringButton).toHaveAttribute('aria-pressed', 'true');

    const badges = page.locator('.result-item .ring-badge');
    await expect(badges.first()).toBeVisible();
    for (const text of await badges.allTextContents()) {
      expect(text.trim()).toBe(ringLabels[0]);
    }
  });
});

test.describe('Blip detail page', () => {
  test('renders blip title and metadata', async ({ page }) => {
    await page.goto('radar/platforms/hashicorp-nomad');
    await expect(page.getByRole('heading', { name: 'HashiCorp Nomad', level: 1 })).toBeVisible();
    await expect(page.locator('.ring-badge').first()).toBeVisible();
  });

  test('breadcrumb links back to the radar and its quadrant', async ({ page }) => {
    await page.goto('radar/platforms/hashicorp-nomad');
    const crumbs = page.getByRole('navigation', { name: 'Breadcrumb' });
    await expect(crumbs.locator(`a[href="${href('/radar')}"]`)).toBeVisible();

    await crumbs.locator(`a[href="${href('/radar#platforms')}"]`).click();
    await expect(radarSvg(page)).toHaveAttribute('viewBox', toViewBox(quadrantBounds('platforms')));
  });

  test('renders a static quadrant radar with the blip labelled', async ({ page }) => {
    await page.goto('radar/platforms/hashicorp-nomad');
    const sidebar = page.locator('.quadrant-radar-sidebar');
    await expect(sidebar.locator('svg').first()).toBeVisible();
    await expect(sidebar.locator('.radar-label')).toHaveText('HashiCorp Nomad');
    await expect(page.locator('astro-island')).toHaveCount(0);
  });

  test('sidebar list uses the quadrant color', async ({ page }) => {
    await page.goto('radar/platforms/hashicorp-nomad');
    const color = await page.locator('.quadrant-blip-list h3').evaluate((el) => getComputedStyle(el).color);
    expect(color).toBe('rgb(9, 105, 218)');
  });
});

test.describe('FAQ page', () => {
  test('renders FAQ heading', async ({ page }) => {
    await page.goto('radar/faq');
    await expect(page.getByRole('heading', { name: 'Frequently Asked Questions', level: 1 })).toBeVisible();
  });

  test('FAQ items are expandable', async ({ page }) => {
    await page.goto('radar/faq');
    const details = page.locator('details').first();
    await details.locator('summary').click();
    await expect(details).toHaveAttribute('open', '');
  });
});

test.describe('Contribute page', () => {
  test('renders contribute heading', async ({ page }) => {
    await page.goto('radar/contribute');
    await expect(page.getByRole('heading', { name: 'Contribute', level: 1 })).toBeVisible();
  });
});

test.describe('Theme toggle', () => {
  test('toggles dark mode', async ({ page }) => {
    await page.goto('radar');
    const toggleBtn = page.getByRole('button', { name: 'Toggle theme' });
    const html = page.locator('html');
    const wasDark = ((await html.getAttribute('class')) ?? '').includes('dark');

    await toggleBtn.click();
    await expect(html).toHaveClass(wasDark ? /^(?!.*dark)/ : /dark/);

    await toggleBtn.click();
    await expect(html).toHaveClass(wasDark ? /dark/ : /^(?!.*dark)/);
  });
});
