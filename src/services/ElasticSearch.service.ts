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

  // async onApplicationBootstrap() {
  //   await this.reindexAllData();
  // }
  //
  // async reindexAllData() {
  //   const indices = [
  //     'user',
  //     'tags',
  //     'collection_item',
  //     'collection',
  //     'entity_model',
  //     'platforms',
  //   ];
  //   const users = await this.userRepository.find();
  //   const entities = await this.entityRepository.find();
  //   const collections = await this.collectionRepository.find();
  //   const collectionItems = await this.collectionItemRepository.find();
  //   const tags = await this.tagsRepository.find();
  //   const models = {
  //     user: users,
  //     entity_model: entities,
  //     collection: collections,
  //     collection_item: collectionItems,
  //     tags: tags,
  //   };
  //   for (const index of indices) {
  //     const model = models[index];
  //     if (index === 'user') {
  //       const contents = await this.indexDocument(index, {
  //         login: model.login,
  //         user_name: model.user_name,
  //         email: model.email,
  //       });
  //       console.log(contents);
  //     } else if (index === 'entity_model') {
  //       const contents = await this.indexDocument(index, {
  //         name: model.name,
  //         id: model.id,
  //         collection_item: model.collection_item,
  //       });
  //       console.log(contents);
  //     } else {
  //       const contents = await this.indexDocument(index, {
  //         id: model.id,
  //         name: model.name,
  //       });
  //       console.log(contents);
  //     }
  //   }
  // }
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
  async searchInIndices(indices: string[], query: any) {
    const result = await this.client.search({
      index: indices.join(','),
      body: {
        query,
      },
    });
    let objects = [];
    result.hits.hits.forEach((hit) => {
      objects.push(new SearchDto(hit));
    });
    return new SearchResultDto(objects);
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
  async createIndexIfNotExists(indexName: string, mappings: any) {
    const indexExists = await this.client.indices.exists({ index: indexName });

    if (!indexExists) {
      await this.client.indices.create({
        index: indexName,
        body: {
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
                  min_gram: 2, // Минимальная длина n-gram
                  max_gram: 5, // Максимальная длина n-gram
                },
              },
            },
          },
          mappings,
        },
      });
      console.log(
        `Индекс ${indexName} создан с анализатором для частичного поиска.`,
      );
    } else {
      console.log(`Индекс ${indexName} уже существует.`);
    }
  }

  async updateIndex(index: string, id: string, body: any) {
    return this.client.update({
      index,
      id,
      body: {
        doc: body,
      },
    });
  }

  async deleteDocument(index: string, id: string) {
    return this.client.delete({
      index,
      id,
    });
  }
}
