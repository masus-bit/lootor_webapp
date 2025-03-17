import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { ConfigModule } from '@nestjs/config';
import { UserRepository } from '../../repositories/User.repository';
import { User } from '../../entities/User';
import { CollectionItemController } from '../../controllers/CollectionItem.controller';
import { CollectionItemService } from '../../services/CollectionItem.service';
import { CollectionItemRepository } from '../../repositories/CollectionItem.repository';
import { CollectionItem } from '../../entities/CollectionItem';
import { JwtRegisterModule } from './jwt.module';
import { HttpModule } from '@nestjs/axios';
import { EventRepository } from '../../repositories/Event.repository';
import { Event } from '../../entities/Event';
import { EntityRepository } from '../../repositories/Entity.repository';
import { EntityModel } from '../../entities/EntityModel';
import { EntityService } from '../../services/Entity.service';

@Module({
  controllers: [CollectionItemController],
  providers: [
    CollectionItemService,
    CollectionItemRepository,
    UserRepository,
    EventRepository,
    EntityRepository,
    EntityService,
  ],
  imports: [
    TypeOrmModule.forFeature([CollectionItem, User, Event, EntityModel]),
    ConfigModule.forRoot({
      envFilePath: '.env',
    }),
    JwtRegisterModule,
    HttpModule,
  ],
  exports: [CollectionItemService],
})
export class CollectionItemModule {}
