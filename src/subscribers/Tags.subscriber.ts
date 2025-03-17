import {
  EntitySubscriberInterface,
  EventSubscriber,
  InsertEvent,
  RemoveEvent,
  UpdateEvent,
} from 'typeorm';
import { Tags } from '../entities/Tags';
import { ElasticsearchService } from '../services/ElasticSearch.service';

@EventSubscriber()
export class TagsSubscriber implements EntitySubscriberInterface<Tags> {
  constructor(private readonly elasticsearchService: ElasticsearchService) {}

  listenTo() {
    return Tags;
  }

  async afterInsert(event: InsertEvent<Tags>) {
    await this.elasticsearchService.indexDocument('tags', {
      id: event.entity.id,
      name: event.entity.name,
    });
  }

  async afterUpdate(event: UpdateEvent<Tags>) {
    await this.elasticsearchService.updateIndex(
      'tags',
      event.entity.id.toString(),
      {
        name: event.entity.name,
      },
    );
  }

  async afterRemove(event: RemoveEvent<Tags>) {
    await this.elasticsearchService.deleteDocument(
      'tags',
      event.entity.id.toString(),
    );
  }
}
