<script setup lang="ts">
import { computed, toRef } from 'vue';
import type { BlipData, Quadrant } from '../types/radar';
import {
  QUADRANT_CONFIGS,
  RING_CONFIGS,
} from '../types/radar';
import { generateRingArc, quarterArcPath } from '../utils/radarMath';
import { useRadarLayout } from '../composables/useRadarLayout';

const props = defineProps<{
  blips: BlipData[];
  quadrant: Quadrant;
  highlightedBlipId?: string;
}>();

const blipsRef = toRef(props, 'blips');
const { positionedBlips } = useRadarLayout(blipsRef);

const VIEWBOX_TARGETS: Record<Quadrant, string> = {
  techniques: '0 -500 500 500',
  platforms: '0 0 500 500',
  tools: '-500 0 500 500',
  'languages-and-frameworks': '-500 -500 500 500',
};

const quadrantViewBox = computed(() => VIEWBOX_TARGETS[props.quadrant]);

const quadrantBlips = computed(() =>
  positionedBlips.value.filter((b) => b.quadrant === props.quadrant)
);

const highlightedBlip = computed(() =>
  quadrantBlips.value.find((b) => b.id === props.highlightedBlipId)
);

// Ring arcs for this quadrant only
const ringArcs = computed(() => {
  const qConfig = QUADRANT_CONFIGS[props.quadrant];
  return RING_CONFIGS.map((ring, i) => ({
    path: generateRingArc(ring.innerRadius, ring.outerRadius, qConfig.startAngle, qConfig.endAngle),
    ringIndex: i,
  }));
});

// Per-quadrant ring labels
const ringLabels = computed(() => {
  const config = QUADRANT_CONFIGS[props.quadrant];
  const midAngle = (config.startAngle + config.endAngle) / 2;
  return RING_CONFIGS.map((ring) => {
    const midR = (ring.innerRadius + ring.outerRadius) / 2;
    return {
      label: ring.label,
      x: midR * Math.cos(midAngle),
      y: midR * Math.sin(midAngle),
    };
  });
});

function blipAngleFromCenter(x: number, y: number): number {
  return Math.atan2(y, x);
}

function movedArcPath(blip: { x: number; y: number; status: string }): string {
  const angle = blipAngleFromCenter(blip.x, blip.y);
  const centerAngle = blip.status === 'moved-in' ? angle + Math.PI : angle;
  return quarterArcPath(blip.x, blip.y, 14, centerAngle);
}
</script>

<template>
  <div class="quadrant-radar-wrapper">
    <svg
      :viewBox="quadrantViewBox"
      class="quadrant-radar-svg"
      xmlns="http://www.w3.org/2000/svg"
    >
      <!-- Ring backgrounds -->
      <path
        v-for="(arc, i) in ringArcs"
        :key="i"
        :d="arc.path"
        class="ring-arc"
        stroke-width="1"
      />

      <!-- Axis lines -->
      <line x1="-475" y1="0" x2="475" y2="0" class="axis-line" stroke-width="1" />
      <line x1="0" y1="-475" x2="0" y2="475" class="axis-line" stroke-width="1" />

      <!-- Ring boundary arcs -->
      <circle
        v-for="ring in RING_CONFIGS"
        :key="ring.label"
        cx="0"
        cy="0"
        :r="ring.outerRadius"
        fill="none"
        class="ring-boundary"
        stroke-width="1"
      />

      <!-- Ring labels -->
      <text
        v-for="rl in ringLabels"
        :key="rl.label"
        :x="rl.x"
        :y="rl.y"
        class="ring-label"
        text-anchor="middle"
        dominant-baseline="middle"
      >
        {{ rl.label }}
      </text>

      <!-- Blips -->
      <a
        v-for="blip in quadrantBlips"
        :key="blip.id"
        :href="`/radar/${blip.quadrant}/${blip.id}`"
        class="blip-link"
      >
        <g
          class="blip-group"
          :class="{ highlighted: blip.id === highlightedBlipId }"
        >
          <!-- Filled circle base -->
          <circle
            :cx="blip.x"
            :cy="blip.y"
            :r="blip.id === highlightedBlipId ? 12 : 10"
            :fill="QUADRANT_CONFIGS[blip.quadrant].color"
            class="blip-dot"
          />

          <!-- New: outer ring -->
          <circle
            v-if="blip.status === 'new'"
            :cx="blip.x"
            :cy="blip.y"
            :r="blip.id === highlightedBlipId ? 16 : 14"
            fill="none"
            :stroke="QUADRANT_CONFIGS[blip.quadrant].color"
            stroke-width="2"
            stroke-opacity="0.7"
          />

          <!-- Moved in/out: quarter arc -->
          <path
            v-if="blip.status === 'moved-in' || blip.status === 'moved-out'"
            :d="movedArcPath(blip)"
            fill="none"
            :stroke="QUADRANT_CONFIGS[blip.quadrant].color"
            stroke-width="2.5"
            stroke-opacity="0.7"
            stroke-linecap="round"
          />

          <!-- Number -->
          <text
            :x="blip.x"
            :y="blip.y"
            fill="#fff"
            :font-size="blip.id === highlightedBlipId ? 11 : 9"
            font-weight="700"
            text-anchor="middle"
            dominant-baseline="central"
            class="blip-number"
          >
            {{ blip.number }}
          </text>

          <!-- Highlight ring for current blip -->
          <circle
            v-if="blip.id === highlightedBlipId"
            :cx="blip.x"
            :cy="blip.y"
            r="18"
            fill="none"
            class="blip-highlight-ring"
            stroke-width="2.5"
          />
        </g>
      </a>

      <!-- Highlighted blip label (always visible) -->
      <g v-if="highlightedBlip" class="highlighted-label">
        <rect
          :x="highlightedBlip.x + 22"
          :y="highlightedBlip.y - 14"
          :width="highlightedBlip.title.length * 7 + 16"
          height="24"
          rx="4"
          class="tooltip-bg"
          fill-opacity="0.9"
        />
        <text
          :x="highlightedBlip.x + 30"
          :y="highlightedBlip.y + 1"
          class="tooltip-text"
          font-size="11"
          font-weight="600"
          dominant-baseline="middle"
        >
          {{ highlightedBlip.title }}
        </text>
      </g>
    </svg>
  </div>
</template>

<style scoped>
.quadrant-radar-wrapper {
  width: 100%;
}

.quadrant-radar-svg {
  width: 100%;
  height: auto;
}

.ring-arc {
  fill: var(--color-radar-ring-fill);
  stroke: var(--color-radar-ring-stroke);
}

.axis-line {
  stroke: var(--color-radar-ring-stroke);
}

.ring-boundary {
  stroke: var(--color-radar-ring-stroke);
}

.ring-label {
  font-size: 10px;
  fill: var(--color-text-muted);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  pointer-events: none;
}

.blip-highlight-ring {
  stroke: var(--color-radar-highlight-stroke);
}

.tooltip-bg {
  fill: var(--color-radar-tooltip-bg);
}

.tooltip-text {
  fill: var(--color-radar-tooltip-text);
}

.blip-link {
  text-decoration: none;
}

.blip-group {
  cursor: pointer;
  transition: opacity 0.2s;
}

.blip-group.highlighted .blip-dot {
  filter: drop-shadow(0 0 6px rgba(0, 0, 0, 0.4));
}

.blip-number {
  pointer-events: none;
  user-select: none;
}
</style>
