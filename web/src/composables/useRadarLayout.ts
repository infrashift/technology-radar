import { computed, type Ref } from 'vue';
import type { BlipData, PositionedBlip, Quadrant, Ring } from '../types/radar';
import {
  QUADRANT_CONFIGS,
  RING_CONFIGS,
  RING_ORDER,
} from '../types/radar';
import {
  polarToCartesian,
  hashString,
  hashString2,
  clampToRing,
  clampToQuadrant,
  distance,
} from '../utils/radarMath';
import { deriveBlipStatus } from '../utils/blipStatus';

const BLIP_RADIUS = 14;
const COLLISION_ITERATIONS = 50;

function getRingConfig(ring: Ring) {
  const idx = RING_ORDER.indexOf(ring);
  return RING_CONFIGS[idx];
}

function computeInitialPosition(blip: BlipData): { x: number; y: number } {
  const quadrant = QUADRANT_CONFIGS[blip.quadrant];
  const ring = getRingConfig(blip.ring);

  // Deterministic seed from blip id
  const radialSeed = hashString(blip.id);
  const angularSeed = hashString2(blip.id);

  // Map seed to radius within ring bounds (with padding)
  const padding = BLIP_RADIUS + 4;
  const minR = ring.innerRadius + padding;
  const maxR = ring.outerRadius - padding;
  const radius = minR + radialSeed * (maxR - minR);

  // Map seed to angle within quadrant bounds (with padding)
  const anglePadding = 0.08;
  const minAngle = quadrant.startAngle + anglePadding;
  const maxAngle = quadrant.endAngle - anglePadding;
  const angle = minAngle + angularSeed * (maxAngle - minAngle);

  return polarToCartesian(angle, radius);
}

function resolveCollisions(blips: PositionedBlip[]): void {
  for (let iter = 0; iter < COLLISION_ITERATIONS; iter++) {
    let moved = false;
    for (let i = 0; i < blips.length; i++) {
      for (let j = i + 1; j < blips.length; j++) {
        const dist = distance(blips[i].x, blips[i].y, blips[j].x, blips[j].y);
        const minDist = BLIP_RADIUS * 2.5;

        if (dist < minDist && dist > 0.001) {
          moved = true;
          const overlap = (minDist - dist) / 2;
          const dx = (blips[j].x - blips[i].x) / dist;
          const dy = (blips[j].y - blips[i].y) / dist;

          blips[i].x -= dx * overlap;
          blips[i].y -= dy * overlap;
          blips[j].x += dx * overlap;
          blips[j].y += dy * overlap;

          // Clamp both back into their ring/quadrant
          for (const b of [blips[i], blips[j]]) {
            const ring = getRingConfig(b.ring);
            const quadrant = QUADRANT_CONFIGS[b.quadrant];
            const clamped = clampToRing(b.x, b.y, ring.innerRadius, ring.outerRadius, BLIP_RADIUS + 4);
            const final = clampToQuadrant(clamped.x, clamped.y, quadrant.startAngle, quadrant.endAngle);
            b.x = final.x;
            b.y = final.y;
          }
        }
      }
    }
    if (!moved) break;
  }
}

export function useRadarLayout(blips: Ref<BlipData[]>) {
  const positionedBlips = computed<PositionedBlip[]>(() => {
    const positioned: PositionedBlip[] = blips.value.map((blip) => {
      const pos = computeInitialPosition(blip);
      return { ...blip, x: pos.x, y: pos.y, number: 0, status: deriveBlipStatus(blip) };
    });

    resolveCollisions(positioned);

    // Assign per-quadrant sequential numbers (ring order, then alphabetical)
    for (const q of Object.keys(QUADRANT_CONFIGS) as Quadrant[]) {
      const quadrantBlips = positioned
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

    return positioned;
  });

  const blipsByQuadrant = computed(() => {
    const grouped: Record<Quadrant, Record<Ring, PositionedBlip[]>> = {
      techniques: { adopt: [], trial: [], assess: [], hold: [] },
      platforms: { adopt: [], trial: [], assess: [], hold: [] },
      tools: { adopt: [], trial: [], assess: [], hold: [] },
      'languages-and-frameworks': { adopt: [], trial: [], assess: [], hold: [] },
    };

    for (const blip of positionedBlips.value) {
      grouped[blip.quadrant][blip.ring].push(blip);
    }

    return grouped;
  });

  return { positionedBlips, blipsByQuadrant };
}
