import { Controller, Get, Query } from '@nestjs/common';
import { ElasticsearchService } from '../services/ElasticSearch.service';

@Controller('search')
export class SearchController {
  constructor(private readonly elasticsearchService: ElasticsearchService) {}

  @Get()
  async search(@Query() query: { search: string }) {
    const indices = [
      'user',
      'tags',
      'collection_item',
      'collection',
      'entity_model',
    ];
    const result = await this.elasticsearchService.searchInIndices(indices, {
      multi_match: {
        query: query.search,
        fields: ['*'],
      },
    });
    console.log(result);
    return result.hits.hits;
  }
}
