import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { ConfigModule } from '@nestjs/config';
import { UserRepository } from '../../repositories/User.repository';
import { User } from '../../entities/User';
import { JwtRegisterModule } from './jwt.module';
import { HttpModule } from '@nestjs/axios';
import { GamesRepository } from '../../repositories/Games.repository';
import { GamesService } from '../../services/Games.service';
import { Games } from '../../entities/Games';
import { GamesController } from '../../controllers/Games.controller';

@Module({
  controllers: [GamesController],
  providers: [UserRepository, GamesRepository, GamesService],
  imports: [
    TypeOrmModule.forFeature([User, Games]),
    ConfigModule.forRoot({
      envFilePath: '.env',
    }),
    JwtRegisterModule,
    HttpModule,
  ],
  exports: [GamesService],
})
export class GamesModule {}
