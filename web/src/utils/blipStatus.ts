import type { BlipData } from '../types/radar';
import type { BlipStatus } from '../types/radar';
import { RING_ORDER } from '../config/radar';

/**
 * Derive the movement status of a blip from its isNew flag and history.
 *
 * - 'new': first appearance on the radar (isNew === true)
 * - 'moved-in': ring moved closer to center (e.g., assess → trial)
 * - 'moved-out': ring moved farther from center (e.g., trial → assess)
 * - 'no-change': ring unchanged or insufficient history
 */
export function deriveBlipStatus(blip: BlipData): BlipStatus {
  if (blip.isNew) return 'new';

  const history = blip.history;
  if (!history || history.length < 2) return 'no-change';

  // History is ordered most-recent first
  const currentRingIndex = RING_ORDER.indexOf(history[0].ring);
  const previousRingIndex = RING_ORDER.indexOf(history[1].ring);

  if (currentRingIndex < previousRingIndex) return 'moved-in';
  if (currentRingIndex > previousRingIndex) return 'moved-out';

  return 'no-change';
}
