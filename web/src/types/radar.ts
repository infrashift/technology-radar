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

export interface HistoryEntry {
  date: string;
  ring: Ring;
}

export interface PositionedBlip extends BlipData {
  x: number;
  y: number;
  number: number;
  status: BlipStatus;
}

export type Quadrant = 'techniques' | 'platforms' | 'tools' | 'languages-and-frameworks';
export type Ring = 'adopt' | 'trial' | 'assess' | 'hold';

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

export const QUADRANT_CONFIGS: Record<Quadrant, QuadrantConfig> = {
  techniques: {
    label: 'Techniques',
    startAngle: -Math.PI / 2,
    endAngle: 0,
    color: '#2da44e',
  },
  platforms: {
    label: 'Platforms',
    startAngle: 0,
    endAngle: Math.PI / 2,
    color: '#0969da',
  },
  tools: {
    label: 'Tools',
    startAngle: Math.PI / 2,
    endAngle: Math.PI,
    color: '#bf8700',
  },
  'languages-and-frameworks': {
    label: 'Languages & Frameworks',
    startAngle: -Math.PI,
    endAngle: -Math.PI / 2,
    color: '#cf222e',
  },
};

export const RING_CONFIGS: RingConfig[] = [
  { label: 'Adopt', innerRadius: 0, outerRadius: 150 },
  { label: 'Trial', innerRadius: 150, outerRadius: 275 },
  { label: 'Assess', innerRadius: 275, outerRadius: 375 },
  { label: 'Hold', innerRadius: 375, outerRadius: 475 },
];

export const RING_ORDER: Ring[] = ['adopt', 'trial', 'assess', 'hold'];

export interface SearchableBlip extends BlipData {
  description: string;
}

export const RADAR_SIZE = 1000;
export const RADAR_RADIUS = 475;
