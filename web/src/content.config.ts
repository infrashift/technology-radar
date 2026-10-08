import { defineCollection } from 'astro:content';
import { z } from 'astro/zod';
import { glob } from 'astro/loaders';
import radarConfig from './data/radar-config.json';

const quadrantSlugs = radarConfig.quadrants.map((q) => q.slug) as [string, ...string[]];
const ringSlugs = radarConfig.rings.map((r) => r.slug) as [string, ...string[]];

const blips = defineCollection({
  loader: glob({ pattern: '**/*.mdx', base: './src/content/blips' }),
  schema: z.object({
    id: z.string(),
    title: z.string(),
    quadrant: z.enum(quadrantSlugs),
    ring: z.enum(ringSlugs),
    isNew: z.boolean().default(false),
    ciIdentifier: z.string().optional(),
    history: z.array(z.object({
      date: z.string(),
      ring: z.enum(ringSlugs),
    })).optional(),
    relatedBlips: z.array(z.string()).optional(),
  }),
});

export const collections = { blips };
