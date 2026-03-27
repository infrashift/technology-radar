<script setup lang="ts">
import { ref, computed, toRef } from 'vue';
import type { BlipData, Quadrant, Ring } from '../types/radar';
import {
  QUADRANT_CONFIGS,
  RING_CONFIGS,
  RING_ORDER,
  RADAR_SIZE,
  RADAR_RADIUS,
} from '../types/radar';
import { generateRingArc, quarterArcPath } from '../utils/radarMath';
import { useRadarLayout } from '../composables/useRadarLayout';
import { useViewBoxTween } from '../composables/useViewBoxTween';
import { useHoverSync } from '../composables/useHoverSync';

const props = defineProps<{
  blips: BlipData[];
}>();

const blipsRef = toRef(props, 'blips');
const { positionedBlips, blipsByQuadrant } = useRadarLayout(blipsRef);
const { hoveredBlipId, setHover } = useHoverSync();

const DEFAULT_VIEWBOX = `-500 -500 ${RADAR_SIZE} ${RADAR_SIZE}`;
const { currentViewBox, tweenTo } = useViewBoxTween(DEFAULT_VIEWBOX);

const activeQuadrant = ref<Quadrant | null>(null);

const VIEWBOX_TARGETS: Record<string, string> = {
  techniques: '0 -500 500 500',
  platforms: '0 0 500 500',
  tools: '-500 0 500 500',
  'languages-and-frameworks': '-500 -500 500 500',
};

function handleQuadrantClick(quadrant: Quadrant) {
  if (activeQuadrant.value === quadrant) {
    activeQuadrant.value = null;
    tweenTo(DEFAULT_VIEWBOX);
  } else {
    activeQuadrant.value = quadrant;
    tweenTo(VIEWBOX_TARGETS[quadrant]);
  }
}

function handleBackClick() {
  activeQuadrant.value = null;
  tweenTo(DEFAULT_VIEWBOX);
}

// Generate ring arc paths for each quadrant (white backgrounds)
const ringArcs = computed(() => {
  const arcs: { path: string; quadrant: Quadrant; ringIndex: number }[] = [];
  const quadrants = Object.keys(QUADRANT_CONFIGS) as Quadrant[];

  for (const q of quadrants) {
    const qConfig = QUADRANT_CONFIGS[q];
    for (let i = 0; i < RING_CONFIGS.length; i++) {
      const ring = RING_CONFIGS[i];
      arcs.push({
        path: generateRingArc(ring.innerRadius, ring.outerRadius, qConfig.startAngle, qConfig.endAngle),
        quadrant: q,
        ringIndex: i,
      });
    }
  }

  return arcs;
});

// Quadrant title label positions — placed OUTSIDE the radar rings
const quadrantTitleLabels = computed(() => {
  return (Object.keys(QUADRANT_CONFIGS) as Quadrant[]).map((q) => {
    const config = QUADRANT_CONFIGS[q];
    const midAngle = (config.startAngle + config.endAngle) / 2;
    const labelRadius = RADAR_RADIUS + 35;
    return {
      quadrant: q,
      label: config.label,
      x: labelRadius * Math.cos(midAngle),
      y: labelRadius * Math.sin(midAngle),
      color: config.color,
    };
  });
});

// Per-quadrant ring labels — placed at the midpoint of each ring along the quadrant's bisecting axis
const perQuadrantRingLabels = computed(() => {
  const labels: { label: string; x: number; y: number; quadrant: Quadrant }[] = [];
  const quadrants = Object.keys(QUADRANT_CONFIGS) as Quadrant[];

  for (const q of quadrants) {
    const config = QUADRANT_CONFIGS[q];
    const midAngle = (config.startAngle + config.endAngle) / 2;

    for (const ring of RING_CONFIGS) {
      const midR = (ring.innerRadius + ring.outerRadius) / 2;
      labels.push({
        label: ring.label,
        x: midR * Math.cos(midAngle),
        y: midR * Math.sin(midAngle),
        quadrant: q,
      });
    }
  }

  return labels;
});

const quadrantOrder: Quadrant[] = ['techniques', 'platforms', 'tools', 'languages-and-frameworks'];

// Compute the angle from center to blip (for moved-in/out arc direction)
function blipAngleFromCenter(x: number, y: number): number {
  return Math.atan2(y, x);
}

