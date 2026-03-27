<script setup lang="ts">
import { ref, computed } from 'vue';
import type { SearchableBlip, Quadrant, Ring } from '../types/radar';
import { QUADRANT_CONFIGS, RING_ORDER } from '../types/radar';

const props = defineProps<{
  blips: SearchableBlip[];
}>();

const searchQuery = ref('');
const selectedQuadrants = ref<Set<Quadrant>>(new Set());
const selectedRings = ref<Set<Ring>>(new Set());

const ringColors: Record<Ring, string> = {
  adopt: '#2da44e',
  trial: '#0969da',
  assess: '#bf8700',
  hold: '#cf222e',
};

const quadrantEntries = Object.entries(QUADRANT_CONFIGS) as [Quadrant, typeof QUADRANT_CONFIGS[Quadrant]][];

function toggleQuadrant(q: Quadrant) {
  const s = new Set(selectedQuadrants.value);
  if (s.has(q)) s.delete(q);
  else s.add(q);
  selectedQuadrants.value = s;
}

function toggleRing(r: Ring) {
  const s = new Set(selectedRings.value);
  if (s.has(r)) s.delete(r);
  else s.add(r);
  selectedRings.value = s;
}

const filteredBlips = computed(() => {
  const query = searchQuery.value.toLowerCase().trim();
  const quadrants = selectedQuadrants.value;
  const rings = selectedRings.value;

  return props.blips.filter((blip) => {
    if (quadrants.size > 0 && !quadrants.has(blip.quadrant)) return false;
    if (rings.size > 0 && !rings.has(blip.ring)) return false;
    if (query) {
      const inTitle = blip.title.toLowerCase().includes(query);
      const inDesc = blip.description.toLowerCase().includes(query);
      if (!inTitle && !inDesc) return false;
    }
    return true;
  }).sort((a, b) => {
    if (!query) return a.title.localeCompare(b.title);
    const aTitle = a.title.toLowerCase().includes(query);
    const bTitle = b.title.toLowerCase().includes(query);
    if (aTitle && !bTitle) return -1;
    if (!aTitle && bTitle) return 1;
    return a.title.localeCompare(b.title);
  });
});

function truncate(text: string, max: number): string {
  if (text.length <= max) return text;
  return text.slice(0, max).replace(/\s\S*$/, '') + '...';
}
</script>

<template>
  <div class="search-page">
    <h1>Search blips</h1>

    <div class="search-input-wrapper">
      <input
        v-model="searchQuery"
        type="text"
        placeholder="Search by name or description..."
        autofocus
        class="search-input"
      />
    </div>

    <div class="filters">
      <div class="filter-group">
        <span class="filter-label">Quadrant</span>
        <div class="filter-buttons">
          <button
            v-for="[key, config] in quadrantEntries"
            :key="key"
            :class="['filter-btn', { active: selectedQuadrants.has(key) }]"
            :style="{
              '--btn-color': config.color,
              borderColor: config.color,
              background: selectedQuadrants.has(key) ? config.color : 'transparent',
              color: selectedQuadrants.has(key) ? '#fff' : config.color,
            }"
            @click="toggleQuadrant(key)"
          >
            {{ config.label }}
          </button>
        </div>
      </div>

      <div class="filter-group">
        <span class="filter-label">Ring</span>
        <div class="filter-buttons">
          <button
            v-for="ring in RING_ORDER"
            :key="ring"
            :class="['filter-btn', { active: selectedRings.has(ring) }]"
            :style="{
              '--btn-color': ringColors[ring],
              borderColor: ringColors[ring],
              background: selectedRings.has(ring) ? ringColors[ring] : 'transparent',
              color: selectedRings.has(ring) ? '#fff' : ringColors[ring],
            }"
            @click="toggleRing(ring)"
          >
            {{ ring }}
          </button>
        </div>
      </div>
    </div>

    <div class="result-count">{{ filteredBlips.length }} blip{{ filteredBlips.length !== 1 ? 's' : '' }}</div>

    <ul v-if="filteredBlips.length > 0" class="results-list">
      <li v-for="blip in filteredBlips" :key="blip.id" class="result-item">
        <a :href="`/radar/${blip.quadrant}/${blip.id}`" class="result-link">
          <div class="result-header">
            <span class="result-title">{{ blip.title }}</span>
            <div class="result-badges">
              <span class="ring-badge" :style="{ background: ringColors[blip.ring] }">
                {{ blip.ring }}
              </span>
              <span class="quadrant-label">{{ QUADRANT_CONFIGS[blip.quadrant].label }}</span>
              <span v-if="blip.isNew" class="new-badge">NEW</span>
            </div>
          </div>
          <p v-if="blip.description" class="result-description">{{ truncate(blip.description, 200) }}</p>
        </a>
      </li>
    </ul>

    <div v-else class="empty-state">
      No blips match your search.
    </div>
  </div>
