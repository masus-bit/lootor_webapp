import { Module } from '@nestjs/common';
import { UserModule } from './submodules/User.module';
import { LotModule } from './submodules/Lot.module';
import { DraftLot } from '../entities/DraftLot';
import { ExternalModule } from './submodules/External.module';
import { CollectionModule } from './submodules/Collection.module';

@Module({
  controllers: [],
  providers: [],
  imports: [UserModule, LotModule, DraftLot, ExternalModule, CollectionModule],
})
export class SecuredModule {}
