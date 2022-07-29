import { Module } from '@nestjs/common';
import { UserModule } from './submodules/User.module';
import { LotModule } from './submodules/Lot.module';
import { DraftLot } from '../entities/DraftLot';

@Module({
  controllers: [],
  providers: [],
  imports: [UserModule, LotModule, DraftLot],
})
export class SecuredModule {}
