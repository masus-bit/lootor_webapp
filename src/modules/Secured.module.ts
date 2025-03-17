import { Module } from '@nestjs/common';
import { UserModule } from './submodules/User.module';
import { LotModule } from './submodules/Lot.module';
import { DraftLot } from '../entities/DraftLot';
import { CollectionModule } from './submodules/Collection.module';
import { CollectionItemModule } from './submodules/CollectionItem.module';
import { EventModule } from './submodules/Event.module';
import { TagsModule } from './submodules/Tags.module';
import { EntityModule } from './submodules/Entity.module';

@Module({
  controllers: [],
  providers: [],
  imports: [
    UserModule,
    LotModule,
    DraftLot,
    CollectionModule,
    CollectionItemModule,
    EventModule,
    TagsModule,
    EntityModule,
  ],
})
export class SecuredModule {}
