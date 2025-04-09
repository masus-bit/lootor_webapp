export const findUniqueElement = (arr1: string[], arr2: string[]) => {
  const uniqueInArr1 = arr1.filter((item) => !arr2.includes(item));

  const uniqueInArr2 = arr2.filter((item) => !arr1.includes(item));

  return [...uniqueInArr1, ...uniqueInArr2];
};
