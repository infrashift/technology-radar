// Quadrant and ring slugs are configured in .radar.yaml, so they are plain strings.
export type Quadrant = string;
export type Ring = string;

export interface HistoryEntry {
  date: string;
  ring: Ring;
}

export interface BlipData {
  id: string;
  title: string;
  quadrant: Quadrant;
  ring: Ring;
  isNew: boolean;
  ciIdentifier?: string;
  history?: HistoryEntry[];
  relatedBlips?: string[];
}

export type BlipStatus = 'new' | 'moved-in' | 'moved-out' | 'no-change';

export interface PositionedBlip extends BlipData {
  x: number;
  y: number;
  number: number;
  status: BlipStatus;
}

export interface SearchableBlip extends BlipData {
  description: string;
}

export interface QuadrantConfig {
  label: string;
  startAngle: number;
  endAngle: number;
  color: string;
}

export interface RingConfig {
  label: string;
  innerRadius: number;
  outerRadius: number;
}
