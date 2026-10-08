<script setup lang="ts">
import { computed } from 'vue';
import type { Quadrant } from '../../types/radar';
import { QUADRANT_CONFIGS, QUADRANT_ORDER, RING_CONFIGS, RADAR_RADIUS } from '../../config/radar';
import { generateRingArc } from '../../utils/radarMath';

// Ring arcs, axes, ring boundaries, and per-quadrant ring labels.
const props = withDefaults(
  defineProps<{
    quadrants?: Quadrant[];
    activeQuadrant?: Quadrant | null;
    interactive?: boolean;
  }>(),
  { quadrants: () => QUADRANT_ORDER, activeQuadrant: null, interactive: false },
);

const emit = defineEmits<{ select: [quadrant: Quadrant] }>();

const arcs = computed(() =>
  props.quadrants.flatMap((q) => {
    const { startAngle, endAngle } = QUADRANT_CONFIGS[q];
    return RING_CONFIGS.map((ring, i) => ({
      key: `${q}-${i}`,
      quadrant: q,
      path: generateRingArc(ring.innerRadius, ring.outerRadius, startAngle, endAngle),
    }));
  }),
);

// Ring labels sit at each ring's mid-radius along the quadrant's bisector.
const ringLabels = computed(() =>
  props.quadrants.flatMap((q) => {
    const { startAngle, endAngle } = QUADRANT_CONFIGS[q];
    const mid = (startAngle + endAngle) / 2;
    return RING_CONFIGS.map((ring, i) => {
      const r = (ring.innerRadius + ring.outerRadius) / 2;
      return { key: `${q}-${i}`, quadrant: q, label: ring.label, x: r * Math.cos(mid), y: r * Math.sin(mid) };
    });
  }),
);

function dimmed(q: Quadrant): boolean {
  return props.activeQuadrant !== null && props.activeQuadrant !== q;
}

function select(q: Quadrant) {
  if (props.interactive) emit('select', q);
}
</script>

<template>
  <g class="radar-backdrop" :class="{ interactive }" aria-hidden="true">
    <path
      v-for="arc in arcs"
      :key="arc.key"
      :d="arc.path"
      class="ring-arc"
      :class="{ dimmed: dimmed(arc.quadrant) }"
      @click="select(arc.quadrant)"
    />
    <line :x1="-RADAR_RADIUS" y1="0" :x2="RADAR_RADIUS" y2="0" class="radar-line" />
    <line x1="0" :y1="-RADAR_RADIUS" x2="0" :y2="RADAR_RADIUS" class="radar-line" />
    <circle v-for="ring in RING_CONFIGS" :key="ring.label" cx="0" cy="0" :r="ring.outerRadius" class="radar-line" />
    <text
      v-for="rl in ringLabels"
      :key="rl.key"
      :x="rl.x"
      :y="rl.y"
      class="ring-label"
      :class="{ dimmed: dimmed(rl.quadrant) }"
      text-anchor="middle"
      dominant-baseline="middle"
    >
      {{ rl.label }}
    </text>
  </g>
</template>

<style scoped>
.ring-arc {
  fill: var(--color-radar-ring-fill);
  stroke: var(--color-radar-ring-stroke);
  stroke-width: 1;
  transition: opacity 0.2s;
}

.interactive .ring-arc {
  cursor: pointer;
}

.interactive .ring-arc:hover {
  fill: var(--color-bg-hover);
}

.ring-arc.dimmed {
  opacity: 0.4;
}

.radar-line {
  fill: none;
  stroke: var(--color-radar-ring-stroke);
  stroke-width: 1;
}

.ring-label {
  font-size: 11px;
  fill: var(--color-text-muted);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  pointer-events: none;
  transition: opacity 0.2s;
}

.ring-label.dimmed {
  opacity: 0.3;
}
</style>
