import { MigrationInterface, QueryRunner } from 'typeorm';

export class Migration1706272910642 implements MigrationInterface {
  name = 'Migration1706272910642';

  public async up(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(`DROP INDEX "public"."platforms_name_idx"`);
    await queryRunner.query(`DROP INDEX "public"."platforms_alias_idx"`);
    await queryRunner.query(`DROP INDEX "public"."platforms_id_idx"`);
    await queryRunner.query(`DROP INDEX "public"."games_platform_idx"`);
    await queryRunner.query(`DROP INDEX "public"."games_release_date_idx"`);
    await queryRunner.query(`DROP INDEX "public"."games_game_title_idx"`);
    await queryRunner.query(
      `DROP INDEX "public"."games_id_game_title_SOUNDEX_idx"`,
    );
    await queryRunner.query(`DROP INDEX "public"."games_last_updated_idx"`);
    await queryRunner.query(`DROP INDEX "public"."games_SOUNDEX_idx"`);
    await queryRunner.query(`DROP INDEX "public"."games_id_idx"`);
    await queryRunner.query(`DROP INDEX "public"."banners_userid_idx"`);
    await queryRunner.query(`DROP INDEX "public"."banners_games_id_idx"`);
    await queryRunner.query(
      `DROP INDEX "public"."banners_type_side_games_id_idx"`,
    );
    await queryRunner.query(`DROP INDEX "public"."banners_type_side_idx"`);
    await queryRunner.query(`DROP INDEX "public"."banners_type_games_id_idx"`);
    await queryRunner.query(`DROP INDEX "public"."banners_mirrorupdate_idx"`);
    await queryRunner.query(
      `DROP INDEX "public"."games_pubs_games_id_pub_id_idx"`,
    );
    await queryRunner.query(`DROP INDEX "public"."games_pubs_pub_id_idx"`);
    await queryRunner.query(`DROP INDEX "public"."games_pubs_games_id_idx"`);
    await queryRunner.query(`DROP INDEX "public"."games_genre_games_id_idx"`);
    await queryRunner.query(
      `DROP INDEX "public"."games_genre_games_id_genres_id_idx"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."games_devs_games_id_dev_id_idx"`,
    );
    await queryRunner.query(`DROP INDEX "public"."games_devs_dev_id_idx"`);
    await queryRunner.query(`DROP INDEX "public"."games_devs_games_id_idx"`);
    await queryRunner.query(`ALTER TABLE "platforms" DROP COLUMN "youtube"`);
    await queryRunner.query(`ALTER TABLE "games" DROP COLUMN "SOUNDEX"`);
    await queryRunner.query(
      `ALTER TABLE "games" DROP COLUMN "publisher_to_be_removed"`,
    );
    await queryRunner.query(`ALTER TABLE "games" DROP COLUMN "last_updated"`);
    await queryRunner.query(`ALTER TABLE "games" DROP COLUMN "hits"`);
    await queryRunner.query(`ALTER TABLE "games" DROP COLUMN "mirrorupdate"`);
    await queryRunner.query(`ALTER TABLE "games" DROP COLUMN "disabled"`);
    await queryRunner.query(`ALTER TABLE "games" DROP COLUMN "youtube"`);
    await queryRunner.query(`ALTER TABLE "games" DROP COLUMN "os"`);
    await queryRunner.query(`ALTER TABLE "games" DROP COLUMN "processor"`);
    await queryRunner.query(`ALTER TABLE "games" DROP COLUMN "ram"`);
    await queryRunner.query(`ALTER TABLE "games" DROP COLUMN "hdd"`);
    await queryRunner.query(`ALTER TABLE "games" DROP COLUMN "video"`);
    await queryRunner.query(`ALTER TABLE "games" DROP COLUMN "sound"`);
    await queryRunner.query(
      `ALTER TABLE "games" DROP COLUMN "alternates_to_be_removed"`,
    );
    await queryRunner.query(`ALTER TABLE "games" DROP COLUMN "region_id"`);
    await queryRunner.query(`ALTER TABLE "games" DROP COLUMN "country_id"`);
    await queryRunner.query(`ALTER TABLE "banners" DROP COLUMN "userid"`);
    await queryRunner.query(`ALTER TABLE "banners" DROP COLUMN "subkey"`);
    await queryRunner.query(`ALTER TABLE "banners" DROP COLUMN "username"`);
    await queryRunner.query(`ALTER TABLE "banners" DROP COLUMN "dateadded"`);
    await queryRunner.query(`ALTER TABLE "banners" DROP COLUMN "languageid"`);
    await queryRunner.query(`ALTER TABLE "banners" DROP COLUMN "colors"`);
    await queryRunner.query(`ALTER TABLE "banners" DROP COLUMN "artistcolors"`);
    await queryRunner.query(`ALTER TABLE "banners" DROP COLUMN "mirrorupdate"`);
    await queryRunner.query(`ALTER TABLE "games_developers" DROP COLUMN "id"`);
    await queryRunner.query(
      `ALTER TABLE "user" ADD "background_url" character varying(100)`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_8ce4c93ba419b56bd82e533724"`,
    );
    await queryRunner.query(`ALTER TABLE "user" DROP COLUMN "created"`);
    await queryRunner.query(`ALTER TABLE "user" ADD "created" integer`);
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(
      `ALTER TABLE "games_developers" DROP CONSTRAINT "FK_82b22acd3c7cfb07e86f80ac1be"`,
    );
    await queryRunner.query(
      `ALTER TABLE "games_developers" DROP CONSTRAINT "FK_4b652919c1ca66282cd4a3908c6"`,
    );
    await queryRunner.query(
      `ALTER TABLE "games_genres" DROP CONSTRAINT "FK_4eccc3295cc5eaf8819b72d4b11"`,
    );
    await queryRunner.query(
      `ALTER TABLE "games_genres" DROP CONSTRAINT "FK_7a86d48f933b8b50effe3f0895b"`,
    );
    await queryRunner.query(
      `ALTER TABLE "games_publishers" DROP CONSTRAINT "FK_0e6215a3bc6ef84f6cf06248132"`,
    );
    await queryRunner.query(
      `ALTER TABLE "games_publishers" DROP CONSTRAINT "FK_37fbcb4a3045e193a53aa2da881"`,
    );
    await queryRunner.query(
      `ALTER TABLE "banners" DROP CONSTRAINT "FK_3a628aa4fd1aa55ffe301ac21fa"`,
    );
    await queryRunner.query(
      `ALTER TABLE "games" DROP CONSTRAINT "FK_eb30b2077ef701aedf4e582d13b"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_82b22acd3c7cfb07e86f80ac1b"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_4b652919c1ca66282cd4a3908c"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_4eccc3295cc5eaf8819b72d4b1"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_7a86d48f933b8b50effe3f0895"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_0e6215a3bc6ef84f6cf0624813"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_37fbcb4a3045e193a53aa2da88"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_86d7a5f06453192c229dd46f07"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_ea68f6f4986d699870c0ba9535"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_756944c85c84864d2b8d3a569d"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_3aaaace8c25a64f97dcd62de11"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_9ba1b2a1f42a1d6f117a6c0451"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_85997e3bb458a88e99397c34ed"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_a53b78cf6c42a5f4a76627c627"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_73ba1ed6fb00460e1a1152a363"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_c9dc5ccf8bbc751ce4c5dfc9bc"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_4e0bda4733ce391051f6f7038b"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_6e888f1270eebe1d8b2595901b"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_a01d172e2d7c8883b0b6f82134"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_fa13a0bb9ce359a8a785a20b71"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_d35c01e0483e04066a26b3f144"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_3de2990b681cf13377814a4c89"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_1a3d733ce21005e988a1142a4f"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_a3b00f5cfe80c2cba9fc10c017"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_413604b47989fd289e32029341"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_e86dfa5698c2a75677d369996c"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_e82b74d9958a23d80b78d9f44e"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_17d3eb985cd27e378a343e2db4"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_a68f2de573348166a19e1f7f54"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_c6817300e18021cd3fc418c4ef"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_4cb9051aeab137596d7c3d32e4"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_ad82d08569c03090581232c9d2"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_6add27e349b6905c85e016fa2c"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_778e59ce8961a0c7c8e6534f12"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_c665119f34b09fcf2787985c83"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_a47be47b362a39a6d6cb9c619a"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_8677af0a1db414d6997eddc91b"`,
    );
    await queryRunner.query(
      `DROP INDEX "public"."IDX_8ce4c93ba419b56bd82e533724"`,
    );
    await queryRunner.query(
      `ALTER TABLE "games_publishers" DROP CONSTRAINT "PK_7fd676529cfe2d1bd08d7473a58"`,
    );
    await queryRunner.query(
      `ALTER TABLE "games_publishers" ADD CONSTRAINT "PK_37fbcb4a3045e193a53aa2da881" PRIMARY KEY ("gamesId")`,
    );
    await queryRunner.query(
      `ALTER TABLE "games_publishers" DROP COLUMN "publishersId"`,
    );
    await queryRunner.query(
      `ALTER TABLE "games_publishers" ADD "publishersId" bigint NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "games_publishers" DROP CONSTRAINT "PK_37fbcb4a3045e193a53aa2da881"`,
    );
    await queryRunner.query(
      `ALTER TABLE "games_publishers" ADD CONSTRAINT "PK_7fd676529cfe2d1bd08d7473a58" PRIMARY KEY ("publishersId", "gamesId")`,
    );
    await queryRunner.query(
      `ALTER TABLE "games_publishers" DROP CONSTRAINT "PK_7fd676529cfe2d1bd08d7473a58"`,
    );
    await queryRunner.query(
      `ALTER TABLE "games_publishers" ADD CONSTRAINT "PK_0e6215a3bc6ef84f6cf06248132" PRIMARY KEY ("publishersId")`,
    );
    await queryRunner.query(
      `ALTER TABLE "games_publishers" DROP COLUMN "gamesId"`,
    );
    await queryRunner.query(
      `ALTER TABLE "games_publishers" ADD "gamesId" bigint NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "games_publishers" DROP CONSTRAINT "PK_0e6215a3bc6ef84f6cf06248132"`,
    );
    await queryRunner.query(
      `ALTER TABLE "games_publishers" ADD CONSTRAINT "PK_7fd676529cfe2d1bd08d7473a58" PRIMARY KEY ("gamesId", "publishersId")`,
    );
    await queryRunner.query(`ALTER TABLE "banners" DROP COLUMN "games_id"`);
    await queryRunner.query(
      `ALTER TABLE "banners" ADD "games_id" bigint NOT NULL DEFAULT '0'`,
    );
    await queryRunner.query(`ALTER TABLE "banners" DROP COLUMN "resolution"`);
    await queryRunner.query(
      `ALTER TABLE "banners" ADD "resolution" character varying(9)`,
    );
    await queryRunner.query(`ALTER TABLE "banners" DROP COLUMN "filename"`);
    await queryRunner.query(
      `ALTER TABLE "banners" ADD "filename" character varying(255) NOT NULL DEFAULT ''`,
    );
    await queryRunner.query(`ALTER TABLE "banners" DROP COLUMN "side"`);
    await queryRunner.query(
      `ALTER TABLE "banners" ADD "side" character varying(16)`,
    );
    await queryRunner.query(`ALTER TABLE "banners" DROP COLUMN "type"`);
    await queryRunner.query(
      `ALTER TABLE "banners" ADD "type" character varying(16) NOT NULL DEFAULT ''`,
    );
    await queryRunner.query(
      `ALTER TABLE "banners" DROP CONSTRAINT "PK_e9b186b959296fcb940790d31c3"`,
    );
    await queryRunner.query(`ALTER TABLE "banners" DROP COLUMN "id"`);
    await queryRunner.query(
      `ALTER TABLE "banners" ADD "id" bigint GENERATED BY DEFAULT AS IDENTITY NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "banners" ADD CONSTRAINT "banners_pkey" PRIMARY KEY ("id")`,
    );
    await queryRunner.query(`ALTER TABLE "developers" DROP COLUMN "name"`);
    await queryRunner.query(
      `ALTER TABLE "developers" ADD "name" character varying(150) NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "games" ALTER COLUMN "platform" SET NOT NULL`,
    );
    await queryRunner.query(`ALTER TABLE "games" DROP COLUMN "coop"`);
    await queryRunner.query(
      `ALTER TABLE "games" ADD "coop" character varying(10)`,
    );
    await queryRunner.query(`ALTER TABLE "games" DROP COLUMN "rating"`);
    await queryRunner.query(
      `ALTER TABLE "games" ADD "rating" character varying(45)`,
    );
    await queryRunner.query(`ALTER TABLE "games" DROP COLUMN "overview"`);
    await queryRunner.query(`ALTER TABLE "games" ADD "overview" text`);
    await queryRunner.query(`ALTER TABLE "games" DROP COLUMN "release_date"`);
    await queryRunner.query(`ALTER TABLE "games" ADD "release_date" date`);
    await queryRunner.query(`ALTER TABLE "games" DROP COLUMN "players"`);
    await queryRunner.query(`ALTER TABLE "games" ADD "players" smallint`);
    await queryRunner.query(`ALTER TABLE "games" DROP COLUMN "game_title"`);
    await queryRunner.query(
      `ALTER TABLE "games" ADD "game_title" character varying(255) NOT NULL DEFAULT ''`,
    );
    await queryRunner.query(
      `ALTER TABLE "games" DROP CONSTRAINT "PK_c9b16b62917b5595af982d66337"`,
    );
    await queryRunner.query(`ALTER TABLE "games" DROP COLUMN "id"`);
    await queryRunner.query(
      `ALTER TABLE "games" ADD "id" bigint GENERATED BY DEFAULT AS IDENTITY NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "games" ADD CONSTRAINT "games_pkey" PRIMARY KEY ("id")`,
    );
    await queryRunner.query(`ALTER TABLE "platforms" DROP COLUMN "overview"`);
    await queryRunner.query(`ALTER TABLE "platforms" ADD "overview" text`);
    await queryRunner.query(`ALTER TABLE "platforms" DROP COLUMN "display"`);
    await queryRunner.query(`ALTER TABLE "platforms" ADD "display" text`);
    await queryRunner.query(
      `ALTER TABLE "platforms" DROP COLUMN "maxcontrollers"`,
    );
    await queryRunner.query(
      `ALTER TABLE "platforms" ADD "maxcontrollers" text`,
    );
    await queryRunner.query(`ALTER TABLE "platforms" DROP COLUMN "sound"`);
    await queryRunner.query(`ALTER TABLE "platforms" ADD "sound" text`);
    await queryRunner.query(`ALTER TABLE "platforms" DROP COLUMN "graphics"`);
    await queryRunner.query(`ALTER TABLE "platforms" ADD "graphics" text`);
    await queryRunner.query(`ALTER TABLE "platforms" DROP COLUMN "memory"`);
    await queryRunner.query(`ALTER TABLE "platforms" ADD "memory" text`);
    await queryRunner.query(`ALTER TABLE "platforms" DROP COLUMN "cpu"`);
    await queryRunner.query(`ALTER TABLE "platforms" ADD "cpu" text`);
    await queryRunner.query(`ALTER TABLE "platforms" DROP COLUMN "media"`);
    await queryRunner.query(`ALTER TABLE "platforms" ADD "media" text`);
    await queryRunner.query(
      `ALTER TABLE "platforms" DROP COLUMN "manufacturer"`,
    );
    await queryRunner.query(`ALTER TABLE "platforms" ADD "manufacturer" text`);
    await queryRunner.query(`ALTER TABLE "platforms" DROP COLUMN "developer"`);
    await queryRunner.query(`ALTER TABLE "platforms" ADD "developer" text`);
    await queryRunner.query(`ALTER TABLE "platforms" DROP COLUMN "controller"`);
    await queryRunner.query(
      `ALTER TABLE "platforms" ADD "controller" character varying(100)`,
    );
    await queryRunner.query(`ALTER TABLE "platforms" DROP COLUMN "console"`);
    await queryRunner.query(
      `ALTER TABLE "platforms" ADD "console" character varying(100)`,
    );
    await queryRunner.query(`ALTER TABLE "platforms" DROP COLUMN "icon"`);
    await queryRunner.query(
      `ALTER TABLE "platforms" ADD "icon" character varying(100) NOT NULL`,
    );
    await queryRunner.query(`ALTER TABLE "platforms" DROP COLUMN "alias"`);
    await queryRunner.query(
      `ALTER TABLE "platforms" ADD "alias" character varying(100)`,
    );
    await queryRunner.query(`ALTER TABLE "platforms" DROP COLUMN "name"`);
    await queryRunner.query(
      `ALTER TABLE "platforms" ADD "name" character varying(100) NOT NULL DEFAULT ''`,
    );
    await queryRunner.query(
      `ALTER TABLE "platforms" DROP CONSTRAINT "PK_3b879853678f7368d46e52b81c6"`,
    );
    await queryRunner.query(`ALTER TABLE "platforms" DROP COLUMN "id"`);
    await queryRunner.query(
      `ALTER TABLE "platforms" ADD "id" bigint GENERATED BY DEFAULT AS IDENTITY NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "platforms" ADD CONSTRAINT "platforms_pkey" PRIMARY KEY ("id")`,
    );
    await queryRunner.query(`ALTER TABLE "genres" DROP COLUMN "genre"`);
    await queryRunner.query(
      `ALTER TABLE "genres" ADD "genre" character varying(100) NOT NULL DEFAULT ''`,
    );
    await queryRunner.query(
      `ALTER TABLE "genres" DROP CONSTRAINT "PK_80ecd718f0f00dde5d77a9be842"`,
    );
    await queryRunner.query(`ALTER TABLE "genres" DROP COLUMN "id"`);
    await queryRunner.query(
      `ALTER TABLE "genres" ADD "id" bigint GENERATED BY DEFAULT AS IDENTITY NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "genres" ADD CONSTRAINT "genres_pkey" PRIMARY KEY ("id")`,
    );
    await queryRunner.query(`ALTER TABLE "publishers" DROP COLUMN "logo"`);
    await queryRunner.query(
      `ALTER TABLE "publishers" ADD "logo" character varying(512) NOT NULL`,
    );
    await queryRunner.query(`ALTER TABLE "publishers" DROP COLUMN "keyword"`);
    await queryRunner.query(
      `ALTER TABLE "publishers" ADD "keyword" character varying(255) NOT NULL`,
    );
    await queryRunner.query(`ALTER TABLE "user" DROP COLUMN "created"`);
    await queryRunner.query(`ALTER TABLE "user" ADD "created" TIMESTAMP`);
    await queryRunner.query(
      `CREATE INDEX "IDX_8ce4c93ba419b56bd82e533724" ON "user" ("created") `,
    );
    await queryRunner.query(
      `ALTER TABLE "games_developers" DROP CONSTRAINT "PK_16d8e01c6a9cdcc6b5ddfe73c75"`,
    );
    await queryRunner.query(
      `ALTER TABLE "games_genres" DROP CONSTRAINT "PK_17791367551a00fc8fe97b6710c"`,
    );
    await queryRunner.query(
      `ALTER TABLE "games_publishers" DROP CONSTRAINT "PK_7fd676529cfe2d1bd08d7473a58"`,
    );
    await queryRunner.query(`ALTER TABLE "user" DROP COLUMN "background_url"`);
    await queryRunner.query(
      `ALTER TABLE "games_developers" ADD "id" integer GENERATED BY DEFAULT AS IDENTITY NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "games_developers" ADD CONSTRAINT "games_devs_pkey" PRIMARY KEY ("id")`,
    );
    await queryRunner.query(
      `ALTER TABLE "games_genres" ADD "id" integer GENERATED BY DEFAULT AS IDENTITY NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "games_genres" ADD CONSTRAINT "games_genre_pkey" PRIMARY KEY ("id")`,
    );
    await queryRunner.query(
      `ALTER TABLE "games_publishers" ADD "id" integer GENERATED BY DEFAULT AS IDENTITY NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "games_publishers" ADD CONSTRAINT "games_pubs_pkey" PRIMARY KEY ("id")`,
    );
    await queryRunner.query(
      `ALTER TABLE "banners" ADD "mirrorupdate" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP`,
    );
    await queryRunner.query(
      `ALTER TABLE "banners" ADD "artistcolors" character varying(255)`,
    );
    await queryRunner.query(
      `ALTER TABLE "banners" ADD "colors" character varying(255)`,
    );
    await queryRunner.query(
      `ALTER TABLE "banners" ADD "languageid" integer NOT NULL DEFAULT '7'`,
    );
    await queryRunner.query(`ALTER TABLE "banners" ADD "dateadded" bigint`);
    await queryRunner.query(
      `ALTER TABLE "banners" ADD "username" character varying(45)`,
    );
    await queryRunner.query(
      `ALTER TABLE "banners" ADD "subkey" character varying(16)`,
    );
    await queryRunner.query(
      `ALTER TABLE "banners" ADD "userid" bigint NOT NULL DEFAULT '1'`,
    );
    await queryRunner.query(
      `ALTER TABLE "games" ADD "country_id" integer NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "games" ADD "region_id" integer NOT NULL`,
    );
    await queryRunner.query(
      `ALTER TABLE "games" ADD "alternates_to_be_removed" text`,
    );
    await queryRunner.query(
      `ALTER TABLE "games" ADD "sound" character varying(255)`,
    );
    await queryRunner.query(
      `ALTER TABLE "games" ADD "video" character varying(255)`,
    );
    await queryRunner.query(
      `ALTER TABLE "games" ADD "hdd" character varying(255)`,
    );
    await queryRunner.query(
      `ALTER TABLE "games" ADD "ram" character varying(255)`,
    );
    await queryRunner.query(
      `ALTER TABLE "games" ADD "processor" character varying(255)`,
    );
    await queryRunner.query(
      `ALTER TABLE "games" ADD "os" character varying(255)`,
    );
    await queryRunner.query(
      `ALTER TABLE "games" ADD "youtube" character varying(255)`,
    );
    await queryRunner.query(
      `ALTER TABLE "games" ADD "disabled" character varying(3) NOT NULL DEFAULT 'No'`,
    );
    await queryRunner.query(
      `ALTER TABLE "games" ADD "mirrorupdate" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP`,
    );
    await queryRunner.query(
      `ALTER TABLE "games" ADD "hits" bigint DEFAULT '0'`,
    );
    await queryRunner.query(`ALTER TABLE "games" ADD "last_updated" TIMESTAMP`);
    await queryRunner.query(
      `ALTER TABLE "games" ADD "publisher_to_be_removed" character varying(100)`,
    );
    await queryRunner.query(
      `ALTER TABLE "games" ADD "SOUNDEX" character varying(255)`,
    );
    await queryRunner.query(
      `ALTER TABLE "platforms" ADD "youtube" character varying(255)`,
    );
    await queryRunner.query(
      `CREATE INDEX "games_devs_games_id_idx" ON "games_developers" ("gamesId") `,
    );
    await queryRunner.query(
      `CREATE INDEX "games_devs_dev_id_idx" ON "games_developers" ("developersId") `,
    );
    await queryRunner.query(
      `CREATE UNIQUE INDEX "games_devs_games_id_dev_id_idx" ON "games_developers" ("gamesId", "developersId") `,
    );
    await queryRunner.query(
      `CREATE UNIQUE INDEX "games_genre_games_id_genres_id_idx" ON "games_genres" ("gamesId", "genresId") `,
    );
    await queryRunner.query(
      `CREATE INDEX "games_genre_games_id_idx" ON "games_genres" ("gamesId") `,
    );
    await queryRunner.query(
      `CREATE INDEX "games_pubs_games_id_idx" ON "games_publishers" ("gamesId") `,
    );
    await queryRunner.query(
      `CREATE INDEX "games_pubs_pub_id_idx" ON "games_publishers" ("publishersId") `,
    );
    await queryRunner.query(
      `CREATE UNIQUE INDEX "games_pubs_games_id_pub_id_idx" ON "games_publishers" ("gamesId", "publishersId") `,
    );
    await queryRunner.query(
      `CREATE INDEX "banners_mirrorupdate_idx" ON "banners" ("mirrorupdate") `,
    );
    await queryRunner.query(
      `CREATE INDEX "banners_type_games_id_idx" ON "banners" ("type", "games_id") `,
    );
    await queryRunner.query(
      `CREATE INDEX "banners_type_side_idx" ON "banners" ("type", "side") `,
    );
    await queryRunner.query(
      `CREATE INDEX "banners_type_side_games_id_idx" ON "banners" ("type", "side", "games_id") `,
    );
    await queryRunner.query(
      `CREATE INDEX "banners_games_id_idx" ON "banners" ("games_id") `,
    );
    await queryRunner.query(
      `CREATE INDEX "banners_userid_idx" ON "banners" ("userid") `,
    );
    await queryRunner.query(`CREATE INDEX "games_id_idx" ON "games" ("id") `);
    await queryRunner.query(
      `CREATE INDEX "games_SOUNDEX_idx" ON "games" ("SOUNDEX") `,
    );
    await queryRunner.query(
      `CREATE INDEX "games_last_updated_idx" ON "games" ("last_updated") `,
    );
    await queryRunner.query(
      `CREATE INDEX "games_id_game_title_SOUNDEX_idx" ON "games" ("id", "game_title", "SOUNDEX") `,
    );
    await queryRunner.query(
      `CREATE INDEX "games_game_title_idx" ON "games" ("game_title") `,
    );
    await queryRunner.query(
      `CREATE INDEX "games_release_date_idx" ON "games" ("release_date") `,
    );
    await queryRunner.query(
      `CREATE INDEX "games_platform_idx" ON "games" ("platform") `,
    );
    await queryRunner.query(
      `CREATE INDEX "platforms_id_idx" ON "platforms" ("id") `,
    );
    await queryRunner.query(
      `CREATE INDEX "platforms_alias_idx" ON "platforms" ("alias") `,
    );
    await queryRunner.query(
      `CREATE UNIQUE INDEX "platforms_name_idx" ON "platforms" ("name") `,
    );
  }
}
