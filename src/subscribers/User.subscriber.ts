import {
  EntitySubscriberInterface,
  EventSubscriber,
  InsertEvent,
  RemoveEvent,
  UpdateEvent,
} from 'typeorm';
import { User } from '../entities/User';
import { ElasticsearchService } from '../services/ElasticSearch.service';

@EventSubscriber()
export class UserSubscriber implements EntitySubscriberInterface<User> {
  constructor(private readonly elasticsearchService: ElasticsearchService) {}

  listenTo() {
    return User;
  }

  async afterInsert(event: InsertEvent<User>) {
    await this.elasticsearchService.indexDocument('user', {
      login: event.entity.login,
      user_name: event.entity.user_name,
      email: event.entity.email,
    });
  }

  async afterUpdate(event: UpdateEvent<User>) {
    await this.elasticsearchService.updateIndex(
      'user',
      event.entity.id.toString(),
      {
        user_name: event.entity.user_name,
        email: event.entity.email,
      },
    );
  }

  async afterRemove(event: RemoveEvent<User>) {
    await this.elasticsearchService.deleteDocument(
      'user',
      event.entity.login.toString(),
    );
  }
}
