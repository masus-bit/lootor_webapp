import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { ConfigModule } from '@nestjs/config';
import { SecuredModule } from './modules/Secured.module';
import { User } from './entities/User';
import { AuthModule } from './modules/Auth.module';
import { Lot } from './entities/Lot';
import { DraftLot } from './entities/DraftLot';
import { CollectionItem } from './entities/CollectionItem';
import { CollectionImage } from './entities/CollectionImage';
import { LotImage } from './entities/LotImage';
import { Collection } from './entities/Collection';
import { Event } from './entities/Event';

@Module({
  controllers: [],
  providers: [],
  imports: [
    ConfigModule.forRoot({
      envFilePath: '.env',
    }),
    TypeOrmModule.forRoot({
      type: 'postgres',
      host: process.env.POSTGRES_HOST,
      port: Number(process.env.POSTGRES_PORT),
      username: process.env.POSTGRES_USER,
      password: process.env.POSTGRES_PASSWORD,
      database: process.env.POSTGRES_DB,
      schema: process.env.POSTGRES_SCHEMA,
      entities: [
        User,
        Lot,
        DraftLot,
        CollectionItem,
        CollectionImage,
        Collection,
        LotImage,
        Event,
      ],
      synchronize: true,
      autoLoadEntities: true,
    }),
    SecuredModule,
    AuthModule,
  ],
})
export class AppModule {}
