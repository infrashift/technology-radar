import type { BlipData, PositionedBlip, Quadrant, Ring } from '../types/radar';
import { QUADRANT_CONFIGS, QUADRANT_ORDER, RING_CONFIGS, RING_ORDER } from '../config/radar';
import { polarToCartesian, hashString, hashString2, clampToRing, clampToQuadrant, distance } from './radarMath';
import { deriveBlipStatus } from './blipStatus';

const BLIP_RADIUS = 14;
const RING_PADDING = BLIP_RADIUS + 4;
const COLLISION_ITERATIONS = 50;

function ringConfig(ring: Ring) {
  return RING_CONFIGS[RING_ORDER.indexOf(ring)];
}

function initialPosition(blip: BlipData): { x: number; y: number } {
  const quadrant = QUADRANT_CONFIGS[blip.quadrant];
  const ring = ringConfig(blip.ring);

  // Deterministic seeds from the blip id keep positions stable across builds.
  const minR = ring.innerRadius + RING_PADDING;
  const maxR = ring.outerRadius - RING_PADDING;
  const radius = minR + hashString(blip.id) * (maxR - minR);

  const anglePadding = 0.08;
  const minAngle = quadrant.startAngle + anglePadding;
  const maxAngle = quadrant.endAngle - anglePadding;
  const angle = minAngle + hashString2(blip.id) * (maxAngle - minAngle);

  return polarToCartesian(angle, radius);
}

function clampBlip(b: PositionedBlip): void {
  const ring = ringConfig(b.ring);
  const quadrant = QUADRANT_CONFIGS[b.quadrant];
  const inRing = clampToRing(b.x, b.y, ring.innerRadius, ring.outerRadius, RING_PADDING);
  const final = clampToQuadrant(inRing.x, inRing.y, quadrant.startAngle, quadrant.endAngle);
  b.x = final.x;
  b.y = final.y;
}

function resolveCollisions(blips: PositionedBlip[]): void {
  const minDist = BLIP_RADIUS * 2.5;
  for (let iter = 0; iter < COLLISION_ITERATIONS; iter++) {
    let moved = false;
    for (let i = 0; i < blips.length; i++) {
      for (let j = i + 1; j < blips.length; j++) {
        const a = blips[i];
        const b = blips[j];
        const dist = distance(a.x, a.y, b.x, b.y);
        if (dist >= minDist || dist <= 0.001) continue;

        moved = true;
        const overlap = (minDist - dist) / 2;
        const dx = (b.x - a.x) / dist;
        const dy = (b.y - a.y) / dist;
        a.x -= dx * overlap;
        a.y -= dy * overlap;
        b.x += dx * overlap;
        b.y += dy * overlap;
        clampBlip(a);
        clampBlip(b);
      }
    }
    if (!moved) break;
  }
}

function byRingThenTitle(a: BlipData, b: BlipData): number {
  return RING_ORDER.indexOf(a.ring) - RING_ORDER.indexOf(b.ring) || a.title.localeCompare(b.title);
}

/**
 * Place blips on the radar: deterministic positions inside their ring and
 * quadrant, spread apart to avoid overlaps, then numbered per quadrant in
 * ring order and alphabetically. Runs at build time; the result is passed to
 * the radar components as props.
 */
export function layoutBlips(blips: BlipData[]): PositionedBlip[] {
  const positioned: PositionedBlip[] = blips.map((blip) => ({
    ...blip,
    ...initialPosition(blip),
    number: 0,
    status: deriveBlipStatus(blip),
  }));

  resolveCollisions(positioned);

  for (const q of QUADRANT_ORDER) {
    positioned
      .filter((b) => b.quadrant === q)
      .sort(byRingThenTitle)
      .forEach((b, i) => {
        b.number = i + 1;
      });
  }

  return positioned;
}

export interface RingGroup {
  ring: Ring;
  label: string;
  blips: PositionedBlip[];
}

/** Blips in one quadrant grouped by ring (empty rings omitted), in number order. */
export function groupByRing(blips: PositionedBlip[], quadrant: Quadrant): RingGroup[] {
  const inQuadrant = blips.filter((b) => b.quadrant === quadrant).sort((a, b) => a.number - b.number);
  return RING_ORDER.map((ring, i) => ({
    ring,
    label: RING_CONFIGS[i].label,
    blips: inQuadrant.filter((b) => b.ring === ring),
  })).filter((g) => g.blips.length > 0);
}
