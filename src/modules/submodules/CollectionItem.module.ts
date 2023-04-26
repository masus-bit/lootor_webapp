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

@Module({
  controllers: [CollectionItemController],
  providers: [
    CollectionItemService,
    CollectionItemRepository,
    UserRepository,
    EventRepository,
  ],
  imports: [
    TypeOrmModule.forFeature([CollectionItem, User, Event]),
    ConfigModule.forRoot({
      envFilePath: '.env',
    }),
    JwtRegisterModule,
    HttpModule,
  ],
  exports: [CollectionItemService],
})
export class CollectionItemModule {}
