import { test, expect } from '@playwright/test';
import type { BlipData } from '../src/types/radar';
import { QUADRANT_CONFIGS, QUADRANT_ORDER, RING_CONFIGS, RING_ORDER } from '../src/config/radar';
import { layoutBlips } from '../src/utils/radarLayout';
import { clampToQuadrant, polarToCartesian } from '../src/utils/radarMath';

// Pure layout checks (no browser): every blip must stay inside its own ring
// and quadrant, even when crowded. Quadrants that straddle the ±π wrap of
// atan2 used to snap blips to the opposite edge.

function crowdedBlips(perCell: number): BlipData[] {
  const blips: BlipData[] = [];
  for (const quadrant of QUADRANT_ORDER) {
    for (const ring of RING_ORDER) {
      for (let i = 0; i < perCell; i++) {
        const id = `${quadrant}-${ring}-${i}`;
        blips.push({ id, title: id, quadrant, ring, isNew: false });
      }
    }
  }
  return blips;
}

function angularOffset(x: number, y: number, mid: number): number {
  const a = Math.atan2(y, x);
  return Math.atan2(Math.sin(a - mid), Math.cos(a - mid));
}

test('clampToQuadrant keeps inside points and clamps outside points to the nearest edge', () => {
  const r = 300;
  const padding = 0.05;
  for (const quadrant of QUADRANT_ORDER) {
    const { startAngle, endAngle } = QUADRANT_CONFIGS[quadrant];

    for (const t of [0.1, 0.5, 0.9]) {
      const p = polarToCartesian(startAngle + t * (endAngle - startAngle), r);
      const c = clampToQuadrant(p.x, p.y, startAngle, endAngle, padding);
      expect(Math.hypot(c.x - p.x, c.y - p.y), `${quadrant} inside t=${t}`).toBeLessThan(1e-6);
    }

    for (const [label, angle, edge] of [
      ['before start', startAngle - 0.02, startAngle + padding],
      ['after end', endAngle + 0.02, endAngle - padding],
    ] as const) {
      const p = polarToCartesian(angle, r);
      const c = clampToQuadrant(p.x, p.y, startAngle, endAngle, padding);
      const want = polarToCartesian(edge, r);
      expect(Math.hypot(c.x - want.x, c.y - want.y), `${quadrant} ${label}`).toBeLessThan(1e-6);
    }
  }
});

test('crowded blips stay inside their ring and quadrant', () => {
  const blips = layoutBlips(crowdedBlips(8));

  for (const b of blips) {
    const ring = RING_CONFIGS[RING_ORDER.indexOf(b.ring)];
    const r = Math.hypot(b.x, b.y);
    expect(r, `${b.id} radius`).toBeGreaterThanOrEqual(ring.innerRadius);
    expect(r, `${b.id} radius`).toBeLessThanOrEqual(ring.outerRadius);

    const { startAngle, endAngle } = QUADRANT_CONFIGS[b.quadrant];
    const offset = angularOffset(b.x, b.y, (startAngle + endAngle) / 2);
    expect(Math.abs(offset), `${b.id} angle`).toBeLessThanOrEqual((endAngle - startAngle) / 2);
  }
});

test('blips are numbered 1..n per quadrant in ring order', () => {
  const blips = layoutBlips(crowdedBlips(2));

  for (const quadrant of QUADRANT_ORDER) {
    const sorted = blips.filter((b) => b.quadrant === quadrant).sort((a, b) => a.number - b.number);
    expect(sorted.map((b) => b.number)).toEqual(sorted.map((_, i) => i + 1));
    const ringIndexes = sorted.map((b) => RING_ORDER.indexOf(b.ring));
    expect(ringIndexes).toEqual([...ringIndexes].sort((a, b) => a - b));
  }
});

test('layout is deterministic', () => {
  const blips = crowdedBlips(3);
  expect(layoutBlips(blips)).toEqual(layoutBlips(blips));
});
