import {
  Controller,
  Delete,
  Param,
  Post,
  Query,
  UploadedFile,
  UseGuards,
  UseInterceptors,
} from '@nestjs/common';
import { FileInterceptor } from '@nestjs/platform-express';
import { S3Service } from '../services/s3.service';
import { AuthGuard } from '../guards/Auth.guard';
import { ApiQuery } from '@nestjs/swagger';

@Controller('secured/files')
export class FilesController {
  constructor(private readonly s3Service: S3Service) {}

  @Post('upload')
  @UseGuards(AuthGuard)
  @UseInterceptors(FileInterceptor('file'))
  @ApiQuery({ name: 'width', type: String, required: false })
  @ApiQuery({ name: 'height', type: String, required: false })
  async upload(
    @UploadedFile() file: any,
    @Query() query: { width: string; height: string },
  ) {
    const url = await this.s3Service.uploadOptimizedImage(
      file.buffer,
      file.originalname,
      {
        width: Number(query?.width),
        height: Number(query?.height),
        quality: 100,
      },
    );

    return { url };
  }

  @Delete(':key')
  @UseGuards(AuthGuard)
  async delete(@Param('key') key: string) {
    await this.s3Service.deleteFile(key);
    return { success: true };
  }
}
