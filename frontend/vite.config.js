import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig, loadEnv } from 'vite';

export default defineConfig(({ mode }) => {
    const env = loadEnv(mode, process.cwd(), '');
    const isDevBuild = env.BUILD_MODE === 'development';

    return {
        plugins: [sveltekit()],
        build: {
            minify: 'terser',
            terserOptions: {
                mangle: !isDevBuild, // Only mangle in production
                compress: {
                    keep_fnames: isDevBuild, // Keep function names in development
                    keep_classnames: isDevBuild // Keep class names in development
                },
                format: {
                    comments: false
                }
            }
        }
    };
});
