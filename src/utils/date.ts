export const getUnixDate = () => Math.floor(+new Date() / 1000);

export const getIsoDate = () => new Date().toISOString() as unknown as Date;
