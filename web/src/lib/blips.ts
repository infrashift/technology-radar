import { getCollection, type CollectionEntry } from 'astro:content';
import type { BlipData, PositionedBlip } from '../types/radar';
import { layoutBlips } from '../utils/radarLayout';

export type BlipEntry = CollectionEntry<'blips'>;

export interface RadarData {
  entries: BlipEntry[];
  /** Positioned and numbered blips, in collection order. */
  blips: PositionedBlip[];
  byId: Map<string, PositionedBlip>;
}

let cached: Promise<RadarData> | undefined;

/**
 * Load the blip collection and lay out the radar once per build. In dev the
 * cache is skipped so content edits show up on reload.
 */
export function getRadarData(): Promise<RadarData> {
  if (!import.meta.env.PROD) return load();
  return (cached ??= load());
}

async function load(): Promise<RadarData> {
  const entries = await getCollection('blips');
  const blips = layoutBlips(entries.map((e) => toBlipData(e)));
  return { entries, blips, byId: new Map(blips.map((b) => [b.id, b])) };
}

function toBlipData({ data }: BlipEntry): BlipData {
  const { id, title, quadrant, ring, isNew, ciIdentifier, history, relatedBlips } = data;
  return { id, title, quadrant, ring, isNew, ciIdentifier, history, relatedBlips };
}
