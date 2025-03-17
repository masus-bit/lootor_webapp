export const indices = {
  collection: {
    settings: {
      index: {
        max_ngram_diff: 10,
      },
      analysis: {
        analyzer: {
          common_analyzer: {
            type: 'custom',
            tokenizer: 'standard',
            filter: ['lowercase', 'my_ngram_filter'],
          },
        },
        filter: {
          my_ngram_filter: {
            type: 'ngram',
            min_gram: 2,
            max_gram: 5,
          },
        },
      },
    },
    mappings: {
      properties: {
        name: {
          type: 'text',
          analyzer: 'common_analyzer',
        },
      },
    },
  },
  user: {
    settings: {
      analysis: {
        index: {
          max_ngram_diff: 10,
        },
        analyzer: {
          common_analyzer: {
            type: 'custom',
            tokenizer: 'standard',
            filter: ['lowercase', 'my_ngram_filter'],
          },
        },
        filter: {
          my_ngram_filter: {
            type: 'ngram',
            min_gram: 2,
            max_gram: 5,
          },
        },
      },
    },
    mappings: {
      properties: {
        login: {
          type: 'text',
          analyzer: 'common_analyzer',
        },
        email: {
          type: 'text',
          analyzer: 'common_analyzer',
        },
        user_name: {
          type: 'text',
          analyzer: 'common_analyzer',
        },
      },
    },
  },
  tags: {
    settings: {
      index: {
        max_ngram_diff: 10,
      },
      analysis: {
        analyzer: {
          common_analyzer: {
            type: 'custom',
            tokenizer: 'standard',
            filter: ['lowercase', 'my_ngram_filter'],
          },
        },
        filter: {
          my_ngram_filter: {
            type: 'ngram',
            min_gram: 2,
            max_gram: 5,
          },
        },
      },
    },
    mappings: {
      properties: {
        name: {
          type: 'text',
          analyzer: 'common_analyzer',
        },
      },
    },
  },
  collection_item: {
    settings: {
      index: {
        max_ngram_diff: 10,
      },
      analysis: {
        analyzer: {
          common_analyzer: {
            type: 'custom',
            tokenizer: 'standard',
            filter: ['lowercase', 'my_ngram_filter'],
          },
        },
        filter: {
          my_ngram_filter: {
            type: 'ngram',
            min_gram: 2,
            max_gram: 5,
          },
        },
      },
    },
    mappings: {
      properties: {
        name: {
          type: 'text',
          analyzer: 'common_analyzer',
        },
      },
    },
  },
  entity_model: {
    settings: {
      index: {
        max_ngram_diff: 10,
      },
      analysis: {
        analyzer: {
          common_analyzer: {
            type: 'custom',
            tokenizer: 'standard',
            filter: ['lowercase', 'my_ngram_filter'],
          },
        },
        filter: {
          my_ngram_filter: {
            type: 'ngram',
            min_gram: 2,
            max_gram: 5,
          },
        },
      },
    },
    mappings: {
      properties: {
        name: {
          type: 'text',
          analyzer: 'common_analyzer',
        },
      },
    },
  },
};
