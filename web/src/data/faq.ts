import { site } from './site';

export interface FaqItem {
  question: string;
  answer: string;
}

export const faqItems: FaqItem[] = [
  {
    question: "What is the Technology Radar?",
    answer: `<p>The Technology Radar is an opinionated guide to technology choices. It captures the collective experience and observations of our technology practitioners to help teams make informed decisions about the technologies they adopt.</p>
<p>The radar is not a market analysis or a comprehensive catalog of every available technology. Instead, it highlights items that we believe are particularly noteworthy — things teams should be paying attention to right now.</p>`,
  },
  {
    question: "How is the radar organized?",
    answer: `<p>The radar is divided into four quadrants, each representing a different category of technology:</p>
<ul>
  <li><strong>Techniques</strong> — Elements of the software development process, such as architectural patterns, design approaches, and methodologies.</li>
  <li><strong>Platforms</strong> — Things that we build software on top of, including cloud providers, operating systems, and container platforms.</li>
  <li><strong>Tools</strong> — Software that supports development, testing, deployment, and operations.</li>
  <li><strong>Languages &amp; Frameworks</strong> — Programming languages and the frameworks built around them.</li>
</ul>`,
  },
  {
    question: "What do the rings mean?",
    answer: `<p>Each item on the radar is placed in one of four rings that indicate our confidence in recommending it:</p>
<ul>
  <li><strong>Adopt</strong> — We have high confidence in these technologies. They are mature, proven, and recommended for broad use. Teams should be using these where appropriate.</li>
  <li><strong>Trial</strong> — These technologies are worth pursuing. They have shown clear value, and we recommend teams invest in understanding and experimenting with them on real projects.</li>
  <li><strong>Assess</strong> — These are worth exploring with the goal of understanding how they might affect your work. Teams should look into these to understand their potential impact, but there is typically higher risk or uncertainty involved.</li>
  <li><strong>Hold</strong> — Proceed with caution. These technologies are either too new and unproven, or they have known issues that make us reluctant to recommend them for new projects.</li>
</ul>`,
  },
  {
    question: "How are items placed on the radar?",
    answer: `<p>Items are placed on the radar based on the collective judgment of our technology practitioners. The process typically involves:</p>
<ul>
  <li>Identifying technologies that are gaining relevance or shifting in maturity.</li>
  <li>Discussing real-world experience — what has worked well, what hasn't, and in what contexts.</li>
  <li>Reaching consensus on the appropriate ring placement based on confidence level and practical experience.</li>
</ul>
<p>The radar represents a snapshot in time. Placements can change as technologies mature, new evidence emerges, or the landscape evolves.</p>`,
  },
  {
    question: 'What does "new" or "moved" mean?',
    answer: `<p>Items on the radar are marked to show how they've changed since the last update:</p>
<ul>
  <li><strong>New</strong> — The item is appearing on the radar for the first time.</li>
  <li><strong>Moved in</strong> — The item has moved closer to the center (e.g., from Assess to Trial), indicating growing confidence.</li>
  <li><strong>Moved out</strong> — The item has moved further from the center (e.g., from Trial to Assess), indicating decreasing confidence or new concerns.</li>
  <li><strong>No change</strong> — The item remains in the same ring as the last update.</li>
</ul>`,
  },
  {
    question: "How often is the radar updated?",
    answer: `<p>The radar is updated periodically to reflect shifts in the technology landscape and new insights from our practitioners. Each update may add new items, move existing items between rings, or retire items that are no longer noteworthy.</p>
<p>Check back regularly to see the latest changes, or use the search feature to find specific technologies.</p>`,
  },
  {
    question: 'Should I stop using technologies in "Hold"?',
    answer: `<p>Not necessarily. "Hold" does not mean "stop using immediately." It means we recommend caution when considering these technologies for <em>new</em> projects or significant new investment.</p>
<p>If a Hold technology is already well-established in your stack and working effectively, there may be no urgent reason to migrate away from it. The key consideration is whether to continue investing in it for new work versus exploring alternatives.</p>`,
  },
  {
    question: "How should I use the radar?",
    answer: `<p>The radar is a starting point for discussion, not a mandate. Here are some ways to use it effectively:</p>
<ul>
  <li><strong>Guide conversations</strong> — Use the radar to spark team discussions about technology choices and direction.</li>
  <li><strong>Inform decisions</strong> — When evaluating new technologies, check the radar for relevant context and recommendations.</li>
  <li><strong>Stay informed</strong> — Review new and moved items to stay aware of emerging trends and shifting recommendations.</li>
  <li><strong>Build your own perspective</strong> — Your context matters. Consider the radar's recommendations alongside your team's specific needs, constraints, and experience.</li>
</ul>`,
  },
  {
    question: "Can I suggest a technology for the radar?",
    answer: `<p>Yes! The radar benefits from diverse perspectives and real-world experience. If you've been working with a technology that you believe should be on the radar — or if you think an existing item should be reconsidered — we welcome your input.</p>
<p>Open an issue in the <a href="${site.repoUrl}" target="_blank" rel="noopener noreferrer">source repository</a> with your suggestion, along with context about your experience and why you think it's noteworthy. See the <a href="/radar/contribute">Contribute</a> page for more details.</p>`,
  },
  {
    question: "What is the difference between the quadrants?",
    answer: `<p>The quadrants categorize technologies by their nature, not by their importance. Some items could arguably fit in more than one quadrant — for example, Kubernetes could be seen as both a tool and a platform. In these cases, we place items in the quadrant that best reflects how most teams interact with them.</p>
<p>The quadrant boundaries are a practical organizational tool, not a strict taxonomy. Focus on the item's ring placement and description rather than its quadrant for the most actionable insight.</p>`,
  },
  {
    question: "Where is the source code?",
    answer: `<p>The Technology Radar is open source. You can find the source code, report issues, and contribute in the <a href="${site.repoUrl}" target="_blank" rel="noopener noreferrer">source repository</a>.</p>
<p>See the <a href="/radar/contribute">Contribute</a> page for guidance on how to get involved.</p>`,
  },
];
