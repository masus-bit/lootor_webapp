import { Module } from '@nestjs/common';
import { UserModule } from './submodules/User.module';
import { LotModule } from './submodules/Lot.module';
import { DraftLot } from '../entities/DraftLot';
import { ExternalModule } from './submodules/External.module';

@Module({
  controllers: [],
  providers: [],
  imports: [UserModule, LotModule, DraftLot, ExternalModule],
})
export class SecuredModule {}
