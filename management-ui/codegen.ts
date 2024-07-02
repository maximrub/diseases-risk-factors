import { CodegenConfig } from '@graphql-codegen/cli';
import { config as dotenvConfig } from 'dotenv';
dotenvConfig();

const config: CodegenConfig = {
  schema: process.env.REACT_APP_API_URL,
  documents: ['src/**/*.tsx'],
  generates: {
    './src/__generated__/': {
      preset: 'client',
      plugins: [],
      presetConfig: {
        gqlTagName: 'gql',
      }
    }
  },
  ignoreNoDocuments: true,
};

export default config;