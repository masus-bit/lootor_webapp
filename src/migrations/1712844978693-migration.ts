import { MigrationInterface, QueryRunner } from 'typeorm';

export class Migration1712844978693 implements MigrationInterface {
  name = 'Migration1712844978693';

  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(`ALTER TABLE "event" DROP COLUMN "date"`);
    await queryRunner.query(`ALTER TABLE "event" ADD "date" bigint`);
  }

  public async down(queryRunner: QueryRunner): Promise<void> {}
}