// Quarter arc path for moved blips
function movedArcPath(blip: { x: number; y: number; status: string }): string {
  const angle = blipAngleFromCenter(blip.x, blip.y);
  // moved-in: arc faces toward center (opposite direction from blip)
  // moved-out: arc faces away from center (same direction as blip)
  const centerAngle = blip.status === 'moved-in' ? angle + Math.PI : angle;
  return quarterArcPath(blip.x, blip.y, 14, centerAngle);
}
</script>

<template>
  <div class="radar-container">
    <div class="radar-svg-wrapper">
      <svg
        :viewBox="currentViewBox"
        class="radar-svg"
        xmlns="http://www.w3.org/2000/svg"
      >
        <!-- Ring backgrounds — white with subtle borders -->
        <g class="ring-backgrounds">
          <path
            v-for="(arc, i) in ringArcs"
            :key="i"
            :d="arc.path"
            stroke-width="1"
            class="ring-arc"
            :class="{ dimmed: activeQuadrant && activeQuadrant !== arc.quadrant }"
            @click="handleQuadrantClick(arc.quadrant)"
          />
        </g>

        <!-- Axis lines -->
        <line x1="-475" y1="0" x2="475" y2="0" class="axis-line" stroke-width="1" />
        <line x1="0" y1="-475" x2="0" y2="475" class="axis-line" stroke-width="1" />

        <!-- Ring boundary circles -->
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

        <!-- Per-quadrant ring labels -->
        <g class="ring-labels">
          <text
            v-for="rl in perQuadrantRingLabels"
            :key="`${rl.quadrant}-${rl.label}`"
            :x="rl.x"
            :y="rl.y"
            class="ring-label"
            :class="{ dimmed: activeQuadrant && activeQuadrant !== rl.quadrant }"
            text-anchor="middle"
            dominant-baseline="middle"
          >
            {{ rl.label }}
          </text>
        </g>

        <!-- Blips -->
        <g class="blips">
          <a
            v-for="blip in positionedBlips"
            :key="blip.id"
            :href="`/radar/${blip.quadrant}/${blip.id}`"
            class="blip-link"
          >
            <g
              class="blip-group"
              :class="{
                highlighted: hoveredBlipId === blip.id,
                dimmed: activeQuadrant && activeQuadrant !== blip.quadrant,
              }"
              @mouseenter="setHover(blip.id)"
              @mouseleave="setHover(null)"
            >
              <!-- All blips have a filled circle base -->
              <circle
                :cx="blip.x"
                :cy="blip.y"
                r="10"
                :fill="QUADRANT_CONFIGS[blip.quadrant].color"
                :fill-opacity="hoveredBlipId === blip.id ? 1 : 0.85"
                class="blip-dot"
              />

              <!-- New: outer ring around the circle -->
              <circle
                v-if="blip.status === 'new'"
                :cx="blip.x"
                :cy="blip.y"
                r="14"
                fill="none"
                :stroke="QUADRANT_CONFIGS[blip.quadrant].color"
                stroke-width="2"
                :stroke-opacity="hoveredBlipId === blip.id ? 1 : 0.7"
                class="blip-ring"
              />

              <!-- Moved in/out: quarter arc indicator -->
              <path
                v-if="blip.status === 'moved-in' || blip.status === 'moved-out'"
                :d="movedArcPath(blip)"
                fill="none"
                :stroke="QUADRANT_CONFIGS[blip.quadrant].color"
                stroke-width="2.5"
                :stroke-opacity="hoveredBlipId === blip.id ? 1 : 0.7"
                stroke-linecap="round"
                class="blip-arc"
              />

              <!-- Number text -->
              <text
                :x="blip.x"
                :y="blip.y"
                fill="#fff"
                font-size="9"
                font-weight="700"
                text-anchor="middle"
                dominant-baseline="central"
                class="blip-number"
              >
                {{ blip.number }}
              </text>

              <!-- Hover highlight ring -->
              <circle
                v-if="hoveredBlipId === blip.id"
                :cx="blip.x"
                :cy="blip.y"
                r="12"
                fill="none"
                stroke-width="2"
                class="blip-hover-ring"
              />

              <!-- Tooltip on hover -->
              <g
                v-if="hoveredBlipId === blip.id"
                class="blip-tooltip"
              >
                <rect
                  :x="blip.x + 16"
                  :y="blip.y - 14"
                  :width="blip.title.length * 7 + 16"
                  height="24"
                  rx="4"
                  class="tooltip-bg"
                  fill-opacity="0.9"
                />
                <text
                  :x="blip.x + 24"
                  :y="blip.y + 1"
                  class="tooltip-text"
                  font-size="11"
                  dominant-baseline="middle"
                >
                  {{ blip.title }}
                </text>
              </g>
            </g>
          </a>
        </g>

        <!-- Quadrant title labels — outside the rings on light gray background -->
        <g
          v-for="ql in quadrantTitleLabels"
          :key="ql.quadrant"
          class="quadrant-title-group"
          :class="{ dimmed: activeQuadrant && activeQuadrant !== ql.quadrant }"
          @click="handleQuadrantClick(ql.quadrant)"
        >
          <rect
            :x="ql.x - 80"
            :y="ql.y - 12"
            width="160"
            height="24"
            rx="4"
            class="quadrant-title-bg"
          />
          <text
            :x="ql.x"
            :y="ql.y"
            :fill="ql.color"
            class="quadrant-title-text"
            text-anchor="middle"
            dominant-baseline="middle"
          >
            {{ ql.label }}
          </text>
        </g>
      </svg>

      <!-- Legend -->
      <div class="radar-legend">
        <span class="legend-item">
          <svg width="20" height="20" viewBox="0 0 20 20">
            <circle cx="10" cy="10" r="5" fill="currentColor" />
            <circle cx="10" cy="10" r="8" fill="none" stroke="currentColor" stroke-width="1.5" />
          </svg>
          New
        </span>
        <span class="legend-item">
          <svg width="20" height="20" viewBox="0 0 20 20">
            <circle cx="10" cy="10" r="5" fill="currentColor" />
            <path d="M 10 2 A 8 8 0 0 1 18 10" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
          </svg>
          Moved in/out
        </span>
        <span class="legend-item">
          <svg width="20" height="20" viewBox="0 0 20 20">
            <circle cx="10" cy="10" r="5" fill="currentColor" />
          </svg>
          No change
        </span>
      </div>

      <!-- Back button when zoomed -->
      <button
        v-if="activeQuadrant"
        class="back-button"
        @click="handleBackClick"
      >
        &larr; Back to full view
      </button>
    </div>

    <!-- Blip lists sidebar -->
    <div class="blip-lists">
      <div
        v-for="q in quadrantOrder"
        :key="q"
        class="quadrant-section"
        :class="{ active: activeQuadrant === q, dimmed: activeQuadrant && activeQuadrant !== q }"
      >
        <h3
          :style="{ color: QUADRANT_CONFIGS[q].color }"
          @click="handleQuadrantClick(q)"
        >
          {{ QUADRANT_CONFIGS[q].label }}
        </h3>

        <div
          v-for="ring in RING_ORDER"
          :key="ring"
          class="ring-group"
        >
          <template v-if="blipsByQuadrant[q][ring].length > 0">
            <h4>{{ ring }}</h4>
            <ul>
              <li
                v-for="blip in blipsByQuadrant[q][ring]"
                :key="blip.id"
                :class="{ highlighted: hoveredBlipId === blip.id }"
                @mouseenter="setHover(blip.id)"
                @mouseleave="setHover(null)"
              >
                <a :href="`/radar/${blip.quadrant}/${blip.id}`">
                  <span
                    class="blip-number-badge"
                    :style="{ background: QUADRANT_CONFIGS[q].color }"
                  >
                    {{ blip.number }}
                  </span>
                  {{ blip.title }}
                  <span class="blip-status-indicator">
                    <svg v-if="blip.status === 'new'" width="12" height="12" viewBox="0 0 12 12">
                      <circle cx="6" cy="6" r="3" :fill="QUADRANT_CONFIGS[q].color" />
                      <circle cx="6" cy="6" r="5" fill="none" :stroke="QUADRANT_CONFIGS[q].color" stroke-width="1" />
                    </svg>
                    <svg v-else-if="blip.status === 'moved-in' || blip.status === 'moved-out'" width="12" height="12" viewBox="0 0 12 12">
                      <circle cx="6" cy="6" r="3" :fill="QUADRANT_CONFIGS[q].color" />
                      <path d="M 6 1 A 5 5 0 0 1 11 6" fill="none" :stroke="QUADRANT_CONFIGS[q].color" stroke-width="1" stroke-linecap="round" />
                    </svg>
                  </span>
                </a>
              </li>
            </ul>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.radar-container {
  display: grid;
  grid-template-columns: 1fr 360px;
  gap: 2rem;
  align-items: start;
}

