import { Injectable } from '@nestjs/common';
import { Client } from '@elastic/elasticsearch';
import { SearchDto, SearchResultDto } from '../dto/search/SearchDto';

@Injectable()
export class ElasticsearchService {
  private readonly client: Client;

  constructor() {
    // private readonly tagsRepository: Repository<Tags>, // @InjectRepository(Tags) // private readonly collectionItemRepository: Repository<CollectionItem>, // @InjectRepository(CollectionItem) // private readonly collectionRepository: Repository<Collection>, // @InjectRepository(Collection) // private readonly userRepository: Repository<User>, // @InjectRepository(User) // private entityRepository: EntityRepository, // @Inject(EntityRepository)
    this.client = new Client({
      node: `http://${process.env.ELASTIC_URL}:9200`,
    });
  }

  /**
   * Поиск по одному индексу
   * @param index - Название индекса
   * @param query - Поисковый запрос
   */
  async search(index: string, query: any) {
    return this.client.search({
      index,
      body: {
        query,
      },
    });
  }

  /**
   * Поиск по нескольким индексам
   * @param indices - Массив индексов
   * @param query - Поисковый запрос
   */
  async searchInIndices(
    indices: string[],
    query: string,
  ): Promise<SearchResultDto> {
    if (!query?.trim()) return new SearchResultDto([]);

    const lowerQuery = query.toLowerCase();

    const searchQuery = {
      bool: {
        should: [
          {
            term: {
              'name.keyword': {
                value: lowerQuery,
              },
            },
          },
          {
            match: {
              'name.prefix': {
                query: lowerQuery,
              },
            },
          },
          {
            match: {
              'name.full': {
                query: lowerQuery,
              },
            },
          },
          {
            wildcard: {
              'name.keyword': {
                value: `*${lowerQuery}*`,
                case_insensitive: true,
              },
            },
          },
        ],
        minimum_should_match: 1,
      },
    };

    const result = await this.client.search({
      index: indices.join(','),
      body: { query: searchQuery },
      size: 100,
    });

    return new SearchResultDto(
      result.hits.hits.map((hit) => new SearchDto(hit._index, hit._source)),
    );
  }

  async indexDocument(index: string, body: any) {
    return this.client.index({
      index,
      body,
    });
  }

  /**
   * Создает индекс с анализатором, если он не существует
   * @param indexName - Название индекса
   * @param mappings - Маппинг для индекса
   */

  Copy;

  async createIndexIfNotExists(indexName: string, indexConfig: any) {
    const indexExists = await this.client.indices.exists({ index: indexName });

    if (!indexExists) {
      const settings = indexConfig.settings || {
        analysis: {
          analyzer: {
            default: {
              type: 'standard',
            },
          },
        },
      };

      await this.client.indices.create({
        index: indexName,
        body: {
          settings,
          mappings: indexConfig.mappings,
        },
      });
      console.log(`Индекс ${indexName} успешно создан`);
    }
  }

  async upsertDocument(index: string, id: string, body: any) {
    return this.client.update({
      index,
      id,
      body: {
        doc: body,
        doc_as_upsert: true,
      },
      refresh: 'wait_for',
    });
  }

  async safeIndexDocument(index: string, id: string, body: any) {
    const exists = await this.client.exists({
      index,
      id,
    });

    if (!exists) {
      return this.client.index({
        index,
        id,
        body,
        refresh: 'wait_for',
      });
    }
    throw new Error(`Document with id ${id} already exists`);
  }

  async deleteByQuery(index: string, query: any) {
    return this.client.deleteByQuery({
      index,
      body: {
        query,
      },
      refresh: true,
    });
  }

  async updateIndex(index: string, id: string, body: any) {
    return this.client.update({
      index,
      id,
      body: {
        doc: body,
      },
      refresh: 'wait_for',
    });
  }

  async deleteDocument(index: string, id: string) {
    try {
      return await this.client.delete({
        index,
        id,
        refresh: 'wait_for',
      });
    } catch (e) {
      if (e.meta?.statusCode === 404) {
        return { result: 'not_found' };
      }
      throw e;
    }
  }
}
