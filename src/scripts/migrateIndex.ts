async function migrateIndex(
  elasticsearchService,
  indexName,
  settings,
  mappings,
) {
  const newIndexName = `${indexName}_v2`;
  const aliasName = indexName;

  try {
    const newIndexExists = await elasticsearchService.client.indices.exists({
      index: newIndexName,
    });

    if (newIndexExists) {
      await elasticsearchService.client.indices.delete({ index: newIndexName });
      console.log(`Старый индекс ${newIndexName} удален.`);
    }

    await elasticsearchService.client.indices.create({
      index: newIndexName,
      body: {
        settings,
        mappings,
      },
    });

    console.log(`Новый индекс ${newIndexName} создан.`);

    let scrollId = null;
    const bulkActions = [];

    const initialScrollSearch = await elasticsearchService.client.search({
      index: indexName,
      scroll: '1m',
      size: 100,
      body: {
        query: {
          match_all: {},
        },
      },
    });

    scrollId = initialScrollSearch._scroll_id;
    let documents = initialScrollSearch.hits.hits;

    if (documents.length > 0) {
      bulkActions.push(
        ...documents.flatMap((doc) => [
          { index: { _index: newIndexName, _id: doc._id } },
          doc._source,
        ]),
      );
    }

    while (documents.length > 0) {
      const scrollSearch = await elasticsearchService.client.scroll({
        scroll_id: scrollId,
        scroll: '1m',
      });

      documents = scrollSearch.hits.hits;
      if (documents.length > 0) {
        bulkActions.push(
          ...documents.flatMap((doc) => [
            { index: { _index: newIndexName, _id: doc._id } },
            doc._source,
          ]),
        );
      }

      scrollId = scrollSearch._scroll_id;
    }

    if (bulkActions.length > 0) {
      await elasticsearchService.client.bulk({
        body: bulkActions,
      });
    }

    console.log(`Данные перенесены в новый индекс ${newIndexName}.`);

    const oldIndexExists = await elasticsearchService.client.indices.exists({
      index: indexName,
    });

    if (oldIndexExists) {
      await elasticsearchService.client.indices.delete({ index: indexName });
      console.log(`Старый индекс ${indexName} удален.`);
    } else {
      console.log(`Старый индекс ${indexName} не существует.`);
    }

    await elasticsearchService.client.indices.updateAliases({
      body: {
        actions: [{ add: { index: newIndexName, alias: aliasName } }],
      },
    });

    console.log(`Алиас ${aliasName} обновлен.`);
  } catch (error) {
    console.error(`Ошибка при миграции индекса ${indexName}:`, error);
  }
}

// async function migrateAllIndices() {
//   const elasticsearchService = new ElasticsearchService();
//
//   try {
//     for (const [indexName, { settings, mappings }] of Object.entries(indices)) {
//       await migrateIndex(elasticsearchService, indexName, settings, mappings);
//     }
//   } catch (error) {
//     console.error('Ошибка при миграции индексов:', error);
//   } finally {
//   }
// }
//
// migrateAllIndices();