@media (max-width: 1024px) {
  .radar-container {
    grid-template-columns: 1fr;
  }
}

.radar-svg-wrapper {
  position: relative;
}

.radar-svg {
  width: 100%;
  height: auto;
  max-height: 80vh;
  cursor: default;
}

.ring-arc {
  fill: var(--color-radar-ring-fill);
  stroke: var(--color-radar-ring-stroke);
  cursor: pointer;
  transition: opacity 0.2s;
}

.ring-arc:hover {
  fill: var(--color-bg-hover) !important;
}

.ring-arc.dimmed {
  opacity: 0.4;
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
  transition: opacity 0.2s;
}

.ring-label.dimmed {
  opacity: 0.3;
}

.blip-link {
  text-decoration: none;
}

.blip-group {
  cursor: pointer;
  transition: opacity 0.2s;
}

.blip-group.dimmed {
  opacity: 0.15;
}

.blip-group.highlighted .blip-dot {
  filter: drop-shadow(0 0 4px rgba(0, 0, 0, 0.3));
}

.blip-dot {
  transition: fill-opacity 0.15s;
}

.blip-number {
  pointer-events: none;
  user-select: none;
}

.quadrant-title-group {
  cursor: pointer;
  transition: opacity 0.2s;
}

.quadrant-title-group.dimmed {
  opacity: 0.3;
}

