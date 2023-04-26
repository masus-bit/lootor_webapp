import { Module } from '@nestjs/common';
import { UserModule } from './submodules/User.module';
import { LotModule } from './submodules/Lot.module';
import { DraftLot } from '../entities/DraftLot';
import { CollectionModule } from './submodules/Collection.module';
import { CollectionItemModule } from './submodules/CollectionItem.module';
import { EventModule } from './submodules/Event.module';

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
  ],
})
export class SecuredModule {}
