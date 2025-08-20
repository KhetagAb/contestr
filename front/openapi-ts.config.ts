import { defineConfig } from '@hey-api/openapi-ts';

export default defineConfig({
  input: 'front/SportComponents/openapi.yaml',
  output: 'front/SportComponents/client',
  plugins: [
    "@hey-api/client-fetch",
    "@tanstack/react-query"
  ]


});