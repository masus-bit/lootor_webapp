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

let variant = 'stage';

@Module({
  controllers: [],
  providers: [],
  imports: [
    ConfigModule.forRoot({
      envFilePath: '.env',
    }),
    TypeOrmModule.forRoot({
      type: 'postgres',
      host: variant === 'dev' ? process.env.POSTGRES_HOST : 'localhost:5432',
      port: Number(process.env.POSTGRES_PORT),
      username: 'auc',
      password: 'aucer',
      database: 'aucdb',
      schema: 'public',
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
