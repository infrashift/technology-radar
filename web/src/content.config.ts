import { defineCollection, z } from 'astro:content';
import { glob } from 'astro/loaders';

const blips = defineCollection({
  loader: glob({ pattern: '**/*.mdx', base: './src/content/blips' }),
  schema: z.object({
    id: z.string(),
    title: z.string(),
    quadrant: z.enum(['techniques', 'platforms', 'tools', 'languages-and-frameworks']),
    ring: z.enum(['adopt', 'trial', 'assess', 'hold']),
    isNew: z.boolean().default(false),
    ciIdentifier: z.string().optional(),
    history: z.array(z.object({
      date: z.string(),
      ring: z.enum(['adopt', 'trial', 'assess', 'hold']),
    })).optional(),
    relatedBlips: z.array(z.string()).optional(),
  }),
});

export const collections = { blips };
