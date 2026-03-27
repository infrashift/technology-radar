import type { BlipData, Quadrant } from '../types/radar';
import { RING_ORDER } from '../types/radar';

export interface NumberedBlip extends BlipData {
  number: number;
}

/**
 * Assign sequential per-quadrant numbers to blips.
 * Ordered by ring (adopt first) then alphabetically by title.
 * This function is framework-agnostic and can be used in both Vue and Astro contexts.
 */
export function assignBlipNumbers(blips: BlipData[]): NumberedBlip[] {
  const result: NumberedBlip[] = blips.map((b) => ({ ...b, number: 0 }));

  const quadrants = [...new Set(blips.map((b) => b.quadrant))] as Quadrant[];

  for (const q of quadrants) {
    const quadrantBlips = result
      .filter((b) => b.quadrant === q)
      .sort((a, b) => {
        const ringDiff = RING_ORDER.indexOf(a.ring) - RING_ORDER.indexOf(b.ring);
        if (ringDiff !== 0) return ringDiff;
        return a.title.localeCompare(b.title);
      });
    quadrantBlips.forEach((b, i) => {
      b.number = i + 1;
    });
  }

  return result;
}
