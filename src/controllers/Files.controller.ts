import {
  Controller,
  Delete,
  Get,
  Header,
  Param,
  Post,
  Query,
  Res,
  UploadedFiles,
  UseGuards,
  UseInterceptors,
} from '@nestjs/common';
import { FilesInterceptor } from '@nestjs/platform-express';
import { S3Service } from '../services/s3.service';
import { AuthGuard } from '../guards/Auth.guard';
import { ApiQuery } from '@nestjs/swagger';
import { Response } from 'express';

@Controller('secured/images')
export class FilesController {
  constructor(private readonly s3Service: S3Service) {}

  @Post('upload')
  @UseGuards(AuthGuard)
  @UseInterceptors(FilesInterceptor('files', 10))
  @ApiQuery({ name: 'width', type: String, required: false })
  @ApiQuery({ name: 'height', type: String, required: false })
  async upload(
    @UploadedFiles() files: any[],
    @Query() query: { width: string; height: string },
  ) {
    const urls = await this.s3Service.uploadOptimizedImages(files, {
      width: Number(query?.width),
      height: Number(query?.height),
      quality: 100,
    });

    return { keys: urls };
  }

  @Get(':key')
  @Header('Content-Type', 'image/webp')
  async getFile(@Param('key') key: string, @Res() res: Response) {
    const file = await this.s3Service.getFile(key);
    res.end(file);
  }

  @Delete(':key')
  @UseGuards(AuthGuard)
  async delete(@Param('key') key: string) {
    await this.s3Service.deleteFile(key);
    return { success: true };
  }
}
