import {
  EntitySubscriberInterface,
  EventSubscriber,
  InsertEvent,
  RemoveEvent,
  UpdateEvent,
} from 'typeorm';
import { CollectionItem } from '../entities/CollectionItem';
import { ElasticsearchService } from '../services/ElasticSearch.service';

@EventSubscriber()
export class CollectionItemSubscriber
  implements EntitySubscriberInterface<CollectionItem>
{
  constructor(private readonly elasticsearchService: ElasticsearchService) {}

  listenTo() {
    return CollectionItem;
  }

  async afterInsert(event: InsertEvent<CollectionItem>) {
    await this.elasticsearchService.indexDocument('collection_item', {
      id: event.entity.id,
      name: event.entity.name,
    });
  }

  async afterUpdate(event: UpdateEvent<CollectionItem>) {
    await this.elasticsearchService.updateIndex(
      'collection_item',
      event.entity.id.toString(),
      {
        name: event.entity.name,
      },
    );
  }

  async afterRemove(event: RemoveEvent<CollectionItem>) {
    await this.elasticsearchService.deleteDocument(
      'collection_item',
      event.entity.id.toString(),
    );
  }
}
