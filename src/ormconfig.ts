import { DataSource } from 'typeorm';

export const ormconfig = new DataSource({
  migrationsTableName: 'migrations',
  type: 'postgres',
  host: 'host.docker.internal',
  port: 5432,
  username: 'auc',
  password: 'aucer',
  database: 'aucdb',
  logging: false,
  synchronize: false,
  name: 'default',
  entities: ['src/entities/*{.ts,.js}'],
  migrations: ['src/migrations/**/*{.ts,.js}'],
  // subscribers: ['src/subscriber/**/*{.ts,.js}'],
});

// export default new DataSource(ormconfig);
