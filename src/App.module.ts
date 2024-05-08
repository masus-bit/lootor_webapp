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
import { Tags } from './entities/Tags';
import { Games } from './entities/Games';
import { Publishers } from './entities/Publishers';
import { Genres } from './entities/Genres';
import { Banners } from './entities/Banners';
import { Platforms } from './entities/Platforms';
import { Developers } from './entities/Developers';

console.log(
  process.env.POSTGRES_HOST,
  process.env.POSTGRES_PORT,
  process.env.POSTGRES_USER,
  process.env.POSTGRES_PASSWORD,
  process.env.POSTGRES_DB,
);
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
      migrations: ['src/migrations/**/*{.ts,.js}'],
      migrationsRun: true,
      entities: [
        User,
        Lot,
        DraftLot,
        CollectionItem,
        CollectionImage,
        Collection,
        LotImage,
        Event,
        Tags,
        Games,
        Publishers,
        Genres,
        Banners,
        Platforms,
        Developers,
      ],
      synchronize: false,
      autoLoadEntities: true,
    }),
    SecuredModule,
    AuthModule,
  ],
})
export class AppModule {}
