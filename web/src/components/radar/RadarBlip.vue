<script setup lang="ts">
import { computed } from 'vue';
import type { PositionedBlip } from '../../types/radar';
import { QUADRANT_CONFIGS } from '../../config/radar';
import { quarterArcPath } from '../../utils/radarMath';

// One blip: dot and number, plus a ring when new or a quarter arc when it
// moved rings. `emphasized` marks the current blip on its detail page.
const props = withDefaults(
  defineProps<{
    blip: PositionedBlip;
    href: string;
    highlighted?: boolean;
    dimmed?: boolean;
    emphasized?: boolean;
  }>(),
  { highlighted: false, dimmed: false, emphasized: false },
);

const emit = defineEmits<{ hover: [id: string | null] }>();

const color = computed(() => QUADRANT_CONFIGS[props.blip.quadrant].color);
const active = computed(() => props.highlighted || props.emphasized);
const scale = computed(() => (props.emphasized ? 1.2 : 1));

// The arc faces the center for moved-in blips and away from it for moved-out.
const movedArc = computed(() => {
  const { x, y, status } = props.blip;
  if (status !== 'moved-in' && status !== 'moved-out') return null;
  const angle = Math.atan2(y, x);
  return quarterArcPath(x, y, 14 * scale.value, status === 'moved-in' ? angle + Math.PI : angle);
});
</script>

<template>
  <a
    :href="href"
    class="blip-link"
    :aria-label="`${blip.number}. ${blip.title}`"
    @mouseenter="emit('hover', blip.id)"
    @mouseleave="emit('hover', null)"
    @focus="emit('hover', blip.id)"
    @blur="emit('hover', null)"
  >
    <g class="blip-group" :class="{ highlighted: active, dimmed }">
      <circle :cx="blip.x" :cy="blip.y" :r="10 * scale" :fill="color" :fill-opacity="active ? 1 : 0.85" class="blip-dot" />
      <circle
        v-if="blip.status === 'new'"
        :cx="blip.x"
        :cy="blip.y"
        :r="14 * scale"
        fill="none"
        :stroke="color"
        stroke-width="2"
        :stroke-opacity="active ? 1 : 0.7"
      />
      <path
        v-if="movedArc"
        :d="movedArc"
        fill="none"
        :stroke="color"
        stroke-width="2.5"
        :stroke-opacity="active ? 1 : 0.7"
        stroke-linecap="round"
      />
      <text
        :x="blip.x"
        :y="blip.y"
        fill="#fff"
        :font-size="9 * scale"
        font-weight="700"
        text-anchor="middle"
        dominant-baseline="central"
        class="blip-number"
      >
        {{ blip.number }}
      </text>
      <circle
        v-if="active"
        :cx="blip.x"
        :cy="blip.y"
        :r="emphasized ? 18 : 12"
        fill="none"
        :stroke-width="emphasized ? 2.5 : 2"
        class="blip-highlight-ring"
      />
    </g>
  </a>
</template>

<style scoped>
.blip-link {
  text-decoration: none;
  outline: none;
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

.blip-link:focus-visible .blip-highlight-ring {
  stroke-width: 3;
}

.blip-dot {
  transition: fill-opacity 0.15s;
}

.blip-number {
  pointer-events: none;
  user-select: none;
}

.blip-highlight-ring {
  stroke: var(--color-radar-highlight-stroke);
}
</style>
