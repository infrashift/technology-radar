import type { BlipData } from '../types/radar';

/**
 * Prefix a site-absolute path with the configured base (e.g. "/technology-radar"
 * on a GitHub Pages project site), so links work wherever the site is hosted.
 */
export function withBase(path: string): string {
  const base = import.meta.env.BASE_URL.replace(/\/+$/, '');
  return `${base}/${path.replace(/^\/+/, '')}`;
}

export function blipHref(blip: Pick<BlipData, 'quadrant' | 'id'>): string {
  return withBase(`/radar/${blip.quadrant}/${blip.id}`);
}
