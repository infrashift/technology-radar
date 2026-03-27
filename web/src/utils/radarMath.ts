import { arc } from 'd3-shape';

/**
 * Convert polar coordinates (angle, radius) to cartesian (x, y).
 * Angle is in radians, measured clockwise from the positive y-axis.
 */
export function polarToCartesian(angle: number, radius: number): { x: number; y: number } {
  return {
    x: radius * Math.cos(angle),
    y: radius * Math.sin(angle),
  };
}

/**
 * Simple string hash (djb2) returning a float in [0, 1).
 */
export function hashString(s: string): number {
  let hash = 5381;
  for (let i = 0; i < s.length; i++) {
    hash = ((hash << 5) + hash + s.charCodeAt(i)) >>> 0;
  }
  return (hash % 10000) / 10000;
}

/**
 * A second independent hash for angular variation.
 */
export function hashString2(s: string): number {
  let hash = 0;
  for (let i = 0; i < s.length; i++) {
    hash = ((hash * 31) + s.charCodeAt(i)) >>> 0;
  }
  return (hash % 10000) / 10000;
}

/**
 * Clamp a point to stay within a ring (between inner and outer radii).
 */
export function clampToRing(
  x: number,
  y: number,
  innerRadius: number,
  outerRadius: number,
  padding: number = 12
): { x: number; y: number } {
  const dist = Math.sqrt(x * x + y * y);
  const minR = innerRadius + padding;
  const maxR = outerRadius - padding;

  if (dist < minR) {
    const scale = minR / Math.max(dist, 0.001);
    return { x: x * scale, y: y * scale };
  }
  if (dist > maxR) {
    const scale = maxR / Math.max(dist, 0.001);
    return { x: x * scale, y: y * scale };
  }
  return { x, y };
}

/**
 * Clamp a point to stay within a quadrant's angular range.
 */
export function clampToQuadrant(
  x: number,
  y: number,
  startAngle: number,
  endAngle: number,
  padding: number = 0.05
): { x: number; y: number } {
  const angle = Math.atan2(y, x);
  const dist = Math.sqrt(x * x + y * y);
  const minAngle = startAngle + padding;
  const maxAngle = endAngle - padding;

  if (angle < minAngle) {
    return polarToCartesian(minAngle, dist);
  }
  if (angle > maxAngle) {
    return polarToCartesian(maxAngle, dist);
  }
  return { x, y };
}

/**
 * Generate an SVG arc path for a ring segment within a quadrant.
 */
export function generateRingArc(
  innerRadius: number,
  outerRadius: number,
  startAngle: number,
  endAngle: number
): string {
  const generator = arc<void>()
    .innerRadius(innerRadius)
    .outerRadius(outerRadius)
    .startAngle(startAngle + Math.PI / 2)  // d3 arc measures from 12 o'clock
    .endAngle(endAngle + Math.PI / 2);

  return generator(undefined as any) || '';
}

/**
 * Distance between two points.
 */
export function distance(x1: number, y1: number, x2: number, y2: number): number {
  const dx = x1 - x2;
  const dy = y1 - y2;
  return Math.sqrt(dx * dx + dy * dy);
}

/**
 * Generate a quarter-arc (90°) SVG path centered at a given angle.
 * Used for moved-in/moved-out blip indicators.
 * @param cx - center x of the blip
 * @param cy - center y of the blip
 * @param r - radius of the arc (slightly larger than blip dot)
 * @param centerAngle - the angle the arc is centered on (radians)
 */
export function quarterArcPath(cx: number, cy: number, r: number, centerAngle: number): string {
  const startAngle = centerAngle - Math.PI / 4;
  const endAngle = centerAngle + Math.PI / 4;
  const x1 = cx + r * Math.cos(startAngle);
  const y1 = cy + r * Math.sin(startAngle);
  const x2 = cx + r * Math.cos(endAngle);
  const y2 = cy + r * Math.sin(endAngle);
  return `M ${x1} ${y1} A ${r} ${r} 0 0 1 ${x2} ${y2}`;
}
