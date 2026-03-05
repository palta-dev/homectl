/**
 * JSON Schema for homectl configuration
 * Used for validation and documentation
 */

export const configSchema = {
  $schema: 'http://json-schema.org/draft-07/schema#',
  type: 'object',
  required: ['version', 'groups'],
  properties: {
    version: {
      type: 'integer',
      minimum: 1,
      description: 'Configuration schema version',
    },
    settings: {
      type: 'object',
      properties: {
        title: { type: 'string' },
        theme: { type: 'string', enum: ['dark', 'light'] },
        allowHosts: {
          type: 'array',
          items: { type: 'string' },
        },
        blockPrivateMetaIPs: { type: 'boolean' },
        requestTimeout: { type: 'string' },
        cache: {
          type: 'object',
          properties: {
            defaultTTL: { type: 'integer', minimum: 0 },
            maxEntries: { type: 'integer', minimum: 1 },
            widgetTTL: { type: 'object', additionalProperties: { type: 'integer' } },
          },
        },
        auth: {
          type: 'object',
          properties: {
            enabled: { type: 'boolean' },
            provider: { type: 'string', enum: ['local', 'github', 'google'] },
            session: {
              type: 'object',
              properties: {
                maxAge: { type: 'string' },
              },
            },
            github: {
              type: 'object',
              properties: {
                clientId: { type: 'string' },
                clientSecret: { type: 'string' },
                allowedUsers: {
                  type: 'array',
                  items: { type: 'string' },
                },
              },
            },
          },
        },
      },
    },
    groups: {
      type: 'array',
      items: {
        type: 'object',
        required: ['name', 'services'],
        properties: {
          name: { type: 'string' },
          layout: { type: 'string', enum: ['grid', 'list', 'compact'] },
          collapsed: { type: 'boolean' },
          services: {
            type: 'array',
            items: {
              type: 'object',
              required: ['name', 'url'],
              properties: {
                name: { type: 'string' },
                url: { type: 'string', format: 'uri-reference' },
                icon: { type: 'string' },
                description: { type: 'string' },
                tags: {
                  type: 'array',
                  items: { type: 'string' },
                },
                newTab: { type: 'boolean' },
                pingEnabled: { type: 'boolean' },
                checks: {
                  type: 'array',
                  items: {
                    oneOf: [
                      {
                        type: 'object',
                        required: ['type', 'url'],
                        properties: {
                          type: { const: 'http' },
                          url: { type: 'string' },
                          method: { type: 'string' },
                          expectStatus: { type: 'integer' },
                          expectBodyContains: { type: 'string' },
                          headers: { type: 'object' },
                          timeout: { type: 'string' },
                          intervalSeconds: { type: 'integer', minimum: 1 },
                          retries: { type: 'integer', minimum: 0 },
                        },
                      },
                      {
                        type: 'object',
                        required: ['type', 'host', 'port'],
                        properties: {
                          type: { const: 'tcp' },
                          host: { type: 'string' },
                          port: { type: 'integer', minimum: 1, maximum: 65535 },
                          timeout: { type: 'string' },
                          intervalSeconds: { type: 'integer', minimum: 1 },
                        },
                      },
                      {
                        type: 'object',
                        required: ['type', 'host'],
                        properties: {
                          type: { const: 'ping' },
                          host: { type: 'string' },
                          count: { type: 'integer', minimum: 1 },
                          intervalSeconds: { type: 'integer', minimum: 1 },
                        },
                      },
                    ],
                  },
                },
                widgets: {
                  type: 'array',
                  items: {
                    oneOf: [
                      {
                        type: 'object',
                        required: ['type', 'url', 'jsonPath'],
                        properties: {
                          type: { const: 'httpJson' },
                          url: { type: 'string' },
                          jsonPath: { type: 'string' },
                          label: { type: 'string' },
                          format: { type: 'string', enum: ['status', 'bytes', 'duration', 'raw', 'percent'] },
                          cacheTTL: { type: 'integer', minimum: 0 },
                        },
                      },
                      {
                        type: 'object',
                        required: ['type', 'url', 'selector'],
                        properties: {
                          type: { const: 'httpHtml' },
                          url: { type: 'string' },
                          selector: { type: 'string' },
                          attribute: { type: 'string' },
                          label: { type: 'string' },
                        },
                      },
                      {
                        type: 'object',
                        required: ['type', 'host', 'port'],
                        properties: {
                          type: { const: 'tcpPort' },
                          host: { type: 'string' },
                          port: { type: 'integer', minimum: 1, maximum: 65535 },
                          label: { type: 'string' },
                          showLatency: { type: 'boolean' },
                        },
                      },
                    ],
                  },
                },
              },
            },
          },
        },
      },
    },
    icons: {
      type: 'object',
      required: ['sources'],
      properties: {
        sources: {
          type: 'array',
          items: {
            oneOf: [
              {
                type: 'object',
                required: ['type', 'path'],
                properties: {
                  type: { const: 'local' },
                  path: { type: 'string' },
                },
              },
              {
                type: 'object',
                required: ['type'],
                properties: {
                  type: { const: 'simpleicons' },
                  cache: { type: 'boolean' },
                  cacheTTL: { type: 'integer', minimum: 0 },
                },
              },
              {
                type: 'object',
                required: ['type', 'baseUrl'],
                properties: {
                  type: { const: 'url' },
                  baseUrl: { type: 'string' },
                  pathTemplate: { type: 'string' },
                },
              },
            ],
          },
        },
      },
    },
  },
} as const;

export type ConfigSchema = typeof configSchema;
