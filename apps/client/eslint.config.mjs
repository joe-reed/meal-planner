import nxEslintPlugin from '@nx/eslint-plugin';
import nextConfig from 'eslint-config-next';
import globals from 'globals';
import rootConfig from '../../eslint.config.mjs';

const config = [
  { ignores: ['.next/**/*', 'apps/client/.next/**/*'] },
  ...rootConfig,
  ...nxEslintPlugin.configs['flat/react-typescript'],
  ...nextConfig,
  {
    files: ['**/*.ts', '**/*.tsx', '**/*.js', '**/*.jsx'],
    languageOptions: {
      globals: { ...globals.jest },
    },
    rules: {
      '@next/next/no-html-link-for-pages': ['error', 'apps/client/pages'],
    },
  },
];

export default config;
