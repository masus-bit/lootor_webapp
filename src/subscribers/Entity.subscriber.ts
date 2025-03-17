import {
  EntitySubscriberInterface,
  EventSubscriber,
  InsertEvent,
  RemoveEvent,
  UpdateEvent,
} from 'typeorm';
import { EntityModel } from '../entities/EntityModel';
import { ElasticsearchService } from '../services/ElasticSearch.service';

@EventSubscriber()
export class EntitySubscriber
  implements EntitySubscriberInterface<EntityModel>
{
  constructor(private readonly elasticsearchService: ElasticsearchService) {}

  listenTo() {
    return EntityModel;
  }

  async afterInsert(event: InsertEvent<EntityModel>) {
    await this.elasticsearchService.indexDocument('entity_model', {
      id: event.entity.id,
      name: event.entity.name,
      collection_item: event.entity.collection_item,
    });
  }

  async afterUpdate(event: UpdateEvent<EntityModel>) {
    await this.elasticsearchService.updateIndex(
      'entity_model',
      event.entity.id.toString(),
      {
        name: event.entity.name,
        collection_item: event.entity.collection_item,
      },
    );
  }

  async afterRemove(event: RemoveEvent<EntityModel>) {
    await this.elasticsearchService.deleteDocument(
      'entity_model',
      event.entity.id.toString(),
    );
  }
}
