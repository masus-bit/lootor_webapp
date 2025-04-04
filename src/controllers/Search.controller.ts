import { Controller, Get, Query } from '@nestjs/common';
import { ElasticsearchService } from '../services/ElasticSearch.service';
import { HttpBadRequestError } from '../errors/HttpBadRequestError';

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
    console.log(query.search);
    try {
      return await this.elasticsearchService.searchInIndices(
        indices,
        query.search,
      );
    } catch (e) {
      console.log(e);
      throw new HttpBadRequestError(e);
    }
  }
}
