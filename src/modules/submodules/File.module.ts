import { Module } from '@nestjs/common';
import { MulterModule } from '@nestjs/platform-express';
import { S3Service } from '../../services/s3.service';
import { FilesController } from '../../controllers/Files.controller';
import { UserRepository } from '../../repositories/User.repository';
import { TypeOrmModule } from '@nestjs/typeorm';
import { User } from '../../entities/User';
import { ConfigModule } from '@nestjs/config';
import { JwtRegisterModule } from './jwt.module';
import { HttpModule } from '@nestjs/axios';
import { ElasticsearchService } from '../../services/ElasticSearch.service';
import { CacheRedisModule } from './Cache.module';

@Module({
  imports: [
    MulterModule.register(),
    TypeOrmModule.forFeature([User]),
    ConfigModule.forRoot({
      envFilePath: '.env',
    }),
    JwtRegisterModule,
    HttpModule,
    CacheRedisModule,
  ],
  controllers: [FilesController],
  providers: [S3Service, UserRepository, ElasticsearchService],
})
export class FileModule {}
