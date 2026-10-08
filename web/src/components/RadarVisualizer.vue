<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import type { PositionedBlip, Quadrant } from '../types/radar';
import {
  QUADRANT_CONFIGS,
  QUADRANT_ORDER,
  RADAR_RADIUS,
  FULL_BOUNDS,
  quadrantBounds,
  toViewBox,
} from '../config/radar';
import { groupByRing } from '../utils/radarLayout';
import { blipHref } from '../utils/url';
import { useViewBoxTween } from '../composables/useViewBoxTween';
import RadarBackdrop from './radar/RadarBackdrop.vue';
import RadarBlip from './radar/RadarBlip.vue';
import RadarLabel from './radar/RadarLabel.vue';

const props = defineProps<{
  blips: PositionedBlip[];
}>();

const hoveredId = ref<string | null>(null);
const hoveredBlip = computed(() => props.blips.find((b) => b.id === hoveredId.value));

const activeQuadrant = ref<Quadrant | null>(null);
const viewBounds = computed(() => (activeQuadrant.value ? quadrantBounds(activeQuadrant.value) : FULL_BOUNDS));
const { currentViewBox, tweenTo } = useViewBoxTween(toViewBox(FULL_BOUNDS));

// Zoom to a quadrant (or back out) and mirror it in the URL hash so the view
// can be linked to, e.g. from a blip page's breadcrumb.
function setActive(quadrant: Quadrant | null, animate = true) {
  activeQuadrant.value = quadrant;
  tweenTo(toViewBox(viewBounds.value), animate ? 400 : 0);
  history.replaceState(null, '', quadrant ? `#${quadrant}` : location.pathname + location.search);
}

function toggleQuadrant(quadrant: Quadrant) {
  setActive(activeQuadrant.value === quadrant ? null : quadrant);
}

function isDimmed(quadrant: Quadrant): boolean {
  return activeQuadrant.value !== null && activeQuadrant.value !== quadrant;
}

onMounted(() => {
  const fromHash = decodeURIComponent(location.hash.slice(1));
  if (QUADRANT_ORDER.includes(fromHash)) setActive(fromHash, false);
});

// Quadrant titles sit just outside the outer ring on each quadrant's bisector,
// in pills sized to the label (approximate glyph width of the 14px uppercase
// title) and kept inside the full view.
const TITLE_CHAR_WIDTH = 10.2;
const quadrantTitles = QUADRANT_ORDER.map((quadrant) => {
  const { startAngle, endAngle, label, color } = QUADRANT_CONFIGS[quadrant];
  const mid = (startAngle + endAngle) / 2;
  const r = RADAR_RADIUS + 35;
  const width = label.length * TITLE_CHAR_WIDTH + 24;
  const maxX = FULL_BOUNDS.x + FULL_BOUNDS.size - width / 2;
  const x = Math.max(-maxX, Math.min(maxX, r * Math.cos(mid)));
  return { quadrant, label, color, width, x, y: r * Math.sin(mid) };
});

const groups = computed(() =>
  QUADRANT_ORDER.map((quadrant) => ({ quadrant, ...QUADRANT_CONFIGS[quadrant], rings: groupByRing(props.blips, quadrant) })),
);
</script>

