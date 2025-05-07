// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
	// Set the base path for GitHub Pages deployment
	// Use '/fedramp-datamesh-refarch/' for project pages
	// Remove this line if you're deploying to a custom domain or using organization pages
	base: '/fedramp-datamesh-refarch/',

	integrations: [
		starlight({
			title: 'FedRAMP High Event-Driven Data Mesh',
			social: [{ icon: 'github', label: 'GitHub', href: 'https://github.com/YARP-cloud/fedramp-datamesh-refarch' }],
			sidebar: [
				{
					label: 'Architecture',
					items: [
						{ label: 'Overview', slug: 'architecture/overview' },
					],
				},
				{
					label: 'Security',
					items: [{ label: 'Compliance', slug: 'security/fedramp-compliance' }],
				},
				{
					label: 'Developers',
					items: [
						{ label: 'Getting Started', slug: 'developers/getting-started' },
						{ label: 'Deployment', slug: 'developers/deployment' }
					],
				},
			],
		}),
	],
});
