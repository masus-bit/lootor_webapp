import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { ConfigModule } from '@nestjs/config';
import { UserRepository } from '../../repositories/User.repository';
import { User } from '../../entities/User';
import { EventController } from '../../controllers/Event.controller';
import { EventService } from '../../services/Event.service';
import { EventRepository } from '../../repositories/Event.repository';
import { Event } from '../../entities/Event';
import { Collection } from '../../entities/Collection';
import { CollectionRepository } from '../../repositories/Collection.repository';
import { JwtRegisterModule } from './jwt.module';
import { HttpModule } from '@nestjs/axios';
import { ElasticsearchService } from '../../services/ElasticSearch.service';

@Module({
  controllers: [EventController],
  providers: [
    EventService,
    EventRepository,
    UserRepository,
    CollectionRepository,
    ElasticsearchService,
  ],
  imports: [
    TypeOrmModule.forFeature([Event, User, Collection]),
    ConfigModule.forRoot({
      envFilePath: '.env',
    }),
    JwtRegisterModule,
    HttpModule,
  ],
  exports: [EventService],
})
export class EventModule {}
