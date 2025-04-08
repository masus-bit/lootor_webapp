import { Module } from '@nestjs/common';
import { MulterModule } from '@nestjs/platform-express';
import { UserRepository } from '../../repositories/User.repository';
import { TypeOrmModule } from '@nestjs/typeorm';
import { User } from '../../entities/User';
import { ConfigModule } from '@nestjs/config';
import { JwtRegisterModule } from './jwt.module';
import { HttpModule } from '@nestjs/axios';
import { ElasticsearchService } from '../../services/ElasticSearch.service';
import { TrelloService } from '../../services/Trello.service';
import { FeedbackController } from '../../controllers/Feedback.controller';

@Module({
  imports: [
    MulterModule.register(),
    TypeOrmModule.forFeature([User]),
    ConfigModule.forRoot({
      envFilePath: '.env',
    }),
    JwtRegisterModule,
    HttpModule,
  ],
  controllers: [FeedbackController],
  providers: [TrelloService, UserRepository, ElasticsearchService],
})
export class FeedbackModule {}
