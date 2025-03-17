import {
  EntitySubscriberInterface,
  EventSubscriber,
  InsertEvent,
  RemoveEvent,
  UpdateEvent,
} from 'typeorm';
import { Collection } from '../entities/Collection';
import { ElasticsearchService } from '../services/ElasticSearch.service';

@EventSubscriber()
export class CollectionSubscriber
  implements EntitySubscriberInterface<Collection>
{
  constructor(private readonly elasticsearchService: ElasticsearchService) {}

  listenTo() {
    return Collection;
  }

  async afterInsert(event: InsertEvent<Collection>) {
    await this.elasticsearchService.indexDocument('collection', {
      id: event.entity.id,
      name: event.entity.name,
    });
    console.log('indexed');
  }

  async afterUpdate(event: UpdateEvent<Collection>) {
    await this.elasticsearchService.updateIndex(
      'collection',
      event.entity.id.toString(),
      {
        name: event.entity.name,
      },
    );
  }

  async afterRemove(event: RemoveEvent<Collection>) {
    await this.elasticsearchService.deleteDocument(
      'collection',
      event.entity.id.toString(),
    );
  }
}