<template>
  <div class="radar-container">
    <div class="radar-stage">
      <svg
        :viewBox="currentViewBox"
        class="radar-svg"
        xmlns="http://www.w3.org/2000/svg"
        role="group"
        aria-label="Technology radar"
      >
        <RadarBackdrop :active-quadrant="activeQuadrant" interactive @select="toggleQuadrant" />

        <g>
          <RadarBlip
            v-for="blip in blips"
            :key="blip.id"
            :blip="blip"
            :href="blipHref(blip)"
            :highlighted="hoveredId === blip.id"
            :dimmed="isDimmed(blip.quadrant)"
            @hover="hoveredId = $event"
          />
        </g>

        <g
          v-for="qt in quadrantTitles"
          :key="qt.quadrant"
          class="quadrant-title"
          :class="{ dimmed: isDimmed(qt.quadrant) }"
          role="button"
          tabindex="0"
          :aria-label="activeQuadrant === qt.quadrant ? `Show full radar` : `Zoom to ${qt.label}`"
          :aria-pressed="activeQuadrant === qt.quadrant"
          @click="toggleQuadrant(qt.quadrant)"
          @keydown.enter.prevent="toggleQuadrant(qt.quadrant)"
          @keydown.space.prevent="toggleQuadrant(qt.quadrant)"
        >
          <rect :x="qt.x - qt.width / 2" :y="qt.y - 15" :width="qt.width" height="30" rx="5" class="quadrant-title-bg" />
          <text :x="qt.x" :y="qt.y" :fill="qt.color" class="quadrant-title-text" text-anchor="middle" dominant-baseline="middle">
            {{ qt.label }}
          </text>
        </g>

        <!-- Drawn last so it is never covered by other blips. -->
        <RadarLabel
          v-if="hoveredBlip"
          :x="hoveredBlip.x"
          :y="hoveredBlip.y"
          :text="hoveredBlip.title"
          :right-edge="viewBounds.x + viewBounds.size"
        />
      </svg>

      <div class="radar-legend" aria-hidden="true">
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

      <button v-if="activeQuadrant" type="button" class="back-button" @click="setActive(null)">
        &larr; Back to full view
      </button>
    </div>

    <nav class="blip-lists" aria-label="Blips by quadrant">
      <section
        v-for="group in groups"
        :key="group.quadrant"
        class="quadrant-section"
        :class="{ dimmed: isDimmed(group.quadrant) }"
      >
        <h3>
          <button
            type="button"
            class="quadrant-heading"
            :style="{ color: group.color }"
            :aria-pressed="activeQuadrant === group.quadrant"
            @click="toggleQuadrant(group.quadrant)"
          >
            {{ group.label }}
          </button>
        </h3>

        <div v-for="rg in group.rings" :key="rg.ring" class="ring-group">
          <h4>{{ rg.label }}</h4>
          <ul>
            <li
              v-for="blip in rg.blips"
              :key="blip.id"
              :class="{ highlighted: hoveredId === blip.id }"
              @mouseenter="hoveredId = blip.id"
              @mouseleave="hoveredId = null"
            >
              <a :href="blipHref(blip)" @focus="hoveredId = blip.id" @blur="hoveredId = null">
                <span class="blip-number-badge" :style="{ background: group.color }">{{ blip.number }}</span>
                {{ blip.title }}
                <span class="blip-status-indicator" aria-hidden="true">
                  <svg v-if="blip.status === 'new'" width="12" height="12" viewBox="0 0 12 12">
                    <circle cx="6" cy="6" r="3" :fill="group.color" />
                    <circle cx="6" cy="6" r="5" fill="none" :stroke="group.color" stroke-width="1" />
                  </svg>
                  <svg v-else-if="blip.status !== 'no-change'" width="12" height="12" viewBox="0 0 12 12">
                    <circle cx="6" cy="6" r="3" :fill="group.color" />
                    <path d="M 6 1 A 5 5 0 0 1 11 6" fill="none" :stroke="group.color" stroke-width="1" stroke-linecap="round" />
                  </svg>
                </span>
              </a>
            </li>
          </ul>
        </div>
      </section>
    </nav>
  </div>
</template>

<style scoped>
.radar-container {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 340px;
  gap: 2rem;
  align-items: start;
}

.radar-stage {
  position: relative;
}

.radar-svg {
  display: block;
  width: 100%;
  height: auto;
  /* Clip to the viewBox: when zoomed, the other quadrants must not spill out. */
  overflow: hidden;
}

.quadrant-title {
  cursor: pointer;
  transition: opacity 0.2s;
  outline: none;
}

.quadrant-title.dimmed {
  opacity: 0.3;
}

.quadrant-title-text {
  font-size: 14px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.quadrant-title-bg {
  fill: var(--color-radar-label-bg);
  transition: fill 0.15s;
}

.quadrant-title:hover .quadrant-title-bg {
  fill: var(--color-radar-label-bg-hover);
}

.quadrant-title:focus-visible .quadrant-title-bg {
  stroke: var(--color-radar-highlight-stroke);
  stroke-width: 2;
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

.back-button {
  position: absolute;
  top: 0;
  left: 0;
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
  margin-bottom: 0.5rem;
}

.quadrant-heading {
  display: block;
  width: 100%;
  font: inherit;
  font-size: 1rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  text-align: left;
  background: none;
  cursor: pointer;
  border: 1.5px solid currentColor;
  border-radius: 0.375rem;
  padding: 0.35rem 0.6rem;
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

@media (max-width: 1024px) {
  .radar-container {
    grid-template-columns: 1fr;
  }
}

/* Desktop: fill the space the page gives us; the sidebar scrolls on its own. */
@media (min-width: 1025px) and (min-height: 601px) {
  .radar-container {
    flex: 1;
    min-height: 0;
    grid-template-rows: minmax(0, 1fr);
    align-items: stretch;
  }

  /* A size container lets the SVG be the largest square that fits, so its
     viewport always matches the (square) viewBox and a zoomed quadrant never
     shows its neighbours in letterbox space. */
  .radar-stage {
    min-height: 0;
    container-type: size;
    display: grid;
    place-items: center;
  }

  .radar-svg {
    width: min(100cqw, 100cqh);
    height: min(100cqw, 100cqh);
  }

  .radar-legend {
    position: absolute;
    left: 0;
    bottom: 0;
    flex-direction: column;
    align-items: flex-start;
    gap: 0.35rem;
    margin: 0;
  }

  .blip-lists {
    min-height: 0;
    overflow-y: auto;
    overscroll-behavior: contain;
    padding-right: 0.5rem;
  }
}
</style>
