import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { ConfigModule } from '@nestjs/config';
import { UserRepository } from '../../repositories/User.repository';
import { User } from '../../entities/User';
import { JwtRegisterModule } from './jwt.module';
import { HttpModule } from '@nestjs/axios';
import { EntityService } from '../../services/Entity.service';
import { EntityController } from '../../controllers/Entity.controller';
import { EntityRepository } from '../../repositories/Entity.repository';
import { CollectionItem } from '../../entities/CollectionItem';
import { CollectionItemRepository } from '../../repositories/CollectionItem.repository';
import { EntityModel } from '../../entities/EntityModel';

@Module({
  controllers: [EntityController],
  providers: [
    EntityRepository,
    EntityService,
    UserRepository,
    CollectionItemRepository,
  ],
  imports: [
    TypeOrmModule.forFeature([User, CollectionItem, EntityModel]),
    ConfigModule.forRoot({
      envFilePath: '.env',
    }),
    JwtRegisterModule,
    HttpModule,
  ],
  exports: [EntityService],
})
export class EntityModule {}
