// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
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
					items: [{ label: 'Deployment', slug: 'developers/deployment' }],
				},
			],
		}),
	],
});
