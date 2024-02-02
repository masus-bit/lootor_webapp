import { Collection } from '../entities/Collection';

export const defineShareString = (
  authUserLogin: string,
  collectionUserLogin: string,
  collection: Collection,
) => {
  if (authUserLogin === collectionUserLogin && collection.is_private) {
    return collection.share_string;
  } else {
    return (collection.share_string = '');
  }
};
