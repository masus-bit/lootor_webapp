import { Module } from '@nestjs/common';
import { UserModule } from './submodules/User.module';

@Module({
  controllers: [],
  providers: [],
  imports: [UserModule],
})
export class SecuredModule {}
