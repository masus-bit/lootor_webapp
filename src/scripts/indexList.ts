export const indices = {
  collection: {
    settings: {
      analysis: {
        analyzer: {
          exact_match: {
            type: 'custom',
            tokenizer: 'keyword',
            filter: ['lowercase'],
          },
          prefix_search: {
            type: 'custom',
            tokenizer: 'standard',
            filter: ['lowercase', 'edge_ngram_filter'],
          },
          full_text: {
            type: 'standard',
            filter: ['lowercase'],
          },
        },
        filter: {
          edge_ngram_filter: {
            type: 'edge_ngram',
            min_gram: 3,
            max_gram: 15,
          },
        },
        normalizer: {
          lowercase: {
            type: 'custom',
            filter: ['lowercase'],
          },
        },
      },
    },
    mappings: {
      properties: {
        name: {
          type: 'text',
          fields: {
            exact: {
              type: 'text',
              analyzer: 'exact_match',
            },
            prefix: {
              type: 'text',
              analyzer: 'prefix_search',
            },
            full: {
              type: 'text',
              analyzer: 'full_text',
            },
            keyword: {
              type: 'keyword',
              normalizer: 'lowercase',
            },
          },
        },
      },
    },
  },
  user: {
    settings: {
      analysis: {
        analyzer: {
          exact_match: {
            type: 'custom',
            tokenizer: 'keyword',
            filter: ['lowercase'],
          },
          prefix_search: {
            type: 'custom',
            tokenizer: 'standard',
            filter: ['lowercase', 'edge_ngram_filter'],
          },
          full_text: {
            type: 'standard',
            filter: ['lowercase'],
          },
        },
        filter: {
          edge_ngram_filter: {
            type: 'edge_ngram',
            min_gram: 3,
            max_gram: 15,
          },
        },
        normalizer: {
          lowercase: {
            type: 'custom',
            filter: ['lowercase'],
          },
        },
      },
    },
    mappings: {
      properties: {
        name: {
          type: 'text',
          fields: {
            exact: {
              type: 'text',
              analyzer: 'exact_match',
            },
            prefix: {
              type: 'text',
              analyzer: 'prefix_search',
            },
            full: {
              type: 'text',
              analyzer: 'full_text',
            },
            keyword: {
              type: 'keyword',
              normalizer: 'lowercase',
            },
          },
        },
      },
    },
  },
  tags: {
    settings: {
      analysis: {
        analyzer: {
          exact_match: {
            type: 'custom',
            tokenizer: 'keyword',
            filter: ['lowercase'],
          },
          prefix_search: {
            type: 'custom',
            tokenizer: 'standard',
            filter: ['lowercase', 'edge_ngram_filter'],
          },
          full_text: {
            type: 'standard',
            filter: ['lowercase'],
          },
        },
        filter: {
          edge_ngram_filter: {
            type: 'edge_ngram',
            min_gram: 3,
            max_gram: 15,
          },
        },
        normalizer: {
          lowercase: {
            type: 'custom',
            filter: ['lowercase'],
          },
        },
      },
    },
    mappings: {
      properties: {
        name: {
          type: 'text',
          fields: {
            exact: {
              type: 'text',
              analyzer: 'exact_match',
            },
            prefix: {
              type: 'text',
              analyzer: 'prefix_search',
            },
            full: {
              type: 'text',
              analyzer: 'full_text',
            },
            keyword: {
              type: 'keyword',
              normalizer: 'lowercase',
            },
          },
        },
      },
    },
  },
  collection_item: {
    settings: {
      analysis: {
        analyzer: {
          exact_match: {
            type: 'custom',
            tokenizer: 'keyword',
            filter: ['lowercase'],
          },
          prefix_search: {
            type: 'custom',
            tokenizer: 'standard',
            filter: ['lowercase', 'edge_ngram_filter'],
          },
          full_text: {
            type: 'standard',
            filter: ['lowercase'],
          },
        },
        filter: {
          edge_ngram_filter: {
            type: 'edge_ngram',
            min_gram: 3,
            max_gram: 15,
          },
        },
        normalizer: {
          lowercase: {
            type: 'custom',
            filter: ['lowercase'],
          },
        },
      },
    },
    mappings: {
      properties: {
        name: {
          type: 'text',
          fields: {
            exact: {
              type: 'text',
              analyzer: 'exact_match',
            },
            prefix: {
              type: 'text',
              analyzer: 'prefix_search',
            },
            full: {
              type: 'text',
              analyzer: 'full_text',
            },
            keyword: {
              type: 'keyword',
              normalizer: 'lowercase',
            },
          },
        },
      },
    },
  },
  entity_model: {
    settings: {
      analysis: {
        analyzer: {
          exact_match: {
            type: 'custom',
            tokenizer: 'keyword',
            filter: ['lowercase'],
          },
          prefix_search: {
            type: 'custom',
            tokenizer: 'standard',
            filter: ['lowercase', 'edge_ngram_filter'],
          },
          full_text: {
            type: 'standard',
            filter: ['lowercase'],
          },
        },
        filter: {
          edge_ngram_filter: {
            type: 'edge_ngram',
            min_gram: 3,
            max_gram: 15,
          },
        },
        normalizer: {
          lowercase: {
            type: 'custom',
            filter: ['lowercase'],
          },
        },
      },
    },
    mappings: {
      properties: {
        name: {
          type: 'text',
          fields: {
            exact: {
              type: 'text',
              analyzer: 'exact_match',
            },
            prefix: {
              type: 'text',
              analyzer: 'prefix_search',
            },
            full: {
              type: 'text',
              analyzer: 'full_text',
            },
            keyword: {
              type: 'keyword',
              normalizer: 'lowercase',
            },
          },
        },
      },
    },
  },
};
