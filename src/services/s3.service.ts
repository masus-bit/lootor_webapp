import { Inject, Injectable, Logger } from '@nestjs/common';
import axios from 'axios';
import * as aws4 from 'aws4';
import * as sharp from 'sharp';
import { Cache, CACHE_MANAGER } from '@nestjs/cache-manager';
import { v4 as uuidv4 } from 'uuid';
import { DeleteFilesDto } from '../dto/files/DeleteFilesDto';
import { ReturnDataDto } from '../dto/ReturnDataDto';
import { findUniqueElement } from '../utils/unique';

@Injectable()
export class S3Service {
  private readonly endpoint = 'https://lootor.storage.yandexcloud.net';
  private readonly accessKeyId = process.env.YC_ACCESS_KEY;
  private readonly secretAccessKey = process.env.YC_SECRET_KEY;
  private readonly bucketName = process.env.YC_BUCKET_NAME;
  private readonly region = process.env.YC_REGION;
  private readonly cacheTtl = +process.env.REDIS_TTL || 3600;
  private readonly logger = new Logger(S3Service.name);

  constructor(@Inject(CACHE_MANAGER) private cacheManager: Cache) {
    this.testCacheConnection();
  }

  private async testCacheConnection() {
    try {
      await this.cacheManager.set('connection_test', 'works', 36000);
      this.logger.log('Redis cache connected successfully');
    } catch (error) {
      this.logger.error('Failed to connect to Redis cache', error.stack);
    }
  }

  async getFile(key: string): Promise<Buffer> {
    const cached = await this.cacheManager.get<Buffer>(key);
    if (cached) return cached;
    const file = await this.downloadFromS3(key);
    try {
      console.log(key, 'not cached');
      await this.cacheManager.set(key, file, 3600000);
    } catch (e) {
      console.log(e);
    }
    return file;
  }

  private async downloadFromS3(key: string): Promise<Buffer> {
    const signed = await this.signS3Request('GET', `/images/${key}`);
    const response = await axios.get(signed.url, {
      headers: signed.headers,
      responseType: 'arraybuffer',
    });
    return Buffer.from(response.data);
  }

  private async signS3Request(
    method: string,
    path: string,
    headers: Record<string, string> = {},
    body: string | Buffer = '',
  ) {
    const opts = {
      host: 'lootor.storage.yandexcloud.net',
      path,
      method,
      headers: {
        ...headers,
        'x-amz-content-sha256': 'UNSIGNED-PAYLOAD',
        'x-amz-date': new Date().toISOString().replace(/[:-]|\.\d{3}/g, ''),
      },
      body,
      service: 's3',
      region: this.region,
    };

    const signed = aws4.sign(opts, {
      accessKeyId: this.accessKeyId,
      secretAccessKey: this.secretAccessKey,
    });

    return {
      url: `${this.endpoint}${path}`,
      headers: signed.headers,
      method,
      data: body,
    };
  }

  async uploadFile(file: Buffer, key: string, contentType: string) {
    const path = `/${key}`;
    const signed = await this.signS3Request('PUT', path, {
      'Content-Type': contentType,
    });

    const response = await axios.put(signed.url, file, {
      headers: signed.headers,
    });

    return response.data;
  }

  async uploadOptimizedImages(
    files: any,
    options: { width?: number; height?: number; quality?: number } = {},
  ) {
    let keys: string[] = [];

    for (const file of files) {
      const key = await this.uploadOptimizedImage(file.buffer, uuidv4(), {
        width: Number(options?.width),
        height: Number(options?.height),
        quality: 70,
      });
      keys.push(key);
    }

    return await Promise.all(keys);
  }

  async uploadOptimizedImage(
    file: Buffer,
    filename: string,
    options: { width?: number; height?: number; quality?: number } = {},
  ) {
    const { width, height, quality } = options;
    let optimizedImage;
    optimizedImage = await sharp(file)
      .resize(width || 1920, height || null)
      .webp({ quality })
      .toBuffer();

    const key = `images/${Date.now()}-${filename.replace(
      /\.[^/.]+$/,
      '',
    )}.webp`;

    await this.uploadFile(optimizedImage, key, 'image/webp');

    return `${key}`;
  }

  async deleteFile(dto: DeleteFilesDto): Promise<ReturnDataDto> {
    const result = [];
    for (const key of dto.keys) {
      const path = `/${key}`;
      const signed = await this.signS3Request('DELETE', path);

      await axios.delete(signed.url, {
        headers: signed.headers,
      });
      result.push(key);
    }
    if (result?.length === dto.keys?.length) {
      return new ReturnDataDto({ success: true });
    }
    const notFound = findUniqueElement(dto.keys, result);

    return new ReturnDataDto({
      success: false,
      error: `not found keys: ${notFound.join(',')}`,
    });
  }
}
