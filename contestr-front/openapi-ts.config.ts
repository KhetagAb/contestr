import { defineConfig } from '@hey-api/openapi-ts';

export default defineConfig({
  input: 'match-integration/openapi.yaml',
  output: 'match-integration/client',
  plugins: [
    "@hey-api/client-fetch",
    "@tanstack/react-query"
  ]
});