.quadrant-title-text {
  font-size: 12px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.quadrant-title-bg {
  fill: var(--color-radar-label-bg);
  transition: fill 0.15s;
}

.quadrant-title-group:hover .quadrant-title-bg {
  fill: var(--color-radar-label-bg-hover);
}

.radar-legend {
  display: flex;
  justify-content: center;
  gap: 1.5rem;
  margin-top: 0.75rem;
  font-size: 0.8rem;
  color: var(--color-radar-legend);
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 0.35rem;
}

.blip-hover-ring {
  stroke: var(--color-radar-highlight-stroke);
}

.tooltip-bg {
  fill: var(--color-radar-tooltip-bg);
}

.tooltip-text {
  fill: var(--color-radar-tooltip-text);
}

.back-button {
  position: absolute;
  top: 1rem;
  left: 1rem;
  background: var(--color-button-bg);
  color: var(--color-text-on-dark);
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 0.5rem;
  font-size: 0.85rem;
  cursor: pointer;
  font-family: inherit;
}

.back-button:hover {
  background: var(--color-button-bg-hover);
}

/* Sidebar lists */
.blip-lists {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.quadrant-section {
  transition: opacity 0.2s;
}

.quadrant-section.dimmed {
  opacity: 0.3;
}

.quadrant-section h3 {
  font-size: 1rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 0.5rem;
  cursor: pointer;
  border: 1.5px solid currentColor;
  border-radius: 0.375rem;
  padding: 0.35rem 0.6rem;
  text-align: left;
}

.ring-group + .ring-group {
  border-top: 1px solid var(--color-border-light);
  margin-top: 0.25rem;
  padding-top: 0.25rem;
}

.ring-group h4 {
  font-size: 0.75rem;
  color: var(--color-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.08em;
  margin: 0.5rem 0 0.25rem;
  font-weight: 600;
}

.ring-group ul {
  list-style: none;
  margin: 0;
  padding: 0;
}

.ring-group li {
  padding: 0.25rem 0.5rem;
  border-radius: 0.25rem;
  transition: background 0.15s;
  font-size: 0.875rem;
}

.ring-group li.highlighted {
  background: var(--color-highlight-bg);
}

.ring-group li a {
  color: var(--color-text-primary);
  text-decoration: none;
  display: flex;
  align-items: center;
  gap: 0.4rem;
}

.ring-group li a:hover {
  text-decoration: underline;
}

.blip-number-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  color: #fff;
  font-size: 0.65rem;
  font-weight: 700;
  flex-shrink: 0;
}

.blip-status-indicator {
  margin-left: auto;
  display: flex;
  align-items: center;
  flex-shrink: 0;
}
</style>
