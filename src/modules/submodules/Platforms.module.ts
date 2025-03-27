import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { ConfigModule } from '@nestjs/config';
import { JwtRegisterModule } from './jwt.module';
import { HttpModule } from '@nestjs/axios';
import { PlatformsRepository } from '../../repositories/Platforms.repository';
import { PlatformsService } from '../../services/Platforms.service';
import { PlatformsController } from '../../controllers/Platforms.controller';
import { Platforms } from '../../entities/Platforms';
import { UserRepository } from '../../repositories/User.repository';
import { User } from '../../entities/User';
import { ElasticsearchService } from '../../services/ElasticSearch.service';

@Module({
  controllers: [PlatformsController],
  providers: [
    PlatformsRepository,
    PlatformsService,
    UserRepository,
    ElasticsearchService,
  ],
  imports: [
    TypeOrmModule.forFeature([Platforms, User]),
    ConfigModule.forRoot({
      envFilePath: '.env',
    }),
    JwtRegisterModule,
    HttpModule,
  ],
  exports: [PlatformsService],
})
export class PlatformsModule {}