</template>

<style scoped>
.search-page {
  max-width: 800px;
  margin: 0 auto;
}

h1 {
  font-size: 2rem;
  font-weight: 700;
  margin-bottom: 1.5rem;
  letter-spacing: -0.02em;
}

.search-input-wrapper {
  margin-bottom: 1.25rem;
}

.search-input {
  width: 100%;
  padding: 0.75rem 1rem;
  font-size: 1rem;
  font-family: inherit;
  border: 1px solid var(--color-border);
  border-radius: 0.5rem;
  background: var(--color-bg-surface);
  color: var(--color-text-primary);
  outline: none;
  transition: border-color 0.15s;
}

.search-input:focus {
  border-color: var(--color-link);
}

.search-input::placeholder {
  color: var(--color-text-tertiary);
}

.filters {
  display: flex;
  flex-wrap: wrap;
  gap: 1.5rem;
  margin-bottom: 1rem;
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.filter-label {
  font-size: 0.8rem;
  font-weight: 600;
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.filter-buttons {
  display: flex;
  gap: 0.375rem;
  flex-wrap: wrap;
}

.filter-btn {
  padding: 0.25rem 0.625rem;
  border: 1.5px solid;
  border-radius: 1rem;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.03em;
  cursor: pointer;
  transition: background 0.15s, color 0.15s;
  font-family: inherit;
}

.filter-btn:hover {
  opacity: 0.85;
}

.result-count {
  font-size: 0.85rem;
  color: var(--color-text-tertiary);
  margin-bottom: 1rem;
}

.results-list {
  list-style: none;
  padding: 0;
}

.result-item {
  border-bottom: 1px solid var(--color-border-light);
}

.result-link {
  display: block;
  padding: 1rem 0.75rem;
  text-decoration: none;
  color: inherit;
  border-radius: 0.375rem;
  transition: background 0.1s;
}

.result-link:hover {
  background: var(--color-bg-hover);
}

.result-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.result-title {
  font-size: 1.05rem;
  font-weight: 600;
  color: var(--color-link);
}

.result-badges {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  flex-shrink: 0;
}

.ring-badge {
  color: #fff;
  padding: 0.15rem 0.5rem;
  border-radius: 1rem;
  font-size: 0.7rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.quadrant-label {
  color: var(--color-text-secondary);
  font-size: 0.8rem;
}

.new-badge {
  background: var(--color-new-badge);
  color: #fff;
  padding: 0.15rem 0.5rem;
  border-radius: 0.25rem;
  font-size: 0.7rem;
  font-weight: 700;
  letter-spacing: 0.05em;
}

.result-description {
  margin-top: 0.375rem;
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  line-height: 1.5;
}

.empty-state {
  padding: 3rem 1rem;
  text-align: center;
  color: var(--color-text-tertiary);
  font-size: 1rem;
}

@media (max-width: 600px) {
  .result-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .filters {
    flex-direction: column;
    gap: 0.75rem;
  }
}
</style>
