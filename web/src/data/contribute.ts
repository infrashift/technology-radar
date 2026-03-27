import { site } from './site';

export interface ContributeSection {
  title: string;
  content: string;
}

export const contributeItems: ContributeSection[] = [
  {
    title: "Suggest a Technology",
    content: `<p>Have you been working with a technology that you think deserves a place on the radar? We welcome suggestions for new blips across all four quadrants:</p>
<ul>
  <li><strong>Techniques</strong> — Development practices, architectural patterns, or methodologies.</li>
  <li><strong>Platforms</strong> — Infrastructure, cloud services, or runtime environments.</li>
  <li><strong>Tools</strong> — Software that supports development, testing, deployment, or operations.</li>
  <li><strong>Languages &amp; Frameworks</strong> — Programming languages or the frameworks built around them.</li>
</ul>
<p>When suggesting a technology, it helps to include:</p>
<ol>
  <li>The name of the technology and which quadrant it belongs to.</li>
  <li>A brief description of what it does and why it's noteworthy.</li>
  <li>Your experience with it — what worked well, what didn't, and in what context.</li>
  <li>A suggested ring placement (Adopt, Trial, Assess, or Hold) with reasoning.</li>
</ol>
<p><a href="${site.repoUrl}/issues/new" target="_blank" rel="noopener noreferrer">Open an issue</a> in the repository to submit your suggestion.</p>`,
  },
  {
    title: "Request a Reassessment",
    content: `<p>Technologies evolve, and so should their placement on the radar. If you believe an existing item should move to a different ring — whether closer to the center or further out — let us know.</p>
<p>A reassessment request is most useful when it includes:</p>
<ul>
  <li>The current blip and its ring placement.</li>
  <li>The ring you believe it should move to.</li>
  <li>What has changed — new releases, shifting community adoption, your team's experience, or emerging concerns.</li>
</ul>
<p><a href="${site.repoUrl}/issues/new" target="_blank" rel="noopener noreferrer">Open an issue</a> in the repository to request a reassessment.</p>`,
  },
  {
    title: "Report a Bug",
    content: `<p>Found something broken on the radar? Whether it's a rendering issue, incorrect data, a broken link, or a problem with the interactive visualization, we want to know about it.</p>
<p>A good bug report includes:</p>
<ul>
  <li>What you expected to happen.</li>
  <li>What actually happened.</li>
  <li>Steps to reproduce the issue.</li>
  <li>Your browser and device (if it's a visual or interaction issue).</li>
</ul>
<p><a href="${site.repoUrl}/issues/new" target="_blank" rel="noopener noreferrer">Report a bug</a> in the repository.</p>`,
  },
  {
    title: "Fix a Bug",
    content: `<p>If you've found a bug and know how to fix it, we welcome pull requests. The project is built with Astro, Vue, and TypeScript on the frontend, with a Go CLI for content validation and generation.</p>
<p>To get started:</p>
<ol>
  <li>Fork the <a href="${site.repoUrl}" target="_blank" rel="noopener noreferrer">repository</a> and create a branch for your fix.</li>
  <li>Make your changes, ensuring they follow the existing code style and patterns.</li>
  <li>Test your changes locally by running <code>bun run dev</code> in the <code>web/</code> directory.</li>
  <li>Submit a pull request with a clear description of the bug and your fix.</li>
</ol>`,
  },
  {
    title: "Enhance the Radar",
    content: `<p>Have an idea to improve the radar experience? We're open to enhancements that make the radar more useful, accessible, or easier to maintain. Some areas where contributions are especially welcome:</p>
<ul>
  <li><strong>Accessibility</strong> — Improving screen reader support, keyboard navigation, or color contrast.</li>
  <li><strong>Visualization</strong> — Enhancements to the interactive radar chart or blip detail pages.</li>
  <li><strong>Search and filtering</strong> — Better ways to find and explore technologies.</li>
  <li><strong>Content tooling</strong> — Improvements to the CUE validation schemas or the Go CLI.</li>
</ul>
<p>For larger changes, consider <a href="${site.repoUrl}/issues/new" target="_blank" rel="noopener noreferrer">opening an issue</a> first to discuss the approach before investing time in implementation.</p>`,
  },
  {
    title: "How to Contribute",
    content: `<p>All contributions — whether suggestions, bug reports, or code changes — follow the same general process:</p>
<ol>
  <li><strong><a href="${site.repoUrl}/issues/new" target="_blank" rel="noopener noreferrer">Open an issue</a></strong> describing the suggestion, bug, or enhancement.</li>
  <li><strong>Discuss</strong> the approach with maintainers if the change is non-trivial.</li>
  <li><strong>Submit a <a href="${site.repoUrl}/pulls" target="_blank" rel="noopener noreferrer">pull request</a></strong> if you're contributing code or content changes.</li>
</ol>
<p>We appreciate contributions of all sizes. Even small improvements like fixing a typo in a blip description or clarifying a ring placement help make the radar better for everyone.</p>`,
  },
];
