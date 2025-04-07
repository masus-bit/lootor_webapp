import { redisStore } from 'cache-manager-redis-store';
import { Module } from '@nestjs/common';
import { CacheModule } from '@nestjs/cache-manager';
import { ConfigModule } from '@nestjs/config';

@Module({
  imports: [
    ConfigModule.forRoot({
      envFilePath: '.env',
    }),
    CacheModule.registerAsync({
      useFactory: async () => ({
        store: redisStore,
        host: process.env.REDIS_HOST,
        port: process.env.REDIS_PORT,
        ttl: +process.env.REDIS_TTL,
      }),
      isGlobal: true,
    }),
  ],
  exports: [CacheModule],
})
export class CacheRedisModule {}
